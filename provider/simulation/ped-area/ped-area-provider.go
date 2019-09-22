package main

import (
	//"context"
	"flag"
	"log"
	"sync"
	"math"
	//"math/rand"
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
	areaId    = flag.Int("areaId", 0, "Area Id")	// Area B
	agentType    = flag.Int("agentType", 0, "Agent Type")	// PEDESTRIAN
	idlist     []uint64
	dmMap      map[uint64]*sxutil.DemandOpts
	spMap		map[uint64]*sxutil.SupplyOpts
	selection 	bool
	mu	sync.Mutex
	ch 			chan *pb.Supply
	sclientArea *sxutil.SMServiceClient
	sclientAgent *sxutil.SMServiceClient
	sclientClock *sxutil.SMServiceClient
	sclientParticipant *sxutil.SMServiceClient
	data *Data
)

func init() {
	idlist = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
	selection = false
	ch = make(chan *pb.Supply)
	data = new(Data)
	data.AgentsInfo = make([]*agent.AgentInfo, 0)
}


type Data struct {
	AreaInfo *area.AreaInfo
	ClockInfo *clock.ClockInfo
	AgentsInfo []*agent.AgentInfo
}


func setArea(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setArea")
	dmArgOneof := dm.GetArg_AreaDemand()

	areaDemand := area.AreaDemand{
		Time: dmArgOneof.Time,
		AreaId: uint32(*areaId), // 
		DemandType: 1, // GET
		StatusType: 2, // NONE
		Meta: "",
	}
		
	nm := "getArea order by ped-area-provider"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, AreaDemand: &areaDemand}
	dmMap, idlist = simutil.SendDemand(sclientArea, opts, dmMap, idlist)

	callback := func(sp *pb.Supply){
		log.Println("getArea response")
		spArgOneof := sp.GetArg_AreaInfo()
		areaInfo := &area.AreaInfo{
			Time: spArgOneof.Time,
			AreaId: spArgOneof.AreaId, // A
			AreaName: spArgOneof.AreaName,
			Map: spArgOneof.Map,
			SupplyType: 0, // RES_SET
			StatusType: 0, // OK
			Meta: "",
		}
		// store AreaInfo data
		data.AreaInfo = areaInfo
		log.Printf("data.AreaInfo %v\n\n", data)
	
		nm2 := "setArea respnse by ped-area-provider"
		js2 := ""
		opts2 := &sxutil.SupplyOpts{
			Target: dm.GetId(),
			Name: nm2, 
			JSON: js2, 
			AreaInfo: areaInfo,
		}

		spMap, idlist = simutil.SendProposeSupply(sclientArea, opts2, spMap, idlist)
	}

	go func(){		
		for {
			select {
				case sp := <- ch:
					log.Printf("getArea response: ", sp.GetArg_AreaInfo())
					spArgOneof := sp.GetArg_AreaInfo()
					if spArgOneof.AreaId == uint32(*areaId){
						callback(sp)
						return
				}
			}
		}

		}()
	
}

// if agent type and coord satisfy, return true
func isAgentInArea(agentInfo *agent.AgentInfo) bool{
	lat := agentInfo.Route.Coord.Lat
	lon := agentInfo.Route.Coord.Lon
	slat := data.AreaInfo.Map.Coord.StartLat
	elat := data.AreaInfo.Map.Coord.EndLat
	slon := data.AreaInfo.Map.Coord.StartLon
	elon := data.AreaInfo.Map.Coord.EndLon
	if agentInfo.AgentType.String() == agent.AgentType_name[int32(*agentType)] && slat <= lat && lat <= elat &&  slon <= lon && lon <= elon {
		return true
	}else{
		log.Printf("agent type and coord is not match...\n\n")
		return false
	}
}

