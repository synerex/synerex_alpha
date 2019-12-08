package provider

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/sxutil"
)

var (
	mu sync.Mutex
	//waitCh chan *pb.Supply
)

const (
	GET_PARTICIPANT_DEMAND = "GET_PARTICIPANT_DEMAND"
	GET_PARTICIPANT_SUPPLY = "GET_PARTICIPANT_SUPPLY"
	SET_PARTICIPANT_DEMAND = "SET_PARTICIPANT_DEMAND"
	SET_PARTICIPANT_SUPPLY = "SET_PARTICIPANT_SUPPLY"
	SET_AGENTS_DEMAND      = "SET_AGENTS_DEMAND"
	SET_AGENTS_SUPPLY      = "SET_AGENTS_SUPPLY"
	FORWARD_CLOCK_DEMAND   = "FORWARD_CLOCK_DEMAND"
	FORWARD_CLOCK_SUPPLY   = "FORWARD_CLOCK_SUPPLY"
	FORWARD_AGENTS_DEMAND  = "FORWARD_AGENTS_DEMAND"
	FORWARD_AGENTS_SUPPLY  = "FORWARD_AGENTS_SUPPLY"
)

// SynerexSimulator :
type SynerexProvider struct {
	DmMap              map[uint64]*sxutil.DemandOpts
	SpMap              map[uint64]*sxutil.SupplyOpts
	SameAreaIdList     []uint64
	NeighborAreaIdList []uint64
	AgentClient        *sxutil.SMServiceClient
	ClockClient        *sxutil.SMServiceClient
	AreaClient         *sxutil.SMServiceClient
	RouteClient        *sxutil.SMServiceClient
	ParticipantClient  *sxutil.SMServiceClient
}

