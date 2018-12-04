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

// routing provider got several information to support user.
type RSInfo struct {
	supplyCh chan *api.Supply
	askDmo *sxutil.DemandOpts
	routes []*rideshare.Route
	info string   // special for Kota experiment
}

var (
	serverAddr            = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv                     = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	spMap             map[uint64]*sxutil.SupplyOpts
	rideShareMap      map[uint64]*RSInfo // subChannels
	rideShareMu       sync.RWMutex
	mu                sync.Mutex
	TrainTrainsitTime = 5 * time.Minute //  5 min
	BusTrainsitTime = 2 * time.Minute // 2 min

)

func init(){
	rideShareMap = make(map[uint64]*RSInfo )//
}


func getOnemileRoute(clt *sxutil.SMServiceClient, from *common.Point, to *common.Point, ft *common.Time, tt *common.Time) (uint64, *sxutil.DemandOpts){
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
	return id,&dmo
}

// start thinking from Train
func trainAndOnemile(clt *sxutil.SMServiceClient, dm *api.Demand) {


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
			ch := make(chan *api.Supply)
			rideShareMu.Lock()
			id, dmo := getOnemileRoute(clt, dpt, aimiPt, dt, nil)
			rts := []*rideshare.Route{}
			rideShareMap[id]= &RSInfo{ch, dmo, rts,"ToAimi,1st"}
			rideShareMu.Unlock()
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

			ch := make(chan *api.Supply)
			rideShareMu.Lock()
			id, dmo := getOnemileRoute(clt, aimiPt, apt, arTTime, nil)
			rts := []*rideshare.Route{route}
			rideShareMap[id]= &RSInfo{ch, dmo, rts,"FromAimi,2nd"}
			rideShareMu.Unlock()

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

			routeBus := getBusRouteFromTime(dpt, aimiPt,int32(dtTime.Minute()+ dtTime.Hour()*60), 0)
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
			id, dmo := getOnemileRoute(clt, dpt, aimiPt, dt, oneMileArriveTime)
			rts := []*rideshare.Route{routeBus}

			rideShareMap[id]= &RSInfo{ch, dmo, rts,"ToAimiBus,1st"}
			rideShareMu.Unlock()
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

			ch := make(chan *api.Supply)
			rideShareMu.Lock()
			id, dmo := getOnemileRoute(clt, aimiPt, apt, arTTime, nil)
			rts := []*rideshare.Route{route}

			rideShareMap[id]= &RSInfo{ch, dmo, rts,"FromAimi,2nd"}
			rideShareMu.Unlock()

		}

	}else { // arrival time only...
		log.Printf("arrival time.. is not yet implemented")
	}



}


func rideshareSupplyCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {

}

// callback for each Demand
func rideshareDemandCallback(clt *sxutil.SMServiceClient, dm *api.Demand) {
	// check if demand is match with my supply.
	log.Println("Got rideshare demand callback on Routing")

	// we need to start a new "Multiple Routing Suggestion".
	// for each Demand, we start go routine for that!

	// currently fix the idea for multi routing

	if dm.TargetId != 0 { // this is Select!
		log.Println("getSelect for Route!")
		// which select ?
		mu.Lock()
		if _, ok := spMap[dm.TargetId]; ok {
			log.Printf("Send Confirm !")
			 // we should start contract with them..
			clt.Confirm(sxutil.IDType(dm.Id))
			delete(spMap, dm.TargetId)
		}else{
			log.Printf("Can't find select spply ID")
		}
		mu.Unlock()

	}else { // not SelectSupply
		if sxutil.IDType(dm.SenderId) == clt.ClientID {

		}else { // not Demand From me.
			// select any ride share demand!
			// should check the type of ride..

			// we have several choices.

			// train + onemile
			// train + bus + onemile
			//  currently we do not consider walk.

			go trainAndOnemile(clt, dm)
			go trainAndBusAndOnemile(clt, dm)
		}
	}
}

// wait for rideshare demand.
func subscribeRideshareDemand(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeDemand(ctx, rideshareDemandCallback)
	// comes here if channel closed
	log.Printf("Server closed... on Routing provider")
}

// wait for rideshare demand.
func subscribeRideshareSupply(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeSupply(ctx, rideshareSupplyCallback)
	// comes here if channel closed
	log.Printf("SupplyServer closed... on Routing provider")
}



func main() {
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
	argJson := fmt.Sprintf("{Client:Routing}")
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_RIDE_SHARE,argJson)

	wg.Add(1)
	go subscribeRideshareDemand(sclient)

	go subscribeRideshareSupply(sclient)


	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
