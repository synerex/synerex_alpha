package main

import (
	"flag"
	"log"
	"sync"

	"fmt"

	"context"
	"net"
	"time"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/daemon"
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
	"github.com/synerex/synerex_alpha/provider/simulation/scenario/communicator"
	"github.com/synerex/synerex_alpha/provider/simulation/scenario/simulator"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

// プロバイダ順番関係なく
// rvo2適用
// daemonでprovider起動setup

// scenario以外のプロバイダを自由に起動可能/scenarioは参加者管理ができている
// scenarioを途中で起動可能/ 各プロバイダは参加者登録、クロック同期ができている
// 起動停止(ctl+c)を二回押すことをなくす

var (
	serverAddr       = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv          = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	port             = flag.Int("port", 10080, "HarmoVis Provider Listening Port")
	daemonPort       = ":9996"
	clockTime        = flag.Int("time", 1, "Time")
	version          = "0.01"
	startCollectId   bool
	isStop           bool
	isStart          bool
	isSetClock       bool
	isSetArea        bool
	isSetAgent       bool
	isGetParticipant bool
	mu               sync.Mutex
	startSync        bool
	com              *communicator.ScenarioCommunicator
	sim              *simulator.ScenarioSimulator
)

func init() {
	isStop = false
	isSetClock = false
	isSetAgent = false
	isSetArea = false
	isStart = false
	isGetParticipant = false
	startSync = false
}

// startClock:
func startClock(stepNum uint64) {

	com.ForwardClockRequest(stepNum)

	com.WaitForwardClockResponse()

	// calc next time
	sim.ForwardStep()
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock forwarded \n Time:  %v \x1b[0m\n", sim.GlobalTime)

	// 待機
	time.Sleep(time.Duration(sim.TimeStep) * time.Second)

	// 次のサイクルを行う
	if isStop != true {
		startClock(stepNum)
	} else {
		log.Printf("\x1b[30m\x1b[47m \n Finish: Clock stopped \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.GlobalTime, sim.TimeStep)
		isStop = false
		isStart = false
		// exit goroutin
		return
	}

}

// stopClock: Clockを停止する
func stopClock() (bool, error) {
	isStop = true
	return true, nil
}

// setAgents: agentをセットするDemandを出す関数
func setAgents(agents []*agent.Agent) (bool, error) {

	// エージェントを設置するリクエスト
	com.SetAgentsRequest(agents)

	// 同期のため待機
	com.WaitSetAgentsResponse()

	isSetAgent = true
	log.Printf("\x1b[30m\x1b[47m \n Finish: Agents set \n Add: %v \x1b[0m\n", len(agents))
	return true, nil
}

// ClearAgents: agentを消去するDemandを出す関数
func clearAgents() (bool, error) {

	// エージェントを設置するリクエスト
	com.ClearAgentsRequest()

	// 同期のため待機
	com.WaitClearAgentsResponse()

	log.Printf("\x1b[30m\x1b[47m \n Finish: Agents cleared. \x1b[0m\n")
	return true, nil
}

// setClock : クロック情報をDaemonから受け取りセットする
func setClock(globalTime float64, timeStep float64) (bool, error) {
	// クロックをセット
	sim.SetGlobalTime(globalTime)
	sim.SetTimeStep(timeStep)

	// クロック情報をプロバイダに送信
	clockInfo := sim.GetClock()
	com.SetClockRequest(clockInfo)
	// Responseを待機
	com.WaitSetClockResponse()
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.GlobalTime, sim.TimeStep)
	return true, nil
}

// collectParticipants : 起動時に、他プロバイダの参加者情報を集める
func collectParticipants() {

	// 情報をプロバイダに送信
	com.CollectParticipantsRequest()
}

// callbackRegistParticipantRequest: 新規参加者を登録するための関数
func callbackRegistParticipantRequest(dm *pb.Demand) {

	participant := dm.GetSimDemand().GetRegistParticipantRequest().GetParticipant()
	// IdListに保存
	com.AddParticipant(participant)

	// 同期するためのIdListを作成
	com.CreateWaitIdList()

	// 新規の参加者情報を参加プロバイダに送信
	com.SetParticipantsRequest()

	// SetParticipantsResponseの待機
	err := com.WaitSetParticipantsResponse()
	if err != nil{
		log.Printf("\x1b[30m\x1b[47m \n Error: %v \x1b[0m\n", err)
	}

	// 新規参加者に登録完了通知
	com.RegistParticipantResponse(dm)
	log.Printf("\x1b[30m\x1b[47m \n Finish: New participant registed  \x1b[0m\n")

}

// callbackDeleteParticipantRequest: 参加者を削除するための関数
func callbackDeleteParticipantRequest(dm *pb.Demand) {

	participant := dm.GetSimDemand().GetDeleteParticipantRequest().GetParticipant()
	// IdListに保存
	com.DeleteParticipant(participant)
	
	// 同期するためのIdListを作成
	com.CreateWaitIdList()

	// 新規の参加者情報を参加プロバイダに送信
	com.SetParticipantsRequest()
	
	// SetParticipantsResponseの待機
	com.WaitSetParticipantsResponse()
	
	// 新規参加者に登録完了通知
	com.DeleteParticipantResponse(dm)
	log.Printf("\x1b[30m\x1b[47m \n Finish: Participant deleted \x1b[0m\n")
}

// callbackGetClockRequest: 新規参加者がクロック情報を取得する関数
func callbackGetClockRequest(dm *pb.Demand) {
	// Clock情報を取得
	clock := sim.GetClock()

	// Clock情報を新規参加者へ送る
	com.GetClockResponse(dm, clock)
}

// notifyDownScenario: 新規参加者がクロック情報を取得する関数
func notifyDownScenario() {

	// 参加取り消しをするRequest
	com.DownScenarioRequest()

	// Responseの待機
	com.WaitDownScenarioResponse()

	log.Printf("\x1b[30m\x1b[47m \n Finish: Scenario-provider notified clash-infomation. \x1b[0m\n")
}

// Supplyのコールバック関数
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	switch sp.GetSimSupply().SupplyType {
	case synerex.SupplyType_SET_PARTICIPANTS_RESPONSE:
		com.SendToSetParticipantsResponse(sp)
	case synerex.SupplyType_SET_AGENTS_RESPONSE:
		com.SendToSetAgentsResponse(sp)
	case synerex.SupplyType_CLEAR_AGENTS_RESPONSE:
		com.SendToClearAgentsResponse(sp)
	case synerex.SupplyType_SET_CLOCK_RESPONSE:
		com.SendToSetClockResponse(sp)
	case synerex.SupplyType_FORWARD_CLOCK_RESPONSE:
		com.SendToForwardClockResponse(sp)
	case synerex.SupplyType_DOWN_SCENARIO_RESPONSE:
		com.SendToDownScenarioResponse(sp)
	default:
		//fmt.Println("order is invalid")
	}
}

// Demandのコールバック関数
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	switch dm.GetSimDemand().DemandType {
	case synerex.DemandType_REGIST_PARTICIPANT_REQUEST:
		callbackRegistParticipantRequest(dm)
	case synerex.DemandType_DELETE_PARTICIPANT_REQUEST:
		callbackDeleteParticipantRequest(dm)
	case synerex.DemandType_GET_CLOCK_REQUEST:
		callbackGetClockRequest(dm)
	default:
		//fmt.Println("order is invalid")
	}

}

