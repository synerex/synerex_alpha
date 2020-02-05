package communicator

import (
	"fmt"
	"context"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
	"github.com/synerex/synerex_alpha/sxutil"
)

var (
	mu                  sync.Mutex
)


type Clients struct {
	AgentClient       *sxutil.SMServiceClient
	ClockClient       *sxutil.SMServiceClient
	AreaClient        *sxutil.SMServiceClient
	RouteClient       *sxutil.SMServiceClient
	ParticipantClient *sxutil.SMServiceClient
}

// SynerexCommunicator :
type SynerexCommunicator struct {
	DmMap        map[uint64]*sxutil.DemandOpts
	SpMap        map[uint64]*sxutil.SupplyOpts
	Participants []*participant.Participant
	//Todo: Clientsを追加
	MyClients *Clients
}

func NewSynerexCommunicator() *SynerexCommunicator {
	s := &SynerexCommunicator{
		DmMap: make(map[uint64]*sxutil.DemandOpts),
		SpMap: make(map[uint64]*sxutil.SupplyOpts),
	}
	return s
}

func (s *SynerexCommunicator) SetClients(clients *Clients) {
	s.MyClients = clients
}

func (s *SynerexCommunicator) AddParticipant(participant *participant.Participant) {
	s.Participants = append(s.Participants, participant)
	fmt.Printf("debug add participant", s.Participants)
}

func (s *SynerexCommunicator) DeleteParticipant(targetPar *participant.Participant) {
	newParticipants := make([]*participant.Participant, 0)

	for _, par := range s.Participants {
		// ターゲット以外を再代入
		if par.GetChannelId().GetParticipantChannelId() != targetPar.GetChannelId().GetParticipantChannelId() {
			newParticipants = append(newParticipants, par)
		}
	}
	fmt.Printf("debug delete participant", newParticipants)
	s.SetParticipants(newParticipants)
}

func (s *SynerexCommunicator) SetParticipants(participants []*participant.Participant) {
	s.Participants = participants
}

func (s *SynerexCommunicator) GetParticipants() []*participant.Participant {
	return s.Participants
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (s *SynerexCommunicator) SubscribeAll(demandCallback func(*sxutil.SMServiceClient, *pb.Demand), supplyCallback func(*sxutil.SMServiceClient, *pb.Supply), wg2 *sync.WaitGroup) {

	// SubscribeDemand, SubscribeSupply
	wg := sync.WaitGroup{}
	wg.Add(1)
	go s.SubscribeDemand(s.MyClients.AgentClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeDemand(s.MyClients.ClockClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeDemand(s.MyClients.AreaClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeDemand(s.MyClients.ParticipantClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeDemand(s.MyClients.RouteClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.MyClients.ClockClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.MyClients.AreaClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.MyClients.AgentClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.MyClients.ParticipantClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.MyClients.RouteClient, supplyCallback, &wg)
	wg.Wait()

	wg2.Done()
}

// Wait: 同期が完了するまで待機する関数
func (s *SynerexCommunicator) Wait(idList []uint64, waitCh chan *pb.Supply) map[uint64]*pb.Supply {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	wg.Add(1)
	pspMap := make(map[uint64]*pb.Supply)
	go func() {
		for {
			select {
			case psp := <-waitCh:


				mu.Lock()
				// spのidがidListに入っているか
				if isSpInIdList(psp, idList){
				pspMap[psp.SenderId] = psp
				// 同期が終了したかどうか
				if isFinishSync(pspMap, idList) {

					mu.Unlock()
					wg.Done()
					return
				}
			}
				mu.Unlock()

			}
		}
	}()
	wg.Wait()
	return pspMap
}

// SendToWait: Waitしているchに受け取ったsupplyを送る
func (s *SynerexCommunicator) SendToWait(sp *pb.Supply, waitCh chan *pb.Supply) {
	waitCh <- sp
}

// isSpInIdList : spのidがidListに入っているか
func isSpInIdList(sp *pb.Supply, idlist []uint64) bool {
	senderId := sp.SenderId
	for _, id := range idlist {
		if senderId == id{
			return true
		}
	}
	return false
}

// isFinishSync : 必要な全てのSupplyを受け取り同期が完了したかどうか
func isFinishSync(pspMap map[uint64]*pb.Supply, idlist []uint64) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			senderId := uint64(sp.SenderId)
			if uint64(id) == senderId {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}

func (s *SynerexCommunicator) IsSupplyTarget(sp *pb.Supply) bool {
	spid := sp.TargetId
	for id, _ := range s.DmMap {
		if id == spid {
			return true
		}
	}
	return false
}

// getAreaRequest :　エリア情報を取得する
func (s *SynerexCommunicator) GetAreaRequest(areaId uint64) {
	getAreaRequest := &area.GetAreaRequest{
		Id: areaId,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_GET_AREA_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_GetAreaRequest{getAreaRequest},
	}

	nm := ""
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}
	s.sendDemand(s.MyClients.AreaClient, opts)
}



// CLEAR OK
// getAreaSupply :　エリア情報を取得する
func (s *SynerexCommunicator) GetAreaResponse(tid uint64, areaInfo *area.Area2) {
	getAreaResponse := &area.GetAreaResponse{
		Area: areaInfo,
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_GET_AREA_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_GetAreaResponse{getAreaResponse},
	}

	nm := "getArea respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    tid,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.AreaClient, opts)
}

// getClockRequest :　エリア情報を取得する
func (s *SynerexCommunicator) GetClockRequest() {
	getClockRequest := &clock.GetClockRequest{}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_GET_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_GetClockRequest{getClockRequest},
	}

	nm := ""
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}
	s.sendDemand(s.MyClients.ClockClient, opts)
}

// getClockSupply :　エリア情報を取得する
func (s *SynerexCommunicator) GetClockResponse(dm *pb.Demand, clockInfo *clock.Clock) {
	getClockResponse := &clock.GetClockResponse{
		Clock: clockInfo,
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_GET_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_GetClockResponse{getClockResponse},
	}

	nm := "getClock respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    dm.GetId(),
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ClockClient, opts)
}

// SetClockRequest : クロック情報をセットする
func (s *SynerexCommunicator) SetClockRequest(clockInfo *clock.Clock) {
	setClockRequest := &clock.SetClockRequest{
		Clock: clockInfo,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_SET_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_SetClockRequest{setClockRequest},
	}

	nm := ""
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}
	s.sendDemand(s.MyClients.ClockClient, opts)
}

