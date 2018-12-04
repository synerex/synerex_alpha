package main

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"github.com/synerex/synerex_alpha/sxutil"
	"github.com/synerex/synerex_alpha/api/routing"
	"time"

	"context"
	"log"
	"sync"
)

var (
	MaxDistance = 10000.0 // 10 km for max distance for onemile mobility
	cltRide, cltRoute *sxutil.SMServiceClient
	routingMap = make(map[uint64]chan *routing.RoutingService)
	routingMu sync.RWMutex
	supplyMap = make(map[uint64]chan *rideshare.RideShare)
	supplyMu sync.RWMutex
	statusStr = []string{"free","pickup","ride","full"}
)

func rideshareToMission(share *rideshare.RideShare) *mission{
	mst := &mission{}
	mst.MissionId = "mis01"
	mst.Title = "お迎え"
	mst.Detail = "XXから相見駅"
	evs := make([]event,0)
	for j, r := range share.Routes{
		rts := make([][2]float64, len(r.Points))
		for i, p := range r.Points {
			rts[i][0]=p.GetLongitude()
			rts[i][1]=p.GetLatitude()
		}

		ev := event{
			EventId: fmt.Sprintf("evt%02d",j+1),
			EventType: statusStr[int(r.StatusType)],
			StartTime: r.GetDepartTime().GetTimestamp().GetSeconds()*1000,
			EndTime: r.GetArriveTime().GetTimestamp().GetSeconds()*1000,
			Route:rts,
		}
		evs = append(evs,ev)
	}
	mst.Events = evs

	return mst
}


func onemileHandleRoutingSupply(clt  *sxutil.SMServiceClient, sp *api.Supply) {
	rs := sp.GetArg_RoutingService()
	if rs == nil {
		log.Printf("Unknown Routing Supply %v",sp)
		return
	}

	if sp.TargetId != 0 { // ProposeSupply!
		routingMu.Lock()
		ch, ok :=routingMap[sp.TargetId]
		if ok {
			delete(routingMap,sp.TargetId)
			ch <- rs
		}else{
			log.Printf("Can't find TargetID %d in mapSize %d", sp.TargetId, len(routingMap))
		}
		routingMu.Unlock()
	}else{
		log.Printf("No target Routing Supply %v",sp)
	}
}


