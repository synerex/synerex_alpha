package main

import (
	"flag"
	"log"
	"sync"
	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"google.golang.org/grpc"
	"time"
	"fmt"
	"github.com/mtfelian/golang-socketio/transport"
	"github.com/mtfelian/golang-socketio"
	"os"
	"net/http"
	"path/filepath"
//	"os/exec"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	port       = flag.Int("port", 10080, "HarmoVis Provider Listening Port")
	version    = "0.01"
	idlist     []uint64
	dmMap      map[uint64]*sxutil.DemandOpts
	spMap		map[uint64]*pb.Supply
	pspMap		map[uint64]*pb.Supply
	ch 			chan *pb.Supply
	startSync 	bool
	startCollectId bool
	isStop bool
	isStart bool
	isSetClock bool
	isSetArea bool
	isSetAgent bool
	isGetParticipant bool
	participantIdList  []uint64
	Order string
	mu	sync.Mutex
	sclientArea *sxutil.SMServiceClient
	sclientAgent *sxutil.SMServiceClient
	sclientClock *sxutil.SMServiceClient
	sclientParticipant *sxutil.SMServiceClient
	sioClient *gosocketio.Client
	idListByChannel *simutil.IdListByChannel
	data *Data
	assetsDir  http.FileSystem
	ioserv     *gosocketio.Server
)

func init() {
	Order = ""
	idlist = make([]uint64, 0)
	participantIdList = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*pb.Supply)
	pspMap = make(map[uint64]*pb.Supply)
	startSync = false
	isStop = false
	isSetClock = false
	isSetAgent = false
	isSetArea = false
	isStart = false
	isGetParticipant = false
	ch = make(chan *pb.Supply)
	data = new(Data)
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


func syncProposeSupply(sp *pb.Supply, idList []uint64, callback func(pspMap map[uint64]*pb.Supply)){
	go func() { 
		log.Println("Send Supply")
		ch <- sp
		return
	}()
	if !startSync {
		log.Println("Start Sync")
		startSync = true

		go func(){		
		for {
			select {
				case psp := <- ch:
					log.Println("recieve ProposeSupply")
					pspMap[psp.SenderId] = psp
//					log.Printf("waitidList %v %v", pspMap, idList)
			
				if simutil.IsFinishSync(pspMap, idList){
					fmt.Printf("Finish Sync\n")
					// if you need, return response
					callback(pspMap)
					
					// init pspMap
					pspMap = make(map[uint64]*pb.Supply)
					startSync = false
					
					return
				}
			}
		}

		}()
	}	
}

func setClock(clockInfo simutil.ClockInfo){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else{
		firstTime := uint32(1)

		clockInfo := clock.ClockInfo{
			Time: firstTime,
			StatusType: 0, // OK
			Meta: "",
		}
		data.ClockInfo = &clockInfo
	
		clockDemand := clock.ClockDemand{
			Time: firstTime,
			DemandType: 2, // SET
			CycleNum: uint32(1),
			CycleDuration: uint32(1),
			CycleInterval: uint32(1),
			StatusType: 2, // NONE
			Meta: "",
		}
		
		nm := "setClock order by scenario"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, ClockDemand: &clockDemand}
	
		dmMap, idlist = simutil.SendDemand(sclientClock, opts, dmMap, idlist)
	}
}

func startClock(){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else if isSetAgent == false{
		fmt.Printf("Error... please order setAgent")
	}else if isSetArea == false{
		fmt.Printf("Error... please order setArea")
	}else if isSetClock == false{
		fmt.Printf("Error... please order setClock")
	}else{
		isStart = true
		// forward clock
		clockDemand := clock.ClockDemand{
			Time: data.ClockInfo.Time,
			DemandType: 0, // Forward
			CycleNum: uint32(1),
			CycleDuration: uint32(1),
			CycleInterval: uint32(1),
			StatusType: 2, // NONE
			Meta: "",
		}
		
		nm := "startClock order by scenario"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, ClockDemand: &clockDemand}
		
		dmMap, idlist = simutil.SendDemand(sclientClock, opts, dmMap, idlist)
		
		// calc next time
		nextTime := data.ClockInfo.Time + 1
		clockInfo := clock.ClockInfo{
			Time: uint32(nextTime),
			StatusType: 0, // OK
			Meta: "",
		}
		data.ClockInfo = &clockInfo
	}
}

func stopClock(){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else{
		if isStart == true{
			isStop = true
		}else{
			fmt.Printf("Clock has not been started. ")
		}
	}
}

