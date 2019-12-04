package main

import (
	//"context"
	"flag"
	"log"
	"sync"

	//"math/rand"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/provider"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"

	//"time"
	"fmt"
	"math/rand"
	"sort"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	sprovider  *provider.SynerexProvider
)

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
		ProviderType: 4, // Route
	}

	sprovider.GetParticipantSupply(dm.GetId(), participantInfo)
}

// Finish Fix
func setParticipant(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("setParticipant")

	participantInfo := &participant.ParticipantInfo{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint32(sprovider.ParticipantClient.ClientID),
			AreaChannelId:        uint32(sprovider.AreaClient.ClientID),
			AgentChannelId:       uint32(sprovider.AgentClient.ClientID),
			ClockChannelId:       uint32(sprovider.ClockClient.ClientID),
			RouteChannelId:       uint32(sprovider.RouteClient.ClientID),
		},
		ProviderType: 4, // Route
	}

	sprovider.SetParticipantSupply(dm.GetId(), participantInfo)
}

func calcRoutes(sLat float32, sLon float32, gLat float32, gLon float32, agentType int32) *agent.RouteInfo {
	dLat := gLat - sLat
	dLon := gLon - sLon
	transitNum := 0
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
	//通常用
	if transitNum == 0 {
		coord := &agent.Coord{
			Lat: float32(gLat),
			Lon: float32(gLon),
		}
		transitPoint = append(transitPoint, coord)
		nextTransit = coord
	}

	// 壁を通過するテスト用
	/*if transitNum == 0 {
		coord := &agent.Coord{
			Lat: float32(35.156208),
			Lon: float32(136.982500),
		}
		transitPoint = append(transitPoint, coord)
		nextTransit = coord
	}*/
	for i := 0; i < transitNum; i++ {
		coord := &agent.Coord{
			Lat: float32(latArray[i]),
			Lon: float32(lonArray[i]),
		}
		transitPoint = append(transitPoint, coord)
		if i == 0 {
			nextTransit = coord
		}
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

	sprovider.GetAgentsRouteSupply(dm.GetId(), newAgentsInfo)
}

// Finish Fix
func setAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("setAgent")
	sprovider.SetAgentsSupply(dm.GetId(), 0, 0, agent.AgentType(0))
}

// Finish Fix
func forwardClock(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("forwardClock")
	forwardClockDemand := dm.GetArg_ForwardClockDemand()
	time := forwardClockDemand.Time
	nextTime := time + 1
	// calculation  area here
	log.Printf("\x1b[30m\x1b[47m \n FORWARD_CLOCK_FINISH \n TIME: %v \x1b[0m\n", time)

	// propose next clock
	nextClockInfo := &clock.ClockInfo{
		Time: nextTime,
	}

	sprovider.ForwardClockSupply(dm.GetId(), nextClockInfo)

	sprovider.ForwardAgentsSupply(dm.GetId(), 0, 0, []*agent.AgentInfo{}, agent.AgentType(0))
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	demandType := sprovider.CheckDemandType(dm)
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
		//log.Println("demand callback is invalid.")
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

	// Clientとして登録
	sprovider = provider.NewSynerexProvider()

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	sprovider.SetupProvider(client, argJson, demandCallback, func(clt *sxutil.SMServiceClient, sp *pb.Supply) {}, &wg)
	wg.Wait()

	wg.Add(1)
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
