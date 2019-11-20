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
	idListByChannel     *simutil.IdListByChannel
	sameAreaIdList      []uint32
	neighborAreaIdList  []uint32
	isGetParticipant    bool
	pspMap              map[uint64]*pb.Supply
	neighborPspMap      map[uint64]*pb.Supply
	samePspMap          map[uint64]*pb.Supply
	selection           bool
	startCollectId      bool
	startSync           bool
	mu                  sync.Mutex
	ch                  chan *pb.Supply
	syncSameCh          chan *pb.Supply
	syncNeighborCh      chan *pb.Supply
	participantsInfo    []*participant.ParticipantInfo
	isFinishSync        bool
	CHANNEL_BUFFER_SIZE int
	sprovider           *simutil.SynerexProvider
	sim                 *simutil.SynerexSimulator
)

func init() {
	CHANNEL_BUFFER_SIZE = 10
	isGetParticipant = false
	pspMap = make(map[uint64]*pb.Supply)
	neighborPspMap = make(map[uint64]*pb.Supply)
	samePspMap = make(map[uint64]*pb.Supply)
	selection = false
	startCollectId = false
	startSync = false
	isFinishSync = false
	ch = make(chan *pb.Supply)
	syncNeighborCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	syncSameCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
}

func isContainNeighborMap(areaId uint32) bool {
	neighborMap := sim.Area.NeighborArea
	for _, neighborId := range neighborMap {
		if areaId == neighborId {
			return true
		}
	}
	return false
}

func wait(pspMap map[uint64]*pb.Supply, idList []uint32, syncCh chan *pb.Supply) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case psp := <-syncCh:
				mu.Lock()
				pspMap[psp.SenderId] = psp
				if simutil.CheckFinishSync(pspMap, idList) {

					mu.Unlock()
					wg.Done()
					return
				}
				mu.Unlock()

			}
		}
	}()
	wg.Wait()
}

// Finish Fix
// when start up,
func setArea() {
	log.Println("setArea")

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

	sp := <-ch
	//wait()
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

	sprovider.SetAgentsSupply(dm.GetId(), 1, uint32(*areaId))
}

// Finish Fix
func setClock() {
	log.Println("setClock")

	// store AgentInfo data
	sim.TimeStep = 1
	sim.GlobalTime = float64(*clockTime)
	//data.ClockInfo = clockInfo
}

// Finish Fix
func calcNextRoute(areaInfo *area.AreaInfo, agentInfo *agent.AgentInfo, otherAgentsInfo []*agent.AgentInfo) *agent.Route {

	route := agentInfo.Route
	speed := route.Speed
	currentLocation := route.Coord
	nextTransit := route.RouteInfo.NextTransit
	transitPoint := route.RouteInfo.TransitPoint
	destination := route.Destination
	// passed all transit point
	if nextTransit != nil {
		destination = nextTransit
	}

	direction, distance := simutil.CalcDirectionAndDistance(currentLocation.Lat, currentLocation.Lon, destination.Lat, destination.Lon)
	//newLat, newLon := simutil.CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, speed*1000/3600, direction)
	newLat, newLon := simutil.CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, destination.Lat, destination.Lon, distance, speed)

	// upate next trasit point
	if distance < 10 {
		if nextTransit != nil {
			nextTransit2 := nextTransit
			for i, tPoint := range transitPoint {
				if tPoint.Lon == nextTransit2.Lon && tPoint.Lat == nextTransit2.Lat {
					if i+1 == len(transitPoint) {
						// pass all transit point
						nextTransit = nil
					} else {
						// go to next transit point
						nextTransit = transitPoint[i+1]
					}
				}
			}
		} else {
			log.Printf("\x1b[30m\x1b[47m Arrived Destination! \x1b[0m\n")
		}

	}

	nextCoord := &agent.Coord{
		Lat: currentLocation.Lat,
		Lon: currentLocation.Lon,
	}
	//TODO: Fix this
	if newLat < 40 && newLat > 0 && newLon < 150 && newLon > 0 {
		nextCoord = &agent.Coord{
			Lat: newLat,
			Lon: newLon,
		}
	} else {
		log.Printf("\x1b[30m\x1b[47m LOCATION CULC ERROR %v \x1b[0m\n", nextCoord)

	}

	routeInfo := &agent.RouteInfo{
		TransitPoint:  transitPoint,
		NextTransit:   nextTransit,
		TotalDistance: route.RouteInfo.TotalDistance,
		RequiredTime:  route.RouteInfo.RequiredTime,
	}
	nextRoute := &agent.Route{
		Coord:       nextCoord,
		Direction:   float32(direction),
		Speed:       float32(speed),
		Destination: route.Destination,
		Departure:   route.Departure,
		RouteInfo:   routeInfo,
	}
	return nextRoute
}

