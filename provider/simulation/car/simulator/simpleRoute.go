package simulator

import (
	"fmt"
	"math"

	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/common"
)

type SimpleRoute struct {
	TimeStep       float64
	GlobalTime     float64
	Area           *area.Area2
	Agents         []*agent.Agent
	AgentType      common.AgentType
	SameAreaAgents []*agent.Agent
}

func NewSimpleRoute(timeStep float64, globalTime float64, area *area.Area2, agents []*agent.Agent, agentType common.AgentType) *SimpleRoute {
	r := &SimpleRoute{
		TimeStep:   timeStep,
		GlobalTime: globalTime,
		Area:       area,
		Agents:     agents,
		AgentType:  agentType,
	}
	return r
}

func (simple *SimpleRoute) CalcDirectionAndDistance(startCoord *common.Coord, goalCoord *common.Coord) (float64, float64) {

	r := 6378137 // equatorial radius
	sLat := startCoord.Latitude * math.Pi / 180
	sLon := startCoord.Longitude * math.Pi / 180
	gLat := goalCoord.Latitude * math.Pi / 180
	gLon := goalCoord.Longitude * math.Pi / 180
	dLon := gLon - sLon
	dLat := gLat - sLat
	cLat := (sLat + gLat) / 2
	dx := float64(r) * float64(dLon) * math.Cos(float64(cLat))
	dy := float64(r) * float64(dLat)

	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	direction := float64(0)
	if dx != 0 && dy != 0 {
		direction = math.Atan2(dy, dx) * 180 / math.Pi
	}

	return direction, distance
}

// TODO: Why Calc Error ? newLat=nan and newLon = inf
func (simple *SimpleRoute) CalcMovedPosition(currentPosition *common.Coord, goalPosition *common.Coord, distance float64, speed float64) *common.Coord {

	sLat := currentPosition.Latitude
	sLon := currentPosition.Longitude
	gLat := goalPosition.Latitude
	gLon := goalPosition.Longitude
	// 割合
	x := speed * 1000 / 3600 / distance

	newLat := sLat + (gLat-sLat)*x
	newLon := sLon + (gLon-sLon)*x

	nextPosition := &common.Coord{
		Latitude:  newLat,
		Longitude: newLon,
	}

	return nextPosition
}

// DecideNextTransit: 次の経由地を決める関数
func (simple *SimpleRoute) DecideNextTransit(nextTransit *common.Coord, transitPoint []*common.Coord, distance float64, destination *common.Coord) *common.Coord {
	// 距離が5m以下の場合
	if distance < 5 {
		if nextTransit != destination {
			for i, tPoint := range transitPoint {
				if tPoint.Longitude == nextTransit.Longitude && tPoint.Latitude == nextTransit.Latitude {
					if i+1 == len(transitPoint) {
						// すべての経由地を通った場合、nilにする
						nextTransit = destination
					} else {
						// 次の経由地を設定する
						nextTransit = transitPoint[i+1]
					}
				}
			}
		} else {
			fmt.Printf("\x1b[30m\x1b[47m Arrived Destination! \x1b[0m\n")
		}
	}
	return nextTransit
}

