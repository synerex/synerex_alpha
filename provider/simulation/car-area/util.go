package main

import (  
    //pb "github.com/synerex/synerex_alpha/api"
	//"github.com/synerex/synerex_alpha/sxutil"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
    //"fmt"
    //"time"
	"log"
	"math"
    //"sync"
    //"context"
)

var (
    //ch 			chan *pb.Supply
    //startSync 	bool
    //mu	sync.Mutex
)

func init() {
	//startSync = false
	//ch = make(chan *pb.Supply)
}

func isAgentInArea(agentInfo *agent.AgentInfo, data *simutil.Data, agentType int) bool{
	lat := agentInfo.Route.Coord.Lat
	lon := agentInfo.Route.Coord.Lon
	slat := data.AreaInfo.Map.Coord.StartLat
	elat := data.AreaInfo.Map.Coord.EndLat
	slon := data.AreaInfo.Map.Coord.StartLon
	elon := data.AreaInfo.Map.Coord.EndLon
	log.Printf("lat lon is %v, %v, %v, %v, %v, %v\n\n", lat, slat, elat, lon, slon, elon)
	if agentInfo.AgentType.String() == agent.AgentType_name[int32(agentType)] && slat <= lat && lat <= elat &&  slon <= lon && lon <= elon {
		return true
	}else{
		log.Printf("agent type and coord is not match...\n\n")
		return false
	}
}

func calcNextRoute(areaInfo *area.AreaInfo, agentInfo *agent.AgentInfo, otherAgentsInfo *agent.AgentsInfo) *agent.Route{

	route := agentInfo.Route
	nextCoord := &agent.Route_Coord{
		Lat: float32(float64(route.Coord.Lat) + float64(0.0001) * 1 * math.Cos(float64(route.Direction))),
		Lon: float32(float64(route.Coord.Lon) + float64(0.0001) * 1 * math.Sin(float64(route.Direction))), 
	}

	nextRoute := &agent.Route{
		Coord: nextCoord,
		Direction: route.Direction,
		Speed: route.Speed,
		Destination: float32(10),
		Departure: float32(100),
	}
	return nextRoute
}
