package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"github.com/synerex/synerex_alpha/api/routing"
	"github.com/synerex/synerex_alpha/sxutil"
)

var (
	MaxDistance       = 10000.0 // 10 km for max distance for onemile mobility
	cltRide, cltRoute *sxutil.SMServiceClient
	routingMap        = make(map[uint64]chan *routing.RoutingService)
	routingMu         sync.RWMutex
	supplyMap         = make(map[uint64]chan uint64)
	supplyMu          sync.RWMutex
	statusStr         = []string{"free", "pickup", "ride", "full"}
)

func rideshareToMission(share *rideshare.RideShare) *mission {
	mst := &mission{}
	mst.MissionId = "mis01"
	mst.Title = "お迎え"
	mst.Detail = "XXから相見駅"
	evs := make([]event, 0)
	for j, r := range share.Routes {
		//		rts := make([][2]float64, len(r.Points))
		/*		log.Println(j,":Size=", len(r.Points))
				rts := "["
				for i, p := range r.Points {
					rts += strconv.FormatFloat(p.GetLongitude(),'f',-1,64)+","
					rts += strconv.FormatFloat(p.GetLatitude(),'f',-1,64)+"]"
					if i+1 != len(r.Points){
						rts+=","
					}
				}
				rts+="]"
		*/
		rts := make([][2]float64, len(r.Points))
		for i, p := range r.Points {
			rts[i][0] = p.GetLatitude()
			rts[i][1] = p.GetLongitude()
		}

		ev := event{
			EventId:     fmt.Sprintf("evt%02d", j+1),
			EventType:   statusStr[int(r.StatusType)],
			StartTime:   r.GetDepartTime().GetTimestamp().GetSeconds() * 1000,
			EndTime:     r.GetArriveTime().GetTimestamp().GetSeconds() * 1000,
			Destination: r.GetArrivePoint().String(),
			Route:       rts,
			Status:      "none",
		}
		evs = append(evs, ev)
	}
	mst.Events = evs

	return mst
}

func onemileHandleRoutingSupply(clt *sxutil.SMServiceClient, sp *api.Supply) {
	rs := sp.GetArg_RoutingService()
	if rs == nil {
		log.Printf("Unknown Routing Supply %v", sp)
		return
	}

	if sp.TargetId != 0 { // ProposeSupply!
		routingMu.Lock()
		ch, ok := routingMap[sp.TargetId]
		if ok {
			delete(routingMap, sp.TargetId)
			ch <- rs
		} else {
			log.Printf("Can't find TargetID %d in mapSize %d", sp.TargetId, len(routingMap))
		}
		routingMu.Unlock()
	} else {
		log.Printf("No target Routing Supply %v", sp)
	}
}

/*
func onemileHandleRideShareSupply(clt  *sxutil.SMServiceClient, sp *api.Supply) {
	log.Println("Got ProposeSupply", *sp)

	if sp.TargetId != 0 { // select supply!
		// SelectSupply.
//		log.Printf("Got ProposeSupply from %d -> %d, ", sp.SenderId, sp.TargetId%1000000)
	}

	// ignore other rideshare supply
}
*/

