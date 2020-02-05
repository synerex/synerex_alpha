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

// ClockCommunicator :
type ClockCommunicator struct {
	ProviderId uint64
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
	DownClockCh                        chan *pb.Supply
	DownClockIdList                    []uint64
	SetClockCh                        chan *pb.Supply
	SetClockIdList                    []uint64
	RegistParticipantCh               chan *pb.Supply
}

// NewClockCommunicator:
func NewClockCommunicator(pid uint64) *ClockCommunicator {

	communicator := &ClockCommunicator{
		ProviderId: pid,
		SynerexCommunicator:   communicator.NewSynerexCommunicator(),
		SetClockIdList:        make([]uint64, 0),
		SetClockCh:            make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		GetClockIdList:        make([]uint64, 0),
		GetClockCh:            make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		DownClockIdList:        make([]uint64, 0),
		DownClockCh:            make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		SetAgentsIdList:       make([]uint64, 0),
		ClearAgentsIdList:     make([]uint64, 0),
		ForwardClockIdList:    make([]uint64, 0),
		SetParticipantsIdList: make([]uint64, 0),
		ForwardClockCh:        make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		SetParticipantsCh:     make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		SetAgentsCh:           make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		ClearAgentsCh:         make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
		RegistParticipantCh:         make(chan *pb.Supply, CHANNEL_BUFFER_SIZE),
	}

	return communicator
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (s *ClockCommunicator) RegistClients(client pb.SynerexClient, argJson string) {

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
func (s *ClockCommunicator) CreateWaitIdList() {
	setAgentsIdList := make([]uint64, 0)
	clearAgentsIdList := make([]uint64, 0)
	setParticipantsIdList := make([]uint64, 0)
	forwardClockIdList := make([]uint64, 0)
	getClockIdList := make([]uint64, 0)
	setClockIdList := make([]uint64, 0)
	downClockIdList := make([]uint64, 0)
	for _, participantInfo := range s.Participants {
		providerType := participantInfo.GetProviderType()
		agentChannelId := participantInfo.GetChannelId().GetAgentChannelId()
		clockChannelId := participantInfo.GetChannelId().GetClockChannelId()
		participantChannelId := participantInfo.GetChannelId().GetParticipantChannelId()

		if providerType != participant.ProviderType_SCENARIO && providerType != participant.ProviderType_AREA {
			forwardClockIdList = append(forwardClockIdList, clockChannelId)
			getClockIdList = append(getClockIdList, clockChannelId)
			setClockIdList = append(setClockIdList, clockChannelId)
			downClockIdList = append(downClockIdList, participantChannelId)
		}
		if providerType != participant.ProviderType_SCENARIO && providerType != participant.ProviderType_AREA && providerType != participant.ProviderType_VISUALIZATION {
			setParticipantsIdList = append(setParticipantsIdList, participantChannelId)
			setAgentsIdList = append(setAgentsIdList, agentChannelId)
			clearAgentsIdList = append(clearAgentsIdList, agentChannelId)
		}
	}
	s.SetAgentsIdList = setAgentsIdList
	s.ClearAgentsIdList = clearAgentsIdList
	s.SetParticipantsIdList = setParticipantsIdList
	s.ForwardClockIdList = forwardClockIdList
	s.GetClockIdList = getClockIdList
	s.SetClockIdList = setClockIdList
	s.DownClockIdList = downClockIdList
}

func (s *ClockCommunicator) GetMyParticipant() *participant.Participant {
	participant := &participant.Participant{
		Id: s.ProviderId,
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint64(s.MyClients.ParticipantClient.ClientID),
			AreaChannelId:        uint64(s.MyClients.AreaClient.ClientID),
			AgentChannelId:       uint64(s.MyClients.AgentClient.ClientID),
			ClockChannelId:       uint64(s.MyClients.ClockClient.ClientID),
			RouteChannelId:       uint64(s.MyClients.RouteClient.ClientID),
		},
		ProviderType: participant.ProviderType_CLOCK,
	}
	return participant
}

// WaitSetAgentsResponse : SetAgentsResponseを待機する
func (s *ClockCommunicator) WaitSetAgentsResponse() {
	// channelの初期化
	s.SetAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.SetAgentsIdList, s.SetAgentsCh)
}

// SendToSetAgentsResponse : SetAgentsResponseを送る
func (s *ClockCommunicator) SendToSetAgentsResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.SetAgentsCh)
}

// WaitClearAgentsResponse : ClearAgentsResponseを待機する
func (s *ClockCommunicator) WaitClearAgentsResponse() {
	// channelの初期化
	s.ClearAgentsCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.ClearAgentsIdList, s.ClearAgentsCh)
}

// SendToClearAgentsResponse : ClearAgentsResponseを送る
func (s *ClockCommunicator) SendToClearAgentsResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.ClearAgentsCh)
}

// WaitDownClockResponse : DownClockResponseを待機する
func (s *ClockCommunicator) WaitDownClockResponse() {
	if len(s.DownClockIdList) != 0{
		// channelの初期化
		s.DownClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
		s.Wait(s.DownClockIdList, s.DownClockCh)
	}
}

// SendToDownClockResponse : DownClockResponseを送る
func (s *ClockCommunicator) SendToDownClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.DownClockCh)
}

// WaitGetClockResponse : GetClockResponseを待機する
func (s *ClockCommunicator) WaitGetClockResponse() {
	// channelの初期化
	s.GetClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.GetClockIdList, s.GetClockCh)
}

// SendToGetClockResponse : GetClockResponseを送る
func (s *ClockCommunicator) SendToGetClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.GetClockCh)
}

// WaitSetClockResponse : SetClockResponseを待機する
func (s *ClockCommunicator) WaitSetClockResponse() {
	// channelの初期化
	s.SetClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.SetClockIdList, s.SetClockCh)
}

// SendToSetClockResponse : SetClockResponseを送る
func (s *ClockCommunicator) SendToSetClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.SetClockCh)
}

// WaitSetParticipantsResponse : SetParticipantsResponseを待機する
func (s *ClockCommunicator) WaitSetParticipantsResponse() error{
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
func (s *ClockCommunicator) SendToSetParticipantsResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.SetParticipantsCh)
}

// WaitForwardClockResponse : ForwardClockResponseを待機する
func (s *ClockCommunicator) WaitForwardClockResponse() {
	// channelの初期化
	s.ForwardClockCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	s.Wait(s.ForwardClockIdList, s.ForwardClockCh)
}

// SendToForwardClockResponse : ForwardClockResponseを送る
func (s *ClockCommunicator) SendToForwardClockResponse(sp *pb.Supply) {
	s.SendToWait(sp, s.ForwardClockCh)
}

// WaitRegistParticipantResponse : RegistParticipantResponseを待機する
func (s *ClockCommunicator) WaitRegistParticipantResponse() error{
	// channelの初期化
	s.RegistParticipantCh = make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)

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
	case <- s.RegistParticipantCh:
		return nil
	}

}


// SendToRegistParticipantResponse : RegistParticipantResponseを送る
func (s *ClockCommunicator) SendToRegistParticipantResponse(sp *pb.Supply) {
	s.RegistParticipantCh <- sp
}
