package simulator


import (
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	ped "github.com/synerex/synerex_alpha/provider/simulation/pedestrian/agent"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator"
	"github.com/paulmach/orb/geojson"
)

var (
	geoInfo                      *geojson.FeatureCollection
)

// SynerexSimulator :
type PedSimulator struct {
	*simulator.SynerexSimulator //埋め込み
	Area                        *area.Area
	AgentType                   common.AgentType
	Agents                      []*agent.Agent
}

// NewSenerexSimulator:
func NewPedSimulator(timeStep float64, globalTime float64) *PedSimulator {

	sim := &PedSimulator{
		SynerexSimulator: &simulator.SynerexSimulator{
			TimeStep:   timeStep,
			GlobalTime: globalTime,
		},
		Area:      &area.Area{},
		AgentType: common.AgentType_PEDESTRIAN,
		Agents:    make([]*agent.Agent, 0),
	}

	return sim
}

// SetObstacles :　Obstaclesを追加する関数
func (sim *PedSimulator) SetGeoInfo(_geoInfo *geojson.FeatureCollection) {
	geoInfo = _geoInfo
}

// SetArea :　Areaを追加する関数
func (sim *PedSimulator) SetArea(areaInfo *area.Area) {
	sim.Area = areaInfo
}

// GetArea :　Areaを取得する関数
func (sim *PedSimulator) GetArea() *area.Area {
	return sim.Area
}

// AddAgents :　Agentsを追加する関数
func (sim *PedSimulator) AddAgents(agentsInfo []*agent.Agent) {
	newAgents := make([]*agent.Agent, 0)
	for _, agentInfo := range agentsInfo {
		if agentInfo.Type == common.AgentType_PEDESTRIAN {
			pedInfo := ped.NewPedestrian(agentInfo)
			if pedInfo.IsInArea(sim.Area.DuplicateArea){
				newAgents = append(newAgents, agentInfo)
			}
		}
	}
	sim.Agents = append(sim.Agents, newAgents...)
}

// SetAgents :　Agentsをセットする関数
func (sim *PedSimulator) SetAgents(agentsInfo []*agent.Agent) {
	newAgents := make([]*agent.Agent, 0)
	for _, agentInfo := range agentsInfo {
		if agentInfo.Type == common.AgentType_PEDESTRIAN {
			newAgents = append(newAgents, agentInfo)
		}
	}
	sim.Agents = newAgents
}

// ClearAgents :　Agentsを追加する関数
func (sim *PedSimulator) ClearAgents() {
	sim.Agents = make([]*agent.Agent, 0)
}

// GetAgents :　Agentsを取得する関数
func (sim *PedSimulator) GetAgents() []*agent.Agent {
	return sim.Agents
}

// UpdateDuplicateAgents :　重複エリアのエージェントを更新する関数
func (sim *PedSimulator) UpdateDuplicateAgents(nextControlAgents []*agent.Agent, neighborAgents []*agent.Agent) []*agent.Agent {
	nextAgents := nextControlAgents
	for _, neighborAgent := range neighborAgents {
		neighborPed := ped.NewPedestrian(neighborAgent)
		//　隣のエージェントが自分のエリアにいてかつ自分のエリアのエージェントと被ってない場合更新
		if len(nextControlAgents) == 0 {
			//
			if neighborPed.IsInArea(sim.Area.DuplicateArea) {
				nextAgents = append(nextAgents, neighborAgent)
			}
		} else {
			isAppendAgent := false
			for _, sameAreaAgent := range nextControlAgents {
				// 自分の管理しているエージェントではなく管理エリアに入っていた場合更新する
				if neighborAgent.Id != sameAreaAgent.Id && neighborPed.IsInArea(sim.Area.ControlArea) {
					isAppendAgent = true
				}
			}
			if isAppendAgent {
				nextAgents = append(nextAgents, neighborAgent)
			}
		}
	}
	return nextAgents
}

// ForwardStep :　次の時刻のエージェントを計算する関数
func (sim *PedSimulator) ForwardStep(sameAreaAgents []*agent.Agent) []*agent.Agent {
	IsRVO2 := true
	nextControlAgents := sim.GetAgents()

	if IsRVO2 {
		// RVO2
		rvo2route := NewRVO2Route(sim.TimeStep, sim.GlobalTime, sim.Area, sim.Agents, sim.AgentType)
		// Agent計算
		nextControlAgents = rvo2route.CalcNextAgents()

	} else {
		// 干渉なしで目的地へ進む
		//simpleRoute := NewSimpleRoute(sim.TimeStep, sim.GlobalTime, sim.Map, sim.Agents, sim.AgentType)
		//nextControlAgents = simpleRoute.CalcNextAgentsBySimple()

	}
	return nextControlAgents
}
