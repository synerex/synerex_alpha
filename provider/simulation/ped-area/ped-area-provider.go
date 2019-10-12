package main

import (
	//"context"
	"flag"
	"log"
	"sync"
	"math"
	"math/rand"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"time"
	"fmt"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	areaId    = flag.Int("areaId", 1, "Area Id")	// Area B
	agentType    = flag.Int("agentType", 0, "Agent Type")	// PEDESTRIAN
	dmIdList     []uint64
	spIdList     []uint64
	myChannelIdList []uint64
	idListByChannel *simutil.IdListByChannel
	isGetParticipant bool
	dmMap      map[uint64]*sxutil.DemandOpts
	spMap		map[uint64]*sxutil.SupplyOpts
	pspMap		map[uint64]*pb.Supply
	selection 	bool
	startCollectId bool
    startSync 	bool
	mu	sync.Mutex
	ch 			chan *pb.Supply
	pspCh		chan map[uint64]*pb.Supply
	sclientArea *sxutil.SMServiceClient
	sclientAgent *sxutil.SMServiceClient
	sclientClock *sxutil.SMServiceClient
	sclientParticipant *sxutil.SMServiceClient
	data *Data
)

func init() {
	spIdList = make([]uint64, 0)
	dmIdList = make([]uint64, 0)
	myChannelIdList = make([]uint64, 0)
	isGetParticipant = false
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
	pspMap = make(map[uint64]*pb.Supply)
	selection = false
	startCollectId = false
	startSync = false
	ch = make(chan *pb.Supply)
	pspCh = make(chan map[uint64]*pb.Supply)
	data = new(Data)
	data.AgentsInfo = make([]*agent.AgentInfo, 0)
}


type Data struct {
	AreaInfo *area.AreaInfo
	ClockInfo *clock.ClockInfo
	AgentsInfo []*agent.AgentInfo
}

func isContainNeighborMap(target uint32) bool{
	tmap := data.AreaInfo.Map.Neighbor
	for _, t := range tmap{
		if target == t{
			return true
		}
	}
	return false
}

// IsFinishSync is a helper function to check if synchronization finish or not 
func isFinishSync(pspMap map[uint64]*pb.Supply, idlist []uint64) bool {
	
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			senderId := sp.SenderId
			agentsInfo := sp.GetArg_AgentsInfo()
			isNeighborArea := isContainNeighborMap(agentsInfo.AreaId)
			if id == senderId || isNeighborArea{
				log.Printf("match! %v %v",id, senderId)
				isMatch = true
			}
		}
		if isMatch == false {
			log.Printf("false")
			return false
		} 
	}
	return true
}

