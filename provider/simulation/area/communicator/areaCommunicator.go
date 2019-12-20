package communicator

import (
	//"fmt"
	//"log"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/communicator"
	"github.com/synerex/synerex_alpha/sxutil"
)

// AreaCommunicator :
type AreaCommunicator struct {
	*communicator.SynerexCommunicator //埋め込み
}

// NewSenerexCommunicator:
func NewAreaCommunicator() *AreaCommunicator {

	communicator := &AreaCommunicator{
		SynerexCommunicator: communicator.NewSynerexCommunicator(),
	}

	return communicator
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (p *AreaCommunicator) RegistClients(client pb.SynerexClient, argJson string) {

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
		ProviderType: participant.ProviderType_AREA,
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

func (p *AreaCommunicator) GetMyParticipant(areaId uint64) *participant.Participant {
	participant := &participant.Participant{
		ChannelId: &participant.ChannelId{
			ParticipantChannelId: uint64(p.MyClients.ParticipantClient.ClientID),
			AreaChannelId:        uint64(p.MyClients.AreaClient.ClientID),
			AgentChannelId:       uint64(p.MyClients.AgentClient.ClientID),
			ClockChannelId:       uint64(p.MyClients.ClockClient.ClientID),
			RouteChannelId:       uint64(p.MyClients.RouteClient.ClientID),
		},
		ProviderType: participant.ProviderType_AREA,
	}
	return participant
}
