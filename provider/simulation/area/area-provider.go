package main

import (
	//"context"
	"flag"
	"log"
	"sync"
	//"math/rand"

	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	//"time"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist     []uint64
	dmMap      map[uint64]*sxutil.DemandOpts
	spMap		map[uint64]*sxutil.SupplyOpts
	selection 	bool
	mu	sync.Mutex
	sclientArea *sxutil.SMServiceClient
	sclientAgent *sxutil.SMServiceClient
	sclientClock *sxutil.SMServiceClient
	sclientParticipant *sxutil.SMServiceClient
	AreaData []Area
	data *Data
)

func init() {
	idlist = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
	selection = false
	AreaData = make([]Area, 0)
	data = new(Data)
}

func readAreaData() []Area{
	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile("area_data.json")
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	var areaData []Area
	
	if err := json.Unmarshal(bytes, &areaData); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("areaData is : %v\n", areaData)
	return areaData
}

type Data struct {
	AreaInfo *area.AreaInfo
	ClockInfo *clock.ClockInfo
}

type Coord struct {
	StartLat float32	`json:"slat"`
	EndLat float32	`json:"elat"`
	StartLon float32	`json:"slon"`
	EndLon float32	`json:"elon"`
}

type Area struct {
	Id uint32	`json:"id"`
	Name string	`json:"name"`
	Coord Coord	`json:"coord"`
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
		AreaId: uint32(0), // Area A
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


func setArea(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setArea")
	argOneof := dm.GetArg_AreaDemand()
	for _, areaData := range readAreaData(){
		if(areaData.Id == argOneof.AreaId){
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
}

func setAgent(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("setAgent")
		
	nm := "setAgent respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
	}
	
	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts, spMap, idlist)
}


func getArea(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("getArea")
	argOneof := dm.GetArg_AreaDemand()
	log.Printf("demand: ", argOneof)
	for _, areaData := range readAreaData(){
		if(areaData.Id == argOneof.AreaId){
			//areaData := AreaData[argOneof.AreaId]
			mapInfo := area.Map{
				Coord: &area.Map_Coord{
					StartLat: areaData.Coord.StartLat,
					StartLon: areaData.Coord.StartLon, 
					EndLat: areaData.Coord.EndLat,
					EndLon: areaData.Coord.EndLon, 
				},
				MapInfo: uint32(0),
			}

			// send area info
			areaInfo := &area.AreaInfo{
				Time: argOneof.Time,
				AreaId: argOneof.AreaId, // A
				AreaName: areaData.Name,
				Map: &mapInfo,
				SupplyType: 1, // RES_GET
				StatusType: 0, // OK
				Meta: "",
			}

			// store AreaInfo data
			data.AreaInfo = areaInfo
			log.Printf("data.AreaInfo %v\n\n", data)
		
			nm := "getArea respnse by area-provider"
			js := ""
			opts := &sxutil.SupplyOpts{Name: nm, JSON: js, AreaInfo: areaInfo}
	
			spMap, idlist = simutil.SendSupply(sclientArea, opts, spMap, idlist)
		}
	}
}


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
}

func forwardClock(clt *sxutil.SMServiceClient, dm *pb.Demand){
	log.Println("forwardClock")
	dmArgOneof := dm.GetArg_ClockDemand()
	time := dmArgOneof.Time
	nextTime := time + 1
	// calculation  area here

	// propose next area
	areaInfo := data.AreaInfo
	mapInfo := area.Map{
		Coord: &area.Map_Coord{
			StartLat: areaInfo.Map.Coord.StartLat,
			StartLon: areaInfo.Map.Coord.StartLon, 
			EndLat: areaInfo.Map.Coord.EndLat,
			EndLon: areaInfo.Map.Coord.EndLon, 
		},
		MapInfo: uint32(0),
	}

	nextAreaInfo := &area.AreaInfo{
		Time: nextTime,
		AreaId: areaInfo.AreaId, // A
		AreaName: areaInfo.AreaName,
		Map: &mapInfo,
		SupplyType: 0, // RES_SET
		StatusType: 0, // OK
		Meta: "",
	}

	nm := "forwardClock to AreaCh respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm, 
		JSON: js, 
		AreaInfo: nextAreaInfo,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientArea, opts, spMap, idlist)


	// propose next clock
	nextClockInfo := clock.ClockInfo{
		Time: nextTime,
		SupplyType: 0,
		StatusType: 0, // OK
		Meta: "",
	}
	
	nm2 := "forwardClock to ClockCh respnse by area-provider"
	js2 := ""
	opts2 := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm2, 
		JSON: js2, 
		ClockInfo: &nextClockInfo,
	}

	spMap, idlist = simutil.SendProposeSupply(sclientClock, opts2, spMap, idlist)
	
	nm3 := "forwardClock to AgentCh respnse by area-provider"
	js3 := ""
	opts3 := &sxutil.SupplyOpts{
		Target: dm.GetId(),
		Name: nm3, 
		JSON: js3, 
	}

	spMap, idlist = simutil.SendProposeSupply(sclientAgent, opts3, spMap, idlist)
}


// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	demandType := simutil.CheckDemandArgOneOf(dm)
	switch demandType{
		case "GET_PARTICIPANT": getParticipant(clt, dm)
		case "SET_CLOCK": setClock(clt, dm)
		case "FORWARD_CLOCK": forwardClock(clt, dm)
		case "SET_AREA": setArea(clt, dm)
		case "GET_AREA": getArea(clt, dm)
		case "SET_AGENT": setAgent(clt, dm)
		//case "START_CLOCK": startClock(clt, dm)
		//case "FORWARD_CLOCK_OK": forwardClockOK(clt, dm)
		default: log.Println("demand callback is invalid.")
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

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