func getAreaInfo() *area.AreaInfo {

	sprovider.GetAreaDemand(uint32(*clockTime), uint32(*areaId))

	sp := <-ch
	//wait()
	log.Println("GET_AREA_FINISH")
	getAreaSupply := sp.GetArg_GetAreaSupply()
	areaInfo := getAreaSupply.AreaInfo

	return areaInfo
}

func getAgentsInfo() []*agent.GetAgentsSupply {

	if len(sameAreaIdList) != 0 {

		// get Agent
		sprovider.GetAgentsDemand(uint32(1), uint32(*areaId))

		wait(samePspMap, sameAreaIdList, syncSameCh)
	} else {
		//log.Println("SAME_AREA_NOTHING")
	}
	getAgentsSupplies := make([]*agent.GetAgentsSupply, 0)
	for _, agentPsp := range samePspMap {
		getAgentsSupply := agentPsp.GetArg_GetAgentsSupply()
		getAgentsSupplies = append(getAgentsSupplies, getAgentsSupply)
	}
	return getAgentsSupplies
}

func updateAgentsInfo(pureNextAgentsInfo []*agent.AgentInfo, neighborPspMap map[uint64]*pb.Supply, currentAreaInfo *area.AreaInfo) []*agent.AgentInfo {
	nextAgentsInfo := pureNextAgentsInfo
	for _, psp := range neighborPspMap {
		forwardAgentsSupply := psp.GetArg_ForwardAgentsSupply()
		agentsInfo := forwardAgentsSupply.AgentsInfo
		for _, neighborAgentInfo := range agentsInfo {
			//　隣のエージェントが自分のエリアにいてかつ自分のエリアのエージェントと被ってない場合更新
			//log.Printf("IS_AGENT_IN_AREA: %v", simutil.IsAgentInArea(neighborAgentInfo, sim.Area, int32(*agentType)))
			if len(pureNextAgentsInfo) == 0 {
				if simutil.IsAgentInArea(neighborAgentInfo, currentAreaInfo, int32(*agentType)) {
					//log.Println("CHANGE_AREA1!!!")
					nextAgentsInfo = append(nextAgentsInfo, neighborAgentInfo)
				}
			} else {
				isAppendAgent := false
				for _, sameAreaAgent := range pureNextAgentsInfo {
					if neighborAgentInfo.AgentId != sameAreaAgent.AgentId && simutil.IsAgentInArea(neighborAgentInfo, currentAreaInfo, int32(*agentType)) {
						isAppendAgent = true
					}
				}
				if isAppendAgent {
					//log.Println("CHANGE_AREA2!!!")
					nextAgentsInfo = append(nextAgentsInfo, neighborAgentInfo)
				}
			}
		}
	}
	return nextAgentsInfo
}

