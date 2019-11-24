package simutil

import (
	//"fmt"
	"log"
	"math"

	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	//"github.com/synerex/synerex_alpha/api/simulation/participant"
	//"github.com/synerex/synerex_alpha/provider/simulation/simutil"
)

type SimpleRoute struct {
	SynSim *SynerexSimulator
	SameAreaAgents []*agent.AgentInfo
	RVOAgents []Agent
}

func NewSimpleRoute(synSim *SynerexSimulator, sameAreaAgents []*agent.AgentInfo) *SimpleRoute {
	r := &SimpleRoute{
		SynSim: synSim,
		SameAreaAgents: sameAreaAgents,
	}
	return r
}

func (simple *SimpleRoute) CalcDirectionAndDistance(sLat float32, sLon float32, gLat float32, gLon float32) (float32, float32) {

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

// TODO: Why Calc Error ? newLat=nan and newLon = inf
func (simple *SimpleRoute) CalcMovedLatLon(sLat float32, sLon float32, gLat float32, gLon float32, distance float32, speed float32) (float32, float32) {

	//r := float64(6378137) // equatorial radius

	// 割合
	x := speed * 1000 / 3600 / distance

	newLat := sLat + (gLat-sLat)*x
	newLon := sLon + (gLon-sLon)*x

	return newLat, newLon
}

// Finish Fix
func (simple *SimpleRoute) CalcNextRoute(agentInfo *agent.AgentInfo, sameAreaAgents []*agent.AgentInfo) *agent.Route {

	route := agentInfo.Route
	speed := route.Speed
	currentLocation := route.Coord
	nextTransit := route.RouteInfo.NextTransit
	transitPoint := route.RouteInfo.TransitPoint
	destination := route.Destination
	// passed all transit point
	if nextTransit != nil {
		destination = nextTransit
	}

	direction, distance := simple.CalcDirectionAndDistance(currentLocation.Lat, currentLocation.Lon, destination.Lat, destination.Lon)
	//newLat, newLon := simutil.CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, speed*1000/3600, direction)
	newLat, newLon := simple.CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, destination.Lat, destination.Lon, distance, speed)

	// upate next trasit point
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
			log.Printf("\x1b[30m\x1b[47m Arrived Destination! \x1b[0m\n")
		}

	}

	nextCoord := &agent.Coord{
		Lat: currentLocation.Lat,
		Lon: currentLocation.Lon,
	}
	//TODO: Fix this
	if newLat < 40 && newLat > 0 && newLon < 150 && newLon > 0 {
		nextCoord = &agent.Coord{
			Lat: newLat,
			Lon: newLon,
		}
	} else {
		log.Printf("\x1b[30m\x1b[47m LOCATION CULC ERROR %v \x1b[0m\n", nextCoord)

	}

	routeInfo := &agent.RouteInfo{
		TransitPoint:  transitPoint,
		NextTransit:   nextTransit,
		TotalDistance: route.RouteInfo.TotalDistance,
		RequiredTime:  route.RouteInfo.RequiredTime,
	}
	nextRoute := &agent.Route{
		Coord:       nextCoord,
		Direction:   float32(direction),
		Speed:       float32(speed),
		Destination: route.Destination,
		Departure:   route.Departure,
		RouteInfo:   routeInfo,
	}
	return nextRoute
}

func (simple *SimpleRoute) IsAgentInControlledArea(agentInfo *agent.AgentInfo, areaInfo *area.AreaInfo, agentType int32) bool {
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

func (simple *SimpleRoute) CalcNextAgentsBySimple() []*agent.AgentInfo{
	pureNextAgents := make([]*agent.AgentInfo, 0)
	for _, agentInfo := range simple.SynSim.Agents {
		// 自エリアにいる場合、次のルートを計算する
		if simple.IsAgentInControlledArea(agentInfo, simple.SynSim.Area, simple.SynSim.AgentType) {

			nextRoute := simple.CalcNextRoute(agentInfo, simple.SameAreaAgents)

			pureNextAgent := &agent.AgentInfo{
				Time:        uint32(simple.SynSim.GlobalTime) + 1,
				AgentId:     agentInfo.AgentId,
				AgentType:   agentInfo.AgentType,
				AgentStatus: agentInfo.AgentStatus,
				Route:       nextRoute,
			}

			pureNextAgents = append(pureNextAgents, pureNextAgent)
		}
	}
	return pureNextAgents
}