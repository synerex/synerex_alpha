package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"log"
	"math"
	"sync"
	"time"
)

var (
	serverAddr                    = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv                                = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	supplyMap        map[uint64]chan uint64 // for confirm
	supplyMu		sync.Mutex
	rideShareMap      map[uint64]chan *api.Supply // subChannels
	rideShareMu       sync.RWMutex
	TrainTrainsitTime  = 5 * time.Minute //  5 min
	BusTrainsitTime    = 2 * time.Minute // 2 min

	demandClient *sxutil.SMServiceClient
	supplyClient *sxutil.SMServiceClient

)

func init(){
	rideShareMap = make(map[uint64]chan *api.Supply )//
	supplyMap = make(map[uint64]chan uint64)
}


func getOnemileRoute(clt *sxutil.SMServiceClient, from *common.Point, to *common.Point, ft *common.Time, tt *common.Time)  uint64{
	rs := new(rideshare.RideShare)
	rs.DepartPoint = common.NewPlace().WithPoint(from)
	rs.DepartTime = ft
	rs.ArrivePoint = common.NewPlace().WithPoint(to)
	rs.ArriveTime = tt

	dmo := sxutil.DemandOpts{
		Name: "OneMileDemand",
		JSON: "",
		RideShare:rs,
	}
	id :=clt.RegisterDemand(&dmo)
	return id
}

// start thinking from Train
func trainAndOnemile(clt *sxutil.SMServiceClient, dm *api.Demand) {


	rsInfo := dm.GetArg_RideShare()
	if rsInfo == nil {
		log.Printf("Demand is not for RideShare! [%v]", dm)
		return
	}
	var dp,ap *common.Place
	dp = rsInfo.DepartPoint
	dpt := dp.GetCentralPoint()
	ap = rsInfo.ArrivePoint
	apt := ap.GetCentralPoint()

	// to Nagoya or from Nagoya

	// find close point to aimi station.
	aimiPt := new(common.Point)
	aimiPt.Latitude = 34.887573
	aimiPt.Longitude = 137.160248

	ddist , _ := aimiPt.Distance(dpt)
	adist , _ := aimiPt.Distance(apt)

	var dt, at *common.Time
	dt = rsInfo.DepartTime
	at = rsInfo.ArriveTime


	if at == nil {// departure time only
		if ddist < adist { // to aimi station (end is Nagoya)
			// onemile -> aimi ->
			ch := make(chan *api.Supply)
			rideShareMu.Lock()
			id  := getOnemileRoute(clt, dpt, aimiPt, dt, nil)
			//
			log.Printf("TrainAndOneMile Now set channel ID %d",id)
			rideShareMap[id]= ch
			rideShareMu.Unlock()
			var rs *rideshare.RideShare
			var rssp *api.Supply
			log.Printf("Wait for onemile reply! %v", dpt)
			select {
			case <-time.After(30 *time.Second ):
				log.Printf("Timeout 30 seconds for getting OneMileRoute")
				return
			case rssp = <-ch:
				rs = rssp.GetArg_RideShare()
				log.Printf("Get OnemleRoute")
			}
			avTime, _ := ptypes.Timestamp(rs.ArriveTime.GetTimestamp())
			avTime.Add(TrainTrainsitTime) //
				route := getTrainRouteFromTime(aimiPt,apt, int32(avTime.Minute()+ avTime.Hour()*60), 0)
				var rdSh rideshare.RideShare

			if(route == nil) {
				log.Printf("Train route is nil..umm")
				// now we found 2 routes( for onemile, and train)
				rdSh = rideshare.RideShare{
					Routes: []*rideshare.Route{rs.Routes[0], rs.Routes[1]}, ///route},
				}
			}else{
				rdSh = rideshare.RideShare{
					Routes: []*rideshare.Route{rs.Routes[1],route},
				}

			}
			spo := sxutil.SupplyOpts{
				Target: dm.GetId(),
				JSON:"{SPOPT0}",
				RideShare: &rdSh,
			}
			spch := make(chan uint64)
			spid :=	clt.ProposeSupply(&spo) // rideshare Demand
			supplyMap[spid] = spch
			select {
			case <-time.After(30 * time.Second):
				log.Printf("Timeout 30sec")
				return
			case id:= <-spch: // receive SelectSupply!
				//
				log.Printf("Select from Kota: %d",id)
				// now we need to selectSupply for OneMile.

				mbusid, mberr := supplyClient.SelectSupply(rssp)
				log.Printf("SelectSupply to  %d",mbusid)
				if mberr == nil {
					log.Printf("SelectSupply Success! and MbusID=%d", mbusid)
					clt.Confirm(sxutil.IDType(id))
					log.Printf("Finally confirm %d ", id)

				}else{
					log.Printf("%v error:",mberr)
				}

			}


		}else{ // end is Kota
			//find Nagoya to Aimi train route
			dtTS := dt.GetTimestamp()
			dtTime, _ := ptypes.Timestamp(dtTS)
			route := getTrainRouteFromTime(dpt, aimiPt,int32(dtTime.Minute()+ dtTime.Hour()*60), 0)
			log.Printf("Train route is ...%v",route)
			arTM := route.ArriveTime
			arTS := arTM.GetTimestamp()
			arTime , _ := ptypes.Timestamp(arTS)
			arTime.Add(TrainTrainsitTime) //
			arTimePro , _ := ptypes.TimestampProto(arTime)
			arTTime := common.NewTime().WithTimestamp(arTimePro)

			rdch := make(chan *api.Supply)
			rideShareMu.Lock()
			id := getOnemileRoute(clt, aimiPt, apt, arTTime, nil)
			log.Printf("TrainAndOneMile Reverse Now set channel ID %d",id)

			rideShareMap[id]= rdch
			rideShareMu.Unlock()
			var rs *rideshare.RideShare
			var rssp *api.Supply
			select {
			case <-time.After(30 *time.Second ):
				log.Printf("Timeout 30 seconds for getting OneMileRoute")
				return
			case rssp = <-rdch:
				rs = rssp.GetArg_RideShare()
				log.Printf("Get OnemileRoute")
			}

			// now we found 2 routes( for onemile, and train)
			var rdSh *rideshare.RideShare
			if rs.Routes[1]== nil{
				log.Println("rs.Route", rs.Routes)
				rdSh = &rideshare.RideShare{
					Routes: []*rideshare.Route{route},
				}

			}else {
				log.Println("routes:", route, rs.Routes[1])
				rdSh = &rideshare.RideShare{
					Routes: []*rideshare.Route{route, rs.Routes[1]},
				}
			}
			spo := sxutil.SupplyOpts{
				Target: dm.GetId(),
				JSON:"{SPOPT1}",
				RideShare: rdSh,
			}
			spch := make(chan uint64)
			spid :=	clt.ProposeSupply(&spo)
			supplyMap[spid] = spch
			select {
			case <-time.After(30 * time.Second):
				log.Printf("TrainAndOnemile propose Timeout 30sec")
				return
			case id:= <-spch: // receive SelectSupply!
				log.Printf("Select from Kota 2nd:%d",id)
				mbusid, mberr := supplyClient.SelectSupply(rssp)
				if mberr == nil {
					log.Printf("Select Supply2nd Success and Mbus %d",mbusid)
					clt.Confirm(sxutil.IDType(id))
					log.Printf("finally confirm 2nd %d",id)
				}else{
					log.Printf("err %v",mberr)
				}
			}

		}

	}else { // arrival time only...
		log.Printf("arrival time.. is not yet implemented")
	}


}


