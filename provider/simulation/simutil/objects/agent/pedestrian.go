package agent

import (
	"github.com/synerex/synerex_alpha/api/simulation/agent"
)

type Sex int

const (
	Male Sex = iota
	Female
)

type Status struct {
	Age  string
	Name string
	Sex  Sex
}

type Pedestrian struct {
	*Agent // 埋め込み
	Status *Status
}

func NewPedestrian() *Pedestrian {
	p := &Pedestrian{
		Agent: &Agent{},
	}
	return p
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
	p.Status = &Status{
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
