package main

import (
	"flag"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/provider"
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
	serverAddr          = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv             = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	port                = flag.Int("port", 10080, "HarmoVis Provider Listening Port")
	clockTime           = flag.Int("time", 1, "Time")
	version             = "0.01"
	participantPspMap   map[uint64]*pb.Supply
	syncForwardCh       chan *pb.Supply
	syncParticipantCh   chan *pb.Supply
	syncSetAgentsCh     chan *pb.Supply
	startCollectId      bool
	isStop              bool
	isStart             bool
	isSetClock          bool
	isSetArea           bool
	isSetAgent          bool
	isGetParticipant    bool
	forwardAgentIdList  []uint64
	participantIdList   []uint64
	setAgentIdList      []uint64
	Order               string
	mu                  sync.Mutex
	participantsInfo    []*participant.ParticipantInfo
	data                *Data
	startSync           bool
	assetsDir           http.FileSystem
	ioserv              *gosocketio.Server
	CHANNEL_BUFFER_SIZE int
	sprovider           *provider.ScenarioProvider
)

func init() {
	CHANNEL_BUFFER_SIZE = 10
	Order = ""
	forwardAgentIdList = make([]uint64, 0)
	participantIdList = make([]uint64, 0)
	setAgentIdList = make([]uint64, 0)
	participantPspMap = make(map[uint64]*pb.Supply)
	isStop = false
	isSetClock = false
	isSetAgent = false
	isSetArea = false
	isStart = false
	isGetParticipant = false
	startSync = false
	syncForwardCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	syncParticipantCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	syncSetAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	data = new(Data)
	data.ClockInfo = &clock.ClockInfo{
		Time: uint32(*clockTime),
	}
}

type Data struct {
	ClockInfo *clock.ClockInfo
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
		//time := data.ClockInfo.Time
		cycleNum := 1
		sprovider.ForwardClockDemand(uint64(data.ClockInfo.Time), uint64(cycleNum))

		syncForwardCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
		forwardPspMap := sprovider.Wait(forwardAgentIdList, syncForwardCh)

		sendToSimulator(forwardPspMap)

		// calc next time
		nextTime := data.ClockInfo.Time + 1
		clockInfo := &clock.ClockInfo{
			Time:          uint32(nextTime),
			CycleNum:      uint32(1),
			CycleDuration: uint32(1),
			CycleInterval: uint32(1),
		}
		data.ClockInfo = clockInfo

		time.Sleep(1 * time.Second)

		if isStop == false {
			log.Printf("FORWARD_CLOCK")
			orderStartClock()
		} else {
			fmt.Printf("Stop clock!")
			isStop = false
			isStart = false
		}

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

// orderSetAgents: agentをセットするDemandを出す関数
func orderSetAgents(order *provider.SetAgentsOrder) {
	if isGetParticipant == false {
		fmt.Printf("Error... please order getParticipant")
	} else {

		sprovider.SetAgentsDemand(agentsInfo)

		syncSetAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
		sprovider.Wait(setAgentIdList, syncSetAgentsCh)
		isSetAgent = true
	}
}

// Finish Fix
func orderGetParticipant() {
	participantInfo := &participant.ParticipantInfo{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint32(sprovider.ParticipantClient.ClientID),
			AreaChannelId:        uint32(sprovider.AreaClient.ClientID),
			AgentChannelId:       uint32(sprovider.AgentClient.ClientID),
			ClockChannelId:       uint32(sprovider.ClockClient.ClientID),
			RouteChannelId:       uint32(sprovider.RouteClient.ClientID),
		},
		ProviderType: 0, //Scenario
	}
	sprovider.GetParticipantDemand(participantInfo)
}

// Finish Fix
func orderSetParticipant() {

	// send participantID to participant provider
	sprovider.SetParticipantDemand(participantsInfo)

	syncParticipantCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	sprovider.Wait(participantIdList, syncParticipantCh)
	fmt.Printf("SET_PARTICIPANT_FINISH!")
	participantPspMap = make(map[uint64]*pb.Supply)

}

