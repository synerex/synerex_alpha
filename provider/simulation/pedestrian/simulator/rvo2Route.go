package simulator

import (
	"fmt"
	"math"

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/common"
)

var (
	sim *rvo.RVOSimulator
)

type RVO2Route struct {
	TimeStep   float64
	GlobalTime float64
	Area       *area.Area
	Agents     []*agent.Agent
	AgentType  common.AgentType
}

func NewRVO2Route(timeStep float64, globalTime float64, area *area.Area, agentsInfo []*agent.Agent, agentType common.AgentType) *RVO2Route {

	r := &RVO2Route{
		TimeStep:   timeStep,
		GlobalTime: globalTime,
		Area:       area,
		Agents:     agentsInfo,
		AgentType:  agentType,
	}
	return r
}

func (rvo2route *RVO2Route) CalcDirectionAndDistance(startCoord *common.Coord, goalCoord *common.Coord) (float64, float64) {

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

/*func Scale(agentInfo []*agent.Agent) ([]Agent, []Obstacle) {
	// エリア情報
	areaCoord := rvo2route.Area.DuplicateArea
	height := math.Abs(areaCoord.EndLat - areaCoord.StartLat)
	minLat := math.Min(areaCoord.EndLat, areaCoord.StartLat)
	width := math.Abs(areaCoord.EndLon - areaCoord.StartLon)
	minLon := math.Min(areaCoord.EndLon, areaCoord.StartLon)

	// エージェント情報
	simAgents := rvo2route.Agents

	scaledAgents := make([]Agent, 0)
	for _, agent := range simAgents {
		route := agent.GetPedestrian().Route
		scaledLat := (route.Position.Latitude - minLat) / height
		scaledLon := (route.Position.Longitude - minLon) / width

		//velLat := agent.Velocity.Lat / height
		//velLon := agent.Velocity.Lon / width
		//fmt.Printf("\x1b[30m\x1b[47m direction: %v \x1b[0m\n", route.Direction)
		velLat := math.Sin(float64(route.Direction * math.Pi / 180))
		velLon := math.Cos(float64(route.Direction * math.Pi / 180))

		goalLat := (route.NextTransit.Latitude - minLat) / height
		goalLon := (route.NextTransit.Longitude - minLon) / width

		scaledCoord := Coord{
			Lat: scaledLat,
			Lon: scaledLon,
		}

		scaledVelocity := Coord{
			Lat: velLat,
			Lon: velLon,
		}
		goal := Coord{
			Lat: goalLat,
			Lon: goalLon,
		}

		scaledAgents = append(scaledAgents, Agent{
			Id:       agent.Id,
			Position: scaledCoord,
			Velocity: scaledVelocity,
			Goal:     goal,
			MaxSpeed: agent.GetPedestrian().Route.Speed,
		})

		//fmt.Printf("\x1b[30m\x1b[47m scaledVelocity: %v \x1b[0m\n", scaledVelocity)
		//fmt.Printf("\x1b[30m\x1b[47m goal: %v \x1b[0m\n", goal)
	}

	scaledObstacles := make([]Obstacle, 0)
	obstacle1 := Obstacle{
		Id: 1,
		Position: []Coord{
			Coord{
				Latitude: 35.152476,
				Longitude: 136.982300,
			},
			Coord{
				Latitude: 35.155078,
				Longitude: 136.982300,
			},
			Coord{
				Latitude: 35.152476,
				Longitude: 136.982700,
			},
			Coord{
				Latitude: 35.155078,
				Longitude: 136.982700,
			},
		},
	}
	obstacle2 := Obstacle{
		Id: 1,
		Position: []Coord{
			Coord{
				Latitude: 35.157576,
				Longitude: 136.982300,
			},
			Coord{
				Latitude: 35.160678,
				Longitude: 136.982300,
			},
			Coord{
				Latitude: 35.157576,
				Longitude: 136.982700,
			},
			Coord{
				Latitude: 35.160678,
				Longitude: 136.982700,
			},
		},
	}
	obstacles := []Obstacle{
		//obstacle1,
		//obstacle2,
	}

	for _, obstacle := range obstacles {
		var scaledObstacle Obstacle
		scaledObstacle.Id = obstacle.Id
		for _, coord := range obstacle.Position {
			scaledLat := (coord.Latitude - minLat) / height
			scaledLon := (coord.Longitude - minLon) / width
			scaledObstacle.Position = append(scaledObstacle.Position, Coord{Latitude: scaledLat, Longitude: scaledLon})
		}
		scaledObstacles = append(scaledObstacles, scaledObstacle)
	}

	return scaledAgents, scaledObstacles
}*/

func (rvo2route *RVO2Route) ScaleCoord(coord *common.Coord) *common.Coord {

	// エリア情報
	areaCoord := rvo2route.Area.DuplicateArea
	height := math.Abs(areaCoord.EndLat - areaCoord.StartLat)
	minLat := math.Min(areaCoord.EndLat, areaCoord.StartLat)
	width := math.Abs(areaCoord.EndLon - areaCoord.StartLon)
	minLon := math.Min(areaCoord.EndLon, areaCoord.StartLon)

	// scale化
	scaledLat := (coord.Latitude - minLat) / height
	scaledLon := (coord.Longitude - minLon) / width

	scaledCoord := &common.Coord{
		Latitude:  scaledLat,
		Longitude: scaledLon,
	}

	return scaledCoord
}

func (rvo2route *RVO2Route) InvScaleCoord(coord *common.Coord) *common.Coord {

	// エリア情報
	areaCoord := rvo2route.Area.DuplicateArea
	height := math.Abs(areaCoord.EndLat - areaCoord.StartLat)
	minLat := math.Min(areaCoord.EndLat, areaCoord.StartLat)
	width := math.Abs(areaCoord.EndLon - areaCoord.StartLon)
	minLon := math.Min(areaCoord.EndLon, areaCoord.StartLon)

	lat := coord.Latitude*height + minLat
	lon := coord.Longitude*width + minLon

	invScaleCoord := &common.Coord{
		Latitude:  lat,
		Longitude: lon,
	}

	return invScaleCoord
}

func (rvo2route *RVO2Route) DecideNextTransit(nextTransit *common.Coord, transitPoint []*common.Coord, distance float64, destination *common.Coord) *common.Coord {
	// 距離が5m以下の場合
	if distance < 150 {
		if nextTransit != destination {
			for i, tPoint := range transitPoint {
				if tPoint.Longitude == nextTransit.Longitude && tPoint.Latitude == nextTransit.Latitude {
					if i+1 == len(transitPoint) {
						// すべての経由地を通った場合、nextTransitをdestinationにする
						nextTransit = destination
					} else {
						// 次の経由地を設定する
						nextTransit = transitPoint[i+1]
					}
				}
			}
		} else {
			//fmt.Printf("arrived!")
		}
	}
	return nextTransit
}

/*func (rvo2route *RVO2Route) InvScale() []*agent.Agent {
	// エリア情報
	areaCoord := rvo2route.Area.DuplicateArea
	height := math.Abs(areaCoord.EndLat - areaCoord.StartLat)
	minLat := math.Min(areaCoord.EndLat, areaCoord.StartLat)
	width := math.Abs(areaCoord.EndLon - areaCoord.StartLon)
	minLon := math.Min(areaCoord.EndLon, areaCoord.StartLon)

	nextAgents := make([]*agent.Agent, 0)
	for i, rvoAgent := range rvo2route.RVOAgents {
			lat := rvoAgent.Position.Latitude*height + minLat
			lon := rvoAgent.Position.Longitude*width + minLon

			route := simAgent.GetPedestrian().Route
			//currentLocation := route.Coord
			nextTransit := route.NextTransit
			transitPoint := route.TransitPoints
			destination := route.Destination

			// calc direction, distance
			// 現在の位置とゴールとの距離と角度を求める (度, m))
			direction, distance := rvo2route.CalcDirectionAndDistance(lat, lon, nextTransit.Latitude, nextTransit.Longitude)

			var nextCoord *common.Coord
			nextCoord = &common.Coord{
				Latitude:  lat,
				Longitude: lon,
			}
			// 次の経由地nextTransitを求める
			nextTransit = rvo2route.DecideNextTransit(nextTransit, transitPoint, distance, destination)

			// ゴールに着いたら停止させる
			if nextTransit == destination && distance < 100 {
				fmt.Printf("\x1b[30m\x1b[47m Goal! id: %v \x1b[0m\n", simAgent.AgentId)
				nextCoord = route.Coord
			} else {
				nextCoord = &common.Coord{
					Lat: lat,
					Lon: lon,
				}
				// 次の経由地nextTransitを求める
				nextTransit = rvo2route.DecideNextTransit(nextTransit, transitPoint, distance, destination)
			}

			fmt.Printf("\x1b[30m\x1b[47m nextTransit %v, distance: %v \x1b[0m\n", nextTransit, distance)

			nextRoute := &agent.PedRoute{
				Position:      nextCoord,
				Direction:     direction,
				Speed:         route.Speed,
				Destination:   route.Destination,
				Departure:     route.Departure,
				TransitPoints: route.TransitPoints,
				NextTransit:   nextTransit,
				TotalDistance: route.TotalDistance,
				RequiredTime:  route.RequiredTime,
			}

			nextAgent := &agent.Agent{
				GlobalTime:  rvo2route.GlobalTime + 1,
				AgentId:     simAgent.GetPedestrian().Id,
				AgentType:   simAgent.GetPedestrian().AgentType,
				AgentStatus: simAgent.GetPedestrian().AgentStatus,
				Route:       nextRoute,
			}

			nextAgents = append(nextAgents, nextAgent)
	}

	return nextAgents
}*/

// FIX: lat?lon?  X: Lon: 経度、Y: Lat: 緯度
func (rvo2route *RVO2Route) SetupScenario() {

	// Set Agent
	for _, agent := range rvo2route.Agents {
		ped := agent.GetPedestrian()
		scaledPosition := rvo2route.ScaleCoord(ped.Route.Position)
		scaledVelocity := rvo2route.ScaleCoord(
			&common.Coord{
				Latitude:  math.Sin(float64(ped.Route.Direction * math.Pi / 180)),
				Longitude: math.Cos(float64(ped.Route.Direction * math.Pi / 180)),
			},
		)
		scaledGoal := rvo2route.ScaleCoord(ped.Route.NextTransit)
		position := &rvo.Vector2{X: scaledPosition.Longitude, Y: scaledPosition.Latitude}
		velocity := &rvo.Vector2{X: scaledVelocity.Longitude, Y: scaledVelocity.Latitude}
		goal := &rvo.Vector2{X: scaledGoal.Longitude, Y: scaledGoal.Latitude}
		// Agentを追加
		id, _ := sim.AddDefaultAgent(position)
		// エージェントの速度方向ベクトルを設定
		sim.SetAgentPrefVelocity(id, velocity)
		// 目的地を設定
		sim.SetAgentGoal(id, goal)
		//sim.SetAgentMaxSpeed(id, float64(agent.MaxSpeed))
	}

	// Set Obstacle
	/*if len(scaledObstacles) != 0 {
		for _, obstacle := range scaledObstacles {
			rvoObsPosition := []*rvo.Vector2{}
			for _, position := range obstacle.Position {
				rvoObsPosition = append(rvoObsPosition, &rvo.Vector2{X: float64(position.Lon), Y: float64(position.Lat)})
			}
			sim.AddObstacle(rvoObsPosition)
		}
		sim.ProcessObstacles()
	}*/

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation..\n\n")
}

func (rvo2route *RVO2Route) SetPreferredVelocities() {

	for i, _ := range rvo2route.Agents {

		// setPrefVelocity
		goalVector := sim.GetAgentGoalVector(i)

		if rvo.Sqr(goalVector) > 1 {
			goalVector = rvo.Normalize(goalVector)
		}

		sim.SetAgentPrefVelocity(i, goalVector)

	}
}

func (rvo2route *RVO2Route) CalcNextAgents() []*agent.Agent {

	nextControlAgents := make([]*agent.Agent, 0)
	currentAgents := rvo2route.Agents

	timeStep := rvo2route.TimeStep
	neighborDist := 0.075 // どのくらいの距離の相手をNeighborと認識するか?Neighborとの距離をどのくらいに保つか？ぶつかったと認識する距離？
	maxneighbors := 100   // 周り何体を計算対象とするか
	timeHorizon := 1.5
	timeHorizonObst := 2.0
	radius := 0.025  // エージェントの半径
	maxSpeed := 0.02 // エージェントの最大スピード
	sim = rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})

	rvo2route.SetupScenario()

	// Stepを進める
	sim.DoStep()

	// 速度ベクトルをセットする
	//rvo2route.SetPreferredVelocities()

	//nextAgents := rvo2route.InvScale(nextScaledAgents)
	//nextRVOAgents := sim.Agents

	// 管理エリアのエージェントのみを抽出
	for rvoId, agentInfo := range currentAgents {
		//nextRVOAgent := sim.GetAgent(int(agentInfo.Id))
		currentPedInfo := agentInfo.GetPedestrian()
		// 計算する前に自エリアにいる場合、次のルートを計算する
		if rvo2route.IsAgentInControlArea(agentInfo) {

			destination := currentPedInfo.Route.Destination

			// rvoの位置情報を緯度経度に変換する
			rvoAgentPosition := sim.GetAgentPosition(int(rvoId))
			scaledPosition := &common.Coord{
				Latitude:  rvoAgentPosition.Y,
				Longitude: rvoAgentPosition.X,
			}
			nextCoord := rvo2route.InvScaleCoord(scaledPosition)

			// 現在の位置とゴールとの距離と角度を求める (度, m))
			direction, distance := rvo2route.CalcDirectionAndDistance(nextCoord, currentPedInfo.Route.NextTransit)
			// 次の経由地nextTransitを求める
			nextTransit := rvo2route.DecideNextTransit(currentPedInfo.Route.NextTransit, currentPedInfo.Route.TransitPoints, distance, destination)

			nextRoute := &agent.PedRoute{
				Position:      nextCoord,
				Direction:     direction,
				Speed:         distance,
				Destination:   destination,
				Departure:     currentPedInfo.Route.Departure,
				TransitPoints: currentPedInfo.Route.TransitPoints,
				NextTransit:   nextTransit,
				TotalDistance: currentPedInfo.Route.TotalDistance,
				RequiredTime:  currentPedInfo.Route.RequiredTime,
			}

			ped := &agent.Pedestrian{
				Status: currentPedInfo.Status,
				Route:  nextRoute,
			}

			nextControlAgent := &agent.Agent{
				Id:   agentInfo.Id,
				Type: agentInfo.Type,
				Data: &agent.Agent_Pedestrian{
					Pedestrian: ped,
				},
			}

			nextControlAgents = append(nextControlAgents, nextControlAgent)
		}
	}

	return nextControlAgents
}

func (rvo2route *RVO2Route) IsAgentInControlArea(agentInfo *agent.Agent) bool {
	
	areaInfo := rvo2route.Area
	agentType := rvo2route.AgentType
	ped := agentInfo.GetPedestrian()
	lat := ped.Route.Position.Latitude
	lon := ped.Route.Position.Longitude
	slat := areaInfo.ControlArea.StartLat
	elat := areaInfo.ControlArea.EndLat
	slon := areaInfo.ControlArea.StartLon
	elon := areaInfo.ControlArea.EndLon
	if agentInfo.Type == agentType && slat <= lat && lat < elat && slon <= lon && lon < elon {
		return true
	}
	//log.Printf("agent type and coord is not match...\n\n")
	return false
}