func setAgent(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setAgent")
	dmArgOneof := dm.GetArg_AgentDemand()
	agentInfo := &agent.AgentInfo{
		Time: dmArgOneof.Time,
		AgentId: dmArgOneof.AgentId,
		AgentName: dmArgOneof.AgentName,
		AgentStatus: dmArgOneof.AgentStatus,
		AgentType: dmArgOneof.AgentType,
		Route: dmArgOneof.Route,
		Rule: dmArgOneof.Rule,
		SupplyType: 0, // RES_SET
		StatusType: 0, //OK
		Meta: "",
	}

	// store AgentInfo data
	if isAgentInArea(agentInfo){
		data.AgentsInfo = append(data.AgentsInfo, agentInfo)
		log.Printf("data.AgentInfo %v\n\n", data)
	}
		
	nm := "setAgent respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
		AgentInfo: agentInfo,
	}
	
	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts, spMap, idlist)
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

	// store AgentInfo data
	data.ClockInfo = &clockInfo
	log.Printf("data.ClockInfo %v\n\n", data)
	
	nm := "setClock respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
		ClockInfo: &clockInfo,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientClock, opts, spMap, idlist)
}

func calcNextRoute(areaInfo *area.AreaInfo, agentInfo *agent.AgentInfo, otherAgentsInfo *agent.AgentsInfo) *agent.Route{

	route := agentInfo.Route
	nextCoord := &agent.Route_Coord{
		Lat: float32(float64(route.Coord.Lat) + float64(route.Speed) * 1 * math.Cos(float64(route.Direction))),
		Lon: float32(float64(route.Coord.Lon) + float64(route.Speed) * 1 * math.Sin(float64(route.Direction))), 
	}

	nextRoute := &agent.Route{
		Coord: nextCoord,
		Direction: route.Direction,
		Speed: route.Speed,
		Destination: float32(10),
		Departure: float32(100),
	}
	return nextRoute
}

func forwardClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("forwardClock")
	dmArgOneof := dm.GetArg_ClockDemand()
	time := dmArgOneof.Time
	nextTime := time + 1

	// get area
	areaDemand := area.AreaDemand{
		Time: time,
		AreaId: data.AreaInfo.AreaId, // A
		DemandType: 1, // GET
		StatusType: 2, // NONE
		Meta: "",
	}
		
	nm := "getArea order by ped-area-provider"
	js := ""
	opts := &sxutil.DemandOpts{
		Name: nm, 
		JSON: js, 
		AreaDemand: &areaDemand,
	}
	dmMap, idlist = simutil.SendDemand(sclientArea, opts, dmMap, idlist)

	sp := <- ch
	log.Println("getArea response")
	spArgOneof := sp.GetArg_AreaInfo()


	// get Agent
//	agentsDemand := &agent.AgentsDemand{
//		Time: uint32(1),
//		AreaId: 1,
//		AgentType: 0,
//		DemandType: 1, // GET
//		StatusType: 2, // NONE
//		Meta: "",
//	}
//		
//	nm := "getAgents order by ped-area-provider"
//	js := ""
//	opts := &sxutil.DemandOpts{
//		Name: nm, 
//		JSON: js, 
//		AreaDemand: agentsDemand,
//	}
//	dmMap, idlist = simutil.SendDemand(sclientAgent, opts, dmMap, idlist)
//
//	spAgent := <- ch
//	log.Println("getAgents response")
//	spAgentsArgOneof := sp.GetArg_AgentsInfo()


	/*areaInfo := &area.AreaInfo{
		Time: spArgOneof.Time,
		AgentId: spArgOneof.AreaId, // A
		AreaName: spArgOneof.AreaName,
		Map: spArgOneof.Map,
		SupplyType: 0, // RES_GET
		StatusType: 0, // OK
		Meta: "",
	}*/

	// calc agent
	agentsInfo := data.AgentsInfo
