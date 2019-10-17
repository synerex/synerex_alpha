package main

import (
	"flag"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	"github.com/synerex/synerex_alpha/sxutil"

	//	"github.com/synerex/synerex_alpha/api/simulation/area"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	gosocketio "github.com/mtfelian/golang-socketio"
	"github.com/mtfelian/golang-socketio/transport"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"google.golang.org/grpc"
	//	"os/exec"
)

var (
	serverAddr        = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv           = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	port              = flag.Int("port", 10080, "HarmoVis Provider Listening Port")
	clockTime         = flag.Int("time", 1, "Time")
	version           = "0.01"
	idlist            []uint64
	dmMap             map[uint64]*sxutil.DemandOpts
	spMap             map[uint64]*pb.Supply
	pspMap            map[uint64]*pb.Supply
	agentPspMap       map[uint64]*pb.Supply
	participantPspMap map[uint64]*pb.Supply
	forwardPspMap     map[uint64]*pb.Supply
	ch                chan *pb.Supply
	startCollectId    bool
	isStop            bool
	isStart           bool
	isSetClock        bool
	isSetArea         bool
	isSetAgent        bool
	isFinishSync      bool

	isGetParticipant   bool
	Order              string
	mu                 sync.Mutex
	sclientArea        *sxutil.SMServiceClient
	sclientAgent       *sxutil.SMServiceClient
	sclientClock       *sxutil.SMServiceClient
	sclientParticipant *sxutil.SMServiceClient
	sclientRoute       *sxutil.SMServiceClient
	participantsInfo   []*participant.ParticipantInfo
	sioClient          *gosocketio.Client
	idListByChannel    *simutil.IdListByChannel
	data               *Data
	startSync          bool
	assetsDir          http.FileSystem
	ioserv             *gosocketio.Server
)

func init() {
	Order = ""
	idlist = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*pb.Supply)
	pspMap = make(map[uint64]*pb.Supply)
	agentPspMap = make(map[uint64]*pb.Supply)
	participantPspMap = make(map[uint64]*pb.Supply)
	forwardPspMap = make(map[uint64]*pb.Supply)
	isFinishSync = false
	isStop = false
	isSetClock = false
	isSetAgent = false
	isSetArea = false
	isStart = false
	isGetParticipant = false
	startSync = false
	ch = make(chan *pb.Supply)
	data = new(Data)
	data.ClockInfo = &clock.ClockInfo{
		Time: uint32(*clockTime),
	}
}

type Data struct {
	ClockInfo *clock.ClockInfo
}

type ChannelIdList struct {
	ParticipantChannelId uint32
	AreaChannelId        uint32
	AgentChannelId       uint32
	ClockChannelId       uint32
	RouteChannelId       uint32
}

type MapMarker struct {
	mtype int32   `json:"mtype"`
	id    int32   `json:"id"`
	lat   float32 `json:"lat"`
	lon   float32 `json:"lon"`
	angle float32 `json:"angle"`
	speed int32   `json:"speed"`
	area  int32   `json:"area"`
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
/*func orderSetClock(clockInfo simutil.ClockInfo) {
	if isGetParticipant == false {
		fmt.Printf("Error... please order getParticipant")
	} else {
		firstTime := uint32(1)

		clockInfo := &clock.ClockInfo{
			Time:          firstTime,
			CycleNum:      uint32(1),
			CycleDuration: uint32(1),
			CycleInterval: uint32(1),
		}
		data.ClockInfo = clockInfo

		setClockDemand := &clock.SetClockDemand{
			Time:          firstTime,
			CycleDuration: uint32(1),
			StatusType:    2, // NONE
			Meta:          "",
		}

		nm := "SetClockDemand"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, SetClockDemand: setClockDemand}

		dmMap, idlist = simutil.SendDemand(sclientClock, opts, dmMap, idlist)
	}
}*/

// Finish Fix
func orderStartClock() {
	if isGetParticipant == false {
		fmt.Printf("Error... please order getParticipant")
	} else if isSetAgent == false {
		fmt.Printf("Error... please order setAgent")
	} else {
		isStart = true
		// forward clock
		forwardClockDemand := &clock.ForwardClockDemand{
			Time:       data.ClockInfo.Time,
			CycleNum:   uint32(1),
			StatusType: 2, // NONE
			Meta:       "",
		}

		nm := "ForwardClockDemand"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, ForwardClockDemand: forwardClockDemand}

		dmMap, idlist = simutil.SendDemand(sclientClock, opts, dmMap, idlist)

		// calc next time
		nextTime := data.ClockInfo.Time + 1
		clockInfo := &clock.ClockInfo{
			Time:          uint32(nextTime),
			CycleNum:      uint32(1),
			CycleDuration: uint32(1),
			CycleInterval: uint32(1),
		}
		data.ClockInfo = clockInfo

		fmt.Printf("wait now...")
		wait()

		fmt.Printf("Callback StartClock! Regist and Send Data to Simulator %v", pspMap)
		for _, psp := range pspMap {
			supplyType := simutil.CheckSupplyType(psp)
			switch supplyType {
			case "FORWARD_AGENTS_SUPPLY":
				fmt.Println("registAgents")
				sendToSimulator(psp)
			default:
				fmt.Println("SupplyType is invalid")
			}
		}

		fmt.Printf("Callback StartClock! move next clock")
		time.Sleep(1 * time.Second)
		if isStop == false {
			log.Printf("FORWARD_CLOCK")
			orderStartClock()
		} else {
			fmt.Printf("Stop clock!")
			isStop = false
			isStart = false
		}

		// clear pspMap
		forwardPspMap = make(map[uint64]*pb.Supply)
	}
}