// simServer: scenarioと通信を行うサーバ
type simDaemonServer struct {
}

// SetAgentsOrder
func (s *simDaemonServer) SetAgentsOrder(ctx context.Context, in *daemon.SetAgentsMessage) (*daemon.Response, error) {
	log.Printf("Received:SetAgents")
	agents := in.GetAgents()
	ok, err := setAgents(agents)
	return &daemon.Response{Ok: ok}, err
}

// ClearAgentsOrder
func (s *simDaemonServer) ClearAgentsOrder(ctx context.Context, in *daemon.ClearAgentsMessage) (*daemon.Response, error) {
	log.Printf("Received:ClearAgents")
	ok, err := clearAgents()
	return &daemon.Response{Ok: ok}, err
}

// SetClock
func (s *simDaemonServer) SetClockOrder(ctx context.Context, in *daemon.SetClockMessage) (*daemon.Response, error) {
	log.Printf("Received:SetClock")
	globalTime := in.GetGlobalTime()
	timeStep := in.GetTimeStep()
	ok, err := setClock(globalTime, timeStep)
	return &daemon.Response{Ok: ok}, err
}

// StartClock
func (s *simDaemonServer) StartClockOrder(ctx context.Context, in *daemon.StartClockMessage) (*daemon.Response, error) {
	log.Printf("Received:StartClock")
	stepNum := in.GetStepNum()
	if isStart {
		log.Printf("\x1b[30m\x1b[47m \n Simulator is already started. \x1b[0m\n")
	} else {
		isStart = true
		go startClock(stepNum)
	}
	return &daemon.Response{Ok: true}, nil
}

// StopClock
func (s *simDaemonServer) StopClockOrder(ctx context.Context, in *daemon.StopClockMessage) (*daemon.Response, error) {
	log.Printf("Received:StopClock")
	ok, err := stopClock()
	return &daemon.Response{Ok: ok}, err
}



func runDaemonServer() {
	// Scenarioとの通信サーバ
	lis, err := net.Listen("tcp", daemonPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	daemon.RegisterSimDaemonServer(s, &simDaemonServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
	sxutil.RegisterDeferFunction(func() { notifyDownScenario(); conn.Close() })

	// run daemon server
	go runDaemonServer()
	log.Printf("Running Daemon Server..\n")

	// connect to synerex server
	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Scenario}")

	// simulator
	timeStep := float64(1)
	globalTime := float64(0)
	sim = simulator.NewScenarioSimulator(timeStep, globalTime)

	// Communicator
	com = communicator.NewScenarioCommunicator()

	// Communicatorのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	// channelごとのClientを作成
	com.RegistClients(client, argJson)
	// ChannelにSubscribe
	com.SubscribeAll(demandCallback, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	// 起動時、プロバイダがいれば登録する
	collectParticipants()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
