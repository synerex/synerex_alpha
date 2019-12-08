package agent

import (
	"fmt"
	"math"

	"github.com/synerex/synerex_alpha/api/simulation/agent"
)

type Sex int

const (
	Male Sex = iota
	Female
)

type PedStatus struct {
	Age  string
	Name string
	Sex  Sex
}

type Pedestrian struct {
	*Agent // 埋め込み
	Status *PedStatus
	Route  *Route
}

func NewPedestrian() *Pedestrian {
	p := &Pedestrian{
		Agent: &Agent{},
	}
	return p
}

// エージェントがエリアの中にいるかどうか
func (p *agent.AgentInfo) IsInArea(mapCoord *MapCoord) bool {
	lat := p.Route.Position.Latitude
	lon := p.Route.Position.Longitude
	if mapCoord.SLatitude < lat && lat < mapCoord.ELatitude && mapCoord.SLongitude < lon && lon < mapCoord.ELongitude {
		return true
	} else {
		return false
	}
}

// ある座標への距離と角度を返す関数
func (p *Pedestrian) GetDirectionAndDistance(goal *Coord) (float64, float64) {
	var direction, distance float64
	r := 6378137 // equatorial radius (m)
	sLat := p.Route.Position.Latitude * math.Pi / 180
	sLon := p.Route.Position.Longitude * math.Pi / 180
	gLat := goal.Latitude * math.Pi / 180
	gLon := goal.Longitude * math.Pi / 180
	dLon := gLon - sLon
	dLat := gLat - sLat
	cLat := (sLat + gLat) / 2
	dx := float64(r) * float64(dLon) * math.Cos(float64(cLat))
	dy := float64(r) * float64(dLat)

	distance = math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	direction = 0
	if dx != 0 && dy != 0 {
		direction = math.Atan2(dy, dx) * 180 / math.Pi
	}

	return direction, distance
}

// ある座標に到着したかどうか
func (p *Pedestrian) IsReachedGoal(goal *Coord, radius float64) bool {

	_, distance := p.GetDirectionAndDistance(goal)
	// 距離がradius m以下の場合
	if distance < math.Abs(radius) {
		return true
	} else {
		return false
	}
}

// 次の目的地を決める関数
func (p *Pedestrian) DecideNextTransit() {
	// 距離が5m以下の場合
	radius := 5.0
	if p.IsReachedGoal(p.Route.NextTransit, radius) {
		if p.Route.NextTransit != p.Route.Destination {
			for i, point := range p.Route.TransitPoints {
				if point.Longitude == p.Route.NextTransit.Longitude && point.Latitude == p.Route.NextTransit.Latitude {
					if i+1 == len(p.Route.TransitPoints) {
						// すべての経由地を通った場合、nilにする
						p.Route.NextTransit = p.Route.Destination
					} else {
						// 次の経由地を設定する
						p.Route.NextTransit = p.Route.TransitPoints[i+1]
					}
				}
			}
		} else {
			fmt.Printf("\x1b[30m\x1b[47m Arrived Destination! \x1b[0m\n")
		}
	}
}

// grpc用のエージェント情報を格納する
func (p *Pedestrian) SetGrpcAgent(grpcAgent *agent.AgentInfo) {

	transitPoints := []*Coord{}
	for _, tp := range grpcAgent.Route.RouteInfo.TransitPoint {
		transitPoints = append(transitPoints, &Coord{
			Latitude:  float64(tp.Lat),
			Longitude: float64(tp.Lon),
		})
	}

	p.ID = uint64(grpcAgent.AgentId)
	p.Type = Type(grpcAgent.AgentType)
	p.Status = &PedStatus{
		Age:  grpcAgent.AgentStatus.Age,
		Sex:  Sex(int(0)),
		Name: grpcAgent.AgentStatus.Name,
	}
	p.Route = &Route{
		Position: &Coord{
			Latitude:  float64(grpcAgent.Route.Coord.Lat),
			Longitude: float64(grpcAgent.Route.Coord.Lon),
		},
		Direction: float64(grpcAgent.Route.Direction),
		Speed:     float64(grpcAgent.Route.Speed),
		Departure: &Coord{
			Latitude:  float64(grpcAgent.Route.Departure.Lat),
			Longitude: float64(grpcAgent.Route.Departure.Lon),
		},
		Destination: &Coord{
			Latitude:  float64(grpcAgent.Route.Destination.Lat),
			Longitude: float64(grpcAgent.Route.Destination.Lon),
		},
		TransitPoints: transitPoints,
		NextTransit: &Coord{
			Latitude:  float64(grpcAgent.Route.RouteInfo.NextTransit.Lat),
			Longitude: float64(grpcAgent.Route.RouteInfo.NextTransit.Lon),
		},
		TotalDistance: float64(grpcAgent.Route.RouteInfo.TotalDistance),
		RequiredTime:  float64(grpcAgent.Route.RouteInfo.RequiredTime),
	}
}

// grpc用のエージェント情報に変換して返す
func (p *Pedestrian) GetGrpcAgent() *agent.AgentInfo {

	transitPoint := []*agent.Coord{}
	for _, tp := range p.Route.TransitPoints {
		transitPoint = append(transitPoint, &agent.Coord{
			Lat: float32(tp.Latitude),
			Lon: float32(tp.Longitude),
		})
	}

	return &agent.AgentInfo{
		AgentId:   uint32(p.ID),
		AgentType: agent.AgentType(int32(p.Type)),
		AgentStatus: &agent.AgentStatus{
			Age:  p.Status.Age,
			Name: p.Status.Name,
			Sex:  "",
		},
		Route: &agent.Route{
			Coord: &agent.Coord{
				Lat: float32(p.Route.Position.Latitude),
				Lon: float32(p.Route.Position.Longitude),
			},
			Direction: float32(p.Route.Direction),
			Speed:     float32(p.Route.Speed),
			Destination: &agent.Coord{
				Lat: float32(p.Route.Destination.Latitude),
				Lon: float32(p.Route.Destination.Longitude),
			},
			Departure: &agent.Coord{
				Lat: float32(p.Route.Departure.Latitude),
				Lon: float32(p.Route.Departure.Longitude),
			},
			RouteInfo: &agent.RouteInfo{
				TransitPoint: transitPoint,
				NextTransit: &agent.Coord{
					Lat: float32(p.Route.NextTransit.Latitude),
					Lon: float32(p.Route.NextTransit.Longitude),
				},
				TotalDistance: float32(p.Route.TotalDistance),
				RequiredTime:  float32(p.Route.RequiredTime),
			},
		},
	}
}