// Finish Fix
func orderStopClock() {
	if isGetParticipant == false {
		fmt.Printf("Error... please order getParticipant")
	} else {
		if isStart == true {
			isStop = true
		} else {
			fmt.Printf("Clock has not been started. ")
		}
	}
}

// Fix  : this don't need ?
/*func orderSetArea(areaInfo simutil.AreaInfo) {
	if isGetParticipant == false {
		fmt.Printf("Error... please order getParticipant")
	} else {
		areaDemand := area.AreaDemand{
			Time: uint32(1),
			AreaId: areaInfo.Id, // A
			DemandType: 0, // SET
			StatusType: 2, // NONE
			Meta: "",
		}

		nm := "setArea order by scenario"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, AreaDemand: &areaDemand}

		dmMap, idlist = simutil.SendDemand(sclientArea, opts, dmMap, idlist)
	}
}*/

// Finish Fix
func orderSetAgents(agentsInfo []*agent.AgentInfo) {
	if isGetParticipant == false {
		fmt.Printf("Error... please order getParticipant")
	} else {

		setAgentsDemand := &agent.SetAgentsDemand{
			AgentsInfo: agentsInfo,
			StatusType: 2, // NONE
			Meta:       "",
		}

		nm := "SetAgentDemand"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, SetAgentsDemand: setAgentsDemand}

		dmMap, idlist = simutil.SendDemand(sclientAgent, opts, dmMap, idlist)

		wait()
		fmt.Printf("Callback SetAgent! return OK to Simulation Synerex Engine")
		isSetAgent = true
		agentPspMap = make(map[uint64]*pb.Supply)
	}
}

// Finish Fix
func orderGetParticipant() {
	participantInfo := &participant.ParticipantInfo{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint32(sclientParticipant.ClientID),
			AreaChannelId:        uint32(sclientArea.ClientID),
			AgentChannelId:       uint32(sclientAgent.ClientID),
			ClockChannelId:       uint32(sclientClock.ClientID),
			RouteChannelId:       uint32(sclientRoute.ClientID),
		},
		ProviderType: 0, //Scenario
	}
	getParticipantDemand := &participant.GetParticipantDemand{
		ParticipantInfo: participantInfo,
		StatusType:      2, // NONE
		Meta:            "",
	}

	nm := "GetParticipantDemand"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, GetParticipantDemand: getParticipantDemand}

	dmMap, idlist = simutil.SendDemand(sclientParticipant, opts, dmMap, idlist)
}

// Finish Fix
func orderSetParticipant() {

	// send participantID to participant provider
	setParticipantDemand := &participant.SetParticipantDemand{
		ParticipantsInfo: participantsInfo,
		StatusType:       2, // NONE
		Meta:             "",
	}

	nm := "setParticipant order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SetParticipantDemand: setParticipantDemand}

	dmMap, idlist = simutil.SendDemand(sclientParticipant, opts, dmMap, idlist)

	wait()
	fmt.Printf("Callback SetParticipant! return OK to Simulation Synerex Engine")
	participantPspMap = make(map[uint64]*pb.Supply)

}

func (m *MapMarker) GetJson() string {
	s := fmt.Sprintf("{\"mtype\":%d,\"id\":%d,\"lat\":%f,\"lon\":%f,\"angle\":%f,\"speed\":%d,\"area\":%d}",
		m.mtype, m.id, m.lat, m.lon, m.angle, m.speed, m.area)
	return s
}

