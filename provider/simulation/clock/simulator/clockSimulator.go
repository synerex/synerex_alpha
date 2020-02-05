package simulator

import (
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator"
)

// SynerexSimulator :
type ClockSimulator struct {
	*simulator.SynerexSimulator //埋め込み
}

// NewSenerexSimulator:
func NewClockSimulator(timeStep float64, globalTime float64) *ClockSimulator {

	sim := &ClockSimulator{
		simulator.NewSynerexSimulator(timeStep, globalTime),
	}

	return sim
}

// ForwardStep :
func (sim *ClockSimulator) ForwardStep() {
	sim.GlobalTime = sim.GlobalTime + sim.TimeStep
}
