package simutil

import (
	"fmt"
	"math"
	//"strconv"

	//"encoding/json"
	//"io/ioutil"
	//"log"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
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

type RVO2Util struct {
	SynSim *SynerexSimulator
	RVOAgents []Agent
}

func NewRVO2Util(synSim *SynerexSimulator ) *RVO2Util {
	r := &RVO2Util{
		SynSim: synSim,
	}
	return r
}

func (rvoutil *RVO2Util) Scale() []Agent {
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
		velLat := float32(math.Cos(float64(route.Direction)) / float64(height))
		velLon := float32(math.Sin(float64(route.Direction)) / float64(width))
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
		fmt.Printf("scaledLat: %v\n", scaledVelocity)
	}

	return scaledAgents
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
		fmt.Printf("id: %v  %v\n", simAgent.AgentId, scaledAgent.Id)
		if simAgent.AgentId == scaledAgent.Id{
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

			if nextTransit != nil {
				destination = nextTransit
			}

			// calc direction, distance
			// 現在の位置とゴールとの距離と角度を求める
			direction, distance := CalcDirectionAndDistance(nextCoord.Lat, nextCoord.Lon, destination.Lat, destination.Lon)

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
					fmt.Printf("\x1b[30m\x1b[47m Arrived Destination! \x1b[0m\n")
				}
			}

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

func SetupScenario(sim *rvo.RVOSimulator, scaledAgents []Agent) {

	for _, agent := range scaledAgents {
		position := &rvo.Vector2{X: float64(agent.Coord.Lat), Y: float64(agent.Coord.Lon)}
		velocity := &rvo.Vector2{X: float64(agent.Velocity.Lat), Y: float64(agent.Velocity.Lon)}
		id, _ := sim.AddAgentPosition(position)
		sim.SetAgentPrefVelocity(id, velocity)
		//sim.SetAgentMaxSpeed(id, float64(agent.MaxSpeed))
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
			X: float64(agent.Goal.Lat),
			Y: float64(agent.Goal.Lon),
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
				Lat: float32(scaledPosition.X),
				Lon: float32(scaledPosition.Y),
			},
			Goal: scaledAgents[i].Goal,
			Velocity: Coord{
				Lat: float32(scaledVelocity.X),
				Lon: float32(scaledVelocity.Y),
			},
		}

		nextScaledAgents = append(nextScaledAgents, scaledAgent)
	}
	return nextScaledAgents
}

func (rvoutil *RVO2Util) CalcNextAgents() []*agent.AgentInfo{
	timeStep := rvoutil.SynSim.TimeStep
	neighborDist := 1.5
	maxneighbors := 100
	timeHorizon := 1.5
	timeHorizonObst := 2.0
	radius := 0.03
	maxSpeed := 0.01
	sim := rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})


	scaledAgents := rvoutil.Scale()
	fmt.Printf("scaledagents is : %v\n", scaledAgents)
	invScaledAgents := rvoutil.InvScale(scaledAgents)

	for _, agent := range rvoutil.SynSim.Agents{
		fmt.Printf("simAgents is : %v\n", agent.Route.Coord)
	}

	for _, agent := range invScaledAgents{
		fmt.Printf("invAgents is : %v\n", agent.Route.Coord)
	}

	SetupScenario(sim, scaledAgents)

	// Stepを進める
	sim.DoStep()

	// 速度ベクトルをセットする
	nextScaledAgents := rvoutil.SetPreferredVelocities(sim, scaledAgents)
	fmt.Printf("nextScaledAgents is : %v\n", nextScaledAgents)

	nextAgents := rvoutil.InvScale(nextScaledAgents)

	fmt.Printf("nextAgents is : %v\n", nextAgents)

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

func (rvoutil *RVO2Util) CalcNextAgentsByRVO() []*agent.AgentInfo{
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