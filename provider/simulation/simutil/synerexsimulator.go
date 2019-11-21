package simutil

import (
	"fmt"
	"log"

	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
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
func NewSynerexSimulator(timeStep float64, agentType int32, globalTime float64) *SynerexSimulator {

	sim := &SynerexSimulator{
		TimeStep:   timeStep,
		Agents:     make([]*agent.AgentInfo, 0),
		AgentType:  agentType,
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

// Finish Fix
func (sim *SynerexSimulator) CalcNextRoute(agentInfo *agent.AgentInfo, sameAreaAgents []*agent.AgentInfo) *agent.Route {

	route := agentInfo.Route
	speed := route.Speed
	currentLocation := route.Coord
	nextTransit := route.RouteInfo.NextTransit
	transitPoint := route.RouteInfo.TransitPoint
	destination := route.Destination
	// passed all transit point
	if nextTransit != nil {
		destination = nextTransit
	}

	direction, distance := CalcDirectionAndDistance(currentLocation.Lat, currentLocation.Lon, destination.Lat, destination.Lon)
	//newLat, newLon := simutil.CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, speed*1000/3600, direction)
	newLat, newLon := CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, destination.Lat, destination.Lon, distance, speed)

	// upate next trasit point
	if distance < 10 {
		if nextTransit != nil {
			nextTransit2 := nextTransit
			for i, tPoint := range transitPoint {
				if tPoint.Lon == nextTransit2.Lon && tPoint.Lat == nextTransit2.Lat {
					if i+1 == len(transitPoint) {
						// pass all transit point
						nextTransit = nil
					} else {
						// go to next transit point
						nextTransit = transitPoint[i+1]
					}
				}
			}
		} else {
			log.Printf("\x1b[30m\x1b[47m Arrived Destination! \x1b[0m\n")
		}

	}

	nextCoord := &agent.Coord{
		Lat: currentLocation.Lat,
		Lon: currentLocation.Lon,
	}
	//TODO: Fix this
	if newLat < 40 && newLat > 0 && newLon < 150 && newLon > 0 {
		nextCoord = &agent.Coord{
			Lat: newLat,
			Lon: newLon,
		}
	} else {
		log.Printf("\x1b[30m\x1b[47m LOCATION CULC ERROR %v \x1b[0m\n", nextCoord)

	}

	routeInfo := &agent.RouteInfo{
		TransitPoint:  transitPoint,
		NextTransit:   nextTransit,
		TotalDistance: route.RouteInfo.TotalDistance,
		RequiredTime:  route.RouteInfo.RequiredTime,
	}
	nextRoute := &agent.Route{
		Coord:       nextCoord,
		Direction:   float32(direction),
		Speed:       float32(speed),
		Destination: route.Destination,
		Departure:   route.Departure,
		RouteInfo:   routeInfo,
	}
	return nextRoute
}

func (sim *SynerexSimulator) CalcNextAgents(sameAreaAgents []*agent.AgentInfo) []*agent.AgentInfo {
	// calc agent
	//otherAgentsInfo := make([]*agent.AgentInfo, 0) // ??
	pureNextAgents := make([]*agent.AgentInfo, 0)
	for _, agentInfo := range sim.Agents {
		// 自エリアにいる場合、次のルートを計算する
		if IsAgentInControlledArea(agentInfo, sim.Area, sim.AgentType) {

			nextRoute := sim.CalcNextRoute(agentInfo, sameAreaAgents)

			pureNextAgent := &agent.AgentInfo{
				Time:        uint32(sim.GlobalTime) + 1,
				AgentId:     agentInfo.AgentId,
				AgentType:   agentInfo.AgentType,
				AgentStatus: agentInfo.AgentStatus,
				Route:       nextRoute,
			}

			pureNextAgents = append(pureNextAgents, pureNextAgent)
		}
	}

	return pureNextAgents
}

// UpdateDuplicateAgents :　重複エリアのエージェントを更新する関数
func (sim *SynerexSimulator) UpdateDuplicateAgents(pureNextAgentsInfo []*agent.AgentInfo, neighborAgents []*agent.AgentInfo) []*agent.AgentInfo {
	nextAgents := pureNextAgentsInfo
	for _, neighborAgent := range neighborAgents {
		//　隣のエージェントが自分のエリアにいてかつ自分のエリアのエージェントと被ってない場合更新
		if len(pureNextAgentsInfo) == 0 {
			if IsAgentInArea(neighborAgent, sim.Area, sim.AgentType) {
				//log.Println("CHANGE_AREA1!!!")
				nextAgents = append(nextAgents, neighborAgent)
			}
		} else {
			isAppendAgent := false
			for _, sameAreaAgent := range pureNextAgentsInfo {
				if neighborAgent.AgentId != sameAreaAgent.AgentId && IsAgentInArea(neighborAgent, sim.Area, sim.AgentType) {
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

func (sim *SynerexSimulator) IsContainNeighborMap(areaId uint32) bool {
	neighborMap := sim.Area.NeighborArea

	fmt.Printf("idlist3 %v\n", neighborMap)
	for _, neighborId := range neighborMap {
		if areaId == neighborId {
			return true
		}
	}
	return false
}

// create sync id list
func (sim *SynerexSimulator) CreateSyncIdList(participantsInfo []*participant.ParticipantInfo) ([]uint64, []uint64) {
	sameAreaIdList := make([]uint64, 0)
	neighborAreaIdList := make([]uint64, 0)
	fmt.Printf("idlist %v\n", sim.Area)
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
			fmt.Printf("idlist2 %v")
			channelId := participantInfo.ChannelId
			agentChannelId := channelId.AgentChannelId
			sameAreaIdList = append(sameAreaIdList, uint64(agentChannelId))
		}
	}

	return sameAreaIdList, neighborAreaIdList
}