func (m *MapMarker) GetJson() string {
	s := fmt.Sprintf("{\"mtype\":%d,\"id\":%d,\"lat\":%f,\"lon\":%f,\"angle\":%f,\"speed\":%d,\"area\":%d}",
		m.mtype, m.id, m.lat, m.lon, m.angle, m.speed, m.area)
	return s
}

// Fix now
func sendToSimulator(pspMap map[uint64]*pb.Supply) {
	sumAgentsInfo := make([]*agent.AgentInfo, 0)
	for _, psp := range pspMap {
		supplyType := sprovider.CheckSupplyType(psp)
		switch supplyType {
		case "FORWARD_AGENTS_SUPPLY":
			forwardAgentsSupply := psp.GetArg_ForwardAgentsSupply()
			agentsInfo := forwardAgentsSupply.AgentsInfo
			sumAgentsInfo = append(sumAgentsInfo, agentsInfo...)
		default:
			//fmt.Println("SupplyType is invalid")
		}
	}
	log.Printf("\x1b[30m\x1b[47m \n FORWARD_CLOCK_FINISH \n TIME:  %v \n Agents Num: %v \x1b[0m\n", data.ClockInfo.Time, len(sumAgentsInfo))

	if sumAgentsInfo != nil {
		jsonAgentsInfo := make([]string, 0)
		for _, agentInfo := range sumAgentsInfo {
			mm := &MapMarker{
				mtype: int32(agentInfo.AgentType), // depends on type of Ped: 0, Car , 1
				id:    int32(agentInfo.AgentId),
				lat:   float32(agentInfo.Route.Coord.Lat),
				lon:   float32(agentInfo.Route.Coord.Lon),
				angle: float32(agentInfo.Route.Direction),
				speed: int32(agentInfo.Route.Speed),
				area:  int32(agentInfo.ControlArea),
			}
			jsonAgentsInfo = append(jsonAgentsInfo, mm.GetJson())
		}
		mu.Lock()
		ioserv.BroadcastToAll("event", jsonAgentsInfo)
		mu.Unlock()
	}
}

// Finish Fix
func callbackStartClock(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	syncForwardCh <- sp

}

// Finish Fix
func callbackSetAgent(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	syncSetAgentsCh <- sp
}

// Finish Fix
func callbackSetParticipant(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	syncParticipantCh <- sp
}

// create sync id list
func createSyncIdList(participantsInfo []*participant.ParticipantInfo) ([]uint64, []uint64, []uint64) {
	setAgentIdList := make([]uint64, 0)
	forwardAgentIdList := make([]uint64, 0)
	participantIdList := make([]uint64, 0)

	for _, participantInfo := range participantsInfo {
		tProviderType := participantInfo.ProviderType
		isSetAgent := tProviderType.String() == "CAR_AREA" || tProviderType.String() == "PED_AREA"
		isForwardAgent := tProviderType.String() == "CAR_AREA" || tProviderType.String() == "PED_AREA" || tProviderType.String() == "AREA"
		isSetParticipant := tProviderType.String() == "CAR_AREA" || tProviderType.String() == "PED_AREA"
		if isSetAgent {
			channelId := participantInfo.ChannelId
			agentChannelId := channelId.AgentChannelId
			setAgentIdList = append(setAgentIdList, uint64(agentChannelId))
		}
		if isForwardAgent {

			channelId := participantInfo.ChannelId
			agentChannelId := channelId.AgentChannelId
			clockChannelId := channelId.ClockChannelId
			forwardAgentIdList = append(forwardAgentIdList, uint64(agentChannelId))
			forwardAgentIdList = append(forwardAgentIdList, uint64(clockChannelId))
		}
		if isSetParticipant {
			channelId := participantInfo.ChannelId
			participantChannelId := channelId.ParticipantChannelId
			participantIdList = append(participantIdList, uint64(participantChannelId))
		}

	}
	return setAgentIdList, participantIdList, forwardAgentIdList
}

