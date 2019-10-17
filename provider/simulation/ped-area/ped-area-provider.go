package main

import (
	//"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

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
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv            = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	areaId             = flag.Int("areaId", 1, "Area Id")       // Area B
	agentType          = flag.Int("agentType", 0, "Agent Type") // PEDESTRIAN
	clockTime          = flag.Int("time", 1, "Time")
	cycleNum           = flag.Int("num", 1, "Num")
	cycleInterval      = flag.Int("interval", 1, "Interval")
	cycleDuration      = flag.Int("duration", 1, "Duration")
	dmIdList           []uint64
	spIdList           []uint64
	idListByChannel    *simutil.IdListByChannel
	sameAreaIdList     []uint32
	neighborAreaIdList []uint32
	isGetParticipant   bool
	dmMap              map[uint64]*sxutil.DemandOpts
	spMap              map[uint64]*sxutil.SupplyOpts
	pspMap             map[uint64]*pb.Supply
	neighborPspMap     map[uint64]*pb.Supply
	samePspMap         map[uint64]*pb.Supply
	selection          bool
	startCollectId     bool
	startSync          bool
	mu                 sync.Mutex
	ch                 chan *pb.Supply
	sclientArea        *sxutil.SMServiceClient
	sclientAgent       *sxutil.SMServiceClient
	sclientClock       *sxutil.SMServiceClient
	sclientRoute       *sxutil.SMServiceClient
	sclientParticipant *sxutil.SMServiceClient
	participantsInfo   []*participant.ParticipantInfo
	data               *Data
	history            *History
	isFinishSync       bool
)

func init() {
	spIdList = make([]uint64, 0)
	dmIdList = make([]uint64, 0)
	isGetParticipant = false
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
	pspMap = make(map[uint64]*pb.Supply)
	neighborPspMap = make(map[uint64]*pb.Supply)
	samePspMap = make(map[uint64]*pb.Supply)
	selection = false
	startCollectId = false
	startSync = false
	isFinishSync = false
	ch = make(chan *pb.Supply)
	data = new(Data)
	data.AgentsInfo = make([]*agent.AgentInfo, 0)
	history = new(History)
	history.History = make(map[uint32]*Data)
}

type Data struct {
	AreaInfo   *area.AreaInfo
	ClockInfo  *clock.ClockInfo
	AgentsInfo []*agent.AgentInfo
}

type History struct {
	CurrentTime uint32
	History     map[uint32]*Data
}

func isContainNeighborMap(areaId uint32) bool {
	neighborMap := data.AreaInfo.NeighborArea
	for _, neighborId := range neighborMap {
		if areaId == neighborId {
			return true
		}
	}
	return false
}

// Finish Fix
// when start up,
func setArea() {
	log.Println("setArea")

	getAreaDemand := &area.GetAreaDemand{
		Time:       uint32(*clockTime),
		AreaId:     uint32(*areaId),
		StatusType: 2, //NONE
		Meta:       "",
	}

	nm := "getArea order by ped-area-provider"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, GetAreaDemand: getAreaDemand}
	dmMap, dmIdList = simutil.SendDemand(sclientArea, opts, dmMap, dmIdList)

	sp := <-ch
	getAreaSupply := sp.GetArg_GetAreaSupply()
	areaInfo := getAreaSupply.AreaInfo
	data.AreaInfo = areaInfo
	log.Printf("finish setting area: %v", areaInfo)
}

// Finish Fix
func setAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("setAgent")
	setAgentsDemand := dm.GetArg_SetAgentsDemand()
	agentsInfo := setAgentsDemand.AgentsInfo
	for _, agentInfo := range agentsInfo {
		if simutil.IsAgentInControlledArea(agentInfo, data.AreaInfo, int32(*agentType)) {
			agentInfo.ControlArea = uint32(*areaId)
		}

		if simutil.IsAgentInArea(agentInfo, data.AreaInfo, int32(*agentType)) {
			data.AgentsInfo = append(data.AgentsInfo, agentInfo)
		}
	}

	if data.AreaInfo != nil && data.AgentsInfo != nil && data.ClockInfo != nil {
		time := data.ClockInfo.Time
		history.History[time] = data
		history.CurrentTime = time
		log.Printf("\x1b[30m\x1b[47m History is : %v\x1b[0m\n", history)
	}

	setAgentsSupply := &agent.SetAgentsSupply{
		Time:       uint32(1),
		AreaId:     uint32(*areaId),
		AgentType:  0, // Ped
		StatusType: 0,
		Meta:       "",
	}

	nm := "setAgentSupply by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:          dm.GetId(),
		Name:            nm,
		JSON:            js,
		SetAgentsSupply: setAgentsSupply,
	}

	spMap, spIdList = simutil.SendProposeSupply(sclientAgent, opts, spMap, spIdList)
}