// Fix now
func sendToSimulator(psp *pb.Supply) {
	forwardAgentsSupply := psp.GetArg_ForwardAgentsSupply()
	agentsInfo := forwardAgentsSupply.AgentsInfo
	if agentsInfo != nil {
		fmt.Printf("\x1b[30m\x1b[47m agentsInfo is : %v\x1b[0m\n", agentsInfo)
		for _, agentInfo := range agentsInfo {
			mm := &MapMarker{
				mtype: int32(agentInfo.AgentType), // depends on type of Ped: 0, Car , 1, Train, 2, Bycycle 3
				id:    int32(agentInfo.AgentId),
				lat:   float32(agentInfo.Route.Coord.Lat),
				lon:   float32(agentInfo.Route.Coord.Lon),
				angle: float32(agentInfo.Route.Direction),
				speed: int32(agentInfo.Route.Speed),
				area:  int32(agentInfo.ControlArea),
			}
			mu.Lock()
			ioserv.BroadcastToAll("event", mm.GetJson())
			mu.Unlock()
		}
	}
}

// Finish Fix
func callbackStartClock(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	mu.Lock()
	spMap[sp.SenderId] = sp
	forwardPspMap[sp.SenderId] = sp
	mu.Unlock()

	clockIdList := idListByChannel.ClockIdList
	agentIdList := idListByChannel.AgentIdList
	syncIdList := append(clockIdList, agentIdList...)
	if simutil.CheckFinishSync(forwardPspMap, syncIdList) {
		isFinishSync = true
	}

}

// Finish Fix
func callbackSetAgent(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	log.Println("Got for set_agent callback")
	mu.Lock()
	spMap[sp.SenderId] = sp
	agentPspMap[sp.SenderId] = sp
	fmt.Printf("sp is: %v", sp)
	mu.Unlock()

	syncAgentIdList := idListByChannel.AgentIdList
	if simutil.CheckFinishSync(agentPspMap, syncAgentIdList) {
		isFinishSync = true
	}
}

// Finish Fix
func callbackSetParticipant(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	log.Println("Got for set_participant callback")
	mu.Lock()
	log.Println("setParticipant2")
	spMap[sp.SenderId] = sp
	participantPspMap[sp.SenderId] = sp
	mu.Unlock()
	log.Println("setParticipant3")
	syncParticipantIdList := idListByChannel.ParticipantIdList
	if simutil.CheckFinishSync(participantPspMap, syncParticipantIdList) {
		isFinishSync = true
	}
}

// Finish Fix
/*func callbackSetClock(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	log.Println("Got for set_clock callback")
	if clt.IsSupplyTarget(sp, idlist) {
		spMap[sp.SenderId] = sp

		callback := func(pspMap map[uint64]*pb.Supply) {
			fmt.Printf("Callback SetClock! return OK to Simulation Synerex Engine")
			isSetClock = true
		}
		syncClockIdList := idListByChannel.ClockIdList
		syncProposeSupply(sp, syncClockIdList, pspMap, callback)

	} else {
		log.Printf("This is not propose supply \n")
	}
}

func callbackSetArea(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	log.Println("Got for set_area callback")
	if clt.IsSupplyTarget(sp, idlist) {

		//		opts :=	dmMap[sp.TargetId]
		//		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp

		callback := func(pspMap map[uint64]*pb.Supply) {
			fmt.Printf("Callback SetArea! return OK to Simulation Synerex Engine")
			isSetArea = true
			//pspMap = make(map[uint64]*pb.Supply)
		}

		syncAreaIdList := idListByChannel.AreaIdList
		syncProposeSupply(sp, syncAreaIdList, pspMap, callback)

	} else {
		log.Printf("This is not propose supply \n")
	}
}*/

// Finish Fix
func collectParticipantId(clt *sxutil.SMServiceClient, d int) {

	for i := 0; i < d; i++ {
		time.Sleep(1 * time.Second)
		log.Printf("waiting... %v", i+1)
	}

	mu.Lock()
	participantsInfo = simutil.CreateParticipantsInfo(participantPspMap)
	idListByChannel = simutil.CreateIdListByChannel(participantsInfo)
	fmt.Printf("participantsInfo, %v", participantsInfo)
	fmt.Printf("idListByChannel ", idListByChannel)
	mu.Unlock()
	startCollectId = false
	isGetParticipant = true
	participantPspMap = make(map[uint64]*pb.Supply)

	// setParticipant
	Order = "SetParticipant"
	orderSetParticipant()
}

// Finish Fix
func callbackGetParticipant(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	//	log.Println("Got for get_participant callback")
	if idListByChannel != nil {
		log.Println("already corrected IdList!")
	} else {
		mu.Lock()
		if clt.IsSupplyTarget(sp, idlist) {
			spMap[sp.SenderId] = sp
			participantPspMap[sp.SenderId] = sp

			if !startCollectId {
				log.Println("start selection")
				startCollectId = true
				go collectParticipantId(clt, 2)
			}
		}
		mu.Unlock()
	}
}

