package simulator

import (

	//"strconv"

	//"encoding/json"
	//"io/ioutil"

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
)

var (
	goals []*rvo.Vector2
)

func init() {
	goals = make([]*rvo.Vector2, 0)
}

/*type Coord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type areaCoord struct {
	StartLat float64 `json:"slat"`
	EndLat   float64 `json:"elat"`
	StartLon float64 `json:"slon"`
	EndLon   float64 `json:"elon"`
}

type Map struct {
	Id    uint64    `json:"id"`
	Coord areaCoord `json:"coord"`
}

type Agent struct {
	Id       uint64 `json:"id"`
	Velocity Coord  `json:"velocity"`
	Coord    Coord  `json:"coord"`
	Goal     Coord  `json:"goal"`
	MaxSpeed float64
}

type Obstacle struct {
	Id       uint64  `json:"id"`
	Position []Coord `json:"velocity"`
}

type RVO2Route struct {
	TimeStep   float64
	GlobalTime float64
	Area       *area.Area
	Agents     []*agent.Agent
	AgentType  int64
}

func NewRVO2Route(timeStep float64, globalTime float64, area *area.Area, agents []*agent.Agent, agentType int64) *RVO2Route {
	r := &RVO2Route{
		TimeStep:   timeStep,
		GlobalTime: globalTime,
		Area:       area,
		Agents:     agents,
		AgentType:  agentType,
	}
	return r
}

func (rvo2route *RVO2Route) CalcDirectionAndDistance(sLat float64, sLon float64, gLat float64, gLon float64) (float64, float64) {

	r := 6378137 // equatorial radius
	sLat = sLat * math.Pi / 180
	sLon = sLon * math.Pi / 180
	gLat = gLat * math.Pi / 180
	gLon = gLon * math.Pi / 180
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

func (rvo2route *RVO2Route) Scale() ([]Agent, []Obstacle) {
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
			Coord:    scaledCoord,
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
				Lat: 35.152476,
				Lon: 136.982300,
			},
			Coord{
				Lat: 35.155078,
				Lon: 136.982300,
			},
			Coord{
				Lat: 35.152476,
				Lon: 136.982700,
			},
			Coord{
				Lat: 35.155078,
				Lon: 136.982700,
			},
		},
	}
	obstacle2 := Obstacle{
		Id: 1,
		Position: []Coord{
			Coord{
				Lat: 35.157576,
				Lon: 136.982300,
			},
			Coord{
				Lat: 35.160678,
				Lon: 136.982300,
			},
			Coord{
				Lat: 35.157576,
				Lon: 136.982700,
			},
			Coord{
				Lat: 35.160678,
				Lon: 136.982700,
			},
		},
	}
	obstacles := []Obstacle{
		obstacle1,
		obstacle2,
	}

	for _, obstacle := range obstacles {
		var scaledObstacle Obstacle
		scaledObstacle.Id = obstacle.Id
		for _, coord := range obstacle.Position {
			scaledLat := (coord.Lat - minLat) / height
			scaledLon := (coord.Lon - minLon) / width
			scaledObstacle.Position = append(scaledObstacle.Position, Coord{Lat: scaledLat, Lon: scaledLon})
		}
		scaledObstacles = append(scaledObstacles, scaledObstacle)
	}

	return scaledAgents, scaledObstacles
}

func (rvo2route *RVO2Route) DecideNextTransit(nextTransit *common.Coord, transitPoint []*common.Coord, distance float64, destination *common.Coord) *common.Coord {
	// 距離が5m以下の場合
	if distance < 150 {
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
		}
	}
	return nextTransit
}

func (rvo2route *RVO2Route) InvScale(scaledAgents []Agent) []*agent.Agent {
	// エリア情報
	areaCoord := rvo2route.Area.DuplicateArea
	height := math.Abs(areaCoord.EndLat - areaCoord.StartLat)
	minLat := math.Min(areaCoord.EndLat, areaCoord.StartLat)
	width := math.Abs(areaCoord.EndLon - areaCoord.StartLon)
	minLon := math.Min(areaCoord.EndLon, areaCoord.StartLon)

	// エージェント情報
	simAgents := rvo2route.Agents

	nextAgents := make([]*agent.Agent, 0)
	for i, scaledAgent := range scaledAgents {
		simAgent := simAgents[i]
		if simAgent.Id == scaledAgent.Id {
			lat := scaledAgent.Coord.Lat*height + minLat
			lon := scaledAgent.Coord.Lon*width + minLon

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
				TransitPoints:  route.TransitPoints,
				NextTransit:   nextTransit,
				TotalDistance: route.TotalDistance,
				RequiredTime:  route.RequiredTime,
			}

			nextAgent := &agent.Agent{
				GlobalTime:        rvo2route.GlobalTime + 1,
				AgentId:     simAgent.GetPedestrian().Id,
				AgentType:   simAgent.GetPedestrian().AgentType,
				AgentStatus: simAgent.GetPedestrian().AgentStatus,
				Route:       nextRoute,
			}

			nextAgents = append(nextAgents, nextAgent)
		}
	}

	return nextAgents
}

// FIX: lat?lon?  X: Lon: 経度、Y: Lat: 緯度
func SetupScenario(sim *rvo.RVOSimulator, scaledAgents []Agent, scaledObstacles []Obstacle) {

	// Set Agent
	for i, agent := range scaledAgents {
		position := &rvo.Vector2{X: float64(agent.Position.Longitude), Y: float64(common.Position.Latitude)}
		velocity := &rvo.Vector2{X: float64(agent.Velocity.Lon), Y: float64(agent.Velocity.Lat)}
		id, _ := sim.AddAgentPosition(position)
		sim.SetAgentPrefVelocity(id, velocity)
		fmt.Printf("\x1b[30m\x1b[47m position %v, velocity: %v \x1b[0m\n", position, sim.GetAgentVelocity(i))
		//sim.SetAgentMaxSpeed(id, float64(agent.MaxSpeed))
	}

	// Set Obstacle
	if len(scaledObstacles) != 0 {
		for _, obstacle := range scaledObstacles {
			rvoObsPosition := []*rvo.Vector2{}
			for _, position := range obstacle.Position {
				rvoObsPosition = append(rvoObsPosition, &rvo.Vector2{X: float64(position.Lon), Y: float64(position.Lat)})
			}
			sim.AddObstacle(rvoObsPosition)
		}
		sim.ProcessObstacles()
	}

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation..\n\n")
}

func (rvo2route *RVO2Route) SetPreferredVelocities(sim *rvo.RVOSimulator, scaledAgents []Agent) []Agent {

	// エージェント情報
	//simAgents := rvo2route.Agents

	nextScaledAgents := make([]Agent, 0)
	for i, agent := range scaledAgents {

		// setPrefVelocity
		goal := &rvo.Vector2{
			X: float64(agent.Goal.Lon),
			Y: float64(agent.Goal.Lat),
		}
		goalVector := rvo.Sub(goal, sim.GetAgentPosition(i))

		if rvo.Sqr(goalVector) > 1 {
			goalVector = rvo.Normalize(goalVector)
		}

		sim.SetAgentPrefVelocity(i, goalVector)

		// nextScledAgents
		scaledPosition := sim.GetAgentPosition(i)
		scaledVelocity := sim.GetAgentPrefVelocity(i)
		scaledAgent := Agent{
			Id: scaledAgents[i].Id,
			Coord: Coord{
				Lat: scaledPosition.Y,
				Lon: scaledPosition.X,
			},
			Goal: scaledAgents[i].Goal,
			Velocity: Coord{
				Lat: scaledVelocity.Y,
				Lon: scaledVelocity.X,
			},
		}

		nextScaledAgents = append(nextScaledAgents, scaledAgent)
	}
	return nextScaledAgents
}

func (rvo2route *RVO2Route) CalcNextAgents() []*agent.Agent {
	timeStep := rvo2route.TimeStep
	neighborDist := 0.075 // どのくらいの距離の相手をNeighborと認識するか?Neighborとの距離をどのくらいに保つか？ぶつかったと認識する距離？
	maxneighbors := 100   // 周り何体を計算対象とするか
	timeHorizon := 1.5
	timeHorizonObst := 2.0
	radius := 0.025  // エージェントの半径
	maxSpeed := 0.02 // エージェントの最大スピード
	sim := rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})

	scaledAgents, scaledObstacles := rvo2route.Scale()
	//fmt.Printf("scaledagents is : %v\n", scaledAgents)
	//invScaledAgents := rvo2route.InvScale(scaledAgents)

	//for _, agent := range rvo2route.Agents{
	//	fmt.Printf("simAgents is : %v\n", agent.Route.Coord)
	//}

	//for _, agent := range invScaledAgents{
	//	fmt.Printf("invAgents is : %v\n", agent.Route.Coord)
	//}

	SetupScenario(sim, scaledAgents, scaledObstacles)

	// Stepを進める
	sim.DoStep()

	// 速度ベクトルをセットする
	nextScaledAgents := rvo2route.SetPreferredVelocities(sim, scaledAgents)
	//fmt.Printf("nextScaledAgents is : %v\n", nextScaledAgents)

	nextAgents := rvo2route.InvScale(nextScaledAgents)

	//fmt.Printf("nextAgents is : %v\n", nextAgents)

	return nextAgents
}

func (rvo2route *RVO2Route) IsAgentInControlledArea(agentInfo *agent.Agent, areaInfo *area.Area, agentType int64) bool {
	lat := agentInfo.GetPedestrian().Route.Coord.Lat
	lon := agentInfo.GetPedestrian().Route.Coord.Lon
	slat := areaInfo.ControlAreaCoord.StartLat
	elat := areaInfo.ControlAreaCoord.EndLat
	slon := areaInfo.ControlAreaCoord.StartLon
	elon := areaInfo.ControlAreaCoord.EndLon
	if agentInfo.AgentType.String() == agent.AgentType_name[int32(agentType)] && slat <= lat && lat < elat && slon <= lon && lon < elon {
		return true
	}
	//log.Printf("agent type and coord is not match...\n\n")
	return false
}

func (rvo2route *RVO2Route) CalcNextAgentsByRVO() []*agent.Agent {
	pureNextAgents := make([]*agent.Agent, 0)
	//rvo2util := NewRVO2Route(rvo2route)
	currentAgents := rvo2route.Agents
	//nextAgents := rvo2util.CalcNextAgents(sim.Agents, sameAreaAgents)
	nextAgents := rvo2route.CalcNextAgents()
	for i, agentInfo := range currentAgents {
		nextAgent := nextAgents[i]
		// 自エリアにいる場合、次のルートを計算する
		if rvo2route.IsAgentInControlledArea(agentInfo, rvo2route.Area, rvo2route.AgentType) {

			pureNextAgent := &agent.AgentInfo{
				Time:        uint32(rvo2route.GlobalTime) + 1,
				AgentId:     agentInfo.AgentId,
				AgentType:   agentInfo.AgentType,
				AgentStatus: agentInfo.AgentStatus,
				Route:       nextAgent.Route,
			}

			pureNextAgents = append(pureNextAgents, pureNextAgent)
		}
	}

	return pureNextAgents
}*/
