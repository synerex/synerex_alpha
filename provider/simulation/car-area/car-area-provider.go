package main

import (
	"context"
	"flag"
	"log"
	"sync"
	//"math/rand"

	"github.com/synerex/synerex_alpha/api/simulation/clock"
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
	areaId    = flag.Int("areaId", 1, "Area Id")
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

type ClockConfig struct {
	Time	uint32
	NumCycle	uint32
	CycleDuration uint32
	CycleTime uint32
}

func checkDemandArgOneOf(dm *pb.Demand) string {
	//demandType := ""
	log.Printf("demandType1.5 is %v", dm)
	if(dm.GetArg_ClockDemand() != nil){
		argOneof := dm.GetArg_ClockDemand()
		log.Printf("demandType2 is %v", argOneof.DemandType.String())
		switch(argOneof.DemandType.String()){
			case "SET": return "SET_CLOCK"
			case "FORWARD": return "FORWARD_CLOCK"
			case "STOP": return "STOP_CLOCK"
			case "BACK": return "BACK_CLOCK"
			case "START": return "START_CLOCK"
		}
	}
	if(dm.GetArg_AreaDemand() != nil){
		argOneof := dm.GetArg_AreaDemand()
		switch(argOneof.DemandType.String()){
			case "SET": return "SET_AREA"
			case "GET": return "GET_AREA"
		}
	}
	if(dm.GetArg_AgentDemand() != nil){
		argOneof := dm.GetArg_AreaDemand()
		switch(argOneof.DemandType.String()){
			case "SET": return "SET_AGENT"
		}
	}
	return "INVALID_TYPE"
}



func setArea(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setArea")
	//sendSupply(clt, "SEND_AREA", "{Area:[{Latitude:36.5, Longitude:135.6},{Latitude:40.5, Longitude:140.6}]}")
}

func getArea(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("getArea")
	//sendSupply(clt, "SEND_AREA", "{Area:[{Latitude:36.5, Longitude:135.6},{Latitude:40.5, Longitude:140.6}]}")
}

func setClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setClock")
	argOneof := dm.GetArg_ClockDemand()
	clockConfig := &ClockConfig{
		Time: argOneof.Time,
		NumCycle: argOneof.NumCycle,
		CycleDuration: argOneof.CycleDuration,
		CycleTime: argOneof.CycleTime,
	}
	clockInfo := clock.ClockInfo{
		Time: argOneof.Time,
		StatusType: 0, // OK
		Meta: "",
	}
	
	nm := "setClock respnse by car-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{Name: nm, JSON: js, ClockInfo: &clockInfo}

	log.Printf("clockConfig is %v", clockConfig)
	log.Printf("clockInfo is %v", clockInfo)
	sendSupply(sclientClock, opts)
}


func forwardAreaOK(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("forwardArea OK")
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	log.Println("Got demand callback")
	//a := dm.GetArg_AreaDemand()
	log.Printf("demand is %v", dm)
	demandType := checkDemandArgOneOf(dm)
	log.Printf("demandType is %v", demandType)
	switch demandType{
		case "SET_CLOCK": setClock(clt, dm)
		//case "SET_CLOCK_OK": setClockOK(clt, dm)
		//case "START_CLOCK": startClock(clt, dm)
		//case "FORWARD_CLOCK_OK": forwardClockOK(clt, dm)
		default: log.Println("demand callback is invalid.")
	}
}

func subscribeDemand(client *sxutil.SMServiceClient) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func sendSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts) {
	mu.Lock()
	id := sclient.RegisterSupply(opts)
	idlist = append(idlist, id) // my demand list
	spMap[id] = opts            // my demand options
	mu.Unlock()
	log.Printf("Register my supply as id %v, %v",id,idlist)
}

func sendDemand(sclient *sxutil.SMServiceClient, opts *sxutil.DemandOpts) {
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id) // my demand list
	dmMap[id] = opts            // my demand options
	mu.Unlock()
	log.Printf("Register my demand as id %v, %v",id,idlist)
}

func main() {
	flag.Parse()

	sxutil.RegisterNodeName(*nodesrv, "CarAreaProvider", false)

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
	argJson := fmt.Sprintf("{Client:CarArea, AreaId: %d}", *areaId)
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
