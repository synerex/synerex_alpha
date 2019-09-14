package main

import (
	"context"
	"flag"
	"log"
	"sync"
	//"math/rand"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"

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
	sclientParticipant *sxutil.SMServiceClient
)

func init() {
	idlist = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
	selection = false
}

type ClockConfig struct {
	Time	uint32
	CycleNum	uint32
	CycleDuration uint32
	CycleInterval uint32
}

func setClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setClock")
	argOneof := dm.GetArg_ClockDemand()
	/*clockConfig := &ClockConfig{
		Time: argOneof.Time,
		CycleNum: argOneof.CycleNum,
		CycleDuration: argOneof.CycleDuration,
		CycleInterval: argOneof.CycleInterval,
	}*/
	clockInfo := clock.ClockInfo{
		Time: argOneof.Time,
		StatusType: 0, // OK
		Meta: "",
	}
	
	nm := "setClock respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
		ClockInfo: &clockInfo,
	}

	simutil.SendProposeSupply(sclientClock, opts, spMap, idlist)
}


func getParticipant(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("getParticipant")
	//argOneof := dm.GetArg_ParticipantDemand()
	participantInfo := participant.ParticipantInfo{
		ClientParticipantId: uint64(sclientParticipant.ClientID),
		ClientAreaId: uint64(sclientArea.ClientID),
		ClientAgentId: uint64(sclientAgent.ClientID),
		ClientClockId: uint64(sclientClock.ClientID),
		ClientType: 3, // PedArea
		AreaId: uint32(0), // Area A
		AgentType: 0, // Pedestrian
		StatusType: 0, // OK
		Meta: "",
	}
	
	nm := "getParticipant respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(), 
		Name: nm, 
		JSON: js, 
		ParticipantInfo: &participantInfo,
	}

	simutil.SendProposeSupply(sclientParticipant, opts, spMap, idlist)
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	demandType := simutil.CheckDemandArgOneOf(dm)
	switch demandType{
		case "GET_PARTICIPANT": getParticipant(clt, dm)
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

func main() {
	flag.Parse()

	sxutil.RegisterNodeName(*nodesrv, "AreaProvider", false)

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
	argJson := fmt.Sprintf("{Client:Area}")
	sclientAgent = sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE,argJson)
	sclientClock = sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE,argJson)
	sclientArea = sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE,argJson)
	sclientParticipant = sxutil.NewSMServiceClient(client, pb.ChannelType_PARTICIPANT_SERVICE,argJson)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go subscribeDemand(sclientAgent)
	go subscribeDemand(sclientClock)
	go subscribeDemand(sclientArea)
	go subscribeDemand(sclientParticipant)

	
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
