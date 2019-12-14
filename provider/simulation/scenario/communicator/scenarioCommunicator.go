package communicator

import (
	pb "github.com/synerex/synerex_alpha/api"
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

// ScenarioCommunicator :
type ScenarioCommunicator struct {
	*communicator.SynerexCommunicator //埋め込み
	SetAgentsIdList                   []uint64
	ForwardClockIdList                []uint64
	SetParticipantsIdList             []uint64
	ForwardClockCh                    chan *pb.Supply
	SetParticipantsCh                 chan *pb.Supply
	SetAgentsCh                       chan *pb.Supply
}

// NewScenarioCommunicator:
func NewScenarioCommunicator() *ScenarioCommunicator {

	communicator := &ScenarioCommunicator{
		SynerexCommunicator:   communicator.NewSynerexCommunicator(),
		SetAgentsIdList:       make([]uint64, 0),
		ForwardClockIdList:    make([]uint64, 0),
		SetParticipantsIdList: make([]uint64, 0),
		ForwardClockCh:        make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		SetParticipantsCh:     make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		SetAgentsCh:           make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
	}

	return communicator
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (s *ScenarioCommunicator) RegistClients(client pb.SynerexClient, argJson string) {

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
		ProviderType: participant.ProviderType_SCENARIO,
	}

	clients := &communicator.Clients{
		AgentClient:       agentClient,
		ClockClient:       clockClient,
		AreaClient:        areaClient,
		ParticipantClient: participantClient,
		RouteClient:       routeClient,
	}

	s.AddParticipant(participant)
	s.SetClients(clients)

}

func (s *ScenarioCommunicator) SetSetAgentsIdList(setAgentsIdList []uint64) {
	s.SetAgentsIdList = setAgentsIdList
}

func (s *ScenarioCommunicator) GetSetAgentsIdList() []uint64 {
	return s.SetAgentsIdList
}

func (s *ScenarioCommunicator) SetForwardClockIdList(forwardClockIdList []uint64) {
	s.ForwardClockIdList = forwardClockIdList
}

func (s *ScenarioCommunicator) GetForwardClockIdList() []uint64 {
	return s.ForwardClockIdList
}

func (s *ScenarioCommunicator) SetSetParticipantsIdList(setParticipantsIdList []uint64) {
	s.SetParticipantsIdList = setParticipantsIdList
}

func (s *ScenarioCommunicator) GetSetParticipantsIdList() []uint64 {
	return s.SetParticipantsIdList
}

func (s *ScenarioCommunicator) SetSetAgentsCh(setAgentsCh chan *pb.Supply) {
	s.SetAgentsCh = setAgentsCh
}

func (s *ScenarioCommunicator) GetSetAgentsCh() chan *pb.Supply {
	return s.SetAgentsCh
}

func (s *ScenarioCommunicator) SetForwardClockCh(forwardClockCh chan *pb.Supply) {
	s.ForwardClockCh = forwardClockCh
}

func (s *ScenarioCommunicator) GetForwardClockCh() chan *pb.Supply {
	return s.ForwardClockCh
}

func (s *ScenarioCommunicator) SetSetParticipantsCh(setParticipantsCh chan *pb.Supply) {
	s.SetParticipantsCh = setParticipantsCh
}

func (s *ScenarioCommunicator) GetSetParticipantsCh() chan *pb.Supply {
	return s.SetParticipantsCh
}

// CreateWaitIdList : 同期するためのIdListを作成する関数
func (s *ScenarioCommunicator) CreateWaitIdList() {
	setAgentsIdList := make([]uint64, 0)
	setParticipantsIdList := make([]uint64, 0)
	forwardClockIdList := make([]uint64, 0)
	for _, participantInfo := range s.Participants {
		providerType := participantInfo.GetProviderType()
		agentChannelId := participantInfo.GetChannelId().GetAgentChannelId()
		clockChannelId := participantInfo.GetChannelId().GetClockChannelId()
		participantChannelId := participantInfo.GetChannelId().GetParticipantChannelId()

		if providerType != participant.ProviderType_SCENARIO && providerType != participant.ProviderType_AREA {
			forwardClockIdList = append(forwardClockIdList, clockChannelId)
		}
		if providerType != participant.ProviderType_SCENARIO && providerType != participant.ProviderType_AREA && providerType != participant.ProviderType_VISUALIZATION {
			setParticipantsIdList = append(setParticipantsIdList, participantChannelId)
			setAgentsIdList = append(setAgentsIdList, agentChannelId)
		}
	}
	s.SetSetAgentsIdList(setAgentsIdList)
	s.SetSetParticipantsIdList(setParticipantsIdList)
	s.SetForwardClockIdList(forwardClockIdList)
}

// WaitSetAgentsResponse : SetAgentsResponseを待機する
func (s *ScenarioCommunicator) WaitSetAgentsResponse() {
	// channelの初期化
	s.SetSetAgentsCh(make(chan *pb.Supply, CHANNEL_BUFFER_SIZE))
	s.Wait(s.GetSetAgentsIdList(), s.GetSetAgentsCh())
}

// SendToSetAgentsResponse : SetAgentsResponseを送る
func (s *ScenarioCommunicator) SendToSetAgentsResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.GetSetAgentsCh())
}

// WaitSetParticipantsResponse : SetParticipantsResponseを待機する
func (s *ScenarioCommunicator) WaitSetParticipantsResponse() {
	// channelの初期化
	s.SetSetParticipantsCh(make(chan *pb.Supply, CHANNEL_BUFFER_SIZE))
	s.Wait(s.GetSetParticipantsIdList(), s.GetSetParticipantsCh())
}

// SendToSetParticipantsResponse : SetParticipantsResponseを送る
func (s *ScenarioCommunicator) SendToSetParticipantsResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.GetSetParticipantsCh())
}

// WaitForwardClockResponse : ForwardClockResponseを待機する
func (s *ScenarioCommunicator) WaitForwardClockResponse() {
	// channelの初期化
	s.SetForwardClockCh(make(chan *pb.Supply, CHANNEL_BUFFER_SIZE))
	s.Wait(s.GetForwardClockIdList(), s.GetForwardClockCh())
}

// SendToForwardClockResponse : ForwardClockResponseを送る
func (s *ScenarioCommunicator) SendToForwardClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.GetForwardClockCh())
}