// Finish Fix
func collectParticipantId(clt *sxutil.SMServiceClient, d int) {

	for i := 0; i < d; i++ {
		time.Sleep(1 * time.Second)
		log.Printf("waiting... %v", i+1)
	}

	mu.Lock()
	participantsInfo = simutil.CreateParticipantsInfo(participantPspMap)
	setAgentIdList, participantIdList, forwardAgentIdList = createSyncIdList(participantsInfo)

	startCollectId = false
	isGetParticipant = true
	participantPspMap = make(map[uint64]*pb.Supply)

	mu.Unlock()
	// setParticipant
	Order = "SetParticipant"
	orderSetParticipant()
}

// Finish Fix
func callbackGetParticipant(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	log.Println("Got for get_participant callback")

	mu.Lock()
	participantPspMap[sp.SenderId] = sp

	if !startCollectId {
		log.Println("start selection")
		startCollectId = true
		go collectParticipantId(clt, 2)
	}

	mu.Unlock()
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	supplyType := sprovider.CheckSupplyType(sp)
	log.Println("Got supply callback", supplyType)
	if sprovider.IsSupplyTarget(sp) {
		switch supplyType {
		case provider.GET_PARTICIPANT_SUPPLY:
			callbackGetParticipant(clt, sp)
		case provider.SET_PARTICIPANT_SUPPLY:
			sprovider.sendToWait(sp, syncParticipantCh)
			//callbackSetParticipant(clt, sp)
		case provider.SET_AGENTS_SUPPLY:
			sprovider.sendToWait(sp, syncSetAgentsCh)
			//callbackSetAgent(clt, sp)
		case provider.FORWARD_AGENTS_SUPPLY:
			sprovider.sendToWait(sp, syncForwardCh)
			//callbackStartClock(clt, sp)
		case provider.FORWARD_CLOCK_SUPPLY:
			sprovider.sendToWait(sp, syncForwardCh)
			//callbackStartClock(clt, sp)
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
	sioClient, sioErr := gosocketio.Dial("ws://localhost:9995/socket.io/?EIO=3&transport=websocket", transport.DefaultWebsocketTransport())
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
	sioClient.On(provider.SET_AGENTS, func(c *gosocketio.Channel, order *provider.SetAgentsOrder) {
		fmt.Printf("set agent %v, \n", order)
		orderSetAgents(order)
	})

	sioClient.On(provider.GET_PARTICIPANT, func(c *gosocketio.Channel, order *provider.GetParticipantOrder) {
		orderGetParticipant()
	})

	sioClient.On(provider.START_CLOCK, func(c *gosocketio.Channel, order *provider.StartClockOrder) {
		orderStartClock()
	})

	sioClient.On(provider.STOP_CLOCK, func(c *gosocketio.Channel, order *provider.StopClockOrder) {
		orderStopClock()
	})

	/*sioClient.On("scenario", func(c *gosocketio.Channel, order *simutil.OrderTest) {
		log.Printf("get order is: %v\n", (*order.Agents[0].Status))
		Order = order.Type
		switch Order {
		case "GetParticipant":
			//			fmt.Println("getParticipant")
			orderGetParticipant()
		case "SetTime":
			fmt.Println("setClock")
			clockInfo := order.ClockInfo
			orderSetClock(clockInfo)
		case "SetArea":
			fmt.Println("setArea")
			areaInfo := order.AreaInfo
			orderSetArea(areaInfo)
		case "SetAgent":
			fmt.Printf("set agent %v, \n", order.Agents)
			//agentsInfo := simutil.ConvertAgentsInfo(order.AgentsInfo)
			//orderSetAgents(agentsInfo)
		case "Start":
			//			fmt.Println("start clock")
			orderStartClock()
		case "Stop":
			//			fmt.Println("stop clock")
			orderStopClock()
		default:
			fmt.Println("error")
		}

	})*/

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
	sioClient := runClient()
	log.Printf("Running Sio Client..\n")
	if sioClient == nil {
		os.Exit(1)
	}

	// connect to synerex server
	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Scenario}")

	// Clientとして登録
	sprovider = provider.NewScenarioProvider()

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	sprovider.SetupProvider(client, argJson, func(clt *sxutil.SMServiceClient, dm *pb.Demand) {}, supplyCallback, &wg)
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
