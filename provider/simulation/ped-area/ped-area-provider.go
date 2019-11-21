package main

import (
	//"context"
	"flag"
	"fmt"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr          = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv             = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	areaId              = flag.Int("areaId", 1, "Area Id")       // Area B
	agentType           = flag.Int("agentType", 0, "Agent Type") // PEDESTRIAN
	clockTime           = flag.Int("time", 1, "Time")
	cycleNum            = flag.Int("num", 1, "Num")
	cycleInterval       = flag.Int("interval", 1, "Interval")
	cycleDuration       = flag.Int("duration", 1, "Duration")
	mu                  sync.Mutex
	ch                  chan *pb.Supply
	sameCh              chan *pb.Supply
	neighborCh          chan *pb.Supply
	sprovider           *simutil.SynerexProvider
	sim                 *simutil.SynerexSimulator
	CHANNEL_BUFFER_SIZE int
)

func init() {
	CHANNEL_BUFFER_SIZE = 10
	ch = make(chan *pb.Supply)
	neighborCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	sameCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
}

// Finish Fix
// when start up,
func setArea() {

	sprovider.GetAreaDemand(uint32(*clockTime), uint32(*areaId))
	sp := <-ch
	getAreaSupply := sp.GetArg_GetAreaSupply()
	areaInfo := getAreaSupply.AreaInfo
	sim.Area = areaInfo
	log.Printf("finish setting area: %v", areaInfo)
}

// Finish Fix
func setAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("setAgent")
	setAgentsDemand := dm.GetArg_SetAgentsDemand()
	agentsInfo := setAgentsDemand.AgentsInfo

	// get route demand
	sprovider.GetAgentsRouteDemand(agentsInfo)
	log.Println("setAgent2")
	sp := <-ch
	log.Println("GET_AGENTS_ROUTE_FINISH")
	getAgentsRouteSupply := sp.GetArg_GetAgentsRouteSupply()
	agentsInfo = getAgentsRouteSupply.AgentsInfo

	for _, agentInfo := range agentsInfo {
		if simutil.IsAgentInControlledArea(agentInfo, sim.Area, int32(*agentType)) {
			agentInfo.ControlArea = uint32(*areaId)
		}

		if simutil.IsAgentInArea(agentInfo, sim.Area, int32(*agentType)) {
			sim.AddAgent(agentInfo)
			//sim.Agents = append(sim.Agents, agentInfo)
		}
	}

	sprovider.SetAgentsSupply(dm.GetId(), 1, uint32(*areaId), agent.AgentType(*agentType))
}

// Finish Fix
func setClock() {
	log.Println("setClock")

	// store AgentInfo data
	sim.TimeStep = 1
	sim.GlobalTime = float64(*clockTime)
	//data.ClockInfo = clockInfo
}

func getAreaInfo() *area.AreaInfo {

	sprovider.GetAreaDemand(uint32(*clockTime), uint32(*areaId))

	sp := <-ch
	log.Println("GET_AREA_FINISH")
	getAreaSupply := sp.GetArg_GetAreaSupply()
	areaInfo := getAreaSupply.AreaInfo

	return areaInfo
}

func getSameAreaAgents() []*agent.AgentInfo {

	sameAreaAgents := make([]*agent.AgentInfo, 0)
	if len(sprovider.SameAreaIdList) != 0 {

		// get Agent
		sprovider.GetAgentsDemand(uint32(1), uint32(*areaId), agent.AgentType(*agentType))

		sameAreaPspMap := sprovider.Wait(sprovider.SameAreaIdList, sameCh)
		fmt.Printf("sameArea: %v\n", sameAreaPspMap)
		for _, agentPsp := range sameAreaPspMap {
			getAgentsSupply := agentPsp.GetArg_GetAgentsSupply()
			sameAreaAgents = append(sameAreaAgents, getAgentsSupply.AgentsInfo...)
		}
	} else {
		//log.Println("SAME_AREA_NOTHING")
	}
	fmt.Printf("sameArea: %v\n", sameAreaAgents)

	return sameAreaAgents
}

