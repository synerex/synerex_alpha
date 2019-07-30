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

/*// this function waits
func startSelection(clt *sxutil.SMServiceClient,d time.Duration){
	var sid uint64

	for i := 0; i < 5; i++{
		time.Sleep(d / 5)
		log.Printf("waiting... %v",i)
	}
	mu.Lock()
	log.Printf("From now, let's find best taxi..")
	lowest := 99999
	// find most valuable proposal
	for k, sp := range spMap {
		dat := TaxiDemand{}
		js := sp.ArgJson
		err := json.Unmarshal([]byte(js),&dat)
		if err != nil {
			log.Printf("Err JSON %v",err)
		}
		if dat.Price < lowest { // we should check some ...
			sid = k
		}
	}
	mu.Unlock()
	log.Printf("Select supply %v", spMap[sid])
	clt.SelectSupply(spMap[sid])
	// we have to cleanup all info.
}*/

func startUp(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("start up")
	sendSupply(clt, "START_UP_OK", "{Area:{Latitude:40.5, Longitude:140.6}}")

}

func setAgent(clt *sxutil.SMServiceClient, dm *pb.Demand){
	// is agent position between area lonlat?
	log.Println("set agent")
	sendSupply(clt, "AET_AGENT_OK", "{Position:{Latitude:40.5, Longitude:140.6}}")
}


func setClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setClock")
	sendDemand(clt, "SET_CLOCK_OK", "{Date: '2019-7-29T22:32:13.234252Z'")
}

func forwardClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	if(dm.DemandName=="FORWARD_CLOCK"){
	log.Println("forwardClock")
	sendDemand(sclientArea, "GET_AREA", "{Area:{Latitude:36.5, Longitude:135.6}}")
	}else if(dm.DemandName=="SEND_AREA"){
		// here is forward function
		//

		// forward area ok
		sendDemand(sclientArea, "FORWARD_AREA_OK", "{Area:{Latitude:36.5, Longitude:135.6}}")
		sendSupply(sclientArea, "FORWARD_AREA_OK", "{Area:{Latitude:36.5, Longitude:135.6}}")

		// forward agent ok
		sendDemand(sclientAgent, "FORWARD_AGENT_OK", "{Position:{Latitude:36.5, Longitude:135.6}}")
		sendSupply(sclientAgent, "FORWARD_AGENT_OK", "{Position:{Latitude:36.5, Longitude:135.6}}")

		// forward area ok
		sendDemand(sclientClock, "FORWARD_CLOCK_OK", "{Date: '2019-7-29T22:32:13.234252Z'}")
		sendSupply(sclientClock, "FORWARD_CLOCK_OK", "{Date: '2019-7-29T22:32:13.234252Z'}")

	}
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	log.Println("Got demand callback")
	log.Printf("demand is %v",dm.DemandName)
	switch dm.DemandName{
	case "START_UP_B": startUp(clt, dm)
	case "SET_AGENT": setAgent(clt, dm)
	case "SET_CLOCK_ALL": setClock(clt, dm)
	case "FORWARD_CLOCK": forwardClock(clt, dm)
	case "SEND_AREA": forwardClock(clt, dm)
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

	sxutil.RegisterNodeName(*nodesrv, "PedAreaBProvider", false)

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
	argJson := fmt.Sprintf("{Client:PedAreaB}")
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