// CalcNextRoute：次の時刻のRouteを計算する関数
func (simple *SimpleRoute) CalcNextRoute(agentInfo *agent.Agent, sameAreaAgents []*agent.Agent) *agent.CarRoute {

	route := agentInfo.GetCar().Route
	speed := route.Speed
	currentPosition := route.Position
	nextTransit := route.NextTransit
	transitPoints := route.TransitPoints
	destination := route.Destination
	// passed all transit point
	//if nextTransit != nil {
	//	destination = nextTransit
	//}

	// 現在位置と目標位置との距離と角度を計算
	direction, distance := simple.CalcDirectionAndDistance(currentPosition, nextTransit)

	// 次の時刻のPositionを計算
	nextPosition := simple.CalcMovedPosition(currentPosition, nextTransit, distance, speed)

	// 経由地に到着していれば、目標位置を次の経由地に更新する
	nextTransit = simple.DecideNextTransit(nextTransit, transitPoints, distance, destination)

	//fmt.Printf("\x1b[30m\x1b[47m Position %v, NextTransit: %v, NextTransit: %v, Direction: %v, Distance: %v \x1b[0m\n", currentPosition, nextTransit, destination, direction, distance)
	//fmt.Printf("\x1b[30m\x1b[47m 上下:  %v, 左右: %v \x1b[0m\n", nextTransit.Lat-currentPosition.Lat, nextTransit.Lon-currentPosition.Lon)
	/*nextPosition := &common.Coord{
		Latitude: currentPosition.Latitude,
		Lonitude: currentPosition.Longitude,
	}
	//TODO: Fix this
	if newLat < 40 && newLat > 0 && newLon < 150 && newLon > 0 {
		nextPosition = &common.Coord{
			Latitude: newLat,
			Longitude: newLon,
		}*/
	//} else {
	//	log.Printf("\x1b[30m\x1b[47m LOCATION CULC ERROR %v \x1b[0m\n", nextPosition)
	//}

	nextRoute := &agent.CarRoute{
		Position:      nextPosition,
		Direction:     direction,
		Speed:         distance,
		Destination:   route.Destination,
		Departure:     route.Departure,
		TransitPoints: transitPoints,
		NextTransit:   nextTransit,
		TotalDistance: route.TotalDistance,
		RequiredTime:  route.RequiredTime,
	}

	return nextRoute
}

// IsAgentInControlArea: AgentがControlAreaに存在するかどうか
func (simple *SimpleRoute) IsAgentInControlArea(agentInfo *agent.Agent) bool {
	areaInfo := simple.Area
	agentType := simple.AgentType
	car := agentInfo.GetCar()
	lat := car.Route.Position.Latitude
	lon := car.Route.Position.Longitude
	areaCoords := areaInfo.ControlArea
	deg := 0.0
	for i, coord := range areaCoords{
		p2lat := coord.Latitude
		p2lon := coord.Longitude
		p3lat := areaCoords[i+1].Latitude
		p3lon := areaCoords[i+1].Longitude
		if i == len(areaCoords)-1 {
			p3lat = areaCoords[0].Latitude
			p3lon = areaCoords[0].Longitude
		}
		alat := p2lat - lat
		alon := p2lon - lon
		blat := p3lat - lat
		blon := p3lon - lon
		cos := (alat*blat + alon*blon) / (math.Sqrt(alat*alat + alon+alon)*math.Sqrt(blat*blat + blon+blon))
		deg += math.Acos(cos) * float64(180) / math.Pi
	}
	if agentInfo.Type == agentType && math.Round(deg) == 360{
		return true
	} else {
		return false
	}
}

// CalcNextAgents: 次の時刻のエージェントを取得する関数
func (simple *SimpleRoute) CalcNextAgents() []*agent.Agent {

	nextControlAgents := make([]*agent.Agent, 0)

	for _, agentInfo := range simple.Agents {
		// 自エリアにいる場合、次のルートを計算する
		if simple.IsAgentInControlArea(agentInfo) {

			// 現在のCar情報
			currentCarInfo := agentInfo.GetCar()

			// 次の時刻のRouteを計算
			nextRoute := simple.CalcNextRoute(agentInfo, simple.SameAreaAgents)

			car := &agent.Car{
				Status: currentCarInfo.Status,
				Route:  nextRoute,
			}

			nextControlAgent := &agent.Agent{
				Id:   agentInfo.Id,
				Type: agentInfo.Type,
				Data: &agent.Agent_Car{
					Car: car,
				},
			}
			// Agent追加
			nextControlAgents = append(nextControlAgents, nextControlAgent)
		}
	}
	return nextControlAgents
}
