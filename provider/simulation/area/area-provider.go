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
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func readMapData() []Map {
	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile("map.json")
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	var mapData []Map

	if err := json.Unmarshal(bytes, &mapData); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("mapData is : %v\n", mapData)
	return mapData
}

type Data struct {
	AreaInfo  *area.AreaInfo
	ClockInfo *clock.ClockInfo
}

type Coord struct {
	StartLat float32 `json:"slat"`
	EndLat   float32 `json:"elat"`
	StartLon float32 `json:"slon"`
	EndLon   float32 `json:"elon"`
}

type Map struct {
	Id         uint32   `json:"id"`
	Name       string   `json:"name"`
	Coord      Coord    `json:"coord"`
	Controlled Coord    `json:"controlled"`
	Neighbor   []uint32 `json:"neighbor"`
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
		ProviderType: 1, // Area
	}

	getParticipantSupply := &participant.GetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "getParticipant respnse by area-provider"
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
		ProviderType: 1, // Area
	}

	setParticipantSupply := &participant.SetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "SetParticipant respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               dm.GetId(),
		Name:                 nm,
		JSON:                 js,
		SetParticipantSupply: setParticipantSupply,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientParticipant, opts, spMap, idlist)
}

// Finish Fix  , this is unused
/*func setArea(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setArea")
	argOneof := dm.GetArg_AreaDemand()
	for _, mapData := range readMapData(){
		if(mapData.Id == argOneof.AreaId){
			// send area info
			areaInfo := &area.AreaInfo{
				Time: argOneof.Time, // uint32(1)
				StatusType: 2, // NONE
				Meta: "",
			}

			nm := "setArea respnse by area-provider"
			js := ""
			opts := &sxutil.SupplyOpts{
				Target: dm.GetId(),
				Name: nm,
				JSON: js,
				AreaInfo: areaInfo,
			}

			spMap, idlist = simutil.SendProposeSupply(sclientArea, opts, spMap, idlist)
		}
	}
}*/

// Finish Fix
func setAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("setAgent")

	nm := "setAgent respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name:   nm,
		JSON:   js,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts, spMap, idlist)
}

// Finish Fix
func getAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("getAgent")

	getAgentsSupply := &agent.GetAgentsSupply{
		Time:       uint32(1),
		StatusType: 0, //OK
		Meta:       "",
	}

	nm := "getAgent respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:          dm.GetId(),
		Name:            nm,
		JSON:            js,
		GetAgentsSupply: getAgentsSupply,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts, spMap, idlist)
}

// Finish Fix
func getArea(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("getArea")
	getAreaDemand := dm.GetArg_GetAreaDemand()
	log.Printf("demand: ", getAreaDemand)
	for _, mapData := range readMapData() {
		if mapData.Id == getAreaDemand.AreaId {

			// send area info
			areaInfo := &area.AreaInfo{
				Time:         getAreaDemand.Time,
				AreaId:       getAreaDemand.AreaId,
				AreaName:     mapData.Name,
				NeighborArea: mapData.Neighbor,
				AreaCoord: &area.AreaCoord{
					StartLat: mapData.Coord.StartLat,
					StartLon: mapData.Coord.StartLon,
					EndLat:   mapData.Coord.EndLat,
					EndLon:   mapData.Coord.EndLon,
				},
				ControlAreaCoord: &area.AreaCoord{
					StartLat: mapData.Controlled.StartLat,
					StartLon: mapData.Controlled.StartLon,
					EndLat:   mapData.Controlled.EndLat,
					EndLon:   mapData.Controlled.EndLon,
				},
			}

			getAreaSupply := &area.GetAreaSupply{
				AreaInfo:   areaInfo,
				StatusType: 0, //OK
				Meta:       "",
			}

			nm := "getArea respnse by area-provider"
			js := ""
			opts := &sxutil.SupplyOpts{
				Target:        dm.GetId(),
				Name:          nm,
				JSON:          js,
				GetAreaSupply: getAreaSupply,
			}

			spMap, idlist = simutil.SendProposeSupply(sclientArea, opts, spMap, idlist)
		}
	}
}

/* this is unused
func setClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setClock")
	argOneof := dm.GetArg_ClockDemand()

	clockInfo := &clock.ClockInfo{
		Time: argOneof.Time,
		SupplyType: 2,
		StatusType: 0, // OK
		Meta: "",
	}

	// store ClockInfo data
	data.ClockInfo = clockInfo
	log.Printf("data.ClockInfo %v\n\n", data)

	nm := "setClock respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm,
		JSON: js,
		ClockInfo: clockInfo,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientClock, opts, spMap, idlist)
}*/

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

	nm2 := "forwardClock to ClockCh respnse by area-provider"
	js2 := ""
	opts2 := &sxutil.SupplyOpts{
		Target:             dm.GetId(),
		Name:               nm2,
		JSON:               js2,
		ForwardClockSupply: forwardClockSupply,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientClock, opts2, spMap, idlist)

	forwardAgentsSupply := &agent.ForwardAgentsSupply{
		Time:       nextTime,
		StatusType: 0, //OK
		Meta:       "",
	}

	nm3 := "forwardClock to AgentCh respnse by area-provider"
	js3 := ""
	opts3 := &sxutil.SupplyOpts{
		Target:              dm.GetId(),
		Name:                nm3,
		JSON:                js3,
		ForwardAgentsSupply: forwardAgentsSupply,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts3, spMap, idlist)
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	demandType := simutil.CheckDemandType(dm)
	switch demandType {
	case "GET_PARTICIPANT_DEMAND":
		getParticipant(clt, dm)
		//		case "SET_CLOCK_DEMAND": setClock(clt, dm)
	case "SET_PARTICIPANT_DEMAND":
		setParticipant(clt, dm)
	case "FORWARD_CLOCK_DEMAND":
		forwardClock(clt, dm)
		//		case "SET_AREA_DEMAND": setArea(clt, dm)
	case "GET_AREA_DEMAND":
		getArea(clt, dm)
	case "SET_AGENTS_DEMAND":
		setAgents(clt, dm)
	case "GET_AGENTS_DEMAND":
		getAgents(clt, dm)
	default:
		log.Println("demand callback is invalid.")
	}
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
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
