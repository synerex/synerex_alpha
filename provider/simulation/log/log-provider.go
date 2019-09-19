package main

import (
	//"context"
	"flag"
	"log"
	"sync"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
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
	//dataMap map[uint32]*AreaData
	//areaData AreaData
)

func init() {
	idlist = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
	selection = false
	//areaData = *AreaData{}
	//dataMap = make(map[uint32]areaData)
}

//type AreaData struct {
//	AreaInfo map[string][]*agent
//	AgentInfo map[string][]*Agent
//}

type AreaData struct {
	AreaInfo *area.AreaInfo
	AgentInfo map[string][]*agent.AgentInfo
}

func setClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setClock")
	argOneof := dm.GetArg_ClockDemand()
	
	clockInfo := clock.ClockInfo{
		Time: argOneof.Time,
		SupplyType: 2,
		StatusType: 0, // OK
		Meta: "",
	}
	
	nm := "setClock respnse by log-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
		ClockInfo: &clockInfo,
	}

	simutil.SendProposeSupply(sclientClock, opts, spMap, idlist)
}

func forwardClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("forwardClock")
	nm := "forwardClock respnse by log-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
	}
	
	spMap, idlist = simutil.SendProposeSupply(sclientClock, opts, spMap, idlist)
}

func setArea(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setArea")
		
	nm := "setArea respnse by log-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
	}
	
	spMap, idlist = simutil.SendProposeSupply(sclientArea, opts, spMap, idlist)
}

func setAgent(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setAgent")
		
	nm := "setAgent respnse by log-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
	}
	
	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts, spMap, idlist)
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
	
	nm := "getParticipant respnse by log-provider"
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
		case "FORWARD_CLOCK": forwardClock(clt, dm)
		case "SET_AREA": setArea(clt, dm)
		case "SET_AGENT": setAgent(clt, dm)
		//case "SET_CLOCK_OK": setClockOK(clt, dm)
		//case "START_CLOCK": startClock(clt, dm)
		//case "FORWARD_CLOCK_OK": forwardClockOK(clt, dm)
		default: log.Println("demand callback is invalid.")
	}
}

func registArea(clt *sxutil.SMServiceClient, sp *pb.Supply){
	spArgOneof := sp.GetArg_AreaInfo()
	areaInfo := spArgOneof
	time := areaInfo.Time
	
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Time: ", time)
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Desc: ", sp.SupplyName)
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Info: ", areaInfo)
}

func registClock(clt *sxutil.SMServiceClient, sp *pb.Supply){
	spArgOneof := sp.GetArg_ClockInfo()
	clockInfo := spArgOneof
	time := clockInfo.Time
	
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Time: ", time)
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Desc: ", sp.SupplyName)
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Info: ", clockInfo)
}

func registAgent(clt *sxutil.SMServiceClient, sp *pb.Supply){
	spArgOneof := sp.GetArg_AgentInfo()
	agentInfo := spArgOneof
	time := agentInfo.Time
	
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Time: ", time)
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Desc: ", sp.SupplyName)
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Info: ", agentInfo)
}

func registAgents(clt *sxutil.SMServiceClient, sp *pb.Supply){
	spArgOneof := sp.GetArg_AgentsInfo()
	agentsInfo := spArgOneof

	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Desc: ", sp.SupplyName)
	fmt.Printf("\x1b[30m\x1b[47m%s %v\x1b[0m\n", "Info: ", agentsInfo)
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	fmt.Sprintf("getSupplyCallback, Regist supply data now")
	//if clt.IsSupplyTarget(sp, idlist) { 
		supplyType := simutil.CheckSupplyArgOneOf(sp)
//		fmt.Printf("supplyType %v", supplyType)
		switch supplyType{
			case "RES_SET_AREA":
				fmt.Println("registArea\n")
				registArea(clt, sp)
			case "RES_SET_CLOCK":
				registClock(clt, sp)
				fmt.Println("registClock")
			case "RES_SET_AGENT":
				registAgent(clt, sp)
				fmt.Println("registAgent")
			case "RES_SET_AGENTS":
				registAgents(clt, sp)
				fmt.Println("registAgents")
			case "RES_FORWARD_CLOCK":
				registClock(clt, sp)
				fmt.Println("registFowardClock")
			default:
				fmt.Println("SupplyCallback SupplyType is invalid")
		}
	//}else{
	//	fmt.Println("this is not propose supply")
	//}
}

func main() {
	flag.Parse()

	sxutil.RegisterNodeName(*nodesrv, "LogProvider", false)

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
	argJson := fmt.Sprintf("{Client:Log}")
	sclientAgent = sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE,argJson)
	sclientClock = sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE,argJson)
	sclientArea = sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE,argJson)
	sclientParticipant = sxutil.NewSMServiceClient(client, pb.ChannelType_PARTICIPANT_SERVICE,argJson)

	wg := sync.WaitGroup{}

	wg.Add(1)

	go simutil.SubscribeDemand(sclientAgent, demandCallback)
	go simutil.SubscribeDemand(sclientClock, demandCallback)
	go simutil.SubscribeDemand(sclientArea, demandCallback)
	go simutil.SubscribeDemand(sclientParticipant, demandCallback)

	go simutil.SubscribeSupply(sclientAgent, supplyCallback)
	go simutil.SubscribeSupply(sclientClock, supplyCallback)
	go simutil.SubscribeSupply(sclientArea, supplyCallback)
	go simutil.SubscribeSupply(sclientParticipant, supplyCallback)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