func onemileHandleRideShareDemand(clt *sxutil.SMServiceClient, dm *api.Demand) {
	if dm.TargetId != 0 { // select supply!
		// SelectSupply.
		supplyMu.Lock()
		spch, ok := supplyMap[dm.GetTargetId()]
		if ok {
			delete(supplyMap,dm.GetTargetId())
			spch <- dm.GetArg_RideShare()
			clt.Confirm(sxutil.IDType(dm.GetId()))
		}else{
			log.Printf("Can't find select id")
		}
		supplyMu.Unlock()

	} else {
		// First we need to check the distance of dx,dy:
		// 0. Find Nearest free car.
		// 1. send RoutingDemand to Simple-routing
		// 2. recv RoutingSupply, then send ProposeSupply to RideShare.


		rs := dm.GetArg_RideShare()
		fpt := 	rs.DepartPoint.GetCentralPoint()
		tpt := rs.ArrivePoint.GetCentralPoint()
		dist,_ := fpt.Distance(tpt)
		log.Printf("Onemile Get Demand for distance %f: %s", dist, dm.ArgJson)
		if dist >= MaxDistance {
			log.Printf("Onemile ignores long distance demand")
			return
		}

		//  now we find which car should support.
		// currently , onemile do not support sharing.

		minDist := MaxDistance
		var selVc *vehicle
		var selPt *common.Point
		for _,v := range vehicleMap {
			if v.Status != "free" {
				continue
			}
			npt := new(common.Point)
			npt.Latitude = v.Coord[0]
			npt.Longitude = v.Coord[1]
			dt,_ := fpt.Distance(npt)
			log.Printf("Dist %f, %f",dt, minDist)
			if dt < minDist {
				minDist = dt
				selVc = v
				selPt = npt
			}
		}
		if selVc == nil {
			log.Printf("Can't find free onemile vehicle.")
			return
		}

		// routing from current position to start point
		rt0 := &routing.RoutingService{
			DepartPlace: common.NewPlace().WithPoint(selPt),
			ArrivePlace: common.NewPlace().WithPoint(fpt),
			DepartTime: common.NewTime().WithTimestamp(ptypes.TimestampNow()),
		}
		route_dmo := sxutil.DemandOpts{
			Name: "Routing Demand",
			JSON: "{from, to}",
			RoutingService: rt0,
		}
//		selVc.Status = "waitSelect0" // now we book the vehicle

		r0id := cltRoute.RegisterDemand(&route_dmo)
		r0ch := make(chan *routing.RoutingService)
		routingMu.Lock()
		routingMap[r0id] = r0ch
		routingMu.Unlock()
		var rs0 *routing.RoutingService
		select {
		case <-time.After(30 *time.Second):
			log.Printf("OneMile Routing Timeout")
			routingMu.Lock()
			delete(routingMap, r0id) // remove ID
			routingMu.Unlock()
//			selVc.Status = "free"
			return
		case rs0 = <-r0ch:
			log.Printf("Got RoutingService Result")
		}

		// now we can go to DepartPlace by using rs0.
		// let's prepare to next Route

		rt := &routing.RoutingService{
			DepartPlace: common.NewPlace().WithPoint(fpt),
			ArrivePlace: common.NewPlace().WithPoint(tpt),
		}
		if rs.DepartTime != nil {
			// we should check the deadline.
			depart , _ := ptypes.Timestamp(rs.DepartTime.GetTimestamp())
			arrive , _ := ptypes.Timestamp(rs0.ArriveTime.GetTimestamp())
			if arrive.After(depart) {// we cant make it... umm
//				selVc.Status = "free"
				return
			}
			// we should consider the time difference (if we have wide gap)
//			log.Printf("Now we have a duration %.1f mins", depart.Sub(arrive).Minutes())

			rt.DepartTime = rs.DepartTime
		}else if rs.ArriveTime != nil {
			rt.ArriveTime = rs.ArriveTime
		}

		dmo := sxutil.DemandOpts{
			Name: "Routing Demand",
			JSON: "{from, to}",
			RoutingService: rt,
		}

		r1id :=clt.RegisterDemand(&dmo)
		r1ch := make(chan *routing.RoutingService)
		routingMu.Lock()
		routingMap[r1id] = r1ch  // save demand
		routingMu.Unlock()
		var rs1 *routing.RoutingService
		select {
		case <-time.After(30 *time.Second):
			log.Printf("OneMile Routing 2nd Timeout")
			routingMu.Lock()
			delete(routingMap, r1id) // remove ID
			routingMu.Unlock()
//			selVc.Status = "free"
			return
		case rs1 = <-r1ch:
			log.Printf("Got RoutingService 2ndResult")
		}

		// we should check the integrity of 2 routes.

		depart2 , _ := ptypes.Timestamp(rs1.DepartTime.GetTimestamp())
		arrive2 , _ := ptypes.Timestamp(rs0.ArriveTime.GetTimestamp())
		if arrive2.After(depart2) {// we cant make it... umm
			log.Printf("Umm. sorry we can't mak it..")
//			selVc.Status = "free"
			return
		}
		// we should consider the time difference (if we have wide gap)
		log.Printf("Now we have a duration %.1f mins for wait", depart2.Sub(arrive2).Minutes())

		// now we have to keep the vehicle to be confirmed.
//		selVc.Status = "waitSelect"

		rsRoute0 := &rideshare.Route{
			TrafficType: rideshare.TrafficType_TAXI,
			StatusType: rideshare.StatusType_PICKUP,
			TransportName :"Onemile",
			DepartPoint: rs0.DepartPlace,
			DepartTime: rs0.DepartTime,
			ArrivePoint: rs0.ArrivePlace,
			ArriveTime: rs0.ArriveTime,
		}

		rsRoute1 := &rideshare.Route{
			TrafficType: rideshare.TrafficType_TAXI,
			StatusType: rideshare.StatusType_RIDE,
			TransportName :"Onemile",
			DepartPoint: rs1.DepartPlace,
			DepartTime: rs1.DepartTime,
			ArrivePoint: rs1.ArrivePlace,
			ArriveTime: rs1.ArriveTime,
		}

		rideShareSvc := &rideshare.RideShare{
			DepartPoint: rs1.DepartPlace,
			DepartTime: rs1.DepartTime,
			ArrivePoint: rs1.ArrivePlace,
			ArriveTime: rs1.ArriveTime,
			Routes: []*rideshare.Route{rsRoute0, rsRoute1},
		}

		spo := &sxutil.SupplyOpts{
			Target: dm.GetId(),
			RideShare: rideShareSvc,
		}

		psid := clt.ProposeSupply(spo)
		rxch := make(chan *rideshare.RideShare)
		supplyMu.Lock()
		supplyMap[psid] = rxch
		supplyMu.Unlock()
		var rdsh *rideshare.RideShare
		select { // wait for select supply
		case <-time.After(30 *time.Second):
			log.Printf("Timeout! 30 seconds")
			return
		case rdsh = <- rxch:
		}
		// now get selectSupply
		// check car availability
		selVc.mu.Lock()
		if selVc.Status == "free" {
			selVc.Status = "pickup"
		// now we need to add event!
			ms := rideshareToMission(rdsh)
			selVc.mission = ms
		}
		selVc.mu.Unlock()

		if selVc.socket != nil {
			selVc.socket.Emit("clt_request_mission", selVc.mission.toMap())
			log.Printf("emit %s: [ payload: %#v]\n", "clt_request_mission", selVc.mission.toMap())
		}

	}
}

func subscribeRouting(rtClient *sxutil.SMServiceClient){
	ctx := context.Background()
	rtClient.SubscribeSupply(ctx, onemileHandleRoutingSupply) // this is on "onemile-routing.go"
	log.Printf("Server closed... on Onemile Routing SubscribeSupply")
}

// Main Entry point of onemile-routing.go
// subscribe rideshare channel
func subscribeRideShare(rdClient, rtClient *sxutil.SMServiceClient) {
	cltRide = rdClient
	cltRoute = rtClient
	go subscribeRouting(rtClient)
	ctx := context.Background()
	rdClient.SubscribeDemand(ctx, onemileHandleRideShareDemand) // this is on "onemile-routing.go"


	log.Printf("Server closed... on Onemile RideShare SubscribeDemand")
}