// SetClockSupply :　SetClockのResponse
func (s *SynerexCommunicator) SetClockResponse(targetId uint64) {
	setClockResponse := &clock.SetClockResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_SET_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_SetClockResponse{setClockResponse},
	}

	nm := "setClock respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    targetId,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ClockClient, opts)
}

// StartClockRequest : クロックをスタートする
func (s *SynerexCommunicator) StartClockRequest(stepNum uint64) {
	startClockRequest := &clock.StartClockRequest{
		StepNum: stepNum,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_SET_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_StartClockRequest{startClockRequest},
	}

	nm := ""
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}
	s.sendDemand(s.MyClients.ClockClient, opts)
}

// StartClockSupply :　StartClockのResponse
func (s *SynerexCommunicator) StartClockResponse(targetId uint64) {
	startClockResponse := &clock.StartClockResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_SET_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_StartClockResponse{startClockResponse},
	}

	nm := "startClock respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    targetId,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ClockClient, opts)
}

// StopClockRequest : クロックをスタートする
func (s *SynerexCommunicator) StopClockRequest() {
	stopClockRequest := &clock.StopClockRequest{
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_SET_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_StopClockRequest{stopClockRequest},
	}

	nm := ""
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}
	s.sendDemand(s.MyClients.ClockClient, opts)
}

// StopClockSupply :　StopClockのResponse
func (s *SynerexCommunicator) StopClockResponse(targetId uint64) {
	stopClockResponse := &clock.StopClockResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_SET_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_StopClockResponse{stopClockResponse},
	}

	nm := "stopClock respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    targetId,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ClockClient, opts)
}

// getNeighborAreaAgentsResponse :　getNeighborAreaAgentDemandに対する応答
func (s *SynerexCommunicator) GetNeighborAreaAgentsResponse(targetId uint64, agents []*agent.Agent) {
	getNeighborAreaAgentsResponse := &agent.GetNeighborAreaAgentsResponse{
		Agents: agents,
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_GET_NEIGHBOR_AREA_AGENTS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_GetNeighborAreaAgentsResponse{getNeighborAreaAgentsResponse},
	}

	nm := "getNeighborAreaAgentSupply by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    targetId,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.AgentClient, opts)
}

// getSameAreaAgentsResponse :　getSameAreaAgentDemandに対する応答
func (s *SynerexCommunicator) GetSameAreaAgentsResponse(targetId uint64, agents []*agent.Agent) {
	getSameAreaAgentsResponse := &agent.GetSameAreaAgentsResponse{
		Agents: agents,
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_GET_SAME_AREA_AGENTS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_GetSameAreaAgentsResponse{getSameAreaAgentsResponse},
	}

	nm := "getSameAreaAgentSupply by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    targetId,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.AgentClient, opts)
}

