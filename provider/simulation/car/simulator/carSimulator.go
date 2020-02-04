package simulator

//"fmt"

import (
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	car "github.com/synerex/synerex_alpha/provider/simulation/car/agent"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator"
)

// SynerexSimulator :
type CarSimulator struct {
	*simulator.SynerexSimulator //埋め込み
	Area                        *area.Area2
	Agents                      []*agent.Agent
	AgentType                   common.AgentType
}

// NewSenerexSimulator:
func NewCarSimulator(timeStep float64, globalTime float64) *CarSimulator {

	sim := &CarSimulator{
		SynerexSimulator: &simulator.SynerexSimulator{
			TimeStep:   timeStep,
			GlobalTime: globalTime,
		},
		Area:      &area.Area2{},
		AgentType: common.AgentType_CAR,
		Agents:    make([]*agent.Agent, 0),
	}

	return sim
}

// SetArea :　Areaを追加する関数
func (sim *CarSimulator) SetArea(areaInfo *area.Area2) {
	sim.Area = areaInfo
}

// GetArea :　Areaを取得する関数
func (sim *CarSimulator) GetArea() *area.Area2 {
	return sim.Area
}

// AddAgents :　Agentsを追加する関数
func (sim *CarSimulator) AddAgents(agentsInfo []*agent.Agent) {
	newAgents := make([]*agent.Agent, 0)
	for _, agentInfo := range agentsInfo {
		if agentInfo.Type == common.AgentType_CAR {
			carInfo := car.NewCar(agentInfo)
			if carInfo.IsInArea(sim.Area.DuplicateArea){
				newAgents = append(newAgents, agentInfo)
			}
		}
	}
	sim.Agents = append(sim.Agents, newAgents...)
}

// SetAgents :　Agentsをセットする関数
func (sim *CarSimulator) SetAgents(agentsInfo []*agent.Agent) {
	newAgents := make([]*agent.Agent, 0)
	for _, agentInfo := range agentsInfo {
		if agentInfo.Type == common.AgentType_CAR {
			newAgents = append(newAgents, agentInfo)
		}
	}
	sim.Agents = newAgents
}

// ClearAgents :　Agentsを追加する関数
func (sim *CarSimulator) ClearAgents() {
	sim.Agents = make([]*agent.Agent, 0)
}

// GetAgents :　Agentsを取得する関数
func (sim *CarSimulator) GetAgents() []*agent.Agent {
	return sim.Agents
}

// UpdateDuplicateAgents :　重複エリアのエージェントを更新する関数
func (sim *CarSimulator) UpdateDuplicateAgents(pureNextAgentsInfo []*agent.Agent, neighborAgents []*agent.Agent) []*agent.Agent {
	nextAgents := pureNextAgentsInfo
	for _, neighborAgent := range neighborAgents {
		neighborCar := car.NewCar(neighborAgent)
		//　隣のエージェントが自分のエリアにいてかつ自分のエリアのエージェントと被ってない場合更新
		if len(pureNextAgentsInfo) == 0 {
			//
			if neighborCar.IsInArea(sim.Area.DuplicateArea) {
				nextAgents = append(nextAgents, neighborAgent)
			}
		} else {
			isAppendAgent := false
			for _, sameAreaAgent := range pureNextAgentsInfo {
				if neighborAgent.Id != sameAreaAgent.Id && neighborCar.IsInArea(sim.Area.DuplicateArea) {
					isAppendAgent = true
				}
			}
			if isAppendAgent {
				//log.Println("CHANGE_AREA2!!!")
				nextAgents = append(nextAgents, neighborAgent)
			}
		}
	}
	return nextAgents
}

// ForwardStep :　次の時刻のエージェントを計算する関数
func (sim *CarSimulator) ForwardStep(sameAreaAgents []*agent.Agent) []*agent.Agent {
	nextControlAgents := sim.GetAgents()

	// 干渉なしで目的地へ進む
	simpleRoute := NewSimpleRoute(sim.TimeStep, sim.GlobalTime, sim.Area, sim.Agents, sim.AgentType)
	nextControlAgents = simpleRoute.CalcNextAgents()

	return nextControlAgents
}
