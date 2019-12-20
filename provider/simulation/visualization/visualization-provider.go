package main

import (
	"flag"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
	"github.com/synerex/synerex_alpha/provider/simulation/visualization/communicator"
	"github.com/synerex/synerex_alpha/sxutil"

	//	"github.com/synerex/synerex_alpha/api/simulation/area"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	gosocketio "github.com/mtfelian/golang-socketio"
	"github.com/synerex/synerex_alpha/provider/simulation/visualization/simulator"
	"google.golang.org/grpc"
	//	"os/exec"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	port       = flag.Int("port", 10080, "HarmoVis Provider Listening Port")
	version    = "0.01"
	mu         sync.Mutex
	assetsDir  http.FileSystem
	ioserv     *gosocketio.Server
	com        *communicator.VisualizationCommunicator
	sim        *simulator.VisualizationSimulator
	isDownScenario bool
)

func init(){
	isDownScenario = false
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

// GetJson: json化する関数
func (m *MapMarker) GetJson() string {
	s := fmt.Sprintf("{\"mtype\":%d,\"id\":%d,\"lat\":%f,\"lon\":%f,\"angle\":%f,\"speed\":%d,\"area\":%d}",
		m.mtype, m.id, m.lat, m.lon, m.angle, m.speed, m.area)
	return s
}

// sendToHarmowareVis: harmowareVisに情報を送信する関数
func sendToHarmowareVis(sumAgents []*agent.Agent) {

	if sumAgents != nil {
		jsonAgents := make([]string, 0)
		for _, agentInfo := range sumAgents {

			// agentInfoTypeによってエージェントを取得
			switch agentInfo.Type {
			case common.AgentType_PEDESTRIAN:
				ped := agentInfo.GetPedestrian()
				mm := &MapMarker{
					mtype: int32(agentInfo.Type), // depends on type of Ped: 0, Car , 1
					id:    int32(agentInfo.Id),
					lat:   float32(ped.Route.Position.Latitude),
					lon:   float32(ped.Route.Position.Longitude),
					angle: float32(ped.Route.Direction),
					speed: int32(ped.Route.Speed),
				}
				jsonAgents = append(jsonAgents, mm.GetJson())

			case common.AgentType_CAR:
				car := agentInfo.GetCar()
				mm := &MapMarker{
					mtype: int32(agentInfo.Type), // depends on type of Ped: 0, Car , 1
					id:    int32(agentInfo.Id),
					lat:   float32(car.Route.Position.Latitude),
					lon:   float32(car.Route.Position.Longitude),
					angle: float32(car.Route.Direction),
					speed: int32(car.Route.Speed),
				}
				jsonAgents = append(jsonAgents, mm.GetJson())

			case common.AgentType_TRAIN:
				train := agentInfo.GetTrain()
				mm := &MapMarker{
					mtype: int32(agentInfo.Type), // depends on type of Ped: 0, Car , 1
					id:    int32(agentInfo.Id),
					lat:   float32(train.Route.Position.Latitude),
					lon:   float32(train.Route.Position.Longitude),
					angle: float32(train.Route.Direction),
					speed: int32(train.Route.Speed),
				}
				jsonAgents = append(jsonAgents, mm.GetJson())
			}
		}
		mu.Lock()
		ioserv.BroadcastToAll("event", jsonAgents)
		mu.Unlock()
	}
}

// registParticipant: 新規参加登録をする関数
func registParticipant() {
	// 新規参加登録をするRequest
	participant := com.GetMyParticipant()
	com.RegistParticipantRequest(participant)
	// Responseの待機
	err := com.WaitRegistParticipantResponse()

	if err != nil {
		log.Printf("\x1b[30m\x1b[47m \n Error: %v \x1b[0m\n", err)
	}else{
		// クロック情報を取得する
		getClock()
		log.Printf("\x1b[30m\x1b[47m \n Finish: This provider registered in scenario-provider \x1b[0m\n")
	}
	return
}

// deleteParticipant: プロバイダ停止時に参加取り消しをする
func deleteParticipant() {
	if isDownScenario == false{
		// 参加取り消しをするRequest
		participant := com.GetMyParticipant()
		com.DeleteParticipantRequest(participant)
	
		// Responseの待機
		com.WaitDeleteParticipantResponse()
		log.Printf("\x1b[30m\x1b[47m \n Finish: This provider deleted from participants list in scenario-provider. \x1b[0m\n")
	}
}

// callbackSetParticipants: 参加者リストをセットする要求
func callbackSetParticipantsRequest(dm *pb.Demand) {
	participants := dm.GetSimDemand().GetSetParticipantsRequest().GetParticipants()
	targetId := dm.GetId()
	// 参加者情報をセットする
	com.SetParticipants(participants)

	// 同期するためのIdListを作成
	com.CreateWaitIdList()

	// セット完了通知を送る
	com.SetParticipantsResponse(targetId)
}

// getClock: クロック情報を取得する関数
func getClock() {

	// エリアを取得するRequest
	com.GetClockRequest()
	// Responseの待機
	clockInfo := com.WaitGetClockResponse()
	// エリア情報をセット
	sim.SetGlobalTime(clockInfo.GlobalTime)
	sim.SetTimeStep(clockInfo.TimeStep)

	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.GlobalTime, sim.TimeStep)
}

