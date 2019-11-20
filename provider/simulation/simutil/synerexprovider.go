package simutil

import (
	"context"
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
)

func init() {
}

// SynerexSimulator :
type SynerexProvider struct {
	DmMap             map[uint64]*sxutil.DemandOpts
	DmIdList          []uint64
	SpMap             map[uint64]*sxutil.SupplyOpts
	SpIdList          []uint64
	AgentClient       *sxutil.SMServiceClient
	ClockClient       *sxutil.SMServiceClient
	AreaClient        *sxutil.SMServiceClient
	RouteClient       *sxutil.SMServiceClient
	ParticipantClient *sxutil.SMServiceClient
}

func NewSynerexProvider() *SynerexProvider {
	s := &SynerexProvider{
		DmMap:    make(map[uint64]*sxutil.DemandOpts),
		DmIdList: make([]uint64, 0),
		SpMap:    make(map[uint64]*sxutil.SupplyOpts),
		SpIdList: make([]uint64, 0),
	}
	return s
}

// RegisterClient :　ClientとしてNodeServerに登録する
func (s *SynerexProvider) RegisterClient(client pb.SynerexClient, channelTypes []pb.ChannelType, argJson string) {
	for _, channelType := range channelTypes {
		switch channelType {
		case pb.ChannelType_AGENT_SERVICE:
			s.AgentClient = sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE, argJson)
		case pb.ChannelType_CLOCK_SERVICE:
			s.ClockClient = sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE, argJson)
		case pb.ChannelType_AREA_SERVICE:
			s.AreaClient = sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE, argJson)
		case pb.ChannelType_PARTICIPANT_SERVICE:
			s.ParticipantClient = sxutil.NewSMServiceClient(client, pb.ChannelType_PARTICIPANT_SERVICE, argJson)
		case pb.ChannelType_ROUTE_SERVICE:
			s.RouteClient = sxutil.NewSMServiceClient(client, pb.ChannelType_ROUTE_SERVICE, argJson)
		}

	}

}

// AddAgent :　エージェントを追加する関数
func Wait(pspMap map[uint64]*pb.Supply, idList []uint32, syncCh chan *pb.Supply) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	wg.Add(1)
	go func() {
		for {
			select {
			case psp := <-syncCh:
				mu.Lock()
				pspMap[psp.SenderId] = psp
				if isFinishSynerexProvider(pspMap, idList) {

					mu.Unlock()
					wg.Done()
					return
				}
				mu.Unlock()

			}
		}
	}()
	wg.Wait()
}

// isFinishSynerexProvider :
func isFinishSynerexProvider(pspMap map[uint64]*pb.Supply, idlist []uint32) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			senderId := uint32(sp.SenderId)
			if id == senderId {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}

// getAgentsDemand :　同じエリアのエージェント情報を取得する
func (s *SynerexProvider) GetAgentsDemand(time uint32, areaId uint32) {
	getAgentsDemand := &agent.GetAgentsDemand{
		Time:       time,
		AreaId:     areaId,
		AgentType:  0, //Ped
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
func (s *SynerexProvider) SetAgentsSupply(tid uint64, time uint32, areaId uint32) {
	setAgentsSupply := &agent.SetAgentsSupply{
		Time:       time,
		AreaId:     areaId,
		AgentType:  0, // Ped
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

// forwardAgentsSupply :　次の時刻のエージェント情報をSupplyする
func (s *SynerexProvider) ForwardAgentsSupply(tid uint64, time uint32, areaId uint32, agentsInfo []*agent.AgentInfo) {
	forwardAgentsSupply := &agent.ForwardAgentsSupply{
		Time:       time,
		AreaId:     areaId,
		AgentType:  0, //Ped
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

// getAgentsSupply :　同じエリアのエージェント情報を提供する
func (s *SynerexProvider) GetAgentsSupply(tid uint64, time uint32, areaId uint32, agentsInfo []*agent.AgentInfo) {
	getAgentsSupply := &agent.GetAgentsSupply{
		Time:       time,
		AreaId:     areaId,
		AgentType:  0, //Ped
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
	s.SpIdList = append(s.SpIdList, id) // my demand list
	s.SpMap[id] = opts                  // my demand options
	mu.Unlock()
}

func (s *SynerexProvider) sendSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts) {
	mu.Lock()
	id := sclient.RegisterSupply(opts)
	s.SpIdList = append(s.SpIdList, id) // my demand list
	s.SpMap[id] = opts                  // my demand options
	mu.Unlock()
}

func (s *SynerexProvider) sendDemand(sclient *sxutil.SMServiceClient, opts *sxutil.DemandOpts) {
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	s.DmIdList = append(s.DmIdList, id) // my demand list
	s.DmMap[id] = opts                  // my demand options
	mu.Unlock()
}

func SubscribeSupply(client *sxutil.SMServiceClient, supplyCallback func(*sxutil.SMServiceClient, *pb.Supply), wg *sync.WaitGroup) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback, wg)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func SubscribeDemand(client *sxutil.SMServiceClient, demandCallback func(*sxutil.SMServiceClient, *pb.Demand), wg *sync.WaitGroup) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback, wg)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}