func setArea(areaInfo simutil.AreaInfo){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else{
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
}

func setAgent(agentInfo simutil.AgentInfo){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else{
		route := agent.Route{
			Coord: &agent.Route_Coord{
				Lat: float32(agentInfo.Route.Coord.Lat),
				Lon: float32(agentInfo.Route.Coord.Lon), 
			},
			Direction: float32(agentInfo.Route.Direction),
			Speed: float32(agentInfo.Route.Speed),
			Destination: float32(10),
			Departure: float32(100),
		}
		rule := agent.Rule{
			RuleInfo: "nothing",
		}
		agentStatus := agent.AgentStatus{
			Age: "20",
			Sex: "Male",
		}
	
		agentDemand := agent.AgentDemand{
			Time: uint32(1),
			AgentId: agentInfo.Id,
			AgentName: "Agent1",
			AgentType: 0,
			AgentStatus: &agentStatus,
			Route: &route,
			Rule: &rule,
			DemandType: 0, // SET
			StatusType: 2, // NONE
			Meta: "",
		}
		
		nm := "setAgent order by scenario"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, AgentDemand: &agentDemand}
	
		dmMap, idlist = simutil.SendDemand(sclientAgent, opts, dmMap, idlist)
	}
}

func getParticipant(){
	participantDemand := participant.ParticipantDemand {
		ClientId: uint64(sclientParticipant.ClientID),
		DemandType: 0, // GET
		StatusType: 2, // NONE
		Meta: "",
	}
	
	nm := "getParticipant order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, ParticipantDemand: &participantDemand}

	dmMap, idlist = simutil.SendDemand(sclientParticipant, opts, dmMap, idlist)
}

func (m *MapMarker) GetJson() string {
	s := fmt.Sprintf("{\"mtype\":%d,\"id\":%d,\"lat\":%f,\"lon\":%f,\"angle\":%f,\"speed\":%d}",
		m.mtype, m.id, m.lat, m.lon, m.angle, m.speed)
	return s
}

func sendToSimulator(psp *pb.Supply){
	agentsInfo := psp.GetArg_AgentsInfo()
	if agentsInfo != nil{
		fmt.Printf("\x1b[30m\x1b[47m agentsInfo is : %v\x1b[0m\n", agentsInfo)
		for _, agentInfo := range agentsInfo.AgentInfo{
			mm := &MapMarker{
				mtype: int32(agentInfo.AgentType), // depends on type of GTFS: 1 for Subway, 2, for Rail, 3 for bus
				id:    int32(agentInfo.AgentId),
				lat:   float32(agentInfo.Route.Coord.Lat),
				lon:   float32(agentInfo.Route.Coord.Lon),
				angle: float32(agentInfo.Route.Direction),
				speed: int32(agentInfo.Route.Speed),
			}
			mu.Lock()
			ioserv.BroadcastToAll("event", mm.GetJson())
			mu.Unlock()
		}
	}
}

func callbackForStartClock(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for start_clock callback")
	if clt.IsSupplyTarget(sp, idlist) { 

		//opts :=	dmMap[sp.TargetId]
//		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp
		
		callback := func (pspMap map[uint64]*pb.Supply){
			fmt.Printf("Callback StartClock! Regist and Send Data to Simulator")
			for _, psp := range pspMap{
				supplyType := simutil.CheckSupplyArgOneOf(psp)
				switch supplyType{
				case "RES_SET_AGENTS":
					fmt.Println("registAgents")
					sendToSimulator(psp)
				default:
					fmt.Println("SupplyType is invalid")
				}
			}
			

			fmt.Printf("Callback StartClock! move next clock")
			time.Sleep(1* time.Second)
			if isStop == false{
				startClock()
			}else{
				fmt.Printf("Stop clock!")
				isStop = false
				isStart = false
			}
		}

		clockIdList := idListByChannel.ClockIdList
		agentIdList := idListByChannel.AgentIdList
		idList := append(clockIdList, agentIdList...)
		syncProposeSupply(sp, idList, callback)

	}else{
		log.Printf("This is not propose supply \n")
	}
}

func callbackForSetAgent(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for set_agent callback")
	if clt.IsSupplyTarget(sp, idlist) { 

//		opts :=	dmMap[sp.TargetId]
//		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp
		
		callback := func (pspMap map[uint64]*pb.Supply){
			fmt.Printf("Callback SetAgent! return OK to Simulation Synerex Engine")
			isSetAgent = true
		}
		agentIdList := idListByChannel.AgentIdList
		syncProposeSupply(sp, agentIdList, callback)	

	}else{
		log.Printf("This is not propose supply \n")
	}
}

func callbackForSetClock(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for set_clock callback")
	if clt.IsSupplyTarget(sp, idlist) { 

//		opts :=	dmMap[sp.TargetId]
//		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp
		
		callback := func (pspMap map[uint64]*pb.Supply){
			fmt.Printf("Callback SetClock! return OK to Simulation Synerex Engine")
			isSetClock = true
		}
		clockIdList := idListByChannel.ClockIdList
		syncProposeSupply(sp, clockIdList, callback)

	}else{
		log.Printf("This is not propose supply \n")
	}
}