/*

func  trainAndBusAndOnemile(clt *sxutil.SMServiceClient, dm *api.Demand) {

	rsInfo := dm.GetArg_RideShare()
	if rsInfo == nil {
		log.Printf("Demand is not for RideShare! [%v]", dm)
		return
	}
	var dp,ap *common.Place
	dp = rsInfo.DepartPoint
	dpt := dp.GetCentralPoint()
	ap = rsInfo.ArrivePoint
	apt := ap.GetCentralPoint()

	// to Nagoya or from Nagoya

	// find close point to aimi station.
	aimiPt := new(common.Point)
	aimiPt.Latitude = 34.887573
	aimiPt.Longitude = 137.160248

	ddist , _ := aimiPt.Distance(dpt)
	adist , _ := aimiPt.Distance(apt)

	var dt, at *common.Time
	dt = rsInfo.DepartTime
	at = rsInfo.ArriveTime

	if at == nil {// departure time only
		if ddist < adist { // to aimi station (end is Nagoya)
			// onemile -> aimi ->
			dtTS := dt.GetTimestamp()
			dtTime, _ := ptypes.Timestamp(dtTS)
			dtTime = dtTime.Add(9*time.Hour)
			log.Printf("Time is %d:%d ",dtTime.Hour(), dtTime.Minute())

			routeBus := getBusRouteFromTime(dpt, aimiPt,int32(dtTime.Minute()+ dtTime.Hour()*60), 0)
			if routeBus == nil {
				log.Println("Can'f find bus time...",dpt, "->", aimiPt)
				return
			}else{
				log.Println("Got bus...",dpt, "->", aimiPt, routeBus)
			}
// now bus time is decided.
//			lets ask Onemile to come with this bus.
			dpTm := routeBus.DepartTime
			dpTs := dpTm.GetTimestamp()
			dpTime , _ := ptypes.Timestamp(dpTs)
			dpTime = dpTime.Add( - BusTrainsitTime)
			dTimePro , _ := ptypes.TimestampProto(dpTime)
			oneMileArriveTime :=common.NewTime().WithTimestamp(dTimePro)

			ch := make(chan *api.Supply)
			rideShareMu.Lock()
//			id  := getOnemileRoute(clt, dpt, aimiPt, dt, oneMileArriveTime)
			id  := getOnemileRoute(clt, dpt, aimiPt, nil, oneMileArriveTime)
			rideShareMap[id]= ch
			rideShareMu.Unlock()
			// select for rideShare.

			var rs *rideshare.RideShare
			var rssp *api.Supply
			select {
			case <-time.After(30 *time.Second ):
				log.Printf("Timeout 30 seconds for getting OneMileRoute")
				return
			case rssp = <-ch:
				rs = rssp.GetArg_RideShare()
				log.Printf("Get OnemileRoute for BSTR")
			}

//			avTime, _ := ptypes.Timestamp(rs.ArriveTime.GetTimestamp())

//
			avTime, _ := ptypes.Timestamp(routeBus.ArriveTime.GetTimestamp())


			avTime.Add(TrainTrainsitTime) //
			route := getTrainRouteFromTime(aimiPt,apt, int32(avTime.Minute()+ avTime.Hour()*60), 0)
			var rdSh rideshare.RideShare

			if(route == nil) {
				log.Printf("Train route is nil..umm")
				// now we found 2 routes( for onemile, and train)
				rdSh = rideshare.RideShare{
					Routes: []*rideshare.Route{rs.Routes[1], routeBus}, ///route},
				}
			}else{
				rdSh = rideshare.RideShare{
					Routes: []*rideshare.Route{rs.Routes[1],routeBus, route},
				}

			}
			spo := sxutil.SupplyOpts{
				Target: dm.GetId(),
				JSON:"{SPOPT0}",
				RideShare: &rdSh,
			}
			spch := make(chan uint64)
			spid :=	clt.ProposeSupply(&spo) // rideshare Demand
			supplyMap[spid] = spch
			select {
			case <-time.After(30 * time.Second):
				log.Printf("Timeout 30sec")
				return
			case id:= <-spch: // receive SelectSupply!
				//
				log.Printf("Select from Kota: %d",id)
				// now we need to selectSupply for OneMile.

				mbusid, mberr := supplyClient.SelectSupply(rssp)
				log.Printf("SelectSupply to  %d",mbusid)
				if mberr == nil {
					log.Printf("SelectSupply Success! and MbusID=%d", mbusid)
					clt.Confirm(sxutil.IDType(id))
					log.Printf("Finally confirm %d ", id)

				}else{
					log.Printf("%v error:",mberr)
				}

			}




		}else{ // end is Kota So, Train -> Bus -> OneMile
			//find Nagoya to Aimi train route
			dtTS := dt.GetTimestamp()
			dtTime, _ := ptypes.Timestamp(dtTS)
			route := getTrainRouteFromTime(dpt, aimiPt,int32(dtTime.Minute()+ dtTime.Hour()*60), 0)
			arTM := route.ArriveTime
			arTS := arTM.GetTimestamp()
			arTime , _ := ptypes.Timestamp(arTS)
			arTime.Add(TrainTrainsitTime) //
			arTimePro , _ := ptypes.TimestampProto(arTime)
			arTTime := common.NewTime().WithTimestamp(arTimePro)

			rdch := make(chan *api.Supply)
			rideShareMu.Lock()
			id := getOnemileRoute(clt, aimiPt, apt, arTTime, nil)
			rideShareMap[id]= rdch
			rideShareMu.Unlock()

			var rs *rideshare.RideShare
			var rssp *api.Supply
			select {
			case <-time.After(30 *time.Second ):
				log.Printf("Timeout 30 seconds for getting OneMileRoute")
				return
			case rssp = <-rdch:
				rs = rssp.GetArg_RideShare()
				log.Printf("Get OnemileRoute for BSTR")
			}

			//			avTime, _ := ptypes.Timestamp(rs.ArriveTime.GetTimestamp())

			avTime, _ := ptypes.Timestamp(rs.ArriveTime.GetTimestamp())

			routeBus := getBusRouteFromTime(dpt, aimiPt,int32(avTime.Minute()+ avTime.Hour()*60), 0)
			if routeBus == nil {
				log.Println("Can'f find bus time...",dpt, "->", aimiPt)
				return
			}else{
				log.Println("Got bus...",dpt, "->", aimiPt, routeBus)
			}
		}

	}else { // arrival time only...
		log.Printf("arrival time.. is not yet implemented")
	}

}
*/
func NewPoint(lat, lon float64) *common.Point{
	pt := new(common.Point)
	pt.Latitude = lat
	pt.Longitude = lon
	return pt
}


