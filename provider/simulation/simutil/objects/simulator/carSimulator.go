package simulator

import (
	"fmt"

	"github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/agent"
)

// SynerexSimulator :
type CarSimulator struct {
	*SynerexSimulator //埋め込み
	Map               *agent.Map
	Agents            []*agent.Car
	AgentType         *agent.Type
}

//Todo: agentType
// NewSenerexSimulator:
func NewCarSimulator(timeStep float64, agentType int64, globalTime float64) *CarSimulator {

	sim := &CarSimulator{
		SynerexSimulator: &SynerexSimulator{
			TimeStep:   timeStep,
			GlobalTime: globalTime,
		},
		AgentType: agent.Type(int64),
		Agents:    make([]*agent.Car, 0),
		Map:       &agent.Map{},
	}

	return sim
}

// SetMap :　Mapを追加する関数
func (sim *CarSimulator) SetMap(_map *agent.Map) {
	sim.Map = _map
}

// AddAgent :　エージェントを追加する関数
func (sim *CarSimulator) AddAgent(agent *agent.Car) {
	sim.Agents = append(sim.Agents, agent)
}

// UpdateAgent :　エージェントを更新する関数
func (sim *CarSimulator) UpdateAgent(a *agent.Car) error {
	for i, agent := range sim.Agents {
		if agent.ID == a.ID {
			sim.Agents[i] = a
			return nil
		}
	}
	// 更新できなかったとき
	return fmt.Errorf("not exist agent by id")
}

// SetAgents :　エージェントを一括で変換する関数
func (sim *CarSimulator) SetAgents(agents []*agent.Car) {
	for _, agentInfo := range agents {
		// 重複エリア内に入っていれば加える
		if agentInfo.IsInArea(sim.Map.Duplicate) {
			sim.AddAgent(agentInfo)
		}
	}
}

// ForwardStep :　次の時刻のエージェントを計算する関数
func (sim *CarSimulator) ForwardStep(sameAreaAgents []*agent.Car) []*agent.Car {

	pureNextAgents := make([]*agent.Car, 0)

	// 干渉なしで目的地へ進む
	simpleRoute := NewSimpleRoute(sim.TimeStep, sim.GlobalTime, sim.Map, sim.Agents, sim.AgentType)
	pureNextAgents = simpleRoute.CalcNextAgentsBySimple()

	return pureNextAgents
}

// UpdateDuplicateAgents :　重複エリアのエージェントを更新する関数
func (sim *CarSimulator) UpdateDuplicateAgents(pureNextAgentsInfo []*agent.Car, neighborAgents []*agent.Car) []*agent.Car {
	nextAgents := pureNextAgentsInfo
	for _, neighborAgent := range neighborAgents {
		//　隣のエージェントが自分のエリアにいてかつ自分のエリアのエージェントと被ってない場合更新
		if len(pureNextAgentsInfo) == 0 {
			//
			if neighborAgent.IsInArea(sim.Map.Duplicate) {
				nextAgents = append(nextAgents, neighborAgent)
			}
		} else {
			isAppendAgent := false
			for _, sameAreaAgent := range pureNextAgentsInfo {
				if neighborAgent.AgentId != sameAreaAgent.AgentId && neighborAgent.IsInArea(sim.Map.Duplicate) {
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

// IsNeighborMap: mapIdが隣接するエリアのものかを判断する関数
func (sim *CarSimulator) IsNeighborMap(mapId uint64) bool {
	neighborMaps := sim.Map.Neighbors

	//fmt.Printf("idlist3 %v\n", neighborMap)
	for _, neighborMap := range neighborMaps {
		if mapId == neighborMap.ID {
			return true
		}
	}
	return false
}

// IsSameMap: mapIdが同じエリアのものかを判断する関数
func (sim *CarSimulator) IsSameMap(mapId uint64) bool {
	if mapId == sim.Map.ID {
		return true
	}
	return false
}

// IsSameAgentType: agentTypeが同じかを判断する関数
func (sim *CarSimulator) IsSameAgentType(agentType *agent.Type) bool {
	if agentType == sim.AgentType {
		return true
	}
	return false
}