func onemileHandleRideShareDemand(clt *sxutil.SMServiceClient, dm *api.Demand) {
	if dm.TargetId != 0 { // select supply!
		// SelectSupply.
		//		log.Println("Got Select Supply!?",*dm)
		log.Printf("Got SelectSupply from %d -> %d, ", dm.SenderId, dm.TargetId%1000000)
		supplyMu.Lock()
		spch, ok := supplyMap[dm.GetTargetId()]
		if ok {
			delete(supplyMap, dm.GetTargetId())
			spch <- dm.GetId()
			log.Printf("Now prepare for Confirm")
		} else {
			log.Printf("RideShare SelectSupply: Can't find select id %d ", dm.TargetId)
		}
		supplyMu.Unlock()
	} else {
		// First we need to check the distance of dx,dy:
		// 0. Find Nearest free car.
		// 1. send RoutingDemand to Simple-routing
		// 2. recv RoutingSupply, then send ProposeSupply to RideShare.
		log.Println("Got Demand", *dm)

		originRideShare := dm.GetArg_RideShare()
		if originRideShare == nil {
			log.Printf("Got nil rideshare! [%v]", dm)
			return
		}
		fpt := originRideShare.DepartPoint.GetCentralPoint()
		tpt := originRideShare.ArrivePoint.GetCentralPoint()
		dist, _ := fpt.Distance(tpt)
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
		for _, v := range vehicleMap {
			if v.Status != "free" {
				continue
			}
			npt := new(common.Point)
			npt.Latitude = v.Coord[0]
			npt.Longitude = v.Coord[1]
			dt, _ := fpt.Distance(npt)
			log.Printf("Dist %f, %f", dt, minDist)
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
			DepartTime:  common.NewTime().WithTimestamp(ptypes.TimestampNow()),
		}
		route_dmo := sxutil.DemandOpts{
			Name:           "Routing Demand",
			JSON:           "{from, to}",
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
		case <-time.After(30 * time.Second):
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
		if originRideShare.DepartTime != nil {
			// we should check the deadline.
			// we should consider the time difference (if we have wide gap)
			//			log.Printf("Now we have a duration %.1f mins", depart.Sub(arrive).Minutes())

			rt.DepartTime = rs0.ArriveTime

			//			rt.DepartTime = originRideShare.DepartTime
		} else if originRideShare.ArriveTime != nil {

			depart, _ := ptypes.Timestamp(originRideShare.DepartTime.GetTimestamp())
			oneDepart, _ := ptypes.Timestamp(rs0.ArriveTime.GetTimestamp())
			if depart.After(oneDepart) { // we cant make it... umm
				//				selVc.Status = "free"
				log.Printf("Onemile can't support this deadline %v -> %v ", depart, oneDepart)
				return
			}

			rt.ArriveTime = originRideShare.ArriveTime
		}

		dmo := sxutil.DemandOpts{
			Name:           "Routing Demand",
			JSON:           "{2ndOnemile Route}",
			RoutingService: rt,
		}

		log.Printf("Now Register Routing Demand for ride")

		r1id := cltRoute.RegisterDemand(&dmo)
		r1ch := make(chan *routing.RoutingService)
		routingMu.Lock()
		routingMap[r1id] = r1ch // save demand
		routingMu.Unlock()
		var rs1 *routing.RoutingService
		select {
		case <-time.After(30 * time.Second):
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

		depart2, _ := ptypes.Timestamp(rs1.DepartTime.GetTimestamp())
		arrive2, _ := ptypes.Timestamp(rs0.ArriveTime.GetTimestamp())
		//		if depart2.After(arrive2) {// we cant make it... umm
		//			log.Printf("Umm. sorry we can't mak it..Arrive %v, 2ndDepart %v",arrive2, depart2)
		//			selVc.Status = "free"
		//			return
		//		}
		// we should consider the time difference (if we have wide gap)
		log.Printf(" 1st route Arrive %v, 2nd route Depart %v", arrive2, depart2)
		log.Printf("Now we have a duration %.1f mins for wait", depart2.Sub(arrive2).Minutes())

		// now we have to keep the vehicle to be confirmed.
		//		selVc.Status = "waitSelect"

		rsRoute0 := &rideshare.Route{
			TrafficType:   rideshare.TrafficType_TAXI,
			StatusType:    rideshare.StatusType_PICKUP,
			TransportName: "Onemile",
			DepartPoint:   rs0.DepartPlace,
			DepartTime:    rs0.DepartTime,
			ArrivePoint:   rs0.ArrivePlace,
			ArriveTime:    rs0.ArriveTime,
			Points:        rs0.Points,
		}

		rsRoute1 := &rideshare.Route{
			TrafficType:   rideshare.TrafficType_TAXI,
			StatusType:    rideshare.StatusType_RIDE,
			TransportName: "Onemile",
			DepartPoint:   rs1.DepartPlace,
			DepartTime:    rs1.DepartTime,
			ArrivePoint:   rs1.ArrivePlace,
			ArriveTime:    rs1.ArriveTime,
			Points:        rs1.Points,
		}

		rideShareSvc := &rideshare.RideShare{
			DepartPoint: rs1.DepartPlace,
			DepartTime:  rs1.DepartTime,
			ArrivePoint: rs1.ArrivePlace,
			ArriveTime:  rs1.ArriveTime,
			Routes:      []*rideshare.Route{rsRoute0, rsRoute1},
		}

		spo := &sxutil.SupplyOpts{
			Target:    dm.GetId(),
			RideShare: rideShareSvc,
			Name:      "FromOneMile!",
			JSON:      "{2routes}",
		}

		// start
		log.Println("Now Propose Supply :", *spo)
		psid := clt.ProposeSupply(spo)
		log.Printf("Propose Supply ID: %d", psid)

		rxch := make(chan uint64)
		supplyMu.Lock()
		supplyMap[psid] = rxch
		supplyMu.Unlock()

		go proposeSupplyForRouting(psid, rxch, selVc, rideShareSvc)
	}
}

func proposeSupplyForRouting(psid uint64, rxch chan uint64, selVc *vehicle, rs *rideshare.RideShare) {
	var cfid uint64
	select { // wait for select supply
	case <-time.After(300 * time.Second):
		log.Printf("Onemile Propose Supply Timeout! 300 seconds %d", psid)
		//should  remvove!
		return
	case cfid = <-rxch:
		log.Printf("Got SelectSupply for Onemile!")
	}
	// now get selectSupply
	// check car availability
	selVc.mu.Lock()
	if selVc.Status == "free" {
		log.Printf("Now Book a vehicle! [%s] by %d", selVc.VehicleId, cfid)
		cltRide.Confirm(sxutil.IDType(cfid))

		selVc.Status = "pickup"
		// now we need to add event!
		ms := rideshareToMission(rs)
		selVc.Mission = ms
		if selVc.socket != nil {
			selVc.socket.Emit("clt_request_mission", selVc.Mission.toMap())
			//			log.Printf("emit %s: [ payload: %#v]\n", "clt_request_mission", selVc.mission.toMap())
		}
		log.Printf("emit %s: [ payload: %#v]\n", "clt_request_mission", selVc.Mission.toMap())
		//		log.Println(selVc.mission)
		//		buf,_ := json.Marshal(selVc.mission)
		//		log.Println(string(buf))
	} else {
		// not confirm! sorry
		log.Printf("Cannot book a vehicle! [%s]", selVc.VehicleId)
	}
	selVc.mu.Unlock()

}

func subscribeRouting(rtClient *sxutil.SMServiceClient) {
	ctx := context.Background()
	rtClient.SubscribeSupply(ctx, onemileHandleRoutingSupply) // this is on "onemile-routing.go"
	log.Printf("Server closed... on Onemile Routing SubscribeSupply")
}

/*
func subscribeRideShareDemand(rdSpClient *sxutil.SMServiceClient){
	ctx := context.Background()
	rdSpClient.SubscribeSupply(ctx, onemileHandleRideShareSupply)
	log.Printf("Server closed... on Onemile RideShare SubscribeSupply")
}
*/

// Main Entry point of onemile-routing.go
// subscribe rideshare channel
func subscribeRideShare(rdClient, rtClient *sxutil.SMServiceClient) {
	cltRide = rdClient
	cltRoute = rtClient

	//	go subscribeRideShareDemand(rdSpClient)

	go subscribeRouting(rtClient)
	ctx := context.Background()
	rdClient.SubscribeDemand(ctx, onemileHandleRideShareDemand) // this is on "onemile-routing.go"

	log.Printf("Server closed... on Onemile RideShare SubscribeDemand")
}
