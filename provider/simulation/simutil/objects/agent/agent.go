package agent

import (
	"fmt"
	"math"
)

type Coord struct {
	Latitude  float64
	Longitude float64
}

type Route struct {
	Position      *Coord
	Direction     float64
	Speed         float64
	Departure     *Coord
	Destination   *Coord
	TransitPoints []*Coord
	NextTransit   *Coord
	TotalDistance float64
	RequiredTime  float64
}

type Type int64

const (
	PEDESTRIAN Type = iota
	CAR
)

type Agent struct {
	ID   uint64
	Type Type
}

// エージェントがエリアの中にいるかどうか
func (a *Agent) IsInArea(mapCoord *MapCoord) bool {
	lat := a.Route.Position.Latitude
	lon := a.Route.Position.Longitude
	if mapCoord.SLatitude < lat && lat < mapCoord.ELatitude && mapCoord.SLongitude < lon && lon < mapCoord.ELongitude {
		return true
	} else {
		return false
	}
}

// ある座標への距離と角度を返す関数
func (a *Agent) GetDirectionAndDistance(goal *Coord) (float64, float64) {
	var direction, distance float64
	r := 6378137 // equatorial radius (m)
	sLat := a.Route.Position.Latitude * math.Pi / 180
	sLon := a.Route.Position.Longitude * math.Pi / 180
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
func (a *Agent) IsReachedGoal(goal *Coord, radius float64) bool {

	_, distance := a.GetDirectionAndDistance(goal)
	// 距離がradius m以下の場合
	if distance < math.Abs(radius) {
		return true
	} else {
		return false
	}
}

// 次の目的地を決める関数
func (a *Agent) DecideNextTransit() {
	// 距離が5m以下の場合
	radius := 5.0
	if a.IsReachedGoal(a.Route.NextTransit, radius) {
		if a.Route.NextTransit != a.Route.Destination {
			for i, point := range a.Route.TransitPoints {
				if point.Longitude == a.Route.NextTransit.Longitude && point.Latitude == a.Route.NextTransit.Latitude {
					if i+1 == len(a.Route.TransitPoints) {
						// すべての経由地を通った場合、nilにする
						a.Route.NextTransit = a.Route.Destination
					} else {
						// 次の経由地を設定する
						a.Route.NextTransit = a.Route.TransitPoints[i+1]
					}
				}
			}
		} else {
			fmt.Printf("\x1b[30m\x1b[47m Arrived Destination! \x1b[0m\n")
		}
	}
}