// Finish Fix
func setClock() {
	log.Println("setClock")

	clockInfo := &clock.ClockInfo{
		Time:          uint32(*clockTime),
		CycleDuration: uint32(*cycleDuration),
		CycleNum:      uint32(*cycleNum),
		CycleInterval: uint32(*cycleInterval),
	}

	// store AgentInfo data
	data.ClockInfo = clockInfo
}

// Finish Fix
func calcNextRoute(areaInfo *area.AreaInfo, agentInfo *agent.AgentInfo, otherAgentsInfo []*agent.AgentInfo) *agent.Route {

	route := agentInfo.Route
	nextCoord := &agent.Route_Coord{
		Lat: float32(float64(route.Coord.Lat) + float64(0.0001)*1*math.Cos(float64(route.Direction*math.Pi/180))),
		Lon: float32(float64(route.Coord.Lon) + float64(0.0001)*1*math.Sin(float64(route.Direction*math.Pi/180))),
	}

	// change speed and direction
	speedFluct := 2
	directionFluct := 5
	rand.Seed(time.Now().UnixNano())
	speed := route.Speed + float32(rand.Intn(speedFluct)) - float32(rand.Intn(speedFluct)/2)
	direction := route.Direction + float32(rand.Intn(directionFluct)) - float32(rand.Intn(directionFluct)/2)

	nextRoute := &agent.Route{
		Coord:       nextCoord,
		Direction:   float32(direction),
		Speed:       float32(speed),
		Destination: float32(10),
		Departure:   float32(100),
	}
	return nextRoute
}

func getAreaInfo() *area.AreaInfo {
	// get area
	getAreaDemand := &area.GetAreaDemand{
		Time:       uint32(*clockTime),
		AreaId:     uint32(*areaId),
		StatusType: 2, //NONE
		Meta:       "",
	}

	nm := "getArea order by ped-area-provider"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, GetAreaDemand: getAreaDemand}
	dmMap, dmIdList = simutil.SendDemand(sclientArea, opts, dmMap, dmIdList)

	sp := <-ch
	//wait()
	log.Println("GET_AREA_FINISH")
	getAreaSupply := sp.GetArg_GetAreaSupply()
	areaInfo := getAreaSupply.AreaInfo

	return areaInfo
}

func getAgentsInfo() []*agent.GetAgentsSupply {
	// get Agent
	getAgentsDemand := &agent.GetAgentsDemand{
		Time:       uint32(1),
		AreaId:     uint32(*areaId),
		AgentType:  0, //Ped
		StatusType: 2, // NONE
		Meta:       "",
	}

	nm4 := "getAgents order by ped-area-provider"
	js4 := ""
	opts4 := &sxutil.DemandOpts{
		Name:            nm4,
		JSON:            js4,
		GetAgentsDemand: getAgentsDemand,
	}
	log.Println("sendDemand AgentsDemand")
	dmMap, dmIdList = simutil.SendDemand(sclientAgent, opts4, dmMap, dmIdList)

	wait()
	//agentPspMap := <-pspCh
	log.Println("GET_AGENTS_FINISH")
	//spAgentsArgOneof := spAgent.GetArg_AgentsInfo()
	getAgentsSupplies := make([]*agent.GetAgentsSupply, 0)
	for _, agentPsp := range samePspMap {
		getAgentsSupply := agentPsp.GetArg_GetAgentsSupply()
		getAgentsSupplies = append(getAgentsSupplies, getAgentsSupply)
	}
	return getAgentsSupplies
}

func updateAgentsInfo(areaInfo *area.AreaInfo, getAgentsSupplies []*agent.GetAgentsSupply) {

}

func calcAgentsInfo(areaInfo *area.AreaInfo, nextTime uint32) []*agent.AgentInfo {
	// calc agent
	agentsInfo := data.AgentsInfo
	//	otherAgentsInfo := spAgentArgOneof
	otherAgentsInfo := make([]*agent.AgentInfo, 0)
	data.AgentsInfo = make([]*agent.AgentInfo, 0)
	for _, agentInfo := range agentsInfo {
		// calc next agentInfo
		//		route := agentInfo.Route
		nextRoute := calcNextRoute(areaInfo, agentInfo, otherAgentsInfo)

		nextAgentInfo := &agent.AgentInfo{
			Time:        nextTime,
			AgentId:     agentInfo.AgentId,
			AgentType:   agentInfo.AgentType,
			AgentStatus: agentInfo.AgentStatus,
			Route:       nextRoute,
		}

		data.AgentsInfo = append(data.AgentsInfo, nextAgentInfo)
	}

	controlAgentsInfo := make([]*agent.AgentInfo, 0)
	for _, agentInfo := range data.AgentsInfo {
		if simutil.IsAgentInControlledArea(agentInfo, data.AreaInfo, int32(*agentType)) {
			controlAgentsInfo = append(controlAgentsInfo, agentInfo)
		}
	}
	return controlAgentsInfo
}