// callbackForwardClockRequest: クロックを進める関数
func callbackForwardClockRequest(dm *pb.Demand) {
	dm.GetSimDemand().GetForwardClockRequest().GetStepNum()
	targetId := dm.GetId()

	// clockを進める
	sim.ForwardGlobalTime()

	// セット完了通知を送る
	com.ForwardClockResponse(targetId)
	log.Printf("\x1b[30m\x1b[47m \n Finish Forward Clock \n Time:  %v \x1b[0m\n", sim.GlobalTime)
}

// callbackSetClock: Clock情報をセットする要求
func callbackSetClockRequest(dm *pb.Demand) {
	clockInfo := dm.GetSimDemand().GetSetClockRequest().GetClock()
	targetId := dm.GetId()

	// Clock情報をセットする
	sim.SetGlobalTime(clockInfo.GlobalTime)
	sim.SetTimeStep(clockInfo.TimeStep)

	// セット完了通知を送る
	com.SetClockResponse(targetId)
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.GlobalTime, sim.TimeStep)
}

// callbackCollectParticipantsRequest:
func callbackCollectParticipantsRequest(dm *pb.Demand) {
	// 新規参加登録 
	// TODO: Why go-routin ? 
	go registParticipant()
	
	// scenarioが再開された
	isDownScenario = false
}

// callbackDownScenarioRequest:
func callbackDownScenarioRequest(dm *pb.Demand) {
	targetId := dm.GetId()
	// scenarioがダウンした
	isDownScenario = true
	// 返答を返す
	com.DownScenarioResponse(targetId)
	log.Printf("\x1b[30m\x1b[47m \n Error: scenario-provider crashed...\n Please restart scenario-provider.   \x1b[0m\n")
}

// callbackVisualizeAgentsResponse: Agentを可視化する関数
func callbackVisualizeAgentsResponse(sp *pb.Supply) {

	// agentsを取得
	agents := sp.GetSimSupply().GetVisualizeAgentsResponse().GetAgents()

	// Harmowareに送る
	sendToHarmowareVis(agents)

}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	switch dm.GetSimDemand().DemandType {

	case synerex.DemandType_SET_PARTICIPANTS_REQUEST:
		// 参加者リストをセットする要求
		callbackSetParticipantsRequest(dm)
	case synerex.DemandType_FORWARD_CLOCK_REQUEST:
		// クロックを進める要求
		callbackForwardClockRequest(dm)
	case synerex.DemandType_SET_CLOCK_REQUEST:
		// クロックをセットする要求
		callbackSetClockRequest(dm)
	case synerex.DemandType_DOWN_SCENARIO_REQUEST:
		// Scenarioがダウンした場合の要求
		callbackDownScenarioRequest(dm)
	case synerex.DemandType_COLLECT_PARTICIPANTS_REQUEST:
		// シナリオプロバイダが参加者リストを募集する要求
		callbackCollectParticipantsRequest(dm)
	default:
		//log.Println("demand callback is invalid.")
	}
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	switch sp.GetSimSupply().SupplyType {
	case synerex.SupplyType_GET_CLOCK_RESPONSE:
		// Clock情報の取得
		com.SendToGetClockResponse(sp)
	case synerex.SupplyType_REGIST_PARTICIPANT_RESPONSE:
		// 参加者登録完了通知の取得
		com.SendToRegistParticipantResponse(sp)
	case synerex.SupplyType_DELETE_PARTICIPANT_RESPONSE:
		// 参加者削除完了通知の取得
		com.SendToDeleteParticipantResponse(sp)
	case synerex.SupplyType_VISUALIZE_AGENTS_RESPONSE:
		// エージェントを可視化する要求
		callbackVisualizeAgentsResponse(sp)
	default:
		//fmt.Println("order is invalid")
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

	sxutil.RegisterNodeName(*nodesrv, "VisualizationProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	// Clientとして登録
	com = communicator.NewVisualizationCommunicator()

	sxutil.RegisterDeferFunction(func() { deleteParticipant(); conn.Close() })

	// run socket.io server
	ioserv = runServer()
	log.Printf("Running Sio Server..\n")
	if ioserv == nil {
		os.Exit(1)
	}

	// synerex simulator
	sim = simulator.NewVisualizationSimulator(1.0, 0.0)

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Visualization}")


	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	// channelごとのClientを作成
	com.RegistClients(client, argJson)
	// ChannelにSubscribe
	com.SubscribeAll(demandCallback, supplyCallback, &wg)
	wg.Wait()

	// 新規参加登録
	registParticipant()

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

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!
}
