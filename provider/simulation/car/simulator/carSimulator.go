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
	Area                        *area.Area
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
		Area:      &area.Area{},
		AgentType: common.AgentType_CAR,
		Agents:    make([]*agent.Agent, 0),
	}

	return sim
}

// SetArea :　Areaを追加する関数
func (sim *CarSimulator) SetArea(areaInfo *area.Area) {
	sim.Area = areaInfo
}

// SetAgents :　Agentsを追加する関数
func (sim *CarSimulator) SetAgents(agentsInfo []*agent.Agent) {
	sim.Agents = agentsInfo
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
	IsRVO2 := true
	pureNextAgents := make([]*agent.Agent, 0)

	if IsRVO2 {
		// RVO2
		//rvo2route := NewRVO2Route(sim.TimeStep, sim.GlobalTime, sim.Map, sim.Agents, sim.AgentType)
		//pureNextAgents = rvo2route.CalcNextAgentsByRVO()
	} else {
		// 干渉なしで目的地へ進む
		//simpleRoute := NewSimpleRoute(sim.TimeStep, sim.GlobalTime, sim.Map, sim.Agents, sim.AgentType)
		//pureNextAgents = simpleRoute.CalcNextAgentsBySimple()

	}
	return pureNextAgents
}
