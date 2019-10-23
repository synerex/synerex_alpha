package main

import (
	//"context"
	"flag"
	"log"
	"sync"

	//"math/rand"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"

	//"time"
	"fmt"
	"math/rand"
	"sort"
)

var (
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv            = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist             []uint64
	dmMap              map[uint64]*sxutil.DemandOpts
	spMap              map[uint64]*sxutil.SupplyOpts
	selection          bool
	mu                 sync.Mutex
	sclientArea        *sxutil.SMServiceClient
	sclientAgent       *sxutil.SMServiceClient
	sclientClock       *sxutil.SMServiceClient
	sclientRoute       *sxutil.SMServiceClient
	sclientParticipant *sxutil.SMServiceClient
	data               *Data
)

func init() {
	idlist = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
	selection = false
	data = new(Data)
}

type Data struct {
	AreaInfo  *area.AreaInfo
	ClockInfo *clock.ClockInfo
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
		ProviderType: 4, // Route
	}

	getParticipantSupply := &participant.GetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "getParticipant respnse by route-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               dm.GetId(),
		Name:                 nm,
		JSON:                 js,
		GetParticipantSupply: getParticipantSupply,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientParticipant, opts, spMap, idlist)
}

// Finish Fix
func setParticipant(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("setParticipant")

	participantInfo := &participant.ParticipantInfo{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint32(sclientParticipant.ClientID),
			AreaChannelId:        uint32(sclientArea.ClientID),
			AgentChannelId:       uint32(sclientAgent.ClientID),
			ClockChannelId:       uint32(sclientClock.ClientID),
			RouteChannelId:       uint32(sclientRoute.ClientID),
		},
		ProviderType: 4, // Route
	}

	setParticipantSupply := &participant.SetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "SetParticipant respnse by route-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               dm.GetId(),
		Name:                 nm,
		JSON:                 js,
		SetParticipantSupply: setParticipantSupply,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientParticipant, opts, spMap, idlist)
}

func calcRoutes(sLat float32, sLon float32, gLat float32, gLon float32, agentType int32) *agent.RouteInfo {
	dLat := gLat - sLat
	dLon := gLon - sLon
	transitNum := 5
	transitPoint := make([]*agent.Coord, 0)
	latArray := make([]float64, 0)
	lonArray := make([]float64, 0)
	for i := 0; i < transitNum; i++ {
		newLat := float64(sLat) + float64(dLat)*rand.Float64()
		newLon := float64(sLon) + float64(dLon)*rand.Float64()
		latArray = append(latArray, newLat)
		lonArray = append(lonArray, newLon)

	}
	if dLat > 0 {
		sort.SliceStable(latArray, func(i, j int) bool { return latArray[i] < latArray[j] })
	} else {
		sort.SliceStable(latArray, func(i, j int) bool { return latArray[i] > latArray[j] })
	}
	if dLon > 0 {
		sort.SliceStable(lonArray, func(i, j int) bool { return lonArray[i] < lonArray[j] })
	} else {
		sort.SliceStable(lonArray, func(i, j int) bool { return lonArray[i] > lonArray[j] })
	}

	nextTransit := &agent.Coord{}
	for i := 0; i < transitNum; i++ {
		coord := &agent.Coord{
			Lat: float32(latArray[i]),
			Lon: float32(lonArray[i]),
		}
		transitPoint = append(transitPoint, coord)
		if i == 0 {
			nextTransit = coord
		}
		log.Printf("\x1b[30m\x1b[47m coord is :i: %v:  %v\x1b[0m\n", i, coord)
	}

	routeInfo := &agent.RouteInfo{
		TransitPoint:  transitPoint,
		NextTransit:   nextTransit,
		TotalDistance: float32(10),
		RequiredTime:  float32(10),
	}

	return routeInfo
}

// Finish Fix
func getAgentsRoute(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("getRoute")

	getAgentsRouteDemand := dm.GetArg_GetAgentsRouteDemand()
	agentsInfo := getAgentsRouteDemand.AgentsInfo
	newAgentsInfo := make([]*agent.AgentInfo, 0)
	for _, agentInfo := range agentsInfo {
		route := agentInfo.Route
		departure := route.Departure
		destination := route.Destination
		agentType := agentInfo.AgentType

		routeInfo := calcRoutes(departure.Lat, departure.Lon, destination.Lat, destination.Lon, int32(agentType))

		agentInfo.Route.RouteInfo = routeInfo
		newAgentsInfo = append(newAgentsInfo, agentInfo)
	}

	getAgentsRouteSupply := &agent.GetAgentsRouteSupply{
		AgentsInfo: newAgentsInfo,
		StatusType: 0, //OK
		Meta:       "",
	}

	nm := "getRoute respnse by route-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               dm.GetId(),
		Name:                 nm,
		JSON:                 js,
		GetAgentsRouteSupply: getAgentsRouteSupply,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientRoute, opts, spMap, idlist)
}

// Finish Fix
func setAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("setAgent")

	nm := "setAgent respnse by route-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name:   nm,
		JSON:   js,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts, spMap, idlist)
}

// Finish Fix
func forwardClock(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("forwardClock")
	forwardClockDemand := dm.GetArg_ForwardClockDemand()
	time := forwardClockDemand.Time
	nextTime := time + 1
	// calculation  area here

	// propose next clock
	nextClockInfo := &clock.ClockInfo{
		Time: nextTime,
	}

	forwardClockSupply := &clock.ForwardClockSupply{
		ClockInfo:  nextClockInfo,
		StatusType: 0, //OK
		Meta:       "",
	}

	nm2 := "forwardClock to ClockCh respnse by route-provider"
	js2 := ""
	opts2 := &sxutil.SupplyOpts{
		Target:             dm.GetId(),
		Name:               nm2,
		JSON:               js2,
		ForwardClockSupply: forwardClockSupply,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientClock, opts2, spMap, idlist)

	nm3 := "forwardClock to AgentCh respnse by route-provider"
	js3 := ""
	opts3 := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name:   nm3,
		JSON:   js3,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts3, spMap, idlist)
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	demandType := simutil.CheckDemandType(dm)
	switch demandType {
	case "GET_AGENTS_ROUTE_DEMAND":
		getAgentsRoute(clt, dm)
	case "GET_PARTICIPANT_DEMAND":
		getParticipant(clt, dm)
	case "SET_PARTICIPANT_DEMAND":
		setParticipant(clt, dm)
	case "FORWARD_CLOCK_DEMAND":
		forwardClock(clt, dm)
	case "SET_AGENTS_DEMAND":
		setAgents(clt, dm)
	default:
		log.Println("demand callback is invalid.")
	}
}

func main() {
	flag.Parse()

	sxutil.RegisterNodeName(*nodesrv, "RouteProvider", false)

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
	argJson := fmt.Sprintf("{Client:Route}")
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
	go simutil.SubscribeDemand(sclientRoute, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
