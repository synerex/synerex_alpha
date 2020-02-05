package main

import (
	"flag"
	"log"
	"sync"

	"fmt"

	"time"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
	"github.com/synerex/synerex_alpha/provider/simulation/clock/communicator"
	"github.com/synerex/synerex_alpha/provider/simulation/clock/simulator"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

// rvo2適用
// daemonでprovider起動setup
// daemonをサーバとしてscenarioに命令

var (
	serverAddr       = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv          = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	pidFlag     = flag.Int("pid", 1, "Provider Id") 
	port             = flag.Int("port", 10080, "HarmoVis Provider Listening Port")
	isStart          bool
	pid uint64
	mu               sync.Mutex
	com              *communicator.ClockCommunicator
	sim              *simulator.ClockSimulator
)

func init() {
	flag.Parse()
	isStart = false
	pid = uint64(*pidFlag)
}

// registParticipant: 新規参加登録をする関数
func registParticipant() {

	// 新規参加登録をするRequest
	participant := com.GetMyParticipant()
	log.Printf("\x1b[31m\x1b[47m \n ParticipantID: %v \x1b[0m\n", participant.Id)

	com.RegistParticipantRequest(participant)
	// Responseの待機
	err := com.WaitRegistParticipantResponse()
	if err != nil {
		log.Printf("\x1b[31m\x1b[47m \n Error: %v \x1b[0m\n", err)
	}else{
		log.Printf("\x1b[30m\x1b[47m \n Finish: This provider registered in scenario-provider \x1b[0m\n")
	}
	return
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
	if isStart {
		startClock(stepNum)
	} else {
		log.Printf("\x1b[30m\x1b[47m \n Finish: Clock stopped \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.GlobalTime, sim.TimeStep)
		isStart = false
		// exit goroutin
		return
	}

}

// stopClock: Clockを停止する
func stopClock() (bool, error) {
	isStart = false
	return true, nil
}


// setClock : クロック情報をDaemonから受け取りセットする
func setClock(dm *pb.Demand) (bool, error) {
	clock := dm.GetSimDemand().GetSetClockRequest().GetClock()
	// クロックをセット
	sim.SetGlobalTime(clock.GlobalTime)
	sim.SetTimeStep(clock.TimeStep)

	// クロック情報をプロバイダに送信
	/*clockInfo := sim.GetClock()
	com.SetClockRequest(clockInfo)
	// Responseを待機
	com.WaitSetClockResponse()*/
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.GlobalTime, sim.TimeStep)
	return true, nil
}

// callbackGetClockRequest: 新規参加者がクロック情報を取得する関数
func callbackGetClockRequest(dm *pb.Demand) {
	// Clock情報を取得
	clock := sim.GetClock()

	// Clock情報を新規参加者へ送る
	com.GetClockResponse(dm, clock)
}

// Supplyのコールバック関数
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	switch sp.GetSimSupply().SupplyType {
	case synerex.SupplyType_SET_CLOCK_RESPONSE:
		com.SendToSetClockResponse(sp)
	case synerex.SupplyType_FORWARD_CLOCK_RESPONSE:
		com.SendToForwardClockResponse(sp)
	default:
		//fmt.Println("order is invalid")
	}
}

// Demandのコールバック関数
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	switch dm.GetSimDemand().DemandType {
	case synerex.DemandType_GET_CLOCK_REQUEST:
		callbackGetClockRequest(dm)
	case synerex.DemandType_SET_CLOCK_REQUEST:
		setClock(dm)
	default:
		//fmt.Println("order is invalid")
	}

}



func main() {

	log.Printf("\x1b[31m\x1b[47m \n SyneServ: %v, NodeServ: %v, Pid: %v   \x1b[0m\n", *serverAddr, *nodesrv, pid)

	// connect to node server
	sxutil.RegisterNodeName(*nodesrv, "ClockProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}


	// Communicator
	com = communicator.NewClockCommunicator(pid)

	sxutil.RegisterDeferFunction(func() { conn.Close() })


	// simulator
	timeStep := float64(1)
	globalTime := float64(0)
	sim = simulator.NewClockSimulator(timeStep, globalTime)


	// connect to synerex server
	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Clock}")

	// Communicatorのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	// channelごとのClientを作成
	com.RegistClients(client, argJson)
	// ChannelにSubscribe
	com.SubscribeAll(demandCallback, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	// 新規参加登録
	registParticipant()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
