package provider

import (
	"log"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/agent"
)

// PedProvider :
type PedProvider struct {
	*SynerexProvider //埋め込み
	Peds             []*agent.Pedestrian
}

// NewSenerexProvider:
func NewPedProvider() *PedProvider {

	provider := &PedProvider{
		SynerexProvider: NewSynerexProvider(),
		Peds:            make([]*agent.Pedestrian, 0),
	}

	return provider
}

func (p *PedProvider) GetAgents(dm *pb.Demand) []*agent.Pedestrian {

	setAgentsDemand := dm.GetArg_SetAgentsDemand()
	agentsInfo := setAgentsDemand.AgentsInfo

	agents := make([]*agent.Pedestrian, 0)
	for _, agent := range agentsInfo {
		ped := agent.NewPedestrian()
		ped.SetGrpcAgent(agent)
		agents = append(agents, ped)
	}
	return agents
}

// GetAgentsRoute : AgentのRouteを取得する関数
func (p *PedProvider) GetAgentsRoute(agents []*agent.Pedestrian, ch chan *pb.Supply) []*agent.Pedestrian {

	// AgentのRoute情報を取得するDemand
	GetAgentsRouteDemand(agentsInfo)
	// Route情報を取得
	sp := <-ch
	log.Println("GET_AGENTS_ROUTE_FINISH")
	getAgentsRouteSupply := sp.GetArg_GetAgentsRouteSupply()
	agentsInfo = getAgentsRouteSupply.AgentsInfo

	return agents
}

// SendAgentsRouteSupply : AgentのRoute取得したspをチャネルに送る関数
func (p *PedProvider) SendAgentsRouteSupply(sp *pb.Supply, ch chan *pb.Supply) {
	ch <- sp
}
