package simulator

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

// ForwardGlobalTime :　GlobalTimeを進める
func (sim *SynerexSimulator) ForwardGlobalTime(forwardTime float64) {
	sim.GlobalTime = sim.GlobalTime + forwardTime
}

// SetGlobalTime :　GlobalTimeをセットする関数
func (sim *SynerexSimulator) SetGlobalTime(time float64) {
	sim.GlobalTime = time
}

// SetTimeStep :　TimeStepをセットする関数
func (sim *SynerexSimulator) SetTimeStep(timeStep float64) {
	sim.TimeStep = timeStep
}