func wait() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			if isFinishSync {
				isFinishSync = false
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
}

// Finish Fix
func forwardClock(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("forwardClock")
	forwardClockDemand := dm.GetArg_ForwardClockDemand()
	time := forwardClockDemand.Time
	nextTime := time + 1

	// update Area data
	areaInfo := getAreaInfo()
	data.AreaInfo = areaInfo

	// update Agents data
	getAgentsSupplies := getAgentsInfo()

	// update Agents data
	updateAgentsInfo(areaInfo, getAgentsSupplies)

	// calc Agents data
	controlAgentsInfo := calcAgentsInfo(areaInfo, nextTime)

	forwardAgentsSupply := &agent.ForwardAgentsSupply{
		Time:       nextTime,
		AreaId:     uint32(*areaId),
		AgentType:  0, //Ped
		AgentsInfo: controlAgentsInfo,
		StatusType: 0, //OK
		Meta:       "",
	}

	nm2 := "forwardAgents to agentCh respnse by ped-area-provider"
	js2 := ""
	opts2 := &sxutil.SupplyOpts{
		Target:              dm.GetId(),
		Name:                nm2,
		JSON:                js2,
		ForwardAgentsSupply: forwardAgentsSupply,
	}
	spMap, spIdList = simutil.SendProposeSupply(sclientAgent, opts2, spMap, spIdList)

	wait()
	log.Println("FORWARD_AGENTS_FINISH")

	// propose clockInfo
	nextClockInfo := &clock.ClockInfo{
		Time: nextTime,
	}
	data.ClockInfo = nextClockInfo

	forwardClockSupply := &clock.ForwardClockSupply{
		ClockInfo:  nextClockInfo,
		StatusType: 0, // OK
		Meta:       "",
	}

	nm3 := "forwardClock to clockCh respnse by ped-area-provider"
	js3 := ""
	opts3 := &sxutil.SupplyOpts{
		Target:             dm.GetId(),
		Name:               nm3,
		JSON:               js3,
		ForwardClockSupply: forwardClockSupply,
	}
	spMap, spIdList = simutil.SendProposeSupply(sclientClock, opts3, spMap, spIdList)

	if data.AreaInfo != nil && data.AgentsInfo != nil && data.ClockInfo != nil {
		time := data.ClockInfo.Time
		history.History[time] = data
		history.CurrentTime = time
		log.Printf("\x1b[30m\x1b[47m History is : %v\x1b[0m\n", history)
	}

	log.Printf("FORWARD_CLOCK_FINISH\n\n")
}

// Finish Fix
func getParticipant(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("getParticipant")

	participantInfo := &participant.ParticipantInfo{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint32(sclientParticipant.ClientID),
			AreaChannelId:        uint32(sclientArea.ClientID),
			AgentChannelId:       uint32(sclientAgent.ClientID),
			ClockChannelId:       uint32(sclientClock.ClientID),
			RouteChannelId:       uint32(sclientRoute.ClientID),
		},
		ProviderType: 3, // PedArea
		AreaId:       uint32(*areaId),
		AgentType:    0, // Ped
	}

	getParticipantSupply := &participant.GetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "getParticipant respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               dm.GetId(),
		Name:                 nm,
		JSON:                 js,
		GetParticipantSupply: getParticipantSupply,
	}

	spMap, spIdList = simutil.SendProposeSupply(sclientParticipant, opts, spMap, spIdList)

}

// create sync id list
func createSyncIdList(participantsInfo []*participant.ParticipantInfo) ([]uint32, []uint32) {
	sameAreaIdList := make([]uint32, 0)
	neighborAreaIdList := make([]uint32, 0)

	for _, participantInfo := range participantsInfo {
		tAgentType := participantInfo.AgentType
		tAreaId := participantInfo.AreaId
		isNeighborArea := isContainNeighborMap(tAreaId)
		isSameArea := int(tAreaId) == *areaId && int(tAgentType) != *agentType
		if isNeighborArea {
			channelId := participantInfo.ChannelId
			agentChannelId := channelId.AgentChannelId
			neighborAreaIdList = append(neighborAreaIdList, agentChannelId)
		}
		if isSameArea {
			channelId := participantInfo.ChannelId
			agentChannelId := channelId.AgentChannelId
			sameAreaIdList = append(sameAreaIdList, agentChannelId)
		}
	}

	return sameAreaIdList, neighborAreaIdList
}

