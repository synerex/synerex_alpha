package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"github.com/golang/protobuf/ptypes"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// simple-routing provider got several information to support user.


var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	spMap      map[uint64]*sxutil.SupplyOpts
	mu		sync.RWMutex
)

type route struct {
	Routes [][2]float64
}

func init(){
	spMap = make(map[uint64]*sxutil.SupplyOpts)
}

func getRoute(from *common.Place , to *common.Place) [][2]float64 {
	RouteURL := "http://rt.synergic.mobi/api/route"
	// currently we do not have good routing engine, so, just use
	// main point.
	fpt := from.GetCentralPoint()
	tpt := to.GetCentralPoint()

	qstr := fmt.Sprintf("?from=%.6f,%.6f&to=%.6f,%.6f",fpt.Longitude,fpt.Longitude, tpt.Longitude, tpt.Latitude)
	url := RouteURL+qstr
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
//	fmt.Println(string(byteArray)) // htmlをstringで取得
	rt := new(route)
	js := json.Unmarshal(byteArray, &rt)
	if js != nil {
		log.Println("Can't unmarshal routing json")
	}
//	fmt.Printf("got result ! %v, %v\n", js, rt)

	return rt.Routes
}

func calcDistance(pts [][2]float64) float64{
	log.Printf("Calc distance %d",len(pts))
	if len(pts) == 0 {
		return 0.0
	}

	spt := pts[0]
	dist:= float64(0.0)
	for _, pt := range pts{
		dist +=	common.DistanceLonLat(spt[0],spt[1],pt[0],pt[1])
		spt = pt
	}
	log.Printf("Distance is %.2f",dist)
	return dist
}

// callback for each Demand
func routingDemandCallback(clt *sxutil.SMServiceClient, dm *api.Demand) {
	MeterPerSoconds := float64(30 * 1000 / 3600)  // meter per seconds
	// check if demand is match with my supply.
	log.Println("Got routing demand callback on SRouting")

	if dm.TargetId != 0 { // check for my data

		log.Printf("Got SelectSupply %d", dm.TargetId)
		mu.Lock()
		if _, ok := spMap[dm.TargetId]; ok {
			log.Printf("Send Confirm ")
			clt.Confirm(sxutil.IDType(dm.Id))
			delete(spMap, dm.TargetId)
		}else{
			log.Printf("Can't find select spply ID")
		}
		mu.Unlock()
		return
	}

	rt := dm.GetArg_RoutingService() // now got routing service.
	if rt == nil{
		log.Println("Can't get routing service info")
		return
	}
	points := getRoute(rt.GetDepartPlace(), rt.GetArrivePlace())
	dist := calcDistance(points)
	durTime := time.Duration( int64( float64(time.Second) * dist / MeterPerSoconds ))
	dp := ptypes.DurationProto(durTime)

	if len(points) == 0 {
		log.Printf("Can't find route!")
	}else{ // find route so prepare supply
		cpts := make([]*common.Point, len(points))

		for i, pts := range points {
			cpts[i] =&common.Point{Longitude:pts[0], Latitude:pts[1]}
		}

		rt.Points =cpts
		rt.AmountTime = dp

		if rt.DepartTime != nil {
			ts, err := ptypes.Timestamp(rt.DepartTime.GetTimestamp())
			if err != nil {
				log.Printf("Timestamp error")
			}
			at := ts.Add(durTime)
			tspb , _ := ptypes.TimestampProto(at)
			arrive := common.NewTime().WithTimestamp(tspb)
			rt.ArriveTime = arrive
		}else if rt.ArriveTime != nil{
			ts, err := ptypes.Timestamp(rt.ArriveTime.GetTimestamp())
			if err != nil {
				log.Printf("ArriveTime Timestamp error")
			}
			dt := ts.Add(-durTime)
			fspb , _ := ptypes.TimestampProto(dt)
			depart := common.NewTime().WithTimestamp(fspb)
			rt.DepartTime = depart
		}else {
			log.Printf("No time specified!")
		}

		spo := sxutil.SupplyOpts{
			Target: dm.GetId(),
			Name: "Supply Route from SRoute",
			JSON: "{Routing}",
			RoutingService: rt,
		}
		id :=clt.ProposeSupply(&spo)
		log.Printf("Send ProposeSupply %d",id)
		mu.Lock()
		spMap[id] = &spo
		mu.Unlock()
	}

}

// wait for routing demand.
func subscribeRoutingDemand(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeDemand(ctx, routingDemandCallback)
	// comes here if channel closed
	log.Printf("Server closed... on Routing provider")
}



func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "SimpleRoutingProvider", false)

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
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_ROUTING_SERVICE,argJson)

	wg.Add(1)
	go subscribeRoutingDemand(sclient)
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}

