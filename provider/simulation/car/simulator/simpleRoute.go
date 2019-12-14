package simulator

/*type SimpleRoute struct {
	TimeStep       float64
	GlobalTime     float64
	Area           *area.Area
	Agents         []*agent.Agent
	AgentType      int64
	SameAreaAgents []*agent.Agent
}

func NewSimpleRoute(timeStep float64, globalTime float64, area *area.Area, agents []*agent.Agent, agentType int64) *SimpleRoute {
	r := &SimpleRoute{
		TimeStep:   timeStep,
		GlobalTime: globalTime,
		Area:       area,
		Agents:     agents,
		AgentType:  agentType,
	}
	return r
}

func (simple *SimpleRoute) CalcDirectionAndDistance(sLat float32, sLon float32, gLat float32, gLon float32) (float32, float32) {

	r := 6378137 // equatorial radius (m)
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

func (simple *SimpleRoute) DecideNextTransit(nextTransit *common.Coord, transitPoint []*common.Coord, distance float32, destination *agent.Coord) *agent.Coord {
	// 距離が5m以下の場合
	if distance < 5 {
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

// Finish Fix
func (simple *SimpleRoute) CalcNextRoute(agentInfo *agent.Agent, sameAreaAgents []*agent.Agent) *agent.PedRoute {

	route := agentInfo.GetPedestrian().Route
	speed := route.Speed
	currentLocation := route.Coord
	nextTransit := route.RouteInfo.NextTransit
	transitPoint := route.RouteInfo.TransitPoint
	destination := route.Destination
	// passed all transit point
	//if nextTransit != nil {
	//	destination = nextTransit
	//}

	direction, distance := simple.CalcDirectionAndDistance(currentLocation.Lat, currentLocation.Lon, nextTransit.Lat, nextTransit.Lon)
	//newLat, newLon := simutil.CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, speed*1000/3600, direction)
	newLat, newLon := simple.CalcMovedLatLon(currentLocation.Lat, currentLocation.Lon, nextTransit.Lat, nextTransit.Lon, distance, speed)

	// upate next trasit point
	nextTransit = simple.DecideNextTransit(nextTransit, transitPoint, distance, destination)


	//fmt.Printf("\x1b[30m\x1b[47m Position %v, NextTransit: %v, NextTransit: %v, Direction: %v, Distance: %v \x1b[0m\n", currentLocation, nextTransit, destination, direction, distance)
	//fmt.Printf("\x1b[30m\x1b[47m 上下:  %v, 左右: %v \x1b[0m\n", nextTransit.Lat-currentLocation.Lat, nextTransit.Lon-currentLocation.Lon)
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

func (simple *SimpleRoute) IsAgentInControlledArea(agentInfo *agent.Agent, areaInfo *area.Area, agentType int64) bool {
	lat := agentInfo.Route.Coord.Lat
	lon := agentInfo.Route.Coord.Lon
	slat := areaInfo.ControlAreaCoord.StartLat
	elat := areaInfo.ControlAreaCoord.EndLat
	slon := areaInfo.ControlAreaCoord.StartLon
	elon := areaInfo.ControlAreaCoord.EndLon
	if agentInfo.AgentType.String() == agent.AgentType_name[int32(agentType)] && slat <= lat && lat < elat && slon <= lon && lon < elon {
		return true
	}
	//log.Printf("agent type and coord is not match...\n\n")
	return false
}

func (simple *SimpleRoute) CalcNextAgentsBySimple() []*agent.Agent {
	pureNextAgents := make([]*agent.Agent, 0)
	for _, agentInfo := range simple.Agents {
		// 自エリアにいる場合、次のルートを計算する
		if simple.IsAgentInControlledArea(agentInfo, simple.Area, simple.AgentType) {

			nextRoute := simple.CalcNextRoute(agentInfo, simple.SameAreaAgents)

			pureNextAgent := &agent.AgentInfo{
				Time:        uint32(simple.GlobalTime) + 1,
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
*/