func syncProposeSupply(sp *pb.Supply, syncIdList []uint64, pspMap map[uint64]*pb.Supply, callback func(pspMap map[uint64]*pb.Supply)){
	go func() { 
		log.Println("Send Supply")
		ch <- sp
		return
	}()
	log.Printf("StartSync : %v", startSync)
	if !startSync {
		log.Println("Start Sync")
		startSync = true
		pspMap2 := make(map[uint64]*pb.Supply)

		go func(){		
		for {
			select {
				case psp := <- ch:
					log.Println("recieve ProposeSupply")
					pspMap2[psp.SenderId] = psp
//					log.Printf("waitidList %v %v", pspMap, idList)
				if isFinishSync(pspMap2, syncIdList){
					fmt.Printf("Finish Sync\n")

					// if you need, return response
					callback(pspMap2)

					// init pspMap
					//pspMap = make(map[uint64]*pb.Supply)
					startSync = false
					fmt.Printf("startSync to false: %v\n", startSync)
					
					return
				}
			}
		}

		}()
	}	
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
	dmMap, dmIdList = simutil.SendDemand(sclientArea, opts, dmMap, dmIdList)

	sp := <- ch
	log.Println("getArea response")
	spArgOneof := sp.GetArg_AreaInfo()
	if spArgOneof.AreaId == uint32(*areaId){
		areaInfo := &area.AreaInfo{
			Time: spArgOneof.Time,
			AreaId: uint32(*areaId), // A
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
	
		spMap, spIdList = simutil.SendProposeSupply(sclientArea, opts2, spMap, spIdList)
	}

	
}

// if agent type and coord satisfy, return true
func isAgentInArea(agentInfo *agent.AgentInfo, data *Data, agentType int) bool{
	lat := agentInfo.Route.Coord.Lat
	lon := agentInfo.Route.Coord.Lon
	slat := data.AreaInfo.Map.Coord.StartLat
	elat := data.AreaInfo.Map.Coord.EndLat
	slon := data.AreaInfo.Map.Coord.StartLon
	elon := data.AreaInfo.Map.Coord.EndLon
	if agentInfo.AgentType.String() == agent.AgentType_name[int32(agentType)] && slat <= lat && lat <= elat &&  slon <= lon && lon <= elon {
		return true
	}else{
		log.Printf("agent type and coord is not match...\n\n")
		return false
	}
}


// if agent type and coord satisfy, return true
func isAgentInControlledArea(agentInfo *agent.AgentInfo, data *Data, agentType int) bool{
	lat := agentInfo.Route.Coord.Lat
	lon := agentInfo.Route.Coord.Lon
	slat := data.AreaInfo.Map.Controlled.StartLat
	elat := data.AreaInfo.Map.Controlled.EndLat
	slon := data.AreaInfo.Map.Controlled.StartLon
	elon := data.AreaInfo.Map.Controlled.EndLon
	if agentInfo.AgentType.String() == agent.AgentType_name[int32(agentType)] && slat <= lat && lat <= elat &&  slon <= lon && lon <= elon {
		return true
	}
	log.Printf("agent type and coord is not match...\n\n")
	return false
}

func setAgent(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setAgent")
	dmArgOneof := dm.GetArg_AgentDemand()
	log.Printf("set agentInfo is %v\n\n", dmArgOneof)
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

	if isAgentInControlledArea(agentInfo, data, *agentType){
		agentInfo.ControllArea = uint32(*areaId)
	}

	// store AgentInfo data
	if isAgentInArea(agentInfo, data, *agentType){
		data.AgentsInfo = append(data.AgentsInfo, agentInfo)
		log.Printf("data.AgentInfo %v\n\n", data)
	}

	log.Printf("agentInfo is:! %v\n\n", agentInfo)
		
	nm := "setAgent respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target    : dm.GetId(),
		Name      : nm, 
		JSON      : js, 
		AgentInfo : agentInfo,  // FIX: why rpc error ?
	}
	
	spMap, spIdList = simutil.SendProposeSupply(sclientAgent, opts, spMap, spIdList)
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

	spMap, spIdList = simutil.SendProposeSupply(sclientClock, opts, spMap, spIdList)
}

