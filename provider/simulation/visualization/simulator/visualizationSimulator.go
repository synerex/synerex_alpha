package simulator

//"fmt"

import (
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator"
)

// SynerexSimulator :
type VisualizationSimulator struct {
	*simulator.SynerexSimulator //埋め込み
	Area                        *area.Area2
	Agents                      []*agent.Agent
	AgentType                   common.AgentType
}

// NewSenerexSimulator:
func NewVisualizationSimulator(timeStep float64, globalTime float64) *VisualizationSimulator {

	sim := &VisualizationSimulator{
		SynerexSimulator: &simulator.SynerexSimulator{
			TimeStep:   timeStep,
			GlobalTime: globalTime,
		},
	}

	return sim
}