//	otherAgentsInfo := spAgentArgOneof
	otherAgentsInfo := &agent.AgentsInfo{}
	areaInfo := spArgOneof
	data.AgentsInfo = make([]*agent.AgentInfo, 0)
	for k, agentInfo := range agentsInfo{
		// calc next agentInfo
//		route := agentInfo.Route
		nextRoute := calcNextRoute(areaInfo, agentInfo, otherAgentsInfo)

		nextAgentInfo := &agent.AgentInfo{
			Time: nextTime,
			AgentId: agentInfo.AgentId,
			AgentName: agentInfo.AgentName,
			AgentType: agentInfo.AgentType,
			AgentStatus: agentInfo.AgentStatus,
			Route: nextRoute,
			Rule: agentInfo.Rule,
			SupplyType: 0, // RES_SET
			StatusType: 2, // NONE
			Meta: "",
		}

		log.Printf("nextAgentInfo %v %v\n\n", nextAgentInfo, k)

		data.AgentsInfo = append(data.AgentsInfo, nextAgentInfo)
	}


	/*// nextAreaInfo
	nextAreaInfo := &area.AreaInfo{
		Time: nextTime,
		AreaId: spArgOneof.AreaId, // A
		AreaName: spArgOneof.AreaName,
		Map: spArgOneof.Map,
		SupplyType: 0, // RES_GET
		StatusType: 0, // OK
		Meta: "",
	}
	data.AreaInfo = nextAreaInfo*/

	// propose agentInfo
	nextAgentsInfo := &agent.AgentsInfo{
		Time: nextTime,
		AreaId: uint32(*areaId),
		AgentType: 0,
		AgentInfo: data.AgentsInfo,
		SupplyType: 0, // RES_SET
		StatusType: 2,
		Meta: "",
	}

	nm2 := "forwardClock to agentCh respnse by ped-area-provider"
	js2 := ""
	opts2 := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm2, 
		JSON: js2, 
		AgentsInfo: nextAgentsInfo,
	}
	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts2, spMap, idlist)
	
	// propose clockInfo
	nextClockInfo := &clock.ClockInfo{
		Time: nextTime,
		SupplyType: 0, //Forward
		StatusType: 0, // OK
		Meta: "",
	}
	data.ClockInfo = nextClockInfo

	nm3 := "forwardClock to clockCh respnse by ped-area-provider"
	js3 := ""
	opts3 := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm3, 
		JSON: js3, 
		ClockInfo: nextClockInfo,
	}
	spMap, idlist = simutil.SendProposeSupply(sclientClock, opts3, spMap, idlist)
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
		AreaId: uint32(*areaId), // Area A
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

	spMap, idlist = simutil.SendProposeSupply(sclientParticipant, opts, spMap, idlist)
}

func getAgents(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("getAgents")
	dmArgOneof := dm.GetArg_AgentsDemand()

	agentsInfo := &agent.AgentsInfo{
		Time: dmArgOneof.Time,
		AreaId: dmArgOneof.AreaId,
		AgentType: dmArgOneof.AgentType,
		AgentInfo: data.AgentsInfo,
		SupplyType: 1,
		StatusType: 0,
		Meta: "",
	}
		
	nm := "getAgent respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
		AgentsInfo: agentsInfo,
	}
	
	spMap, idlist = simutil.SendSupply(sclientAgent, opts, spMap, idlist)
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
		case "GET_AGENTS": getAgents(clt, dm)
		default: log.Println("demand callback is invalid.")
	}
}


// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	fmt.Sprintf("getSupplyCallback", sp)
	supplyType := simutil.CheckSupplyArgOneOf(sp)
	switch supplyType{
		case "RES_GET_AREA":
			ch <- sp
//			callbackForGetArea(clt, sp)
			fmt.Println("getArea")
		default:
			fmt.Println("error")
	}
}


func main() {
	flag.Parse()
	log.Printf("area id is: %v, agent type is %v",*areaId, *agentType)

	sxutil.RegisterNodeName(*nodesrv, "PedAreaProvider", false)

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
	argJson := fmt.Sprintf("{Client:PedArea, AreaId: %d}", *areaId)
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

	go simutil.SubscribeSupply(sclientArea, supplyCallback)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
