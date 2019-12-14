package simulator

import (
	"github.com/synerex/synerex_alpha/api/simulation/clock"
)

var (
	IsRVO2 bool
)

func init() {
	IsRVO2 = true
}

// SynerexSimulator :
type SynerexSimulator struct {
	TimeStep   float64
	GlobalTime float64
}

// NewSenerexSimulator:
func NewSynerexSimulator(timeStep float64, globalTime float64) *SynerexSimulator {

	sim := &SynerexSimulator{
		TimeStep:   timeStep,
		GlobalTime: globalTime,
	}

	return sim
}

// GetClock :　Clock情報を取得する
func (sim *SynerexSimulator) GetClock() *clock.Clock {
	clock := &clock.Clock{
		GlobalTime: sim.GlobalTime,
		TimeStep:   sim.TimeStep,
	}
	return clock
}

// ForwardGlobalTime :　GlobalTimeを進める
func (sim *SynerexSimulator) ForwardGlobalTime() {
	sim.GlobalTime = sim.GlobalTime + sim.TimeStep
}

// BackGlobalTime :　GlobalTimeを戻す
func (sim *SynerexSimulator) BackGlobalTime() {
	sim.GlobalTime = sim.GlobalTime + sim.TimeStep
}

// SetGlobalTime :　GlobalTimeをセットする関数
func (sim *SynerexSimulator) SetGlobalTime(time float64) {
	sim.GlobalTime = time
}

// SetTimeStep :　TimeStepをセットする関数
func (sim *SynerexSimulator) SetTimeStep(timeStep float64) {
	sim.TimeStep = timeStep
}
