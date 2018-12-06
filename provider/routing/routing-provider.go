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
				JSON:"{SPOPT}",
				RideShare: rdSh,
			}
			spch := make(chan uint64)
			spid :=	clt.ProposeSupply(&spo)
			supplyMap[spid] = spch
			select {
			case <-time.After(30 * time.Second):
				log.Printf("Timeout 30sec")
				return
			case id:= <-spch: // receive SelectSupply!
				clt.Confirm(sxutil.IDType(id))
			}

		}

	}else { // arrival time only...
		log.Printf("arrival time.. is not yet implemented")
	}


}


func  trainAndBusAndOnemile(clt *sxutil.SMServiceClient, dm *api.Demand) {

	rsInfo := dm.GetArg_RideShare()
	if rsInfo == nil {
		log.Printf("Demand is not for RideShare!")
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
				return
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
			id  := getOnemileRoute(clt, dpt, aimiPt, dt, oneMileArriveTime)
			rideShareMap[id]= ch
			rideShareMu.Unlock()
			// select for rideShare.

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

		}

	}else { // arrival time only...
		log.Printf("arrival time.. is not yet implemented")
	}



}


func rideshareSupplyCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	if sp.TargetId == 0{
		log.Printf("Should not come here...")
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
			log.Printf("Supply Calback: Can't find report channel for %d",sp.TargetId)
		}
		rideShareMu.RUnlock()
	}else{
		log.Printf("Cant check suply routing %v",*sp)
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

			go trainAndOnemile(clt, dm)
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
