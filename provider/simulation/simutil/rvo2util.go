package simutil

import (
	"fmt"
	"math"

	//"strconv"

	//"encoding/json"
	//"io/ioutil"

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
)

var (
	goals []*rvo.Vector2
)

func init() {
	goals = make([]*rvo.Vector2, 0)
}

type areaCoord struct {
	StartLat float32 `json:"slat"`
	EndLat   float32 `json:"elat"`
	StartLon float32 `json:"slon"`
	EndLon   float32 `json:"elon"`
}

type Map struct {
	Id    uint32    `json:"id"`
	Coord areaCoord `json:"coord"`
}

type Agent struct {
	Id       uint32 `json:"id"`
	Velocity Coord  `json:"velocity"`
	Coord    Coord  `json:"coord"`
	Goal     Coord  `json:"goal"`
	MaxSpeed float32
}

type Obstacle struct {
	Id       uint32  `json:"id"`
	Position []Coord `json:"velocity"`
}

type RVO2Util struct {
	SynSim    *SynerexSimulator
	RVOAgents []Agent
}

func NewRVO2Util(synSim *SynerexSimulator) *RVO2Util {
	r := &RVO2Util{
		SynSim: synSim,
	}
	return r
}

func (rvoutil *RVO2Util) CalcDirectionAndDistance(sLat float32, sLon float32, gLat float32, gLon float32) (float32, float32) {

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

	distance := float32(math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2)))
	direction := float32(0)
	if dx != 0 && dy != 0 {
		direction = float32(math.Atan2(dy, dx)) * 180 / math.Pi
	}

	return direction, distance
}