// check which place
func findCurrentPlace( mp *common.Point ) int {
	ldb := []*common.Point{
		NewPoint(34.892948, 137.188032),
		NewPoint(34.8770639, 137.1892827),
		NewPoint(34.876837, 137.168259),
	}

	dist := 10000.0

	ix := -1
	for i, p := range ldb {
		d,_:= mp.Distance(p)
		if d < dist {
			dist = d
			ix = i
		}
	}
	if dist < 5000 {
		return ix
	}
	return -1
}




// goto Nagoya?
// check total distance
// if there is no big distance than 15km we just use bus and onemile
func checkTrainDest (rsInfo *rideshare.RideShare) (bool, bool){
	const maxDist = 15000.0
	var dp,ap *common.Place
	dp = rsInfo.DepartPoint
	dpt := dp.GetCentralPoint()
	ap = rsInfo.ArrivePoint
	apt := ap.GetCentralPoint()

	// to Nagoya or from Nagoya
	// find closest point to aimi station.
	aimiPt := new(common.Point)   //
	aimiPt.Latitude = 34.887573
	aimiPt.Longitude = 137.160248

	ddist , _ := aimiPt.Distance(dpt)
	adist , _ := aimiPt.Distance(apt)

	useTrain := false
	if math.Max(ddist, adist) > maxDist {
		useTrain = true
	}
	if ddist > adist {
		return true, useTrain
	}
	return false, useTrain
}