func NewSynerexProvider() *SynerexProvider {
	s := &SynerexProvider{
		DmMap: make(map[uint64]*sxutil.DemandOpts),
		SpMap: make(map[uint64]*sxutil.SupplyOpts),
	}
	return s
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (s *SynerexProvider) SetupProvider(client pb.SynerexClient, argJson string, demandCallback func(*sxutil.SMServiceClient, *pb.Demand), supplyCallback func(*sxutil.SMServiceClient, *pb.Supply), wg2 *sync.WaitGroup) {

	s.AgentClient = sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE, argJson)
	s.ClockClient = sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE, argJson)
	s.AreaClient = sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE, argJson)
	s.ParticipantClient = sxutil.NewSMServiceClient(client, pb.ChannelType_PARTICIPANT_SERVICE, argJson)
	s.RouteClient = sxutil.NewSMServiceClient(client, pb.ChannelType_ROUTE_SERVICE, argJson)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go s.SubscribeDemand(s.AgentClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeDemand(s.ClockClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeDemand(s.AreaClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeDemand(s.ParticipantClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeDemand(s.RouteClient, demandCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.ClockClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.AreaClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.AgentClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.ParticipantClient, supplyCallback, &wg)
	wg.Wait()

	wg.Add(1)
	go s.SubscribeSupply(s.RouteClient, supplyCallback, &wg)
	wg.Wait()

	wg2.Done()
}

func (s *SynerexProvider) SetRelateAreaIdList(sameAreaIdList []uint64, neighborAreaIdList []uint64) {
	s.SameAreaIdList = sameAreaIdList
	s.NeighborAreaIdList = neighborAreaIdList
}

// Wait: 同期が完了するまで待機する関数
func (s *SynerexProvider) Wait(idList []uint64, waitCh chan *pb.Supply) map[uint64]*pb.Supply {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	wg.Add(1)
	pspMap := make(map[uint64]*pb.Supply)
	go func() {
		for {
			select {
			case psp := <-waitCh:

				mu.Lock()
				pspMap[psp.SenderId] = psp
				if isFinishSync(pspMap, idList) {

					mu.Unlock()
					wg.Done()
					return
				}
				mu.Unlock()

			}
		}
	}()
	wg.Wait()
	return pspMap
}

// SendToWait: Waitしているchに受け取ったsupplyを送る
func (s *SynerexProvider) SendToWait(sp *pb.Supply, waitCh chan *pb.Supply) {
	fmt.Printf("sendTowait")
	waitCh <- sp
}

// isFinishSync : 必要な全てのSupplyを受け取り同期が完了したかどうか
func isFinishSync(pspMap map[uint64]*pb.Supply, idlist []uint64) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			senderId := uint32(sp.SenderId)
			if uint32(id) == senderId {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}

func (s *SynerexProvider) IsSupplyTarget(sp *pb.Supply) bool {
	spid := sp.TargetId
	for id, _ := range s.DmMap {
		if id == spid {
			return true
		}
	}
	return false
}

// getAgentsDemand :　同じエリアのエージェント情報を取得する
func (s *SynerexProvider) GetAgentsDemand(time uint32, areaId uint32, agentType agent.AgentType) {
	getAgentsDemand := &agent.GetAgentsDemand{
		Time:       time,
		AreaId:     areaId,
		AgentType:  agentType,
		StatusType: 2, // NONE
		Meta:       "",
	}

	nm := "getAgents order by ped-area-provider"
	js := ""
	opts := &sxutil.DemandOpts{
		Name:            nm,
		JSON:            js,
		GetAgentsDemand: getAgentsDemand,
	}
	s.sendDemand(s.AgentClient, opts)
}

// getAreaDemand :　エリア情報を取得する
func (s *SynerexProvider) GetAreaDemand(time uint32, areaId uint32) {
	getAreaDemand := &area.GetAreaDemand{
		Time:       time,
		AreaId:     areaId,
		StatusType: 2, //NONE
		Meta:       "",
	}

	nm := "getArea order by ped-area-provider"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, GetAreaDemand: getAreaDemand}
	s.sendDemand(s.AreaClient, opts)
}

// getAreaDemand :　エリア情報を取得する
func (s *SynerexProvider) GetAreaSupply(tid uint64, areaInfo *area.AreaInfo) {
	getAreaSupply := &area.GetAreaSupply{
		AreaInfo:   areaInfo,
		StatusType: 0, //OK
		Meta:       "",
	}

	nm := "getArea respnse by area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:        tid,
		Name:          nm,
		JSON:          js,
		GetAreaSupply: getAreaSupply,
	}

	s.sendProposeSupply(s.AreaClient, opts)
}

// getAgentsRouteSupply :　エージェントのルート情報を提供する
func (s *SynerexProvider) GetAgentsRouteSupply(tid uint64, agentsInfo []*agent.AgentInfo) {
	getAgentsRouteSupply := &agent.GetAgentsRouteSupply{
		AgentsInfo: agentsInfo,
		StatusType: 0, //OK
		Meta:       "",
	}

	nm := "getRoute respnse by route-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               tid,
		Name:                 nm,
		JSON:                 js,
		GetAgentsRouteSupply: getAgentsRouteSupply,
	}

	s.sendProposeSupply(s.RouteClient, opts)
}

// getAgentsRouteDemand :　エージェントのルート情報を取得する
func (s *SynerexProvider) GetAgentsRouteDemand(agentsInfo []*agent.AgentInfo) {
	getAgentsRouteDemand := &agent.GetAgentsRouteDemand{
		AgentsInfo: agentsInfo,
		StatusType: 2, //NONE
		Meta:       "",
	}

	nm := "getAgentsRouteDemand order by ped-area-provider"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, GetAgentsRouteDemand: getAgentsRouteDemand}
	s.sendDemand(s.RouteClient, opts)
}

// setAgentsSupply :　setAgentDemandに対する応答
func (s *SynerexProvider) SetAgentsSupply(tid uint64, time uint32, areaId uint32, agentType agent.AgentType) {
	setAgentsSupply := &agent.SetAgentsSupply{
		Time:       time,
		AreaId:     areaId,
		AgentType:  agentType,
		StatusType: 0,
		Meta:       "",
	}

	nm := "setAgentSupply by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:          tid,
		Name:            nm,
		JSON:            js,
		SetAgentsSupply: setAgentsSupply,
	}

	s.sendProposeSupply(s.AgentClient, opts)
}

// setAgentsDemand :　Agentsを設置するDemand
func (s *SynerexProvider) SetAgentsDemand(agentsInfo []*agent.AgentInfo) {
	setAgentsDemand := &agent.SetAgentsDemand{
		AgentsInfo: agentsInfo,
		StatusType: 2, // NONE
		Meta:       "",
	}

	nm := "SetAgentDemand"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SetAgentsDemand: setAgentsDemand}

	s.sendDemand(s.AgentClient, opts)
}

