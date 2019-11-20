package main

import (
	//"context"
	"flag"
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

	//"encoding/json"
	"fmt"
)

var (
	serverAddr          = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv             = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	areaId              = flag.Int("areaId", 1, "Area Id")
	agentType           = flag.Int("agentType", 1, "Agent Type") // CAR
	clockTime           = flag.Int("time", 1, "Time")
	cycleNum            = flag.Int("num", 1, "Num")
	cycleInterval       = flag.Int("interval", 1, "Interval")
	cycleDuration       = flag.Int("duration", 1, "Duration")
	dmIdList            []uint64
	spIdList            []uint64
	myChannelIdList     []uint64
	idListByChannel     *simutil.IdListByChannel
	sameAreaIdList      []uint32
	neighborAreaIdList  []uint32
	isGetParticipant    bool
	dmMap               map[uint64]*sxutil.DemandOpts
	spMap               map[uint64]*sxutil.SupplyOpts
	pspMap              map[uint64]*pb.Supply
	neighborPspMap      map[uint64]*pb.Supply
	samePspMap          map[uint64]*pb.Supply
	startCollectId      bool
	selection           bool
	startSync           bool
	mu                  sync.Mutex
	ch                  chan *pb.Supply
	syncSameCh          chan *pb.Supply
	syncNeighborCh      chan *pb.Supply
	sclientArea         *sxutil.SMServiceClient
	sclientAgent        *sxutil.SMServiceClient
	sclientClock        *sxutil.SMServiceClient
	sclientRoute        *sxutil.SMServiceClient
	sclientParticipant  *sxutil.SMServiceClient
	participantsInfo    []*participant.ParticipantInfo
	data                *Data
	history             *History
	isFinishSync        bool
	CHANNEL_BUFFER_SIZE int
)

func init() {
	CHANNEL_BUFFER_SIZE = 10
	spIdList = make([]uint64, 0)
	dmIdList = make([]uint64, 0)
	myChannelIdList = make([]uint64, 0)
	isGetParticipant = false
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
	pspMap = make(map[uint64]*pb.Supply)
	neighborPspMap = make(map[uint64]*pb.Supply)
	samePspMap = make(map[uint64]*pb.Supply)
	isFinishSync = false
	selection = false
	startCollectId = false
	startSync = false
	ch = make(chan *pb.Supply)
	syncNeighborCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	syncSameCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	data = new(Data)
	data.AgentsInfo = make([]*agent.AgentInfo, 0)
	history = new(History)
	history.History = make(map[uint32]*Data)
}

type Data struct {
	AreaInfo     *area.AreaInfo
	ClockInfo    *clock.ClockInfo
	AgentsInfo   []*agent.AgentInfo
	OccupancyMap [][][]uint32
}

type History struct {
	CurrentTime uint32
	History     map[uint32]*Data
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

	nm := "getArea order by car-area-provider"
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

	// get route demand

	getAgentsRouteDemand := &agent.GetAgentsRouteDemand{
		AgentsInfo: agentsInfo,
		StatusType: 2, //NONE
		Meta:       "",
	}

	nm1 := "getAgentsRouteDemand order by ped-area-provider"
	js1 := ""
	opts1 := &sxutil.DemandOpts{Name: nm1, JSON: js1, GetAgentsRouteDemand: getAgentsRouteDemand}
	dmMap, dmIdList = simutil.SendDemand(sclientRoute, opts1, dmMap, dmIdList)

	sp := <-ch
	//wait()
	log.Println("GET_AGENTS_ROUTE_FINISH")
	getAgentsRouteSupply := sp.GetArg_GetAgentsRouteSupply()
	agentsInfo = getAgentsRouteSupply.AgentsInfo

	for _, agentInfo := range agentsInfo {
		if simutil.IsAgentInControlledArea(agentInfo, data.AreaInfo, int32(*agentType)) {
			agentInfo.ControlArea = uint32(*areaId)
		}

		if simutil.IsAgentInArea(agentInfo, data.AreaInfo, int32(*agentType)) {
			data.AgentsInfo = append(data.AgentsInfo, agentInfo)
		}
	}

	occupancyMap := simutil.CreateOccupancyMap(data.AreaInfo, data.AgentsInfo)
	data.OccupancyMap = occupancyMap

	if data.AreaInfo != nil && data.AgentsInfo != nil && data.ClockInfo != nil {
		time := data.ClockInfo.Time
		history.History[time] = data
		history.CurrentTime = time
		//log.Printf("\x1b[30m\x1b[47m History is : %v\x1b[0m\n", history)
	}

	setAgentsSupply := &agent.SetAgentsSupply{
		Time:       uint32(1),
		AreaId:     uint32(*areaId),
		AgentType:  1, // Car
		StatusType: 0,
		Meta:       "",
	}

	nm := "setAgentSupply by car-area-provider"
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
	log.Println("GET_AREA_FINISH")
	getAreaSupply := sp.GetArg_GetAreaSupply()
	areaInfo := getAreaSupply.AreaInfo

	return areaInfo
}

