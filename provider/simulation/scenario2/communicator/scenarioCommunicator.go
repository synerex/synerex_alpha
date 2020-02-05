package communicator

import (
	"fmt"
	"time"
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
	ClearAgentsIdList                 []uint64
	ForwardClockIdList                []uint64
	SetParticipantsIdList             []uint64
	ForwardClockCh                    chan *pb.Supply
	SetParticipantsCh                 chan *pb.Supply
	SetAgentsCh                       chan *pb.Supply
	ClearAgentsCh                     chan *pb.Supply
	GetClockCh                        chan *pb.Supply
	GetClockIdList                    []uint64
	DownScenarioCh                        chan *pb.Supply
	DownScenarioIdList                    []uint64
	SetClockCh                        chan *pb.Supply
	SetClockIdList                    []uint64
	StartClockCh                        chan *pb.Supply
	StartClockIdList                    []uint64
	StopClockCh                        chan *pb.Supply
	StopClockIdList                    []uint64
}

// NewScenarioCommunicator:
func NewScenarioCommunicator() *ScenarioCommunicator {

	communicator := &ScenarioCommunicator{
		SynerexCommunicator:   communicator.NewSynerexCommunicator(),
		SetClockIdList:        make([]uint64, 0),
		SetClockCh:            make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		StartClockIdList:        make([]uint64, 0),
		StartClockCh:            make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		StopClockIdList:        make([]uint64, 0),
		StopClockCh:            make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		GetClockIdList:        make([]uint64, 0),
		GetClockCh:            make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		DownScenarioIdList:        make([]uint64, 0),
		DownScenarioCh:            make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		SetAgentsIdList:       make([]uint64, 0),
		ClearAgentsIdList:     make([]uint64, 0),
		ForwardClockIdList:    make([]uint64, 0),
		SetParticipantsIdList: make([]uint64, 0),
		ForwardClockCh:        make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		SetParticipantsCh:     make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		SetAgentsCh:           make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		ClearAgentsCh:         make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
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


// CreateWaitIdList : 同期するためのIdListを作成する関数
func (s *ScenarioCommunicator) CreateWaitIdList() {
	setAgentsIdList := make([]uint64, 0)
	clearAgentsIdList := make([]uint64, 0)
	setParticipantsIdList := make([]uint64, 0)
	forwardClockIdList := make([]uint64, 0)
	getClockIdList := make([]uint64, 0)
	setClockIdList := make([]uint64, 0)
	stopClockIdList := make([]uint64, 0)
	startClockIdList := make([]uint64, 0)
	downScenarioIdList := make([]uint64, 0)
	for _, participantInfo := range s.Participants {
		providerType := participantInfo.GetProviderType()
		agentChannelId := participantInfo.GetChannelId().GetAgentChannelId()
		clockChannelId := participantInfo.GetChannelId().GetClockChannelId()
		participantChannelId := participantInfo.GetChannelId().GetParticipantChannelId()

		if providerType != participant.ProviderType_SCENARIO && providerType != participant.ProviderType_AREA {
			forwardClockIdList = append(forwardClockIdList, clockChannelId)
			getClockIdList = append(getClockIdList, clockChannelId)
			setClockIdList = append(setClockIdList, clockChannelId)
			downScenarioIdList = append(downScenarioIdList, participantChannelId)
		}
		if providerType != participant.ProviderType_SCENARIO && providerType != participant.ProviderType_AREA && providerType != participant.ProviderType_VISUALIZATION {
			setParticipantsIdList = append(setParticipantsIdList, participantChannelId)
			setAgentsIdList = append(setAgentsIdList, agentChannelId)
			clearAgentsIdList = append(clearAgentsIdList, agentChannelId)
		}
		// clockのみ
		if providerType == participant.ProviderType_CLOCK{
			startClockIdList = append(startClockIdList, clockChannelId)
			stopClockIdList = append(stopClockIdList, clockChannelId)
		}
	}
	s.SetAgentsIdList = setAgentsIdList
	s.ClearAgentsIdList = clearAgentsIdList
	s.SetParticipantsIdList = setParticipantsIdList
	s.ForwardClockIdList = forwardClockIdList
	s.GetClockIdList = getClockIdList
	s.SetClockIdList = setClockIdList
	s.StartClockIdList = startClockIdList
	s.StopClockIdList = stopClockIdList
	s.DownScenarioIdList = downScenarioIdList
}


// WaitSendParticipantInfo : SetAgentsResponseを待機する
func (s *ScenarioCommunicator) WaitSendParticipantInfo() {
	// channelの初期化
	s.SetAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.SetAgentsIdList, s.SetAgentsCh)
}

// SendToWaitSendParticipant : WaitSendParticipantを送る
func (s *ScenarioCommunicator) SendToWaitSendParticipant(sp *pb.Supply) {
	s.SendToWait(sp, s.SetAgentsCh)
}



// WaitSetAgentsResponse : SetAgentsResponseを待機する
func (s *ScenarioCommunicator) WaitSetAgentsResponse() {
	// channelの初期化
	s.SetAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.SetAgentsIdList, s.SetAgentsCh)
}