func calcNextRoute(areaInfo *area.AreaInfo, agentInfo *agent.AgentInfo, otherAgentsInfo *agent.AgentsInfo) *agent.Route{

	route := agentInfo.Route
	nextCoord := &agent.Route_Coord{
		Lat: float32(float64(route.Coord.Lat) + float64(0.0001) * 1 * math.Cos(float64(route.Direction*math.Pi/180))),
		Lon: float32(float64(route.Coord.Lon) + float64(0.0001) * 1 * math.Sin(float64(route.Direction*math.Pi/180))), 
	}

	// change speed and direction
	speedFluct := 2
	directionFluct := 5
	rand.Seed(time.Now().UnixNano())
	speed := route.Speed + float32(rand.Intn(speedFluct)) - float32(rand.Intn(speedFluct)/2)
	direction := route.Direction + float32(rand.Intn(directionFluct)) - float32(rand.Intn(directionFluct)/2)


	nextRoute := &agent.Route{
		Coord: nextCoord,
		Direction: float32(direction),
		Speed: float32(speed),
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
	dmMap, dmIdList = simutil.SendDemand(sclientArea, opts, dmMap, dmIdList)

	sp := <- ch
	log.Println("getArea response")
	spArgOneof := sp.GetArg_AreaInfo()


	// get Agent
	agentsDemand := &agent.AgentsDemand{
		Time: uint32(1),
		AreaId: uint32(*areaId),
		AgentType: 0,
		DemandType: 1, // GET
		StatusType: 2, // NONE
		Meta: "",
	}
		
	nm4 := "getAgents order by ped-area-provider"
	js4 := ""
	opts4 := &sxutil.DemandOpts{
		Name: nm4, 
		JSON: js4, 
		AgentsDemand: agentsDemand,
	}
	log.Println("sendDemand AgentsDemand")
	dmMap, dmIdList = simutil.SendDemand(sclientAgent, opts4, dmMap, dmIdList)

	agentPspMap := <- pspCh
	log.Println("getAgents response")
	//spAgentsArgOneof := spAgent.GetArg_AgentsInfo()
	log.Printf("get Agents is: %v", agentPspMap)


	/*areaInfo := &area.AreaInfo{
		Time: spArgOneof.Time,
		AgentId: spArgOneof.AreaId, // A
		AreaName: spArgOneof.AreaName,
		Map: spArgOneof.Map,
		SupplyType: 0, // RES_GET
		StatusType: 0, // OK
		Meta: "",
	}*/
	
	// update Agents data


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
	sendAgentsInfo := make([]*agent.AgentInfo, 0)
	for _, agentInfo := range data.AgentsInfo{
		if isAgentInControlledArea(agentInfo, data, *agentType){
			sendAgentsInfo = append(sendAgentsInfo, agentInfo)
		}
	}

	// propose agentInfo
	nextAgentsInfo := &agent.AgentsInfo{
		Time: nextTime,
		AreaId: uint32(*areaId),
		AgentType: 0,
		AgentInfo: sendAgentsInfo,
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
	log.Printf("nextAgentsInfo! %v %v\n\n", nextAgentsInfo)
	spMap, spIdList = simutil.SendProposeSupply(sclientAgent, opts2, spMap, spIdList)
	
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
	spMap, spIdList = simutil.SendProposeSupply(sclientClock, opts3, spMap, spIdList)
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

	spMap, spIdList = simutil.SendProposeSupply(sclientParticipant, opts, spMap, spIdList)

	// get participant to same area provider
	if isGetParticipant == false{
		log.Println("getParticipant for same area")
		isGetParticipant = true
		participantDemand := participant.ParticipantDemand {
			ClientId: uint64(sclientParticipant.ClientID),
			DemandType: 0, // GET
			StatusType: 2, // NONE
			Meta: "",
		}
			
		nm := "getParticipant order by car area provider"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, ParticipantDemand: &participantDemand}
		
		dmMap, dmIdList = simutil.SendDemand(sclientParticipant, opts, dmMap, dmIdList)
	}
}

func setParticipant(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setParticipant")
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

	spMap, spIdList = simutil.SendProposeSupply(sclientParticipant, opts, spMap, spIdList)

	// get participant to same area provider
	if isGetParticipant == false{
		log.Println("getParticipant for same area")
		isGetParticipant = true
		participantDemand := participant.ParticipantDemand {
			ClientId: uint64(sclientParticipant.ClientID),
			DemandType: 0, // GET
			StatusType: 2, // NONE
			Meta: "",
		}
			
		nm := "getParticipant order by car area provider"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, ParticipantDemand: &participantDemand}
		
		dmMap, dmIdList = simutil.SendDemand(sclientParticipant, opts, dmMap, dmIdList)
	}
}