// callback for each Supply
func proposeSupplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	//	log.Println("Got supply callback")
	if simutil.IsSupplyTarget(sp, idlist) {
		switch Order {
		case "GetParticipant":
			callbackGetParticipant(clt, sp)
		case "SetParticipant":
			callbackSetParticipant(clt, sp)
		case "SetAgent":
			callbackSetAgent(clt, sp)
		case "Start":
			callbackStartClock(clt, sp)
		default:
			fmt.Println("order is invalid")
		}
	}
}

func runServer() *gosocketio.Server {

	currentRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	d := filepath.Join(currentRoot, "mclient", "build")

	assetsDir = http.Dir(d)
	log.Println("AssetDir:", assetsDir)

	assetsDir = http.Dir(d)
	server := gosocketio.NewServer()

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected from %s as %s", c.IP(), c.Id())
		// do something.
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected from %s as %s", c.IP(), c.Id())
	})

	return server
}

func runClient() *gosocketio.Client {
	log.Println("RUN_CLIENT")
	var sioErr error
	sioClient, sioErr = gosocketio.Dial("ws://localhost:9995/socket.io/?EIO=3&transport=websocket", transport.DefaultWebsocketTransport())
	if sioErr != nil {
		log.Println("se: Error to connect with se-daemon. You have to start se-daemon first.") //,err)
		os.Exit(1)
	} else {
		log.Println("se: connect OK")
	}
	sioClient.On(gosocketio.OnConnection, func(c *gosocketio.Channel, param interface{}) {
		fmt.Println("Go socket.io connected ")
		c.Emit("setCh", "Scenario")
	})

	sioClient.On("scenario", func(c *gosocketio.Channel, order *simutil.Order) {
		log.Printf("get order is: %v\n", order)
		Order = order.Type
		switch Order {
		case "GetParticipant":
			//			fmt.Println("getParticipant")
			orderGetParticipant()
		/*case "SetTime":
			fmt.Println("setClock")
			clockInfo := order.ClockInfo
			orderSetClock(clockInfo)
		case "SetArea":
			fmt.Println("setArea")
			areaInfo := order.AreaInfo
			orderSetArea(areaInfo)*/
		case "SetAgent":
			//			fmt.Println("set agent")
			agentsInfo := simutil.ConvertAgentsInfo(order.AgentsInfo)
			orderSetAgents(agentsInfo)
		case "Start":
			//			fmt.Println("start clock")
			orderStartClock()
		case "Stop":
			//			fmt.Println("stop clock")
			orderStopClock()
		default:
			fmt.Println("error")
		}

	})

	sioClient.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel, param interface{}) {
		fmt.Println("Go socket.io disconnected ", c)
	})

	return sioClient
}

// assetsFileHandler for static Data
func assetsFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}

	file := r.URL.Path
	//	log.Printf("Open File '%s'",file)
	if file == "/" {
		file = "/index.html"
	}
	f, err := assetsDir.Open(file)
	if err != nil {
		log.Printf("can't open file %s: %v\n", file, err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Printf("can't open file %s: %v\n", file, err)
		return
	}
	http.ServeContent(w, r, file, fi.ModTime(), f)
}

func main() {

	flag.Parse()

	// connect to node server
	sxutil.RegisterNodeName(*nodesrv, "ScenarioProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })

	// run socket.io server
	ioserv = runServer()
	log.Printf("Running Sio Server..\n")
	if ioserv == nil {
		os.Exit(1)
	}

	// run socket.io client
	sioClient = runClient()
	log.Printf("Running Sio Client..\n")
	if sioClient == nil {
		os.Exit(1)
	}

	// connect to synerex server
	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Scenario}")
	sclientAgent = sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE, argJson)
	sclientClock = sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE, argJson)
	sclientArea = sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE, argJson)
	sclientParticipant = sxutil.NewSMServiceClient(client, pb.ChannelType_PARTICIPANT_SERVICE, argJson)
	sclientRoute = sxutil.NewSMServiceClient(client, pb.ChannelType_ROUTE_SERVICE, argJson)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go simutil.SubscribeSupply(sclientArea, proposeSupplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sclientAgent, proposeSupplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sclientClock, proposeSupplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go simutil.SubscribeSupply(sclientParticipant, proposeSupplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	serveMux := http.NewServeMux()

	serveMux.Handle("/socket.io/", ioserv)
	serveMux.HandleFunc("/", assetsFileHandler)
	//Order = "GetParticipant"
	//orderGetParticipant()
	log.Printf("Starting Harmoware VIS  Provider %s  on port %d", version, *port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), serveMux)
	if err != nil {
		log.Fatal(err)
	}

	//orderGetParticipant()
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