// getSameAreaAgentsRequest :　Agentsを設置するRequest
func (s *SynerexCommunicator) GetSameAreaAgentsRequest(areaId uint64, agentType common.AgentType) {
	getSameAreaAgentsRequest := &agent.GetSameAreaAgentsRequest{
		AreaId:    areaId,
		AgentType: agentType,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_GET_SAME_AREA_AGENTS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_GetSameAreaAgentsRequest{getSameAreaAgentsRequest},
	}

	nm := "GetSameAreaAgentsRequest"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.AgentClient, opts)
}

// setAgentsResponse :　setAgentDemandに対する応答
func (s *SynerexCommunicator) SetAgentsResponse(targetId uint64) {
	setAgentsResponse := &agent.SetAgentsResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_SET_AGENTS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_SetAgentsResponse{setAgentsResponse},
	}

	nm := "setAgentSupply by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    targetId,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.AgentClient, opts)
}

// setAgentsDemand :　Agentsを設置するDemand
func (s *SynerexCommunicator) SetAgentsRequest(agents []*agent.Agent) {
	setAgentsRequest := &agent.SetAgentsRequest{
		Agents: agents,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_SET_AGENTS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_SetAgentsRequest{setAgentsRequest},
	}

	nm := "SetAgentDemand"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.AgentClient, opts)
}

// ClearAgentsResponse :　clearAgentDemandに対する応答
func (s *SynerexCommunicator) ClearAgentsResponse(targetId uint64) {
	clearAgentsResponse := &agent.ClearAgentsResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_CLEAR_AGENTS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_ClearAgentsResponse{clearAgentsResponse},
	}

	nm := "ClearAgentSupply by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    targetId,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.AgentClient, opts)
}

// ClearAgentsDemand :　Agentsを消去するDemand
func (s *SynerexCommunicator) ClearAgentsRequest() {
	clearAgentsRequest := &agent.ClearAgentsRequest{}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_CLEAR_AGENTS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_ClearAgentsRequest{clearAgentsRequest},
	}

	nm := "ClearAgentDemand"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.AgentClient, opts)
}

// setAgentsDemand :　Agentsを設置するDemand
func (s *SynerexCommunicator) VisualizeAgentsResponse(agents []*agent.Agent, areaId uint64, agentType common.AgentType) {
	visualizeAgentsResponse := &agent.VisualizeAgentsResponse{
		AreaId:    areaId,
		AgentType: agentType,
		Agents:    agents,
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_VISUALIZE_AGENTS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_VisualizeAgentsResponse{visualizeAgentsResponse},
	}

	nm := "VisualizeAgents to agentCh "
	js := ""
	opts := &sxutil.SupplyOpts{
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendSupply(s.MyClients.AgentClient, opts)
}

// forwardClockResponse :　forwardClockDemandに対するResponse
func (s *SynerexCommunicator) ForwardClockResponse(tid uint64) {
	forwardClockResponse := &clock.ForwardClockResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_FORWARD_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_ForwardClockResponse{forwardClockResponse},
	}

	nm := "forwardClock to clockCh "
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    tid,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ClockClient, opts)

}

// forwardClockRequest :　Clockを進めるRequest
func (s *SynerexCommunicator) ForwardClockRequest(stepNum uint64) {
	forwardClockRequest := &clock.ForwardClockRequest{
		StepNum: stepNum,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_FORWARD_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_ForwardClockRequest{forwardClockRequest},
	}

	nm := "ForwardClockDemand"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.ClockClient, opts)
}

// downScenarioResponse :　Scenarioの障害のResponse
func (s *SynerexCommunicator) DownScenarioResponse(tid uint64) {
	downScenarioResponse := &participant.DownScenarioResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_DOWN_SCENARIO_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_DownScenarioResponse{downScenarioResponse},
	}

	nm := "downScenario to clockCh "
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    tid,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ParticipantClient, opts)

}

// downScenarioRequest :　Scenarioの障害を知らせる
func (s *SynerexCommunicator) DownScenarioRequest() {
	downScenarioRequest := &participant.DownScenarioRequest{
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_DOWN_SCENARIO_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_DownScenarioRequest{downScenarioRequest},
	}

	nm := "DownScenarioDemand"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.ParticipantClient, opts)
}


// backClockResponse :　backClockDemandに対するResponse
func (s *SynerexCommunicator) BackClockResponse(tid uint64) {
	backClockResponse := &clock.BackClockResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_BACK_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_BackClockResponse{backClockResponse},
	}

	nm := "backClock to clockCh "
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    tid,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ClockClient, opts)

}

// backClockRequest :　Clockを進めるRequest
func (s *SynerexCommunicator) BackClockRequest(stepNum uint64) {
	backClockRequest := &clock.BackClockRequest{
		StepNum: stepNum,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_BACK_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_BackClockRequest{backClockRequest},
	}

	nm := "BackClockDemand"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.ClockClient, opts)

}