func getAgents(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("getAgents")
	dmArgOneof := dm.GetArg_AgentsDemand()
	log.Printf("info %v, %v, %v, %v", int(dmArgOneof.AreaId), *areaId, int(dmArgOneof.AgentType), *agentType)
	if int(dmArgOneof.AreaId) == *areaId && int(dmArgOneof.AgentType) != *agentType{
		log.Println("getAgents2")
		agentsInfo := &agent.AgentsInfo{
			Time: dmArgOneof.Time,
			AreaId: dmArgOneof.AreaId,
			AgentType: 0,
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
		
		spMap, spIdList = simutil.SendProposeSupply(sclientAgent, opts, spMap, spIdList)
	}
}

func collectParticipantId(clt *sxutil.SMServiceClient, d int){

	for i := 0; i < d; i++{
		time.Sleep(1 * time.Second)
		log.Printf("waiting... %v",i+1)
	}

	mu.Lock()
	idListByChannel = simutil.CreateIdListByChannel(pspMap)
	mu.Unlock()
	log.Printf("finish collecting Id %v", idListByChannel)
	log.Printf("clock Id %v", idListByChannel.ClockIdList)
	log.Printf("area Id %v", idListByChannel.AreaIdList)
	log.Printf("agent Id %v", idListByChannel.AgentIdList)
	startCollectId = false
	pspMap = make(map[uint64]*pb.Supply)
}

/*func isContainNeighborMap(target uint32) bool{
	tmap := data.AreaInfo.Map.Neighbor
	for _, t := range tmap{
		if target == t{
			return true
		}
	}
	return false
}*/

func callbackForGetParticipant(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for get_participant callback")
	
	if idListByChannel != nil{
		log.Println("already corrected IdList!")
	}else{
		mu.Lock()
		pInfo := sp.GetArg_ParticipantInfo()
		isSameArea := pInfo.AreaId == uint32(*areaId)
		isDiffAgentType := int(pInfo.AgentType) != *agentType
		//isNeighborArea := isContainNeighborMap(pInfo.AreaId)
		if clt.IsSupplyTarget(sp, dmIdList) && (pInfo != nil && isSameArea && isDiffAgentType){ 
			pspMap[sp.SenderId] = sp

			if !startCollectId {
				log.Println("start selection")
				startCollectId = true
				go collectParticipantId(clt, 2)
			}
		}else{
			log.Printf("This is not propose supply \n")
		}
		mu.Unlock()
	}
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	demandType := simutil.CheckDemandArgOneOf(dm)
	switch demandType{
		case "GET_PARTICIPANT": getParticipant(clt, dm)
		case "SET_PARTICIPANT": setParticipant(clt, dm)
		case "SET_CLOCK": setClock(clt, dm)
		case "FORWARD_CLOCK": forwardClock(clt, dm)
		case "SET_AREA": setArea(clt, dm)
		case "SET_AGENT": setAgent(clt, dm)
		case "GET_AGENTS": getAgents(clt, dm)
		default: log.Println("demand callback is invalid.")
	}
}


// callback for each Supply
func proposeSupplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	fmt.Sprintln("getSupplyCallback")
	if simutil.IsSupplyTarget(sp, dmIdList) {
		supplyType := simutil.CheckSupplyArgOneOf(sp)
		switch supplyType{
			case "RES_GET_AREA":
				ch <- sp
				fmt.Println("getArea")
			case "RES_GET_AGENTS":
				fmt.Println("getAgents response in callback")
				callback := func (pspMap map[uint64]*pb.Supply){
					fmt.Printf("Callback GetAgents!")
					pspCh <- pspMap
					pspMap = make(map[uint64]*pb.Supply)
				}
				syncAgentIdList := idListByChannel.AgentIdList
				syncProposeSupply(sp, syncAgentIdList, pspMap, callback)

			case "RES_GET_PARTICIPANT":
				fmt.Println("getParticipant response in callback")
				callbackForGetParticipant(clt, sp)
			default:
				fmt.Println("error")
		}
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
	myChannelIdList = append(myChannelIdList, []uint64{uint64(sclientAgent.ClientID), uint64(sclientClock.ClientID), uint64(sclientArea.ClientID), uint64(sclientParticipant.ClientID)}...)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go simutil.SubscribeDemand(sclientAgent, demandCallback)
	go simutil.SubscribeDemand(sclientClock, demandCallback)
	go simutil.SubscribeDemand(sclientArea, demandCallback)
	go simutil.SubscribeDemand(sclientParticipant, demandCallback)

	go simutil.SubscribeSupply(sclientArea, proposeSupplyCallback)
	go simutil.SubscribeSupply(sclientAgent, proposeSupplyCallback)
	go simutil.SubscribeSupply(sclientParticipant, proposeSupplyCallback)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
