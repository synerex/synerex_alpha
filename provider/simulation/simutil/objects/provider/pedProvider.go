package provider

import (
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/agent"
)

// PedProvider :
type PedProvider struct {
	*SynerexProvider //埋め込み
	Peds             []*agent.Pedestrian
}

// NewSenerexProvider:
func NewPedProvider(timeStep float64, agentType int64, globalTime float64) *PedProvider {

	provider := &PedProvider{
		SynerexProvider: NewSynerexProvider(),
		Peds:            make([]*agent.Pedestrian, 0),
	}

	return provider
}