// Finish Fix
func forwardClock(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Printf("FORWARD_CLOCK_START \n\n")
	forwardClockDemand := dm.GetArg_ForwardClockDemand()
	time := forwardClockDemand.Time
	currentClockInfo := &clock.ClockInfo{
		Time:          uint32(sim.GlobalTime),
		CycleDuration: uint32(sim.TimeStep),
		CycleNum:      uint32(*cycleNum),
		CycleInterval: uint32(*cycleInterval),
	}
	nextTime := time + 1

	// 同じエリアにいるエージェントを取得する
	sameAreaAgents := getSameAreaAgents()

	// 次の時間のエージェントを計算する。重複エリアの更新をすませてないのでpureNextAgentsInfoとしている
	pureNextAgents := sim.CalcNextAgents(sameAreaAgents)

	// 計算後のエージェントを同じエリアの異種エージェントプロバイダへ送る
	sprovider.ForwardAgentsSupply(dm.GetId(), nextTime, uint32(*areaId), pureNextAgents, agent.AgentType(*agentType))

	// 同じエリアの異種エージェントプロバイダから計算後のエージェント情報を取得する
	neighborAgents := make([]*agent.AgentInfo, 0)
	if len(sprovider.NeighborAreaIdList) != 0 {
		fmt.Printf("neihg")
		neighborPspMap := sprovider.Wait(sprovider.NeighborAreaIdList, neighborCh)
		for _, psp := range neighborPspMap {
			forwardAgentsSupply := psp.GetArg_ForwardAgentsSupply()
			neighborAgents = append(neighborAgents, forwardAgentsSupply.AgentsInfo...)
		}
	} else {
		//log.Println("NEIGHBOR_AREA_NOTHING")
	}

	// update Agents data
	nextAgentsInfo := sim.UpdateDuplicateAgents(pureNextAgents, neighborAgents)

	// propose clockInfo
	sprovider.ForwardClockSupply(dm.GetId(), currentClockInfo)
	sim.Agents = nextAgentsInfo
	sim.GlobalTime = float64(nextTime)

	log.Printf("\x1b[30m\x1b[47m \n FORWARD_CLOCK_FINISH \n TIME: %v \n AGENT_NUM: %v \x1b[0m\n", time, len(pureNextAgents))

}

// Finish Fix
func getParticipant(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("getParticipant")

	participantInfo := &participant.ParticipantInfo{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint32(sprovider.ParticipantClient.ClientID),
			AreaChannelId:        uint32(sprovider.AreaClient.ClientID),
			AgentChannelId:       uint32(sprovider.AgentClient.ClientID),
			ClockChannelId:       uint32(sprovider.ClockClient.ClientID),
			RouteChannelId:       uint32(sprovider.RouteClient.ClientID),
		},
		ProviderType: 3, // PedArea
		AreaId:       uint32(*areaId),
		AgentType:    0, // Ped
	}

	sprovider.GetParticipantSupply(dm.GetId(), participantInfo)

}

// Finish Fix
func setParticipant(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("setParticipant")
	setParticipantDemand := dm.GetArg_SetParticipantDemand()

	participantsInfo := setParticipantDemand.ParticipantsInfo
	fmt.Printf("participants: %v\n", participantsInfo)
	sameAreaIdList, neighborAreaIdList := sim.CreateSyncIdList(participantsInfo)
	sprovider.SetRelateAreaIdList(sameAreaIdList, neighborAreaIdList)
	fmt.Printf("Same %v", sameAreaIdList)
	fmt.Printf("Neighbor %v", neighborAreaIdList)

	participantInfo := &participant.ParticipantInfo{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint32(sprovider.ParticipantClient.ClientID),
			AreaChannelId:        uint32(sprovider.AreaClient.ClientID),
			AgentChannelId:       uint32(sprovider.AgentClient.ClientID),
			ClockChannelId:       uint32(sprovider.ClockClient.ClientID),
			RouteChannelId:       uint32(sprovider.RouteClient.ClientID),
		},
		ProviderType: 3, // PedArea
		AreaId:       uint32(*areaId),
		AgentType:    0, // Ped
	}

	sprovider.SetParticipantSupply(dm.GetId(), participantInfo)

}

