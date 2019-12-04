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
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/provider"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"

	//"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	sprovider  *provider.SynerexProvider
)

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
			ParticipantChannelId: uint32(sprovider.ParticipantClient.ClientID),
			AreaChannelId:        uint32(sprovider.AreaClient.ClientID),
			AgentChannelId:       uint32(sprovider.AgentClient.ClientID),
			ClockChannelId:       uint32(sprovider.ClockClient.ClientID),
			RouteChannelId:       uint32(sprovider.RouteClient.ClientID),
		},
		ProviderType: 1, // Area
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
		ProviderType: 1, // Area
	}

	sprovider.SetParticipantSupply(dm.GetId(), participantInfo)
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

	sprovider.SetAgentsSupply(dm.GetId(), 0, 0, agent.AgentType(0))
}

// Finish Fix
func getAgents(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	log.Println("getAgent")

	sprovider.GetAgentsSupply(dm.GetId(), 0, 0, []*agent.AgentInfo{}, agent.AgentType(0))

	/*getAgentsSupply := &agent.GetAgentsSupply{
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

	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts, spMap, idlist)*/
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

			sprovider.GetAreaSupply(dm.GetId(), areaInfo)

			/*getAreaSupply := &area.GetAreaSupply{
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

			spMap, idlist = simutil.SendProposeSupply(sclientArea, opts, spMap, idlist)*/
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
		//log.Println("demand callback is invalid.")
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
