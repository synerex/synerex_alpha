package main

import (
	"context"
	"flag"
	"log"
	"sync"
	//"math/rand"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	//"time"
	//"encoding/json"
	"fmt"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist     []uint64
	dmMap      map[uint64]*sxutil.DemandOpts
	spMap		map[uint64]*sxutil.SupplyOpts
	selection 	bool
	mu	sync.Mutex
	sclientArea *sxutil.SMServiceClient
	sclientAgent *sxutil.SMServiceClient
	sclientClock *sxutil.SMServiceClient
)

func init() {
	idlist = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
	selection = false
}

type LonLat struct{
	Latitude	float32
	Longitude	float32
}

type TaxiDemand struct {
	Price	int
	Distance	int
	Arrival int
	Destination int
	Position LonLat
}


func setClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setClock")
	sendDemand(clt, "SET_CLOCK_ALL", "{Date: '2019-7-29T22:32:13.234252Z'")
}

func setClockOK(clt *sxutil.SMServiceClient, dm *pb.Demand){
	// wait untill receive 3 SET_CLOCK_OK
	log.Println("setClockAllOK")
	sendSupply(clt, "SET_CLOCK_ALL_OK", "{Date: '2019-7-29T22:32:13.234252Z'")
}

func startClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("startClock")
	// every 1 cycle
	sendDemand(clt, "FORWARD_CLOCK", "{Forward: 1}")
}

func forwardClockOK(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("forwardClock OK")
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	log.Println("Got demand callback")
	log.Printf("demand is %v",dm.DemandName)
	switch dm.DemandName{
	case "SET_CLOCK": setClock(clt, dm)
	case "SET_CLOCK_OK": setClockOK(clt, dm)
	case "START_CLOCK": startClock(clt, dm)
	case "FORWARD_CLOCK_OK": forwardClockOK(clt, dm)
	default: log.Println("demand callback is valid.")

	}
}

func subscribeDemand(client *sxutil.SMServiceClient) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func sendSupply(sclient *sxutil.SMServiceClient, nm string, js string) {
	opts := &sxutil.SupplyOpts{Name: nm, JSON: js}
	mu.Lock()
	id := sclient.RegisterSupply(opts)
	idlist = append(idlist, id) // my demand list
	spMap[id] = opts            // my demand options
	mu.Unlock()
	log.Printf("Register my supply as id %v, %v",id,idlist)
}

func sendDemand(sclient *sxutil.SMServiceClient, nm string, js string) {
	opts := &sxutil.DemandOpts{Name: nm, JSON: js}
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id) // my demand list
	dmMap[id] = opts            // my demand options
	mu.Unlock()
	log.Printf("Register my demand as id %v, %v",id,idlist)
}

func main() {
	flag.Parse()

	sxutil.RegisterNodeName(*nodesrv, "ClockProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Clock}")
	sclientAgent = sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE,argJson)
	sclientClock = sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE,argJson)
	sclientArea = sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE,argJson)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go subscribeDemand(sclientAgent)
	go subscribeDemand(sclientClock)
	go subscribeDemand(sclientArea)

	/*for {
		sendDemand(sclient, "Share Ride to Home", "{Destination:{Latitude:36.5, Longitude:135.6}, Duration: 1200}")
		time.Sleep(time.Second * time.Duration(10 + rand.Int()%10))
	}*/
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