func calcAgentsInfo(currentAgentsInfo []*agent.AgentInfo, currentAreaInfo *area.AreaInfo, currentTime uint32, sameAreaAgentsSupply []*agent.GetAgentsSupply) []*agent.AgentInfo {
	// calc agent
	otherAgentsInfo := make([]*agent.AgentInfo, 0)
	pureNextAgentsInfo := make([]*agent.AgentInfo, 0)
	for _, agentInfo := range currentAgentsInfo {
		// 自エリアにいる場合、次のルートを計算する
		if simutil.IsAgentInControlledArea(agentInfo, currentAreaInfo, int32(*agentType)) {

			nextRoute := calcNextRoute(currentAreaInfo, agentInfo, otherAgentsInfo)

			pureNextAgentInfo := &agent.AgentInfo{
				Time:        currentTime + 1,
				AgentId:     agentInfo.AgentId,
				AgentType:   agentInfo.AgentType,
				AgentStatus: agentInfo.AgentStatus,
				Route:       nextRoute,
			}

			pureNextAgentsInfo = append(pureNextAgentsInfo, pureNextAgentInfo)
		}
	}

	return pureNextAgentsInfo
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
	currentAreaInfo := sim.Area
	currentAgentsInfo := sim.Agents
	nextTime := time + 1
	// update Agents data
	sameAreaAgentsSupply := getAgentsInfo()

	// 次の時間のエージェントを計算する。重複エリアの更新をすませてないのでpureNextAgentsInfoとしている
	pureNextAgentsInfo := calcAgentsInfo(currentAgentsInfo, currentAreaInfo, time, sameAreaAgentsSupply)

	sprovider.ForwardAgentsSupply(dm.GetId(), nextTime, uint32(*areaId), pureNextAgentsInfo)

	if len(neighborAreaIdList) != 0 {
		//syncNeighborCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
		wait(neighborPspMap, neighborAreaIdList, syncNeighborCh)
	} else {
		//log.Println("NEIGHBOR_AREA_NOTHING")
	}

	// update Agents data
	nextAgentsInfo := updateAgentsInfo(pureNextAgentsInfo, neighborPspMap, currentAreaInfo)

	// propose clockInfo
	sprovider.ForwardClockSupply(dm.GetId(), currentClockInfo)
	sim.Agents = nextAgentsInfo
	sim.GlobalTime = float64(nextTime)

	neighborPspMap = make(map[uint64]*pb.Supply)
	samePspMap = make(map[uint64]*pb.Supply)
	log.Printf("\x1b[30m\x1b[47m \n FORWARD_CLOCK_FINISH \n TIME: %v \n AGENT_NUM: %v \x1b[0m\n", time, len(pureNextAgentsInfo))

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

// create sync id list
func createSyncIdList(participantsInfo []*participant.ParticipantInfo) ([]uint32, []uint32) {
	sameAreaIdList := make([]uint32, 0)
	neighborAreaIdList := make([]uint32, 0)

	for _, participantInfo := range participantsInfo {
		tAgentType := participantInfo.AgentType
		tAreaId := participantInfo.AreaId
		isNeighborArea := isContainNeighborMap(tAreaId) && int(tAgentType) == *agentType
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
	sameAreaIdList, neighborAreaIdList = createSyncIdList(participantsInfo)
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
	log.Println("getAgents")

	getAgentsDemand := dm.GetArg_GetAgentsDemand()
	//log.Printf("info %v, %v, %v, %v", int(getAgentsDemand.AreaId), *areaId, int(getAgentsDemand.AgentType), *agentType)
	isNeighborArea := isContainNeighborMap(getAgentsDemand.AreaId)
	isSameArea := int(getAgentsDemand.AreaId) == *areaId && int(getAgentsDemand.AgentType) != *agentType
	if isSameArea || isNeighborArea {

		agentsInfo := sim.Agents
		sprovider.GetAgentsSupply(dm.GetId(), getAgentsDemand.Time, uint32(*areaId), agentsInfo)
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
func proposeSupplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	supplyType := simutil.CheckSupplyType(sp)
	if simutil.IsSupplyTarget(sp, sprovider.DmIdList) {
		switch supplyType {
		case "GET_AREA_SUPPLY":
			ch <- sp
		case "GET_AGENTS_ROUTE_SUPPLY":
			ch <- sp
		case "GET_AGENTS_SUPPLY":
			getAgentsSupply := sp.GetArg_GetAgentsSupply()
			if getAgentsSupply.AreaId == uint32(*areaId) {
				syncSameCh <- sp
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
			syncNeighborCh <- sp

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
	go simutil.SubscribeSupply(sprovider.AreaClient, proposeSupplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sprovider.AgentClient, proposeSupplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sprovider.ParticipantClient, proposeSupplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sprovider.RouteClient, proposeSupplyCallback, &wg)
	wg.Wait()

	// start up(setArea)
	wg.Add(1)
	setArea()
	setClock()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
