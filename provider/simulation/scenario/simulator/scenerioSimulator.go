package simulator

import (
	"time"

	"github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator"
)

// SynerexSimulator :
type ScenarioSimulator struct {
	*simulator.SynerexSimulator //埋め込み
}

// NewSenerexSimulator:
func NewScenarioSimulator(timeStep float64, globalTime float64) *ScenarioSimulator {

	sim := &ScenarioSimulator{
		simulator.NewSynerexSimulator(timeStep, globalTime),
	}

	return sim
}

// ForwardStep :
func (sim *ScenarioSimulator) ForwardStep() {
	sim.GlobalTime = sim.GlobalTime + sim.TimeStep
	// 待機
	time.Sleep(time.Duration(sim.TimeStep) * time.Second)
}