// for each userID, we fix the route
//
func  expSpecial(clt *sxutil.SMServiceClient, dm *api.Demand) {
	var rdSh rideshare.RideShare
	rsInfo := dm.GetArg_RideShare()
	if rsInfo == nil {
		log.Printf("Demand is not for RideShare! [%v]", dm)
		return
	}
	toStation, useTrain := checkTrainDest(rsInfo)

	var fpt, tpt *common.Point
	var brt ,trt  *rideshare.Route
	var subj int
	dptTime := rsInfo.DepartTime
	if toStation { // to Aimi
		fpt = 	rsInfo.GetDepartPoint().GetCentralPoint()
		subj = findCurrentPlace(fpt)
	}else{
		tpt = 	rsInfo.GetArrivePoint().GetCentralPoint()
		subj = findCurrentPlace(tpt)
	}
	brt = getBusRoute(toStation, subj)
	trt = getTrainRouteFromTimeExp(toStation)

	if !toStation {
		fpt = brt.GetArrivePoint().GetCentralPoint()
		arTM := brt.ArriveTime
		arTS := arTM.GetTimestamp()
		arTime , _ := ptypes.Timestamp(arTS)
		arTime.Add(BusTrainsitTime) //
		arTimePro , _ := ptypes.TimestampProto(arTime)
		dptTime = common.NewTime().WithTimestamp(arTimePro)
	}else{
		tpt = brt.GetDepartPoint().GetCentralPoint()
	}

	ch := make(chan *api.Supply)
	rideShareMu.Lock()
	id  := getOnemileRoute(clt, fpt, tpt, dptTime, nil)
	rideShareMap[id]= ch
	rideShareMu.Unlock()

	var rs *rideshare.RideShare
	var rssp *api.Supply
	select {
	case <-time.After(30 *time.Second ):
		log.Printf("Timeout 30 seconds for getting OneMileRoute")
		return
	case rssp = <-ch:
		rs = rssp.GetArg_RideShare()
		log.Printf("Get OnemileRoute for BSTR")
	}

	if useTrain {
		if toStation {
			rdSh = rideshare.RideShare{
				Routes: []*rideshare.Route{rs.Routes[1],brt, trt },
			}
		}else{
			rdSh = rideshare.RideShare{
				Routes: []*rideshare.Route{ trt, brt, rs.Routes[1] },
			}
		}
	}else{
		if toStation {
			rdSh = rideshare.RideShare{
				Routes: []*rideshare.Route{rs.Routes[1],brt },
			}
		}else {
			rdSh = rideshare.RideShare{
				Routes: []*rideshare.Route{ brt, rs.Routes[1]},
			}
		}
	}

	spo := sxutil.SupplyOpts{
		Target: dm.GetId(),
		JSON:"{SPOPT0}",
		RideShare: &rdSh,
	}

	spch := make(chan uint64)
	spid :=	clt.ProposeSupply(&spo) // rideshare Demand
	supplyMap[spid] = spch
	select {
	case <-time.After(300 * time.Second):
		log.Printf("Timeout 300sec")
		return
	case id:= <-spch: // receive SelectSupply!
		//
		log.Printf("Select from Kota: %d",id)
		// now we need to selectSupply for OneMile.

		mbusid, mberr := supplyClient.SelectSupply(rssp)
		log.Printf("SelectSupply to  %d",mbusid)
		if mberr == nil {
			log.Printf("SelectSupply Success! and MbusID=%d", mbusid)
			clt.Confirm(sxutil.IDType(id))
			log.Printf("Finally confirm %d ", id)

		}else{
			log.Printf("%v error:",mberr)
		}

	}



}




func rideshareSupplyCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	if sp.TargetId == 0{
//		log.Printf("Should not come here...")
		return
	}
	log.Printf("Got RideShare Supply from (may from Onemile) %v",*sp)
	rt := sp.GetArg_RideShare()
	if rt != nil { // get Routing supplu
		rideShareMu.RLock()
		if ch, ok := rideShareMap[sp.TargetId]; ok{

			delete(rideShareMap, sp.TargetId)

			log.Printf("Send SelectSupply %d", sp.Id)
			ch <- sp // .GetArg_RideShare()
			// we have to wait for Select by User.

//			id, err := clt.SelectSupply(sp)
//			if err == nil {
//				log.Printf("SelectSupply Success! and MbusID=%d", id)
//			}

		}else{
			log.Printf("Supply Callback: Can't find report channel for %d",sp.TargetId)
		}
		rideShareMu.RUnlock()
	}else{
		log.Printf("Routing:not for rideshare %v",*sp)
	}
}

// callback for each Demand
func rideshareDemandCallback(clt *sxutil.SMServiceClient, dm *api.Demand) {
	// check if demand is match with my supply.
	log.Println("Got rideshare demand callback on Routing:",dm)

	// we need to start a new "Multiple Routing Suggestion".
	// for each Demand, we start go routine for that!

	// currently fix the idea for multi routing

	if dm.TargetId != 0 { // this is Select!
		log.Println("getSelect for Route!")
		// which select ?
		supplyMu.Lock()
		if ch, ok := supplyMap[dm.TargetId]; ok {
			log.Printf("Send Confirm ! to %d", dm.TargetId)
			 // we should start contract with them..
			delete(supplyMap, dm.TargetId)
			ch<- dm.GetId() // confirm
		}else{
			log.Printf("Can't find select spply ID")
		}
		supplyMu.Unlock()

	}else { // not SelectSupply
		log.Println("No target as :",dm)
		if sxutil.IDType(dm.SenderId) == clt.ClientID {
			log.Println("From me ",dm)
		}else { // not Demand From me.
			// select any ride share demand!
			// should check the type of ride..

			// we have several choices.

			// train + onemile
			// train + bus + onemile
			//  currently we do not consider walk.
			expSpecial(clt,dm)
//			go trainAndOnemile(clt, dm)
//			go trainAndBusAndOnemile(clt, dm)
		}
	}
}

// wait for rideshare demand.
func subscribeRideshareDemand(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	err :=	client.SubscribeDemand(ctx, rideshareDemandCallback)
	// comes here if channel closed
	log.Printf("Server closed... on Routing provider %v", err)
}

// wait for rideshare demand.
func subscribeRideshareSupply(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	log.Printf("Now supporting rideshare Supply")
	client.SubscribeSupply(ctx, rideshareSupplyCallback)
	// comes here if channel closed
	log.Printf("SupplyServer closed... on Routing provider")
}



func main() {
//	time.Local = time.FixedZone("Asia/Tokyo",9*60*60)

	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "RoutingProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Routing:RSDM}")
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_RIDE_SHARE,argJson)

	wg.Add(1)
	demandClient = sclient

	argJson2 := fmt.Sprintf("{Client:Routing:RSSP}")
	sclient2 := sxutil.NewSMServiceClient(client, api.ChannelType_RIDE_SHARE,argJson2)

	supplyClient = sclient2

	go subscribeRideshareDemand(sclient)

	wg.Add(1)

	go subscribeRideshareSupply(sclient2)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
