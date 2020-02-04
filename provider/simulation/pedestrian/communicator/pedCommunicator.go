package communicator

import (
	"fmt"
	"time"
	"log"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/communicator"
	"github.com/synerex/synerex_alpha/sxutil"
)

var (
	CHANNEL_BUFFER_SIZE int
	CHANNEL_BUFFER_SIZE2 int
)

func init() {
	CHANNEL_BUFFER_SIZE = 10
	CHANNEL_BUFFER_SIZE2 = 30
}

// PedCommunicator :
type PedCommunicator struct {
	*communicator.SynerexCommunicator //埋め込み
	GetAreaCh                         chan *pb.Supply
	RegistParticipantCh               chan *pb.Supply
	GetClockIdList                    []uint64
	GetClockCh                        chan *pb.Supply
	DeleteParticipantIdList           []uint64
	DeleteParticipantCh               chan *pb.Supply
	GetSameAreaAgentsIdList           []uint64
	GetSameAreaAgentsCh               chan *pb.Supply
	GetNeighborAreaAgentsIdList       []uint64
	GetNeighborAreaAgentsCh           chan *pb.Supply
}

// NewSenerexCommunicator:
func NewPedCommunicator() *PedCommunicator {

	communicator := &PedCommunicator{
		SynerexCommunicator:         communicator.NewSynerexCommunicator(),
		GetAreaCh:                   make(chan *pb.Supply, CHANNEL_BUFFER_SIZE2),
		DeleteParticipantIdList:     make([]uint64, 0),
		DeleteParticipantCh:         make(chan *pb.Supply, CHANNEL_BUFFER_SIZE2),
		GetClockIdList:              make([]uint64, 0),
		GetClockCh:                  make(chan *pb.Supply, CHANNEL_BUFFER_SIZE2),
		GetSameAreaAgentsIdList:     make([]uint64, 0),
		GetSameAreaAgentsCh:         make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		GetNeighborAreaAgentsIdList: make([]uint64, 0),
		GetNeighborAreaAgentsCh:     make(chan *pb.Supply, CHANNEL_BUFFER_SIZE2),
		RegistParticipantCh:         make(chan *pb.Supply, CHANNEL_BUFFER_SIZE2),
	}

	return communicator
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (p *PedCommunicator) RegistClients(client pb.SynerexClient, argJson string) {

	agentClient := sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE, argJson)
	clockClient := sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE, argJson)
	areaClient := sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE, argJson)
	participantClient := sxutil.NewSMServiceClient(client, pb.ChannelType_PARTICIPANT_SERVICE, argJson)
	routeClient := sxutil.NewSMServiceClient(client, pb.ChannelType_ROUTE_SERVICE, argJson)

	// Participantに追加
	participant := &participant.Participant{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint64(participantClient.ClientID),
			AgentChannelId:       uint64(agentClient.ClientID),
			AreaChannelId:        uint64(areaClient.ClientID),
			ClockChannelId:       uint64(clockClient.ClientID),
			RouteChannelId:       uint64(routeClient.ClientID),
		},
		ProviderType: participant.ProviderType_PEDESTRIAN,
	}

	clients := &communicator.Clients{
		AgentClient:       agentClient,
		ClockClient:       clockClient,
		AreaClient:        areaClient,
		ParticipantClient: participantClient,
		RouteClient:       routeClient,
	}

	p.AddParticipant(participant)
	p.SetClients(clients)

}

// CreateWaitIdList : 同期するためのIdListを作成する関数
func (p *PedCommunicator) CreateWaitIdList(myAgentType common.AgentType, myAreaId uint64, neighborAreaIds []uint64) {
	getClockIdList := make([]uint64, 0)
	deleteParticipantIdList := make([]uint64, 0)
	getSameAreaAgentsIdList := make([]uint64, 0)
	getNeighborAreaAgentsIdList := make([]uint64, 0)
	for _, participantInfo := range p.Participants {
		providerType := participantInfo.GetProviderType()
		//agentType := participantInfo.GetAgentType()
		areaId := participantInfo.GetAreaId()
		clockChannelId := participantInfo.GetChannelId().GetClockChannelId()
		agentChannelId := participantInfo.GetChannelId().GetAgentChannelId()
		participantChannelId := participantInfo.GetChannelId().GetParticipantChannelId()

		if providerType == participant.ProviderType_SCENARIO {
			getClockIdList = append(getClockIdList, clockChannelId)
			deleteParticipantIdList = append(deleteParticipantIdList, participantChannelId)
		}
		if myAreaId == areaId && agentChannelId != uint64(p.MyClients.AgentClient.ClientID) {
			getSameAreaAgentsIdList = append(getSameAreaAgentsIdList, agentChannelId)
		}
		for _, neighborAreaId := range neighborAreaIds {
			if neighborAreaId == areaId{
				getNeighborAreaAgentsIdList = append(getNeighborAreaAgentsIdList, agentChannelId)
			}
		}
	}
	p.DeleteParticipantIdList = deleteParticipantIdList
	p.GetClockIdList = getClockIdList
	p.GetSameAreaAgentsIdList = getSameAreaAgentsIdList
	p.GetNeighborAreaAgentsIdList = getNeighborAreaAgentsIdList
}

func (p *PedCommunicator) GetMyParticipant(areaId uint64) *participant.Participant {
	participant := &participant.Participant{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint64(p.MyClients.ParticipantClient.ClientID),
			AreaChannelId:        uint64(p.MyClients.AreaClient.ClientID),
			AgentChannelId:       uint64(p.MyClients.AgentClient.ClientID),
			ClockChannelId:       uint64(p.MyClients.ClockClient.ClientID),
			RouteChannelId:       uint64(p.MyClients.RouteClient.ClientID),
		},
		ProviderType: participant.ProviderType_PEDESTRIAN,
		AreaId:       areaId,
		AgentType:    common.AgentType_PEDESTRIAN,
	}
	return participant
}

