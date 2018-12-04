package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// simple-routing test provider
//

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	dmMap 	map[uint64]*sxutil.DemandOpts
	mu		sync.RWMutex
)

func init(){
	dmMap = make(map[uint64]*sxutil.DemandOpts)
}

type point struct {
	Point [2]float64
}


func getPoint() *common.Point {
	PointURL := "http://rt.synergic.mobi/api/rndpoint"
	// currently we do not have good routing engine, so, just use
	// main point.

	url := PointURL
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(byteArray)) // htmlをstringで取得
	rt := new(point)
	js := json.Unmarshal(byteArray, &rt)
	if js != nil {
		log.Println("Can't unmarshal routing json")
	}

	ppt := new(common.Point)

	ppt.Longitude = rt.Point[0]
	ppt.Latitude = rt.Point[1]

	return ppt
}

func supplyRideshareCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	rt := sp.GetArg_RideShare()
	if rt != nil { // get Rideshare supplu
		mu.RLock()
		if dmo, ok := dmMap[sp.TargetId]; ok{
//			log.Printf("Supply: %v",sp)
			log.Printf("dmo:%v",*dmo)

			delete(dmMap, sp.TargetId)

			// send select! (required?!)
			log.Printf("Send SelectSupply %d", sp.Id)
			id, err := clt.SelectSupply(sp)
			if err == nil {
				log.Printf("SelectSupply Success! and MbusID=%d", id)
			}

		}
		mu.RUnlock()
	}else{
		log.Printf("Cant check suply routing %v",*sp)
	}
}


// wait for routing demand.
func subscribeRideshareSupply(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeSupply(ctx, supplyRideshareCallback)
	// comes here if channel closed
	log.Printf("Server closed... on Routing provider")
}


func registerRideshareDemand(clt *sxutil.SMServiceClient){
	pt1 := getPoint()
	pt2 := new(common.Point)
	pt2.Longitude = 136.881638
	pt2.Latitude = 35.170694

	tm := time.Now().In(time.Local)
	tsp,_ := ptypes.TimestampProto(tm)

	rs :=rideshare.RideShare{
		DepartPoint: common.NewPlace().WithPoint(pt1),
		DepartTime: common.NewTime().WithTimestamp(tsp),
		ArrivePoint: common.NewPlace().WithPoint(pt2),
	}

	dmo := sxutil.DemandOpts{
		Name: "RideshareDemand",
		JSON: "{from, to}",
		RideShare: &rs,
	}
	id :=clt.RegisterDemand(&dmo)
	mu.Lock()
	dmMap[id] = &dmo  // save demand
	mu.Unlock()
	log.Printf("Sent RideshareDemand %d, %v",id, rs)

}

func main() {
//	time.Local = time.FixedZone("Asia/Tokyo",9*60*60)

	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "FakeKotaProvider", false)

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
	argJson := fmt.Sprintf("{Client:KotaTest}")
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_RIDE_SHARE,argJson)

	wg.Add(1)
	go subscribeRideshareSupply(sclient)

	// send RoutingDemand

//	for {
		log.Printf("Sending Demand!")
		registerRideshareDemand(sclient)
		time.Sleep(10*time.Second)
//	}

	wg.Wait()

	sxutil.CallDeferFunctions() // cleanup!

}

