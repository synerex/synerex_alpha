package simutil

import (
	"fmt"

	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
)

// SynerexSimulator :
type SynerexSimulator struct {
	TimeStep   float64
	Agents     []*agent.AgentInfo
	AgentType  int32
	Area       *area.AreaInfo
	GlobalTime float64
}

//Todo: agentType
// NewSenerexSimulator:
func NewSynerexSimulator(timeStep float64, agentType uint64, globalTime float64) *SynerexSimulator {

	sim := &SynerexSimulator{
		TimeStep:   timeStep,
		Agents:     make([]*agent.AgentInfo, 0),
		AgentType:  0,
		Area:       &area.AreaInfo{},
		GlobalTime: globalTime,
	}

	return sim
}

// AddAgent :　エージェントを追加する関数
func (sim *SynerexSimulator) AddAgent(agent *agent.AgentInfo) {
	sim.Agents = append(sim.Agents, agent)
}

// UpdateAgent :　エージェントを更新する関数
func (sim *SynerexSimulator) UpdateAgent(a *agent.AgentInfo) error {
	for i, agent := range sim.Agents {
		if agent.AgentId == a.AgentId {
			sim.Agents[i] = a
			return nil
		}
	}
	// 更新できなかったとき
	return fmt.Errorf("not exist agent by id")
}

// SetAgents :　エージェントを一括で変換する関数
func (sim *SynerexSimulator) SetAgents(agents []*agent.AgentInfo) {
	sim.Agents = agents
}

// ForwardStep :　次の時刻のエージェントを計算する関数
func (sim *SynerexSimulator) ForwardStep() ([]*agent.AgentInfo, error) {
	nextAgents := make([]*agent.AgentInfo, 0)

	// 同じエリアの次の時刻のエージェントを計算をする

	// forwardAgentsSupply

	// 重複エリアを更新する

	return nextAgents, nil

}

// UpdateDuplicateAgents :　重複エリアのエージェントを更新する関数
func (sim *SynerexSimulator) UpdateDuplicateAgents() {

}