// forwardAgentsSupply :　次の時刻のエージェント情報をSupplyする
func (s *SynerexProvider) ForwardAgentsSupply(tid uint64, time uint32, areaId uint32, agentsInfo []*agent.AgentInfo, agentType agent.AgentType) {
	forwardAgentsSupply := &agent.ForwardAgentsSupply{
		Time:       time,
		AreaId:     areaId,
		AgentType:  agentType,
		AgentsInfo: agentsInfo,
		StatusType: 0, //OK
		Meta:       "",
	}

	opts := &sxutil.SupplyOpts{
		Target:              tid,
		ForwardAgentsSupply: forwardAgentsSupply,
	}

	s.sendProposeSupply(s.AgentClient, opts)

}

// forwardClockSupply :　forwardClockDemandに対するSupply
func (s *SynerexProvider) ForwardClockSupply(tid uint64, clockInfo *clock.ClockInfo) {
	forwardClockSupply := &clock.ForwardClockSupply{
		ClockInfo:  clockInfo,
		StatusType: 0, // OK
		Meta:       "",
	}

	nm := "forwardClock to clockCh respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:             tid,
		Name:               nm,
		JSON:               js,
		ForwardClockSupply: forwardClockSupply,
	}

	s.sendProposeSupply(s.ClockClient, opts)

}

// forwardClockDemand :　Clockを進めるDemand
func (s *SynerexProvider) ForwardClockDemand(time uint64, cyclenum uint64) {
	forwardClockDemand := &clock.ForwardClockDemand{
		Time:       uint32(time),
		CycleNum:   uint32(cyclenum),
		StatusType: 2, // NONE
		Meta:       "",
	}

	nm := "ForwardClockDemand"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, ForwardClockDemand: forwardClockDemand}

	s.sendDemand(s.ClockClient, opts)

}

// getParticipantSupply :　参加する意向をSupply
func (s *SynerexProvider) GetParticipantSupply(tid uint64, participantInfo *participant.ParticipantInfo) {

	getParticipantSupply := &participant.GetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "getParticipant respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               tid,
		Name:                 nm,
		JSON:                 js,
		GetParticipantSupply: getParticipantSupply,
	}

	s.sendProposeSupply(s.ParticipantClient, opts)
}

// getParticipantDemand :　参加者を募集するDemand
func (s *SynerexProvider) GetParticipantDemand(participantInfo *participant.ParticipantInfo) {

	getParticipantDemand := &participant.GetParticipantDemand{
		ParticipantInfo: participantInfo,
		StatusType:      2, // NONE
		Meta:            "",
	}

	nm := "GetParticipantDemand"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, GetParticipantDemand: getParticipantDemand}

	s.sendDemand(s.ParticipantClient, opts)
}

// setParticipantSupply :　setParticipantDemandに対するSupply
func (s *SynerexProvider) SetParticipantSupply(tid uint64, participantInfo *participant.ParticipantInfo) {

	setParticipantSupply := &participant.SetParticipantSupply{
		ParticipantInfo: participantInfo,
		StatusType:      0, // OK
		Meta:            "",
	}

	nm := "SetParticipant respnse by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:               tid,
		Name:                 nm,
		JSON:                 js,
		SetParticipantSupply: setParticipantSupply,
	}

	s.sendProposeSupply(s.ParticipantClient, opts)
}

// setParticipantDemand :　参加者を共有するDEmand
func (s *SynerexProvider) SetParticipantDemand(participantsInfo []*participant.ParticipantInfo) {

	setParticipantDemand := &participant.SetParticipantDemand{
		ParticipantsInfo: participantsInfo,
		StatusType:       2, // NONE
		Meta:             "",
	}

	nm := "setParticipant order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SetParticipantDemand: setParticipantDemand}

	s.sendDemand(s.ParticipantClient, opts)
}

// getAgentsSupply :　同じエリアのエージェント情報を提供する
func (s *SynerexProvider) GetAgentsSupply(tid uint64, time uint32, areaId uint32, agentsInfo []*agent.AgentInfo, agentType agent.AgentType) {
	getAgentsSupply := &agent.GetAgentsSupply{
		Time:       time,
		AreaId:     areaId,
		AgentType:  agentType, //Ped
		AgentsInfo: agentsInfo,
		StatusType: 0, //OK
		Meta:       "",
	}

	nm := "getAgentSupply  by ped-area-provider"
	js := ""
	opts := &sxutil.SupplyOpts{
		Target:          tid,
		Name:            nm,
		JSON:            js,
		GetAgentsSupply: getAgentsSupply,
	}
	s.sendProposeSupply(s.AgentClient, opts)
}

