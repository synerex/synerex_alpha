package routes

import (
	"fmt"
	"log"

	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
)

var (
	isRVO2 bool
)

func init() {
	isRVO2 = false
}

// Finish Fix
func (sim *SynerexSimulator) CalcNextRoute(agentInfo *agent.AgentInfo, sameAreaAgents []*agent.AgentInfo) *agent.Route {

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

	direction, distance := simutil.CalcDirectionAndDistance(currentLocation.Lat, currentLocation.Lon, destination.Lat, destination.Lon)
	//newLat, newLon := simutil.CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, speed*1000/3600, direction)
	newLat, newLon := simutil.CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, destination.Lat, destination.Lon, distance, speed)

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

func calcNextAgentsBySimple(synSim *SynerexSimulator) []*agent.AgentInfo{
	pureNextAgents := make([]*agent.AgentInfo, 0)
	for _, agentInfo := range synSim.Agents {
		// 自エリアにいる場合、次のルートを計算する
		if simutil.IsAgentInControlledArea(agentInfo, synSim.Area, synSim.AgentType) {

			nextRoute := CalcNextRoute(agentInfo, sameAreaAgents)

			pureNextAgent := &agent.AgentInfo{
				Time:        uint32(synSim.GlobalTime) + 1,
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