package main

import (
	//"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"github.com/paulmach/orb/geojson"
	"io/ioutil"
	"encoding/json"
	//"time"
	//"runtime"
	//"encoding/json"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
	"github.com/synerex/synerex_alpha/provider/simulation/pedestrian/communicator"
	"github.com/synerex/synerex_alpha/provider/simulation/pedestrian/simulator"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)


var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	areaIdFlag     = flag.Int("areaId", 1, "Area Id") 
	pidFlag     = flag.Int("pid", 1, "Provider Id") 
	areaJson     = flag.String("area_json", "", "Area Json") 
	areaInfo *area.Area2
	agentType  = common.AgentType_PEDESTRIAN               // PEDESTRIAN
	com        *communicator.PedCommunicator
	sim        *simulator.PedSimulator;
	areaId uint64
	pid uint64
)


func flagToAreaInfo(areaJson string) *area.Area2{
	bytes := []byte(areaJson)
	json.Unmarshal(bytes, &areaInfo)
    return areaInfo
}

func init(){
	flag.Parse()
	areaInfo = flagToAreaInfo(*areaJson)
	log.Printf("\x1b[31m\x1b[47m \nAreaInfo: %v \x1b[0m\n", areaInfo)
	areaId = uint64(*areaIdFlag)
	pid = uint64(*pidFlag)
}



// registParticipant: 新規参加登録をする関数
func registParticipant() {

	// 新規参加登録をするRequest
	participant := com.GetMyParticipant(areaId)

	com.RegistParticipantRequest(participant)
	// Responseの待機
	err := com.WaitRegistParticipantResponse()
	if err != nil {
		log.Printf("\x1b[31m\x1b[47m \n Error: %v \x1b[0m\n", err)
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
	participant := com.GetMyParticipant(areaId)
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
	com.CreateWaitIdList(agentType, areaId, sim.Area.NeighborAreas)

	// セット完了通知を送る
	com.SetParticipantsResponse(targetId)
	log.Printf("\x1b[30m\x1b[47m \n Finish: Set Participants \x1b[0m\n")
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

// callbackSetAgents: Agent情報をセットする要求
func callbackSetAgentsRequest(dm *pb.Demand) {
	agents := dm.GetSimDemand().GetSetAgentsRequest().GetAgents()
	targetId := dm.GetId()

	// Agent情報を追加する
	sim.AddAgents(agents)

	// セット完了通知を送る
	com.SetAgentsResponse(targetId)
	log.Printf("\x1b[30m\x1b[47m \n Finish: Agents information set. \n Total:  %v \n Add: %v \x1b[0m\n", len(sim.GetAgents()), len(agents))
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

// callbackClearAgentsRequest: Agent情報を消去する要求
func callbackClearAgentsRequest(dm *pb.Demand) {
	targetId := dm.GetId()

	// エージェントをクリアする
	sim.ClearAgents()
	// Responseを送る
	com.ClearAgentsResponse(targetId)
	log.Printf("\x1b[30m\x1b[47m \n Finish: Agents cleared.  \n Total:  %v \x1b[0m\n", len(sim.GetAgents()))
}



// callbackForwardClock: Agentを計算し、クロックを進める要求
func callbackForwardClockRequest(dm *pb.Demand) {
	//t_start := time.Now()
	com.InitChannel()
	dm.GetSimDemand().GetForwardClockRequest().GetStepNum()
	targetId := dm.GetId()

	// 同じエリアのAgent情報を取得する
	com.GetSameAreaAgentsRequest(areaId, agentType)
	// Responseの待機
	sameAreaAgents := com.WaitGetSameAreaAgentsResponse()

	// 次の時間のエージェントを計算する
	// agentがDuplicateで計算されている？
	nextControlAgents := sim.ForwardStep(sameAreaAgents)


	// 隣接エリアにエージェントの情報を送信
	com.GetNeighborAreaAgentsResponse(targetId, nextControlAgents)

	// 次の時刻の隣接しているエリアの同じAgentTypeのエージェント情報を取得する
	// 同じデータを取得している？
	neighborAreaAgents := com.WaitGetNeighborAreaAgentsResponse()

	// 重複エリアのエージェントを更新する
	nextDuplicateAgents := sim.UpdateDuplicateAgents(nextControlAgents, neighborAreaAgents)

	
	// Agentsをセットする
	sim.SetAgents(nextDuplicateAgents)

	// クロックを進める
	sim.ForwardGlobalTime()

	// 可視化プロバイダへ送信
	com.VisualizeAgentsResponse(nextControlAgents, areaId, agentType)

	// セット完了通知を送る
	com.ForwardClockResponse(targetId)
	//t_end := time.Now()
	//log.Printf("time: %v\n", t_end.Sub(t_start).Seconds())
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock forwarded. \n Time:  %v \n Agents Num: %v \x1b[0m\n", sim.GlobalTime, len(nextControlAgents))
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	switch dm.GetSimDemand().DemandType {

	case synerex.DemandType_SET_PARTICIPANTS_REQUEST:
		// 参加者リストをセットする要求
		callbackSetParticipantsRequest(dm)
	case synerex.DemandType_SET_AGENTS_REQUEST:
		// Agentをセットする要求
		callbackSetAgentsRequest(dm)
	case synerex.DemandType_CLEAR_AGENTS_REQUEST:
		// Agentをクリアする要求
		callbackClearAgentsRequest(dm)
	case synerex.DemandType_FORWARD_CLOCK_REQUEST:
		// クロックを進める要求
		callbackForwardClockRequest(dm)
	case synerex.DemandType_SET_CLOCK_REQUEST:
		// クロックをセットする要求
		callbackSetClockRequest(dm)
	case synerex.DemandType_GET_SAME_AREA_AGENTS_REQUEST:
		// 同じエリアのエージェントの要求
		callbackGetSameAreaAgentsRequest(dm)
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
	case synerex.SupplyType_GET_SAME_AREA_AGENTS_RESPONSE:
		// 同じエリアの異種エージェント情報の取得
		com.SendToGetSameAreaAgentsResponse(sp)
	case synerex.SupplyType_GET_NEIGHBOR_AREA_AGENTS_RESPONSE:
		// 隣接エリアの同種エージェント情報の取得
		com.SendToGetNeighborAreaAgentsResponse(sp)
	case synerex.SupplyType_REGIST_PARTICIPANT_RESPONSE:
		// 参加者登録完了通知の取得
		com.SendToRegistParticipantResponse(sp)
	case synerex.SupplyType_DELETE_PARTICIPANT_RESPONSE:
		// 参加者削除完了通知の取得
		com.SendToDeleteParticipantResponse(sp)

	default:
		//fmt.Println("order is invalid")
	}

}

func main() {
	log.Printf("\x1b[31m\x1b[47m \n SyneServ: %v, NodeServ: %v, AreaJson: %v Pid: %v   \x1b[0m\n", *serverAddr, *nodesrv, *areaJson, pid)
	//log.Printf("area id is: %v, agent type is %v", areaId, agentType)

	sxutil.RegisterNodeName(*nodesrv, "PedProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	// Communicator
	com = communicator.NewPedCommunicator(pid)

	sxutil.RegisterDeferFunction(func() { deleteParticipant(); conn.Close() })

	// Pedestrian Simulator
	sim = simulator.NewPedSimulator(1.0, 0.0)



	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Ped, AreaId: %d}", areaId)

	// プロバイダのsetup
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
