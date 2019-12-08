package agent

import (
	"fmt"
	"math"

	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/common"
)

type Car struct {
	*agent.Car
}

func NewCar(car *agent.Car) *Car {
	p := &Car{
		Car: car,
	}
	return p
}

// エージェントがエリアの中にいるかどうか
func (p *Car) IsInArea(mapCoord *area.AreaCoord) bool {
	lat := p.Route.Position.Latitude
	lon := p.Route.Position.Longitude
	if mapCoord.StartLat < lat && lat < mapCoord.EndLat && mapCoord.StartLon < lon && lon < mapCoord.EndLon {
		return true
	} else {
		return false
	}
}

// ある座標への距離と角度を返す関数
func (p *Car) GetDirectionAndDistance(goal *common.Coord) (float64, float64) {
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
func (p *Car) IsReachedGoal(goal *common.Coord, radius float64) bool {

	_, distance := p.GetDirectionAndDistance(goal)
	// 距離がradius m以下の場合
	if distance < math.Abs(radius) {
		return true
	} else {
		return false
	}
}

// 次の目的地を決める関数
func (p *Car) DecideNextTransit() {
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