func getAgentsInfo() []*agent.GetAgentsSupply {

	if len(sameAreaIdList) != 0 {

		// get Agent
		getAgentsDemand := &agent.GetAgentsDemand{
			Time:       uint32(1),
			AreaId:     uint32(*areaId),
			AgentType:  1, //CAR
			StatusType: 2, // NONE
			Meta:       "",
		}

		nm4 := "getAgents order by car-area-provider"
		js4 := ""
		opts4 := &sxutil.DemandOpts{
			Name:            nm4,
			JSON:            js4,
			GetAgentsDemand: getAgentsDemand,
		}
		dmMap, dmIdList = simutil.SendDemand(sclientAgent, opts4, dmMap, dmIdList)

		//syncSameCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
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
	//log.Printf("NEIGHBOR_PSP: %v", neighborPspMap)
	//log.Println("PURE_AGENTS! %v\n\n", nextAgentsInfo)
	for _, psp := range neighborPspMap {
		forwardAgentsSupply := psp.GetArg_ForwardAgentsSupply()
		agentsInfo := forwardAgentsSupply.AgentsInfo
		for _, neighborAgentInfo := range agentsInfo {
			//log.Println("NEIGHBOR_AGENTS! %v\n\n", neighborAgentInfo)
			//log.Println("PURE_AGENTS! %v\n\n", nextAgentsInfo)
			//　隣のエージェントが自分のエリアにいてかつ自分のエリアのエージェントと被ってない場合更新
			//log.Printf("IS_AGENT_IN_AREA: %v", simutil.IsAgentInArea(neighborAgentInfo, data.AreaInfo, int32(*agentType)))
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
	//currentAgentsInfo := data.AgentsInfo
	//	otherAgentsInfo := spAgentArgOneof
	otherAgentsInfo := make([]*agent.AgentInfo, 0)
	pureNextAgentsInfo := make([]*agent.AgentInfo, 0)
	for _, agentInfo := range currentAgentsInfo {
		// calc next agentInfo
		//		route := agentInfo.Route
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
	currentTime := forwardClockDemand.Time
	currentClockInfo := data.ClockInfo
	currentAreaInfo := data.AreaInfo
	currentAgentsInfo := data.AgentsInfo
	nextTime := currentTime + 1
	//nextTime := time + 1

	// update Agents data
	sameAreaAgentsSupply := getAgentsInfo()

	// 次の時間のエージェントを計算する。重複エリアの更新をすませてないのでpureNextAgentsInfoとしている
	pureNextAgentsInfo := calcAgentsInfo(currentAgentsInfo, currentAreaInfo, currentTime, sameAreaAgentsSupply)

	forwardAgentsSupply := &agent.ForwardAgentsSupply{
		Time:       nextTime,
		AreaId:     uint32(*areaId),
		AgentType:  1, //Car
		AgentsInfo: pureNextAgentsInfo,
		StatusType: 0, //OK
		Meta:       "",
	}

	nm2 := "forwardClock to agentCh respnse by car-area-provider"
	js2 := ""
	opts2 := &sxutil.SupplyOpts{
		Target:              dm.GetId(),
		Name:                nm2,
		JSON:                js2,
		ForwardAgentsSupply: forwardAgentsSupply,
	}
	spMap, spIdList = simutil.SendProposeSupply(sclientAgent, opts2, spMap, spIdList)

	// update
	if len(neighborAreaIdList) != 0 {
		//syncNeighborCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
		wait(neighborPspMap, neighborAreaIdList, syncNeighborCh)
	} else {
		//log.Println("NEIGHBOR_AREA_NOTHING")
	}

	// update Agents data
	nextAgentsInfo := updateAgentsInfo(pureNextAgentsInfo, neighborPspMap, currentAreaInfo)

	// get nextArea Data   :: nextAreaInfo
	/*wait()
	log.Println("FORWARD_AGENTS_FINISH")
	// update Area data
	nextAreaInfo := getAreaInfo()*/

	// propose clockInfo
	forwardClockSupply := &clock.ForwardClockSupply{
		ClockInfo:  currentClockInfo,
		StatusType: 0, // OK
		Meta:       "",
	}

	nm3 := "forwardClock to clockCh respnse by car-area-provider"
	js3 := ""
	opts3 := &sxutil.SupplyOpts{
		Target:             dm.GetId(),
		Name:               nm3,
		JSON:               js3,
		ForwardClockSupply: forwardClockSupply,
	}
	spMap, spIdList = simutil.SendProposeSupply(sclientClock, opts3, spMap, spIdList)

	// update data and history
	occupancyMap := simutil.UpdateOccupancyMap(data.AreaInfo, nextAgentsInfo)
	data.OccupancyMap = occupancyMap

	nextClockInfo := &clock.ClockInfo{
		Time: nextTime,
	}
	data.AgentsInfo = nextAgentsInfo
	data.ClockInfo = nextClockInfo
	//data.AreaInfo = nextAreaInfo
	history.History[nextTime] = data
	history.CurrentTime = nextTime

	neighborPspMap = make(map[uint64]*pb.Supply)
	samePspMap = make(map[uint64]*pb.Supply)

	log.Printf("\x1b[30m\x1b[47m \n FORWARD_CLOCK_FINISH \n TIME: %v \n AGENT_NUM: %v \x1b[0m\n\n", currentTime, len(pureNextAgentsInfo))

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
		ProviderType: 2, // CarArea
		AreaId:       uint32(*areaId),
		AgentType:    1, // Car
	}

	getParticipantSupply := &participant.GetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "getParticipant respnse by car-area-provider"
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
	//idListByChannel = simutil.CreateIdListByChannel(participantsInfo)
	sameAreaIdList, neighborAreaIdList = createSyncIdList(participantsInfo)
	fmt.Printf("Same %v", sameAreaIdList)
	fmt.Printf("Neighbor %v", neighborAreaIdList)

	participantInfo := &participant.ParticipantInfo{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint32(sclientParticipant.ClientID),
			AreaChannelId:        uint32(sclientArea.ClientID),
			AgentChannelId:       uint32(sclientAgent.ClientID),
			ClockChannelId:       uint32(sclientClock.ClientID),
			RouteChannelId:       uint32(sclientRoute.ClientID),
		},
		ProviderType: 2, // CarArea
		AreaId:       uint32(*areaId),
		AgentType:    1, // Car
	}

	setParticipantSupply := &participant.SetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "SetParticipantSupply respnse by car-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               dm.GetId(),
		Name:                 nm,
		JSON:                 js,
		SetParticipantSupply: setParticipantSupply,
	}

	spMap, spIdList = simutil.SendProposeSupply(sclientParticipant, opts, spMap, spIdList)

}

// Fix now
func getAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("getAgents")

	getAgentsDemand := dm.GetArg_GetAgentsDemand()
	log.Printf("info %v, %v, %v, %v", int(getAgentsDemand.AreaId), *areaId, int(getAgentsDemand.AgentType), *agentType)
	isNeighborArea := isContainNeighborMap(getAgentsDemand.AreaId)
	isSameArea := int(getAgentsDemand.AreaId) == *areaId && int(getAgentsDemand.AgentType) != *agentType
	if isSameArea || isNeighborArea {

		getAgentsSupply := &agent.GetAgentsSupply{
			Time:       getAgentsDemand.Time,
			AreaId:     uint32(*areaId),
			AgentType:  1, //Car
			AgentsInfo: data.AgentsInfo,
			StatusType: 0, //OK
			Meta:       "",
		}

		nm := "getAgentSupply  by car-area-provider"
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
		//log.Println("demand callback is invalid.")
	}
}