// Finish Fix
func getAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Printf("getAgents %v", dm)

	getAgentsDemand := dm.GetArg_GetAgentsDemand()
	//log.Printf("info %v, %v, %v, %v", int(getAgentsDemand.AreaId), *areaId, int(getAgentsDemand.AgentType), *agentType)
	isNeighborArea := sim.IsContainNeighborMap(getAgentsDemand.AreaId)
	isSameArea := int(getAgentsDemand.AreaId) == *areaId && int(getAgentsDemand.AgentType) != *agentType
	if isSameArea || isNeighborArea {

		agentsInfo := sim.Agents
		sprovider.GetAgentsSupply(dm.GetId(), getAgentsDemand.Time, uint32(*areaId), agentsInfo, agent.AgentType(*agentType))
	}

}

// Finish Fix
// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	demandType := simutil.CheckDemandType(dm)
	switch demandType {
	case "GET_PARTICIPANT_DEMAND":
		getParticipant(clt, dm)
	case "SET_PARTICIPANT_DEMAND":
		setParticipant(clt, dm)
	case "FORWARD_CLOCK_DEMAND":
		forwardClock(clt, dm)
	case "SET_AGENTS_DEMAND":
		setAgents(clt, dm)
	case "GET_AGENTS_DEMAND":
		getAgents(clt, dm)
	default:
		//log.Println("demand callback is invalid.")
	}
}

// Finish Fix
// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {

	supplyType := simutil.CheckSupplyType(sp)
	if sprovider.IsSupplyTarget(sp) {
		switch supplyType {
		case "GET_AREA_SUPPLY":
			ch <- sp
		case "GET_AGENTS_ROUTE_SUPPLY":
			ch <- sp
		case "GET_AGENTS_SUPPLY":
			getAgentsSupply := sp.GetArg_GetAgentsSupply()
			if getAgentsSupply.AreaId == uint32(*areaId) {
				//ch <- sp
				sprovider.SendToWait(sp, sameCh)
			}
		default:
			//fmt.Println("error")
		}
	}

	// FORWARD_AGENTS_SUPPLY
	//supplyType := simutil.CheckSupplyType(sp)
	if supplyType == "FORWARD_AGENTS_SUPPLY" {
		//mu.Lock()
		forwardAgentsSupply := sp.GetArg_ForwardAgentsSupply()
		// not equal areaId and not AreaProvider
		if forwardAgentsSupply.AreaId != uint32(*areaId) && forwardAgentsSupply.AreaId != 0 {
			//ch <- sp
			sprovider.SendToWait(sp, neighborCh)
		}
	}

}

func main() {
	flag.Parse()
	log.Printf("area id is: %v, agent type is %v", *areaId, *agentType)

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

	// synerex simulator
	sim = simutil.NewSynerexSimulator(1, 0, 0)

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:PedArea, AreaId: %d}", *areaId)

	// Clientとして登録
	sprovider = simutil.NewSynerexProvider()
	channelTypes := []pb.ChannelType{
		pb.ChannelType_AGENT_SERVICE,
		pb.ChannelType_CLOCK_SERVICE,
		pb.ChannelType_AREA_SERVICE,
		pb.ChannelType_PARTICIPANT_SERVICE,
		pb.ChannelType_ROUTE_SERVICE,
	}
	sprovider.RegisterClient(client, channelTypes, argJson)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go simutil.SubscribeDemand(sprovider.AgentClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeDemand(sprovider.ClockClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeDemand(sprovider.AreaClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeDemand(sprovider.ParticipantClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sprovider.AreaClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sprovider.AgentClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sprovider.ParticipantClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sprovider.RouteClient, supplyCallback, &wg)
	wg.Wait()

	// start up(setArea)
	wg.Add(1)
	setArea()
	setClock()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