// SendToSetAgentsResponse : SetAgentsResponseを送る
func (s *ScenarioCommunicator) SendToSetAgentsResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.SetAgentsCh)
}

// WaitClearAgentsResponse : ClearAgentsResponseを待機する
func (s *ScenarioCommunicator) WaitClearAgentsResponse() {
	// channelの初期化
	s.ClearAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.ClearAgentsIdList, s.ClearAgentsCh)
}

// SendToClearAgentsResponse : ClearAgentsResponseを送る
func (s *ScenarioCommunicator) SendToClearAgentsResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.ClearAgentsCh)
}

// WaitDownScenarioResponse : DownScenarioResponseを待機する
func (s *ScenarioCommunicator) WaitDownScenarioResponse() {
	if len(s.DownScenarioIdList) != 0{
		// channelの初期化
		s.DownScenarioCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
		s.Wait(s.DownScenarioIdList, s.DownScenarioCh)
	}
}

// SendToDownScenarioResponse : DownScenarioResponseを送る
func (s *ScenarioCommunicator) SendToDownScenarioResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.DownScenarioCh)
}

// WaitGetClockResponse : GetClockResponseを待機する
func (s *ScenarioCommunicator) WaitGetClockResponse() {
	// channelの初期化
	s.GetClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.GetClockIdList, s.GetClockCh)
}

// SendToGetClockResponse : GetClockResponseを送る
func (s *ScenarioCommunicator) SendToGetClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.GetClockCh)
}

// WaitSetClockResponse : SetClockResponseを待機する
func (s *ScenarioCommunicator) WaitSetClockResponse() {
	// channelの初期化
	s.SetClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.SetClockIdList, s.SetClockCh)
}

// SendToSetClockResponse : SetClockResponseを送る
func (s *ScenarioCommunicator) SendToSetClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.SetClockCh)
}

// WaitStartClockResponse : StartClockResponseを待機する
func (s *ScenarioCommunicator) WaitStartClockResponse() {
	// channelの初期化
	s.StartClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.StartClockIdList, s.StartClockCh)
}

// SendToStartClockResponse : StartClockResponseを送る
func (s *ScenarioCommunicator) SendToStartClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.StartClockCh)
}

// WaitStopClockResponse : StopClockResponseを待機する
func (s *ScenarioCommunicator) WaitStopClockResponse() {
	// channelの初期化
	s.StopClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.StopClockIdList, s.StopClockCh)
}

// SendToStopClockResponse : StopClockResponseを送る
func (s *ScenarioCommunicator) SendToStopClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.StopClockCh)
}

// WaitSetParticipantsResponse : SetParticipantsResponseを待機する
func (s *ScenarioCommunicator) WaitSetParticipantsResponse() error{
	// Participantが自分しかいない場合、スキップする
	if len(s.SetParticipantsIdList) == 0{
		return nil
	} else {
		// channelの初期化
		s.SetParticipantsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)

		errch := make(chan error, 1)
		go func(){
			s.Wait(s.SetParticipantsIdList, s.SetParticipantsCh)
			errch <- nil
		}()
		// timeout
		go func(){
			time.Sleep(2*time.Second)
			errch <- fmt.Errorf("timeout occor...\n")
			return
		}()
		select {
		case err := <- errch:
			return err
		}
	}
}

// SendToSetParticipantsResponse : SetParticipantsResponseを送る
func (s *ScenarioCommunicator) SendToSetParticipantsResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.SetParticipantsCh)
}

// WaitForwardClockResponse : ForwardClockResponseを待機する
func (s *ScenarioCommunicator) WaitForwardClockResponse() {
	// channelの初期化
	s.ForwardClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.ForwardClockIdList, s.ForwardClockCh)
}

// SendToForwardClockResponse : ForwardClockResponseを送る
func (s *ScenarioCommunicator) SendToForwardClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.ForwardClockCh)
}