// registParticipantResponse :　登録完了通知
func (s *SynerexCommunicator) RegistParticipantResponse(dm *pb.Demand) {

	registParticipantResponse := &participant.RegistParticipantResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_REGIST_PARTICIPANT_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_RegistParticipantResponse{registParticipantResponse},
	}

	nm := "getParticipant respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    dm.GetId(),
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ParticipantClient, opts)
}

// RegistParticipantRequest :　参加者を登録するDemand
func (s *SynerexCommunicator) RegistParticipantRequest(participantInfo *participant.Participant) {

	registParticipantRequest := &participant.RegistParticipantRequest{
		Participant: participantInfo,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_REGIST_PARTICIPANT_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_RegistParticipantRequest{registParticipantRequest},
	}

	nm := "RegistParticipantRequest"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.ParticipantClient, opts)
}

// DeleteParticipantResponse :　参加取り消しをする意向をSupply
func (s *SynerexCommunicator) DeleteParticipantResponse(dm *pb.Demand) {

	deleteParticipantResponse := &participant.DeleteParticipantResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_DELETE_PARTICIPANT_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_DeleteParticipantResponse{deleteParticipantResponse},
	}

	nm := "getParticipant respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    dm.GetId(),
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ParticipantClient, opts)
}

// DeleteParticipantRequest :　参加者を削除するDemand
func (s *SynerexCommunicator) DeleteParticipantRequest(participantInfo *participant.Participant) {

	deleteParticipantRequest := &participant.DeleteParticipantRequest{
		Participant: participantInfo,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_DELETE_PARTICIPANT_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_DeleteParticipantRequest{deleteParticipantRequest},
	}

	nm := "DeleteParticipantRequest"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.ParticipantClient, opts)
}

// setParticipantsResponse :　setParticipantsDemandに対するSupply
func (s *SynerexCommunicator) SetParticipantsResponse(tid uint64) {

	setParticipantsResponse := &participant.SetParticipantsResponse{}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_SET_PARTICIPANTS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_SetParticipantsResponse{setParticipantsResponse},
	}

	nm := "SetParticipants respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:    tid,
		Name:      nm,
		JSON:      js,
		SimSupply: simSupply,
	}

	s.sendProposeSupply(s.MyClients.ParticipantClient, opts)
}

// setParticipantsRequest :　参加者を共有するDemand
func (s *SynerexCommunicator) SetParticipantsRequest() {
	participants := s.Participants

	setParticipantsRequest := &participant.SetParticipantsRequest{
		Participants: participants,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_SET_PARTICIPANTS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_SetParticipantsRequest{setParticipantsRequest},
	}

	nm := "setParticipants order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.ParticipantClient, opts)
}

// CollectParticipantsRequest :　参加者を共有するDemand
func (s *SynerexCommunicator) CollectParticipantsRequest() {

	collectParticipantsRequest := &participant.CollectParticipantsRequest{}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_COLLECT_PARTICIPANTS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_CollectParticipantsRequest{collectParticipantsRequest},
	}

	nm := "collectParticipants order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.ParticipantClient, opts)
}

// NotifyStartUpRequest :　参加者を共有するDemand
func (s *SynerexCommunicator) NotifyStartUpRequest(providerType participant.ProviderType) {

	notifyStartUpRequest := &participant.NotifyStartUpRequest{
		ProviderType: providerType,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_NOTIFY_START_UP_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_NotifyStartUpRequest{notifyStartUpRequest},
	}

	nm := "notifyStartUp"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	s.sendDemand(s.MyClients.ParticipantClient, opts)
}

func (s *SynerexCommunicator) sendProposeSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts) {
	mu.Lock()
	id := sclient.ProposeSupply(opts)
	s.SpMap[id] = opts // my demand options
	mu.Unlock()
}

func (s *SynerexCommunicator) sendSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts) {
	mu.Lock()
	id := sclient.RegisterSupply(opts)
	s.SpMap[id] = opts // my demand options
	mu.Unlock()
}

func (s *SynerexCommunicator) sendDemand(sclient *sxutil.SMServiceClient, opts *sxutil.DemandOpts) {
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	s.DmMap[id] = opts // my demand options
	mu.Unlock()
}

func (s *SynerexCommunicator) SubscribeSupply(client *sxutil.SMServiceClient, supplyCallback func(*sxutil.SMServiceClient, *pb.Supply), wg *sync.WaitGroup) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback, wg)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func (s *SynerexCommunicator) SubscribeDemand(client *sxutil.SMServiceClient, demandCallback func(*sxutil.SMServiceClient, *pb.Demand), wg *sync.WaitGroup) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback, wg)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}