func (s *SynerexProvider) sendProposeSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts) {
	mu.Lock()
	id := sclient.ProposeSupply(opts)
	s.SpMap[id] = opts // my demand options
	mu.Unlock()
}

func (s *SynerexProvider) sendSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts) {
	mu.Lock()
	id := sclient.RegisterSupply(opts)
	s.SpMap[id] = opts // my demand options
	mu.Unlock()
}

func (s *SynerexProvider) sendDemand(sclient *sxutil.SMServiceClient, opts *sxutil.DemandOpts) {
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	s.DmMap[id] = opts // my demand options
	mu.Unlock()
}

func (s *SynerexProvider) SubscribeSupply(client *sxutil.SMServiceClient, supplyCallback func(*sxutil.SMServiceClient, *pb.Supply), wg *sync.WaitGroup) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback, wg)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func (s *SynerexProvider) SubscribeDemand(client *sxutil.SMServiceClient, demandCallback func(*sxutil.SMServiceClient, *pb.Demand), wg *sync.WaitGroup) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback, wg)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func (s *SynerexProvider) CheckDemandType(dm *pb.Demand) string {
	// clock
	if dm.GetArg_SetClockDemand() != nil {
		return "SET_CLOCK_DEMAND"
	}
	if dm.GetArg_ForwardClockDemand() != nil {
		return "FORWARD_CLOCK_DEMAND"
	}
	if dm.GetArg_BackClockDemand() != nil {
		return "BACK_CLOCK_DEMAND"
	}
	// area
	if dm.GetArg_GetAreaDemand() != nil {
		return "GET_AREA_DEMAND"
	}
	// agents
	if dm.GetArg_GetAgentsDemand() != nil {
		return "GET_AGENTS_DEMAND"
	}
	if dm.GetArg_SetAgentsDemand() != nil {
		return "SET_AGENTS_DEMAND"
	}
	// participant
	if dm.GetArg_GetParticipantDemand() != nil {
		return "GET_PARTICIPANT_DEMAND"
	}
	if dm.GetArg_SetParticipantDemand() != nil {
		return "SET_PARTICIPANT_DEMAND"
	}
	// route
	if dm.GetArg_GetAgentRouteDemand() != nil {
		return "GET_AGENT_ROUTE_DEMAND"
	}
	if dm.GetArg_GetAgentsRouteDemand() != nil {
		return "GET_AGENTS_ROUTE_DEMAND"
	}

	return "INVALID_TYPE"
}

func (s *SynerexProvider) CheckSupplyType(sp *pb.Supply) string {
	// clock
	if sp.GetArg_SetClockSupply() != nil {
		return "SET_CLOCK_SUPPLY"
	}
	if sp.GetArg_ForwardClockSupply() != nil {
		return "FORWARD_CLOCK_SUPPLY"
	}
	if sp.GetArg_BackClockSupply() != nil {
		return "BACK_CLOCK_SUPPLY"
	}
	// area
	if sp.GetArg_GetAreaSupply() != nil {
		return "GET_AREA_SUPPLY"
	}
	// agents
	if sp.GetArg_GetAgentsSupply() != nil {
		return "GET_AGENTS_SUPPLY"
	}
	if sp.GetArg_SetAgentsSupply() != nil {
		return "SET_AGENTS_SUPPLY"
	}
	if sp.GetArg_ForwardAgentsSupply() != nil {
		return "FORWARD_AGENTS_SUPPLY"
	}
	// participant
	if sp.GetArg_GetParticipantSupply() != nil {
		return "GET_PARTICIPANT_SUPPLY"
	}
	if sp.GetArg_SetParticipantSupply() != nil {
		return "SET_PARTICIPANT_SUPPLY"
	}
	// route
	if sp.GetArg_GetAgentRouteSupply() != nil {
		return "GET_AGENT_ROUTE_SUPPLY"
	}
	if sp.GetArg_GetAgentsRouteSupply() != nil {
		return "GET_AGENTS_ROUTE_SUPPLY"
	}

	return "INVALID_TYPE"
}