// WaitGetAreaResponse : GetAreaResponseを待機する
func (p *PedCommunicator) WaitGetAreaResponse() (*area.Area2, error) {
	// channelの初期化
	p.GetAreaCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE2)

	errch := make(chan error, 1)
	// timeout
	go func(){
		time.Sleep(2*time.Second)
		errch <- fmt.Errorf("timeout occor...\n area-provider closed ?\n You don't have to restart this provider. \n Please start area-provider.")
		return
	}()
	select {
	case err := <- errch:
		return nil, err
	case sp := <- p.GetAreaCh:
		areaInfo := sp.GetSimSupply().GetGetAreaResponse().GetArea()
		return areaInfo, nil
	}
}

// SendToGetAreaResponse : GetAreaResponseを送る
func (p *PedCommunicator) SendToGetAreaResponse(sp *pb.Supply) {
	p.GetAreaCh <- sp
}

// WaitRegistParticipantResponse : RegistParticipantResponseを待機する
func (p *PedCommunicator) WaitRegistParticipantResponse() error{


	// channelの初期化
	p.RegistParticipantCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE2)

	errch := make(chan error, 1)

	// timeout
	go func(){
		time.Sleep(2*time.Second)
		errch <- fmt.Errorf("timeout occor...\n scenario-provider closed ?\n You don't have to restart this provider. \n Please start scenario-provider.")
		return
	}()
	select {
	case err := <- errch:
		return err
	case <- p.RegistParticipantCh:
		return nil
	}

}


// SendToRegistParticipantResponse : RegistParticipantResponseを送る
func (p *PedCommunicator) SendToRegistParticipantResponse(sp *pb.Supply) {
	p.RegistParticipantCh <- sp
}

// WaitDeleteParticipantResponse : DeleteParticipantResponseを待機する
func (p *PedCommunicator) WaitDeleteParticipantResponse() {
	// channelの初期化
	p.DeleteParticipantCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	// spの待機
	p.Wait(p.DeleteParticipantIdList, p.DeleteParticipantCh)
}

// SendToDeleteParticipantResponse : DeleteParticipantResponseを送る
func (p *PedCommunicator) SendToDeleteParticipantResponse(sp *pb.Supply) {
	p.SendToWait(sp, p.DeleteParticipantCh)
}

// WaitGetClockResponse : GetClockResponseを待機する
func (p *PedCommunicator) WaitGetClockResponse() *clock.Clock {
	// channelの初期化
	p.GetClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE2)
	// spの待機
	spMap := p.Wait(p.GetClockIdList, p.GetClockCh)
	// ClockInfoを取得
	var clockInfo *clock.Clock
	for _, sp := range spMap {
		clockInfo = sp.GetSimSupply().GetGetClockResponse().GetClock()
	}

	return clockInfo
}

// SendToGetClockResponse : GetClockResponseを送る
func (p *PedCommunicator) SendToGetClockResponse(sp *pb.Supply) {
	p.SendToWait(sp, p.GetClockCh)
}

// WaitGetSameAreaAgentsResponse : GetSameAreaAgentsResponseを待機する
func (p *PedCommunicator) WaitGetSameAreaAgentsResponse() []*agent.Agent {
	sameAgents := make([]*agent.Agent, 0)
	// 同じエリアが他にない場合、スキップ
	if len(p.GetSameAreaAgentsIdList) != 0{
		spMap := p.Wait(p.GetSameAreaAgentsIdList, p.GetSameAreaAgentsCh)
		// Agentsを取得
	
		for _, sp := range spMap {
			agents := sp.GetSimSupply().GetGetSameAreaAgentsResponse().GetAgents()
			sameAgents = append(sameAgents, agents...)
		}
	
	}

	// channelの初期化
	p.GetSameAreaAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	return sameAgents
}

// SendToGetSameAreaAgentsResponse : GetSameAreaAgentsResponseを送る
func (p *PedCommunicator) SendToGetSameAreaAgentsResponse(sp *pb.Supply) {
	p.SendToWait(sp, p.GetSameAreaAgentsCh)
}

// WaitGetNeighborAreaAgentsResponse : GetNeighborAreaAgentsResponseを待機する
func (p *PedCommunicator) WaitGetNeighborAreaAgentsResponse() []*agent.Agent {
	neighborAgents := make([]*agent.Agent, 0)
	
	// 隣接エリアが他にない場合、スキップ
	if len(p.GetNeighborAreaAgentsIdList) != 0 {
		//p.GetNeighborAreaAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
		log.Printf("neighborAgents: %v\n", p.GetNeighborAreaAgentsIdList)
		spMap := p.Wait(p.GetNeighborAreaAgentsIdList, p.GetNeighborAreaAgentsCh)
		// Agentsを取得
		for _, sp := range spMap {
			agents := sp.GetSimSupply().GetGetNeighborAreaAgentsResponse().GetAgents()
			neighborAgents = append(neighborAgents, agents...)
		}

		//close(p.GetNeighborAreaAgentsCh)
		//p.GetNeighborAreaAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
		// channelの初期化:　チャネルにすでに情報が入っているため最後に初期化する
		//p.GetNeighborAreaAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	}
	

	return neighborAgents
}

// SendToGetNeighborAreaAgentsResponse : GetNeighborAreaAgentsResponseを送る
func (p *PedCommunicator) SendToGetNeighborAreaAgentsResponse(sp *pb.Supply) {
	p.SendToWait(sp, p.GetNeighborAreaAgentsCh)
}

func (p *PedCommunicator) InitChannel(){
	p.GetNeighborAreaAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE2)
}
