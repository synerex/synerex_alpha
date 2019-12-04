package simulator

import (
	"fmt"

	"github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/agent"
)

// SynerexSimulator :
type AreaSimulator struct {
	*SynerexSimulator //埋め込み
	Map               *agent.Map
	Signals           []*agent.Signals
}

// NewSenerexSimulator:
func NewAreaSimulator(timeStep float64, globalTime float64) *AreaSimulator {

	sim := &AreaSimulator{
		SynerexSimulator: &SynerexSimulator{
			TimeStep:   timeStep,
			GlobalTime: globalTime,
		},
		Map:     *agent.Map{},
		Signals: make([]*agent.Signal, 0),
	}

	return sim
}

// SetMap :　Mapを追加する関数
func (sim *AreaSimulator) SetMap(_map *agent.Map) {
	sim.Map = _map
}

// AddSignal :　エージェントを追加する関数
func (sim *AreaSimulator) AddSignal(signal *agent.Signal) {
	sim.Signals = append(sim.Signals, signal)
}

// UpdateSignal :　エージェントを更新する関数
func (sim *AreaSimulator) UpdateSignal(a *agent.Signal) error {
	for i, signal := range sim.Signals {
		if signal.ID == a.ID {
			sim.Signals[i] = a
			return nil
		}
	}
	// 更新できなかったとき
	return fmt.Errorf("not exist agent by id")
}

// SetSignals :　エージェントを一括で変換する関数
func (sim *AreaSimulator) SetSignals(signals []*agent.Signal) {
	sim.Signals = signals
}
