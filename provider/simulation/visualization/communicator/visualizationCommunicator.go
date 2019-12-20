package communicator

import (
	"fmt"
	"time"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	//"github.com/synerex/synerex_alpha/api/simulation/common"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/communicator"
	"github.com/synerex/synerex_alpha/sxutil"
)

var (
	CHANNEL_BUFFER_SIZE int
)

func init() {
	CHANNEL_BUFFER_SIZE = 10
}

// VisualizationCommunicator :
type VisualizationCommunicator struct {
	*communicator.SynerexCommunicator //埋め込み
	VisualizeAgentsIdList             []uint64
	VisualizeAgentsCh                 chan *pb.Supply
	RegistParticipantCh               chan *pb.Supply
	GetClockIdList                    []uint64
	GetClockCh                        chan *pb.Supply
	DeleteParticipantIdList           []uint64
	DeleteParticipantCh               chan *pb.Supply
}

// NewVisualizationCommunicator:
func NewVisualizationCommunicator() *VisualizationCommunicator {

	communicator := &VisualizationCommunicator{
		SynerexCommunicator:     communicator.NewSynerexCommunicator(),
		VisualizeAgentsIdList:   make([]uint64, 0),
		VisualizeAgentsCh:       make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		RegistParticipantCh:     make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		DeleteParticipantIdList: make([]uint64, 0),
		DeleteParticipantCh:     make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		GetClockIdList:          make([]uint64, 0),
		GetClockCh:              make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
	}

	return communicator
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (p *VisualizationCommunicator) RegistClients(client pb.SynerexClient, argJson string) {

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
		ProviderType: participant.ProviderType_VISUALIZATION,
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
func (p *VisualizationCommunicator) CreateWaitIdList() {
	getClockIdList := make([]uint64, 0)
	deleteParticipantIdList := make([]uint64, 0)
	visualizeAgentsIdList := make([]uint64, 0)
	for _, participantInfo := range p.Participants {
		providerType := participantInfo.GetProviderType()
		//agentType := participantInfo.GetAgentType()
		clockChannelId := participantInfo.GetChannelId().GetClockChannelId()
		agentChannelId := participantInfo.GetChannelId().GetAgentChannelId()
		participantChannelId := participantInfo.GetChannelId().GetParticipantChannelId()

		if providerType == participant.ProviderType_SCENARIO {
			getClockIdList = append(getClockIdList, clockChannelId)
			deleteParticipantIdList = append(deleteParticipantIdList, participantChannelId)
		}
		if providerType != participant.ProviderType_SCENARIO { // agentTypeにNONEを作るべき
			visualizeAgentsIdList = append(visualizeAgentsIdList, agentChannelId)
		}
	}
	p.DeleteParticipantIdList = deleteParticipantIdList
	p.GetClockIdList = getClockIdList
	p.VisualizeAgentsIdList = visualizeAgentsIdList
}

func (p *VisualizationCommunicator) GetMyParticipant() *participant.Participant {
	participant := &participant.Participant{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint64(p.MyClients.ParticipantClient.ClientID),
			AreaChannelId:        uint64(p.MyClients.AreaClient.ClientID),
			AgentChannelId:       uint64(p.MyClients.AgentClient.ClientID),
			ClockChannelId:       uint64(p.MyClients.ClockClient.ClientID),
			RouteChannelId:       uint64(p.MyClients.RouteClient.ClientID),
		},
		ProviderType: participant.ProviderType_VISUALIZATION,
	}
	return participant
}

// WaitRegistParticipantResponse : RegistParticipantResponseを待機する
func (p *VisualizationCommunicator) WaitRegistParticipantResponse() error{
	// channelの初期化
	p.RegistParticipantCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)

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
func (p *VisualizationCommunicator) SendToRegistParticipantResponse(sp *pb.Supply) {
	p.RegistParticipantCh <- sp
}

// WaitDeleteParticipantResponse : DeleteParticipantResponseを待機する
func (p *VisualizationCommunicator) WaitDeleteParticipantResponse() {
	// channelの初期化
	p.DeleteParticipantCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	// spの待機
	p.Wait(p.DeleteParticipantIdList, p.DeleteParticipantCh)
}

// SendToDeleteParticipantResponse : DeleteParticipantResponseを送る
func (p *VisualizationCommunicator) SendToDeleteParticipantResponse(sp *pb.Supply) {
	p.SendToWait(sp, p.DeleteParticipantCh)
}

// WaitGetClockResponse : GetClockResponseを待機する
func (p *VisualizationCommunicator) WaitGetClockResponse() *clock.Clock {
	// channelの初期化
	p.GetClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
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
func (p *VisualizationCommunicator) SendToGetClockResponse(sp *pb.Supply) {
	p.SendToWait(sp, p.GetClockCh)
}
