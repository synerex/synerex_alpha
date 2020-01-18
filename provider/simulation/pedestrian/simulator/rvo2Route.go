package simulator

import (
	"fmt"
	"math"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	//monitor "github.com/RuiHirano/rvo2-go/monitor"

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	"io/ioutil"
	"log"

)

var (
	sim *rvo.RVOSimulator
	fcs *geojson.FeatureCollection
)

func loadGeoJson(fname string) *geojson.FeatureCollection{

	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Print("Can't read file:", err)
		panic("load json")
	}
	fc, _ := geojson.UnmarshalFeatureCollection(bytes)

	return fc
}

type RVO2Route struct {
	TimeStep   float64
	GlobalTime float64
	Area       *area.Area
	Agents     []*agent.Agent
	AgentType  common.AgentType
}

func NewRVO2Route(timeStep float64, globalTime float64, area *area.Area, agentsInfo []*agent.Agent, agentType common.AgentType) *RVO2Route {

	// set obstacle
	fcs = loadGeoJson("higashiyama.geojson")
	
	r := &RVO2Route{
		TimeStep:   timeStep,
		GlobalTime: globalTime,
		Area:       area,
		Agents:     agentsInfo,
		AgentType:  agentType,
	}
	return r
}

// CalcDirectionAndDistance: 目的地までの距離と角度を計算する関数
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


// DecideNextTransit: 次の経由地を求める関数
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

// SetupScenario: Scenarioを設定する関数
func (rvo2route *RVO2Route) SetupScenario() {

	// Set Agent
	for _, agent := range rvo2route.Agents {
		ped := agent.GetPedestrian()

		position := &rvo.Vector2{X: ped.Route.Position.Longitude, Y: ped.Route.Position.Latitude}
		goal := &rvo.Vector2{X: ped.Route.NextTransit.Longitude, Y: ped.Route.NextTransit.Latitude}

		// Agentを追加
		id, _ := sim.AddDefaultAgent(position)

		// 目的地を設定
		sim.SetAgentGoal(id, goal)
		
		// エージェントの速度方向ベクトルを設定
		goalVector := sim.GetAgentGoalVector(id)
		sim.SetAgentPrefVelocity(id, goalVector)
		//sim.SetAgentMaxSpeed(id, float64(agent.MaxSpeed))
	}

	// Set Obstacle
	for _, feature := range fcs.Features {
		multiPosition := feature.Geometry.(orb.MultiLineString)[0]
		//fmt.Printf("geometry: ", multiPosition)
		rvoObstacle := []*rvo.Vector2{}

		for _, positionArray := range multiPosition{
				position := &rvo.Vector2{
					X: positionArray[0],
					Y: positionArray[1],
				}

				rvoObstacle = append(rvoObstacle, position)
		}

		sim.AddObstacle(rvoObstacle)
	}

	sim.ProcessObstacles()

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation..\n\n")
}

func (rvo2route *RVO2Route) CalcNextAgents() []*agent.Agent {

	nextControlAgents := make([]*agent.Agent, 0)
	currentAgents := rvo2route.Agents

	timeStep := rvo2route.TimeStep
	neighborDist := 0.00008 // どのくらいの距離の相手をNeighborと認識するか?Neighborとの距離をどのくらいに保つか？ぶつかったと認識する距離？
	maxneighbors := 10   // 周り何体を計算対象とするか
	timeHorizon := 1.0
	timeHorizonObst := 1.0
	radius := 0.00001  // エージェントの半径
	maxSpeed := 0.00004 // エージェントの最大スピード
	sim = rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})

	// scenario設定
	rvo2route.SetupScenario()

	// Stepを進める
	sim.DoStep()

	// 管理エリアのエージェントのみを抽出
	for rvoId, agentInfo := range currentAgents {
		//nextRVOAgent := sim.GetAgent(int(agentInfo.Id))
		currentPedInfo := agentInfo.GetPedestrian()
		// 計算する前に自エリアにいる場合、次のルートを計算する
		if rvo2route.IsAgentInControlArea(agentInfo) {
			destination := currentPedInfo.Route.Destination

			// rvoの位置情報を緯度経度に変換する
			rvoAgentPosition := sim.GetAgentPosition(int(rvoId))

			nextCoord := &common.Coord{
				Latitude:  rvoAgentPosition.Y,
				Longitude: rvoAgentPosition.X,
			}

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

// IsAgentInControlArea: エージェントが管理エリアにいるかどうか
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
