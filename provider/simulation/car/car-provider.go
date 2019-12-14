package main

import (
	//"context"
	"flag"
	"fmt"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
	"github.com/synerex/synerex_alpha/provider/simulation/car/communicator"
	"github.com/synerex/synerex_alpha/provider/simulation/car/simulator"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	areaId     = uint64(*flag.Int("areaId", 1, "Area Id")) // Area B
	agentType  = common.AgentType_CAR                      // CAR
	com        *communicator.CarCommunicator
	sim        *simulator.CarSimulator
)

// getArea: 起動時にエリアを取得する関数
func getArea() {
	// エリアを取得するRequest
	com.GetAreaRequest(areaId)
	// Responseの待機
	areaInfo := com.WaitGetAreaResponse()
	// エリア情報をセット
	sim.SetArea(areaInfo)
	fmt.Printf("Finish Get Area\n")
}

// registParticipant: 新規参加登録をする関数
func registParticipant() {
	// 新規参加登録をするRequest
	participant := com.GetMyParticipant(areaId)
	com.RegistParticipantRequest(participant)
	// Responseの待機
	com.WaitRegistParticipantResponse()
	fmt.Printf("Finish Regist Participant\n")
}

// deleteParticipant: プロバイダ停止時に参加取り消しをする
func deleteParticipant() {
	// 参加取り消しをするRequest
	participant := com.GetMyParticipant(areaId)
	com.DeleteParticipantRequest(participant)

	// Responseの待機
	com.WaitDeleteParticipantResponse()
	fmt.Printf("Finish Delete Participant\n")
}

// callbackSetParticipants: 参加者リストをセットする要求
func callbackSetParticipantsRequest(dm *pb.Demand) {
	participants := dm.GetSimDemand().GetSetParticipantsRequest().GetParticipants()
	targetId := dm.GetId()
	// 参加者情報をセットする
	com.SetParticipants(participants)

	// 同期するためのIdListを作成
	com.CreateWaitIdList(agentType, areaId)

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
	fmt.Printf("Finish Get Clock\n")
}

// callbackSetAgents: Agent情報をセットする要求
func callbackSetAgentsRequest(dm *pb.Demand) {
	agents := dm.GetSimDemand().GetSetAgentsRequest().GetAgents()
	targetId := dm.GetId()

	// Agent情報をセットする
	sim.SetAgents(agents)

	// セット完了通知を送る
	com.SetAgentsResponse(targetId)
	fmt.Printf("Finish Set Agents\n")
}

// callbackGetSameAreaAgentsRequest: Agent情報をセットする要求
func callbackGetSameAreaAgentsRequest(dm *pb.Demand) {
	areaId := dm.GetSimDemand().GetGetSameAreaAgentsRequest().GetAreaId()
	// agentType := dm.GetSimDemand().GetSameAreaAgentsRequest().GetAgentType()
	targetId := dm.GetId()

	// Areaが等しい場合
	if areaId == sim.Area.Id {
		// Agentを送る
		com.GetSameAreaAgentsResponse(targetId, sim.Agents)
	}
}

// callbackForwardClock: Agentを計算し、クロックを進める要求
func callbackForwardClockRequest(dm *pb.Demand) {
	dm.GetSimDemand().GetForwardClockRequest().GetStepNum()
	targetId := dm.GetId()

	// 同じエリアのAgent情報を取得する
	com.GetSameAreaAgentsRequest(areaId, agentType)
	// Responseの待機
	sameAreaAgents := com.WaitGetSameAreaAgentsResponse()

	// 次の時間のエージェントを計算する
	pureNextAgents := sim.ForwardStep(sameAreaAgents)

	// 次の時刻の隣接しているエリアの同じAgentTypeのエージェント情報を取得する
	neighborAreaAgents := com.WaitGetNeighborAreaAgentsResponse()

	// 重複エリアのエージェントを更新する
	nextAgents := sim.UpdateDuplicateAgents(pureNextAgents, neighborAreaAgents)

	// Agentsをセットする
	sim.SetAgents(nextAgents)

	// クロックを進める
	sim.ForwardGlobalTime()

	// セット完了通知を送る
	com.ForwardClockResponse(targetId)
	fmt.Printf("Finish Forward Clock\n")
}

// CLEAR
// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	switch dm.GetSimDemand().DemandType {

	case synerex.DemandType_SET_PARTICIPANTS_REQUEST:
		// 参加者リストをセットする要求
		callbackSetParticipantsRequest(dm)
	case synerex.DemandType_SET_AGENTS_REQUEST:
		// 参加者リストをセットする要求
		callbackSetAgentsRequest(dm)
	case synerex.DemandType_FORWARD_CLOCK_REQUEST:
		// クロックを進める要求
		callbackForwardClockRequest(dm)
	case synerex.DemandType_GET_SAME_AREA_AGENTS_REQUEST:
		// クロックを進める要求
		callbackGetSameAreaAgentsRequest(dm)
	default:
		//log.Println("demand callback is invalid.")
	}
}

// CLEAR
// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {

	switch sp.GetSimSupply().SupplyType {
	case synerex.SupplyType_GET_AREA_RESPONSE:
		com.SendToGetAreaResponse(sp)
	case synerex.SupplyType_GET_CLOCK_RESPONSE:
		com.SendToGetClockResponse(sp)
	case synerex.SupplyType_GET_SAME_AREA_AGENTS_RESPONSE:
		com.SendToGetSameAreaAgentsResponse(sp)
	case synerex.SupplyType_GET_NEIGHBOR_AREA_AGENTS_RESPONSE:
		com.SendToGetNeighborAreaAgentsResponse(sp)
	case synerex.SupplyType_REGIST_PARTICIPANT_RESPONSE:
		com.SendToRegistParticipantResponse(sp)
	case synerex.SupplyType_DELETE_PARTICIPANT_RESPONSE:
		com.SendToDeleteParticipantResponse(sp)
	default:
		//fmt.Println("order is invalid")
	}

}

func main() {
	flag.Parse()
	log.Printf("area id is: %v, agent type is %v", areaId, agentType)

	sxutil.RegisterNodeName(*nodesrv, "CarAreaProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { deleteParticipant(); conn.Close() })

	// synerex simulator
	sim = simulator.NewCarSimulator(1.0, 0.0)

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:CarArea, AreaId: %d}", areaId)

	// Clientとして登録
	com = communicator.NewCarCommunicator()

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	// channelごとのClientを作成
	com.RegistClients(client, argJson)
	// ChannelにSubscribe
	com.SubscribeAll(demandCallback, supplyCallback, &wg)
	wg.Wait()

	// start up(setArea)
	wg.Add(1)
	// 起動時にエリア情報を取得する
	getArea()
	// 新規参加登録
	registParticipant()
	// クロック情報を取得する
	getClock()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
