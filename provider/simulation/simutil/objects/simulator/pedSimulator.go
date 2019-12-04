package simulator

import (
	"fmt"

	"github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/agent"
)

// SynerexSimulator :
type PedSimulator struct {
	*SynerexSimulator //埋め込み
	Map               *agent.Map
	Agents            []*agent.Pedestrian
	AgentType         *agent.Type
}

// NewSenerexSimulator:
func NewPedSimulator(timeStep float64, agentType int64, globalTime float64) *PedSimulator {

	sim := &PedSimulator{
		SynerexSimulator: &SynerexSimulator{
			TimeStep:   timeStep,
			GlobalTime: globalTime,
		},
		AgentType: agent.Type(int64),
		Agents:    make([]*agent.Pedestrian, 0),
		Map:       &agent.Map{},
	}

	return sim
}

// SetMap :　Mapを追加する関数
func (sim *PedSimulator) SetMap(_map *agent.Map) {
	sim.Map = _map
}

// AddAgent :　エージェントを追加する関数
func (sim *PedSimulator) AddAgent(agent *agent.Pedestrian) {
	sim.Agents = append(sim.Agents, agent)
}

// UpdateAgent :　エージェントを更新する関数
func (sim *PedSimulator) UpdateAgent(a *agent.Pedestrian) error {
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
func (sim *PedSimulator) SetAgents(agents []*agent.Pedestrian) {
	for _, agentInfo := range agents {
		// 重複エリア内に入っていれば加える
		if agentInfo.IsInArea(sim.Map.Duplicate) {
			sim.AddAgent(agentInfo)
		}
	}
}

// ForwardStep :　次の時刻のエージェントを計算する関数
func (sim *PedSimulator) ForwardStep(sameAreaAgents []*agent.Pedestrian) []*agent.Pedestrian {

	pureNextAgents := make([]*agent.Pedestrian, 0)

	if IsRVO2 {
		// RVO2
		rvo2route := NewRVO2Route(sim.TimeStep, sim.GlobalTime, sim.Map, sim.Agents, sim.AgentType)
		pureNextAgents = rvo2route.CalcNextAgentsByRVO()
	} else {
		// 干渉なしで目的地へ進む
		simpleRoute := NewSimpleRoute(sim.TimeStep, sim.GlobalTime, sim.Map, sim.Agents, sim.AgentType)
		pureNextAgents = simpleRoute.CalcNextAgentsBySimple()

	}
	return pureNextAgents
}

// UpdateDuplicateAgents :　重複エリアのエージェントを更新する関数
func (sim *PedSimulator) UpdateDuplicateAgents(pureNextAgentsInfo []*agent.Pedestrian, neighborAgents []*agent.Pedestrian) []*agent.Pedestrian {
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
func (sim *PedSimulator) IsNeighborMap(mapId uint64) bool {
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
func (sim *PedSimulator) IsSameMap(mapId uint64) bool {
	if mapId == sim.Map.ID {
		return true
	}
	return false
}

// IsSameAgentType: agentTypeが同じかを判断する関数
func (sim *PedSimulator) IsSameAgentType(agentType *agent.Type) bool {
	if agentType == sim.AgentType {
		return true
	}
	return false
}

/*func (sim *PedSimulator) IsContainNeighborMap(areaId uint32) bool {
	neighborMap := sim.Area.NeighborArea

	//fmt.Printf("idlist3 %v\n", neighborMap)
	for _, neighborId := range neighborMap {
		if areaId == neighborId {
			return true
		}
	}
	return false
}

// どうしようか、、neighborAreaIdListとうの所在
// create sync id list
func (sim *PedSimulator) CreateSyncIdList(participantsInfo []*participant.ParticipantInfo) ([]uint64, []uint64) {
	sameAreaIdList := make([]uint64, 0)
	neighborAreaIdList := make([]uint64, 0)
	//fmt.Printf("idlist %v\n", sim.Area)
	for _, participantInfo := range participantsInfo {
		tAgentType := participantInfo.AgentType
		tAreaId := participantInfo.AreaId
		isNeighborArea := sim.IsContainNeighborMap(tAreaId) && int(tAgentType) == int(sim.AgentType)
		isSameArea := int(tAreaId) == int(sim.Area.AreaId) && int(tAgentType) != int(sim.AgentType)
		if isNeighborArea {

			channelId := participantInfo.ChannelId
			agentChannelId := channelId.AgentChannelId
			neighborAreaIdList = append(neighborAreaIdList, uint64(agentChannelId))
		}
		if isSameArea {
			//fmt.Printf("idlist2 %v")
			channelId := participantInfo.ChannelId
			agentChannelId := channelId.AgentChannelId
			sameAreaIdList = append(sameAreaIdList, uint64(agentChannelId))
		}
	}

	return sameAreaIdList, neighborAreaIdList
}

// if agent type and coord satisfy, return true
func (sim *PedSimulator) IsAgentInArea(agentInfo *agent.Pedestrian, areaInfo *area.AreaInfo, agentType int64) bool {
	lat := agentInfo.Route.Coord.Lat
	lon := agentInfo.Route.Coord.Lon
	slat := areaInfo.AreaCoord.StartLat
	elat := areaInfo.AreaCoord.EndLat
	slon := areaInfo.AreaCoord.StartLon
	elon := areaInfo.AreaCoord.EndLon
	if agentInfo.AgentType.String() == agent.AgentType_name[int32(agentType)] && slat <= lat && lat < elat && slon <= lon && lon < elon {
		return true
	} else {
		//log.Printf("agent type and coord is not match...\n\n")
		return false
	}
}*/