// Finish Fix
func setParticipant(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("setParticipant")
	setParticipantDemand := dm.GetArg_SetParticipantDemand()

	participantsInfo = setParticipantDemand.ParticipantsInfo
	fmt.Printf("ParticipantsInfo ", participantsInfo)
	idListByChannel = simutil.CreateIdListByChannel(participantsInfo)
	sameAreaIdList, neighborAreaIdList = createSyncIdList(participantsInfo)

	participantInfo := &participant.ParticipantInfo{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint32(sclientParticipant.ClientID),
			AreaChannelId:        uint32(sclientArea.ClientID),
			AgentChannelId:       uint32(sclientAgent.ClientID),
			ClockChannelId:       uint32(sclientClock.ClientID),
			RouteChannelId:       uint32(sclientRoute.ClientID),
		},
		ProviderType: 3, // PedArea
		AreaId:       uint32(*areaId),
		AgentType:    0, // Ped
	}

	setParticipantSupply := &participant.SetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "SetParticipant respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               dm.GetId(),
		Name:                 nm,
		JSON:                 js,
		SetParticipantSupply: setParticipantSupply,
	}

	spMap, spIdList = simutil.SendProposeSupply(sclientParticipant, opts, spMap, spIdList)

}

// Finish Fix
func getAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("getAgents")

	getAgentsDemand := dm.GetArg_GetAgentsDemand()
	log.Printf("info %v, %v, %v, %v", int(getAgentsDemand.AreaId), *areaId, int(getAgentsDemand.AgentType), *agentType)
	isNeighborArea := isContainNeighborMap(getAgentsDemand.AreaId)
	if (int(getAgentsDemand.AreaId) == *areaId && int(getAgentsDemand.AgentType) != *agentType) || isNeighborArea {

		getAgentsSupply := &agent.GetAgentsSupply{
			Time:       getAgentsDemand.Time,
			AreaId:     uint32(*areaId),
			AgentType:  0, //Ped
			AgentsInfo: data.AgentsInfo,
			StatusType: 0, //OK
			Meta:       "",
		}

		nm := "getAgentSupply  by ped-area-provider"
		js := ""
		opts := &sxutil.SupplyOpts{
			Target:          dm.GetId(),
			Name:            nm,
			JSON:            js,
			GetAgentsSupply: getAgentsSupply,
		}
		log.Printf("SEND_AGENTS_INFO to : areaId: %v, agentType: %v", getAgentsDemand.AreaId, getAgentsDemand.AgentType.String())

		spMap, spIdList = simutil.SendProposeSupply(sclientAgent, opts, spMap, spIdList)
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
		log.Println("demand callback is invalid.")
	}
}

// Finish Fix
// callback for each Supply
func proposeSupplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	fmt.Sprintln("getSupplyCallback")
	if simutil.IsSupplyTarget(sp, dmIdList) {
		supplyType := simutil.CheckSupplyType(sp)
		switch supplyType {
		case "GET_AREA_SUPPLY":
			ch <- sp
			fmt.Println("getArea")
		case "GET_AGENTS_SUPPLY":
			fmt.Println("getAgents response in callback")
			samePspMap[sp.SenderId] = sp
			if simutil.CheckFinishSync(samePspMap, sameAreaIdList) {
				isFinishSync = true
			}
		default:
			fmt.Println("error")
		}
	}

	// FORWARD_AGENTS_SUPPLY
	supplyType := simutil.CheckSupplyType(sp)
	if supplyType == "FORWARD_AGENTS_SUPPLY" {
		fmt.Println("forwardAgents response in callback")
		neighborPspMap[sp.SenderId] = sp
		if simutil.CheckFinishSync(neighborPspMap, neighborAreaIdList) {
			isFinishSync = true
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

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:PedArea, AreaId: %d}", *areaId)
	sclientAgent = sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE, argJson)
	sclientClock = sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE, argJson)
	sclientArea = sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE, argJson)
	sclientParticipant = sxutil.NewSMServiceClient(client, pb.ChannelType_PARTICIPANT_SERVICE, argJson)
	sclientRoute = sxutil.NewSMServiceClient(client, pb.ChannelType_ROUTE_SERVICE, argJson)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go simutil.SubscribeDemand(sclientAgent, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeDemand(sclientClock, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeDemand(sclientArea, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeDemand(sclientParticipant, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sclientArea, proposeSupplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sclientAgent, proposeSupplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sclientParticipant, proposeSupplyCallback, &wg)
	wg.Wait()

	// start up(setArea)
	wg.Add(1)
	setArea()
	setClock()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