// Finish Fix
// callback for each Supply
func proposeSupplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	supplyType := simutil.CheckSupplyType(sp)
	if simutil.IsSupplyTarget(sp, dmIdList) {
		//supplyType := simutil.CheckSupplyType(sp)
		switch supplyType {
		case "GET_AREA_SUPPLY":
			ch <- sp
		case "GET_AGENTS_ROUTE_SUPPLY":
			ch <- sp
		case "GET_AGENTS_SUPPLY":
			getAgentsSupply := sp.GetArg_GetAgentsSupply()
			if getAgentsSupply.AreaId == uint32(*areaId) {
				log.Printf("\x1b[30m\x1b[47m 1.3. GET_SAME_AREA_AGENTS\x1b[0m\n")
				syncSameCh <- sp
			}
		default:
			//fmt.Println("error")
		}
	}

	// FORWARD_AGENTS_SUPPLY
	if supplyType == "FORWARD_AGENTS_SUPPLY" {
		//mu.Lock()
		forwardAgentsSupply := sp.GetArg_ForwardAgentsSupply()
		// not equal areaId and not AreaProvider
		if forwardAgentsSupply.AreaId != uint32(*areaId) && forwardAgentsSupply.AreaId != 0 {
			syncNeighborCh <- sp
		}
	}

	/*else if supplyType == "FORWARD_AREA_SUPPLY" {
		fmt.Println("forwardArea response in callback")
		mu.Lock()
		neighborPspMap[sp.SenderId] = sp
		if simutil.CheckFinishSync(neighborPspMap, neighborAreaIdList) {
			isFinishSync = true
		}
		mu.Unlock()
	}*/

}

func main() {
	flag.Parse()
	log.Printf("area id is: %v, agent type is %v", *areaId, *agentType)

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

	wg.Add(1)
	go simutil.SubscribeSupply(sclientRoute, proposeSupplyCallback, &wg)
	wg.Wait()

	// start up(setArea)
	wg.Add(1)
	setArea()
	setClock()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