func callbackForSetArea(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for set_area callback")
	if clt.IsSupplyTarget(sp, idlist) { 

//		opts :=	dmMap[sp.TargetId]
//		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp
		

		callback := func (pspMap map[uint64]*pb.Supply){
			fmt.Printf("Callback SetArea! return OK to Simulation Synerex Engine")
			isSetArea = true
		}
	
		areaIdList := idListByChannel.AreaIdList
		syncProposeSupply(sp, areaIdList, callback)

	}else{
		log.Printf("This is not propose supply \n")
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
	isGetParticipant = true
	pspMap = make(map[uint64]*pb.Supply)
}

func callbackForGetParticipant(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for get_participant callback")
	
	if idListByChannel != nil{
		log.Println("already corrected IdList!")
	}else{
		mu.Lock()
		if clt.IsSupplyTarget(sp, idlist) { 

//			opts :=	dmMap[sp.TargetId]
//			log.Printf("Got Supply for %v as '%v'",opts, sp )
			spMap[sp.SenderId] = sp
			pspMap[sp.SenderId] = sp

			if !startCollectId {
				log.Println("start selection")
				startCollectId = true
				go collectParticipantId(clt, 3)
			}
		}else{
			log.Printf("This is not propose supply \n")
		}
		mu.Unlock()
	}
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got supply callback")
	//log.Printf("supply is %v",sp)
	switch Order {
		case "GetParticipant":
			fmt.Println("getParticipant")
			callbackForGetParticipant(clt, sp)
		case "SetTime":
			fmt.Println("setClock")
			callbackForSetClock(clt, sp)
		case "SetArea":
			fmt.Println("setArea")
			callbackForSetArea(clt, sp)
		case "SetAgent":
			fmt.Println("setAgent")
			callbackForSetAgent(clt, sp)
		case "Start":
			fmt.Println("startClock")
			callbackForStartClock(clt, sp)
		default:
			fmt.Println("order is invalid")
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

func runClient() *gosocketio.Client{
	var sioErr error
	sioClient, sioErr = gosocketio.Dial("ws://localhost:9995/socket.io/?EIO=3&transport=websocket", transport.DefaultWebsocketTransport())
	if sioErr != nil {
		fmt.Println("se: Error to connect with se-daemon. You have to start se-daemon first.") //,err)
		os.Exit(1)
	}else{
		fmt.Println("se: connect OK")
	}
	sioClient.On(gosocketio.OnConnection, func(c *gosocketio.Channel,param interface{}) {
		fmt.Println("Go socket.io connected ")
		c.Emit("setCh", "Scenario")
	})

	sioClient.On("scenario", func(c *gosocketio.Channel, order *simutil.Order) {
		fmt.Printf("get order is: %v\n", order)
		Order = order.Type
		switch Order {
		case "GetParticipant":
			fmt.Println("getParticipant")
			getParticipant()
		case "SetTime":
			fmt.Println("setClock")
			clockInfo := order.ClockInfo
			setClock(clockInfo)
		case "SetArea":
			fmt.Println("setArea")
			areaInfo := order.AreaInfo
			setArea(areaInfo)
		case "SetAgent":
			fmt.Println("set agent")
			agentInfo := order.AgentInfo
			setAgent(agentInfo)
		case "Start":
			fmt.Println("start clock")
			startClock()
		case "Stop":
			fmt.Println("stop clock")
			stopClock()
		default:
			fmt.Println("error")
		}
		
	})

	sioClient.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel,param interface{}) {
		fmt.Println("Go socket.io disconnected ",c)
	})

	return sioClient
}



func main() {
//	fmt.Printf("running test.go\n")
//	time.Sleep(2)
//	cmd := exec.Command("./scenario-provider")
//	cmd.Start()

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
	fmt.Printf("Running Sio Server..\n")
	if ioserv == nil {
		os.Exit(1)
	}

	// run socket.io client
	sioClient = runClient()
	fmt.Printf("Running Sio Client..\n")
	if sioClient == nil {
		os.Exit(1)
	}

	// connect to synerex server
	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Scenario}")
	sclientAgent = sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE,argJson)
	sclientClock = sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE,argJson)
	sclientArea = sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE,argJson)
	sclientParticipant = sxutil.NewSMServiceClient(client, pb.ChannelType_PARTICIPANT_SERVICE,argJson)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go simutil.SubscribeSupply(sclientAgent, supplyCallback)
	go simutil.SubscribeSupply(sclientClock, supplyCallback)
	go simutil.SubscribeSupply(sclientArea, supplyCallback)
	go simutil.SubscribeSupply(sclientParticipant, supplyCallback)

	serveMux := http.NewServeMux()

	serveMux.Handle("/socket.io/", ioserv)
	serveMux.HandleFunc("/", assetsFileHandler)

	log.Printf("Starting Harmoware VIS  Provider %s  on port %d", version, *port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), serveMux)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