func (rvoutil *RVO2Util) Scale() ([]Agent, []Obstacle) {
	// エリア情報
	areaCoord := rvoutil.SynSim.Area.AreaCoord
	height := float32(math.Abs(float64(areaCoord.EndLat - areaCoord.StartLat)))
	minLat := float32(math.Min(float64(areaCoord.EndLat), float64(areaCoord.StartLat)))
	width := float32(math.Abs(float64(areaCoord.EndLon - areaCoord.StartLon)))
	minLon := float32(math.Min(float64(areaCoord.EndLon), float64(areaCoord.StartLon)))

	// エージェント情報
	simAgents := rvoutil.SynSim.Agents

	scaledAgents := make([]Agent, 0)
	for _, agent := range simAgents {
		route := agent.Route
		scaledLat := (route.Coord.Lat - minLat) / height
		scaledLon := (route.Coord.Lon - minLon) / width

		//velLat := agent.Velocity.Lat / height
		//velLon := agent.Velocity.Lon / width
		//fmt.Printf("\x1b[30m\x1b[47m direction: %v \x1b[0m\n", route.Direction)
		velLat := float32(math.Sin(float64(route.Direction * math.Pi / 180)))
		velLon := float32(math.Cos(float64(route.Direction * math.Pi / 180)))

		goalLat := (route.RouteInfo.NextTransit.Lat - minLat) / height
		goalLon := (route.RouteInfo.NextTransit.Lon - minLon) / width

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
			Id:       agent.AgentId,
			Coord:    scaledCoord,
			Velocity: scaledVelocity,
			Goal:     goal,
			MaxSpeed: agent.Route.Speed,
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
				Lon: 136.982500,
			},
			Coord{
				Lat: 35.154578,
				Lon: 136.982500,
			},
		},
	}
	obstacle2 := Obstacle{
		Id: 1,
		Position: []Coord{
			Coord{
				Lat: 35.158576,
				Lon: 136.982500,
			},
			Coord{
				Lat: 35.160678,
				Lon: 136.982500,
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

func (rvoutil *RVO2Util) DecideNextTransit(nextTransit *agent.Coord, transitPoint []*agent.Coord, distance float32, destination *agent.Coord) *agent.Coord {
	// 距離が5m以下の場合
	if distance < 40 {
		if nextTransit != destination {
			for i, tPoint := range transitPoint {
				if tPoint.Lon == nextTransit.Lon && tPoint.Lat == nextTransit.Lat {
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

func (rvoutil *RVO2Util) InvScale(scaledAgents []Agent) []*agent.AgentInfo {
	// エリア情報
	areaCoord := rvoutil.SynSim.Area.AreaCoord
	height := float32(math.Abs(float64(areaCoord.EndLat - areaCoord.StartLat)))
	minLat := float32(math.Min(float64(areaCoord.EndLat), float64(areaCoord.StartLat)))
	width := float32(math.Abs(float64(areaCoord.EndLon - areaCoord.StartLon)))
	minLon := float32(math.Min(float64(areaCoord.EndLon), float64(areaCoord.StartLon)))

	// エージェント情報
	simAgents := rvoutil.SynSim.Agents

	nextAgents := make([]*agent.AgentInfo, 0)
	for i, scaledAgent := range scaledAgents {
		simAgent := simAgents[i]
		//fmt.Printf("id: %v  %v\n", simAgent.AgentId, scaledAgent.Id)
		if simAgent.AgentId == scaledAgent.Id {
			lat := scaledAgent.Coord.Lat*height + minLat
			lon := scaledAgent.Coord.Lon*width + minLon
			//velLat := scaledAgent.Velocity.Lat * height
			//velLon := scaledAgent.Velocity.Lon * width
			//goalLat := scaledAgent.Goal.Lat*height + minLat
			//goalLon := scaledAgent.Goal.Lon*width + minLon

			// calc direction
			//rad := math.Atan2(velLon, velLat)
			//nextDirection := rad * 180 / math.Pi

			nextCoord := &agent.Coord{
				Lat: lat,
				Lon: lon,
			}

			/*velocity := Coord{
				Lat: velLat,
				Lon: velLon,
			}

			goal := Coord{
				Lat: goalLat,
				Lon: goalLon,
			}*/
			route := simAgent.Route
			//currentLocation := route.Coord
			nextTransit := route.RouteInfo.NextTransit
			transitPoint := route.RouteInfo.TransitPoint
			destination := route.Destination

			// 次の経由地があればゴールを次の経由地にする
			//if nextTransit != nil {
			//	destination = nextTransit
			//}

			// calc direction, distance
			// 現在の位置とゴールとの距離と角度を求める (度, m))
			direction, distance := rvoutil.CalcDirectionAndDistance(nextCoord.Lat, nextCoord.Lon, nextTransit.Lat, nextTransit.Lon)

			// 次の経由地nextTransitを求める
			nextTransit = rvoutil.DecideNextTransit(nextTransit, transitPoint, distance, destination)
			/*if distance < 5 {
				if nextTransit != destination {
					nextTransit2 := nextTransit
					for i, tPoint := range transitPoint {
						if tPoint.Lon == nextTransit2.Lon && tPoint.Lat == nextTransit2.Lat {
							if i+1 == len(transitPoint) {
								// pass all transit point
								nextTransit = destination
							} else {
								// go to next transit point
								nextTransit = transitPoint[i+1]
							}
						}
					}
				} else {
					log.Printf("\x1b[30m\x1b[47m Arrived Destination! \x1b[0m\n")
				}

			}*/

			//fmt.Printf("\x1b[30m\x1b[47m Position %v, NextTransit: %v, Destination: %v, Direction: %v, Distance: %v \x1b[0m\n", nextCoord, nextTransit, destination, direction, distance)
			//fmt.Printf("\x1b[30m\x1b[47m 上下:  %v, 左右: %v \x1b[0m\n", nextTransit.Lat-nextCoord.Lat, nextTransit.Lon-nextCoord.Lon)
			nextRouteInfo := &agent.RouteInfo{
				TransitPoint:  route.RouteInfo.TransitPoint,
				NextTransit:   nextTransit,
				TotalDistance: route.RouteInfo.TotalDistance,
				RequiredTime:  route.RouteInfo.RequiredTime,
			}

			nextRoute := &agent.Route{
				Coord:       nextCoord,
				Direction:   float32(direction),
				Speed:       route.Speed,
				Destination: route.Destination,
				Departure:   route.Departure,
				RouteInfo:   nextRouteInfo,
			}

			nextAgent := &agent.AgentInfo{
				Time:        uint32(rvoutil.SynSim.GlobalTime) + 1,
				AgentId:     simAgent.AgentId,
				AgentType:   simAgent.AgentType,
				AgentStatus: simAgent.AgentStatus,
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
	for _, agent := range scaledAgents {
		position := &rvo.Vector2{X: float64(agent.Coord.Lon), Y: float64(agent.Coord.Lat)}
		velocity := &rvo.Vector2{X: float64(agent.Velocity.Lon), Y: float64(agent.Velocity.Lat)}
		id, _ := sim.AddAgentPosition(position)
		sim.SetAgentPrefVelocity(id, velocity)
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

func (rvoutil *RVO2Util) SetPreferredVelocities(sim *rvo.RVOSimulator, scaledAgents []Agent) []Agent {

	// エージェント情報
	//simAgents := rvoutil.SynSim.Agents

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
			Id: uint32(scaledAgents[i].Id),
			Coord: Coord{
				Lat: float32(scaledPosition.Y),
				Lon: float32(scaledPosition.X),
			},
			Goal: scaledAgents[i].Goal,
			Velocity: Coord{
				Lat: float32(scaledVelocity.Y),
				Lon: float32(scaledVelocity.X),
			},
		}

		nextScaledAgents = append(nextScaledAgents, scaledAgent)
	}
	return nextScaledAgents
}

func (rvoutil *RVO2Util) CalcNextAgents() []*agent.AgentInfo {
	timeStep := rvoutil.SynSim.TimeStep
	neighborDist := 1.5
	maxneighbors := 100
	timeHorizon := 1.5
	timeHorizonObst := 2.0
	radius := 0.03
	maxSpeed := 0.01
	sim := rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})

	scaledAgents, scaledObstacles := rvoutil.Scale()
	//fmt.Printf("scaledagents is : %v\n", scaledAgents)
	//invScaledAgents := rvoutil.InvScale(scaledAgents)

	//for _, agent := range rvoutil.SynSim.Agents{
	//	fmt.Printf("simAgents is : %v\n", agent.Route.Coord)
	//}

	//for _, agent := range invScaledAgents{
	//	fmt.Printf("invAgents is : %v\n", agent.Route.Coord)
	//}

	SetupScenario(sim, scaledAgents, scaledObstacles)

	// Stepを進める
	sim.DoStep()

	// 速度ベクトルをセットする
	nextScaledAgents := rvoutil.SetPreferredVelocities(sim, scaledAgents)
	//fmt.Printf("nextScaledAgents is : %v\n", nextScaledAgents)

	nextAgents := rvoutil.InvScale(nextScaledAgents)

	//fmt.Printf("nextAgents is : %v\n", nextAgents)

	return nextAgents
}

func (rvoutil *RVO2Util) IsAgentInControlledArea(agentInfo *agent.AgentInfo, areaInfo *area.AreaInfo, agentType int32) bool {
	lat := agentInfo.Route.Coord.Lat
	lon := agentInfo.Route.Coord.Lon
	slat := areaInfo.ControlAreaCoord.StartLat
	elat := areaInfo.ControlAreaCoord.EndLat
	slon := areaInfo.ControlAreaCoord.StartLon
	elon := areaInfo.ControlAreaCoord.EndLon
	if agentInfo.AgentType.String() == agent.AgentType_name[agentType] && slat <= lat && lat < elat && slon <= lon && lon < elon {
		return true
	}
	//log.Printf("agent type and coord is not match...\n\n")
	return false
}

func (rvoutil *RVO2Util) CalcNextAgentsByRVO() []*agent.AgentInfo {
	pureNextAgents := make([]*agent.AgentInfo, 0)
	//rvo2util := NewRVO2Util(rvoutil.SynSim)
	currentAgents := rvoutil.SynSim.Agents
	//nextAgents := rvo2util.CalcNextAgents(sim.Agents, sameAreaAgents)
	nextAgents := rvoutil.CalcNextAgents()
	for i, agentInfo := range currentAgents {
		nextAgent := nextAgents[i]
		// 自エリアにいる場合、次のルートを計算する
		if rvoutil.IsAgentInControlledArea(agentInfo, rvoutil.SynSim.Area, rvoutil.SynSim.AgentType) {

			pureNextAgent := &agent.AgentInfo{
				Time:        uint32(rvoutil.SynSim.GlobalTime) + 1,
				AgentId:     agentInfo.AgentId,
				AgentType:   agentInfo.AgentType,
				AgentStatus: agentInfo.AgentStatus,
				Route:       nextAgent.Route,
			}

			pureNextAgents = append(pureNextAgents, pureNextAgent)
		}
	}

	return pureNextAgents
}
