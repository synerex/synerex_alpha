package simulator

import (
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator"
)

// SynerexSimulator :
type AreaSimulator struct {
	*simulator.SynerexSimulator //埋め込み
}

// NewSenerexSimulator:
func NewAreaSimulator(timeStep float64, globalTime float64) *AreaSimulator {

	sim := &AreaSimulator{
		SynerexSimulator: &simulator.SynerexSimulator{
			TimeStep:   timeStep,
			GlobalTime: globalTime,
		},
	}

	return sim
}
