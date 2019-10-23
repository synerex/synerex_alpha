package simutil

import (
	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/sxutil"

	//	"github.com/synerex/synerex_alpha/api/simulation/route"

	//"time"
	"context"
	"log"
	"math"
	"sync"
)

var (
	ch        chan *pb.Supply
	startSync bool
	mu        sync.Mutex
)

func init() {
	startSync = false
	ch = make(chan *pb.Supply)
}

type Data struct {
	AreaInfo   *area.AreaInfo
	ClockInfo  *clock.ClockInfo
	AgentsInfo []*agent.AgentInfo
}

type History struct {
	CurrentTime uint32
	History     map[uint32]*Data
}

type IdListByChannel struct {
	ParticipantIdList []uint32
	ClockIdList       []uint32
	AgentIdList       []uint32
	AreaIdList        []uint32
	RouteIdList       []uint32
}

type Order struct {
	Type       string
	ClockInfo  ClockInfo
	AreaInfo   AreaInfo
	AgentsInfo []AgentInfo
}

type Coord struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type Route struct {
	Coord       Coord   `json:"coord"`
	Direction   float32 `json:"direction"`
	Speed       float32 `json:"speed"`
	Departure   Coord   `json:"departure"`
	Destination Coord   `json:"destination"`
}

type Status struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	Sex  string `json:"sex"`
}

type Rule struct {
}

type ClockInfo struct {
	Time string `json:"time"`
}

type AreaInfo struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
}

type AgentInfo struct {
	Id     uint32 `json:"id"`
	Type   string `json:"type"`
	Status Status `json:"status"`
	Route  Route  `json:"route"`
	Rule   Rule   `json:"rule"`
}

type SimData struct {
	Time  string      `json:"time"`
	Area  []AreaInfo  `json:"area"`
	Agent []AgentInfo `json:"agent"`
}

func ConvertAgentsInfo(agentsInfo2 []AgentInfo) []*agent.AgentInfo {
	agentsInfo := make([]*agent.AgentInfo, 0)
	for _, agentInfo := range agentsInfo2 {
		route := &agent.Route{
			Coord: &agent.Coord{
				Lat: float32(agentInfo.Route.Coord.Lat),
				Lon: float32(agentInfo.Route.Coord.Lon),
			},
			Direction: float32(agentInfo.Route.Direction),
			Speed:     float32(agentInfo.Route.Speed),
			Destination: &agent.Coord{
				Lat: float32(agentInfo.Route.Destination.Lat),
				Lon: float32(agentInfo.Route.Destination.Lon),
			},
			Departure: &agent.Coord{
				Lat: float32(agentInfo.Route.Departure.Lat),
				Lon: float32(agentInfo.Route.Departure.Lon),
			},
		}
		agentStatus := &agent.AgentStatus{
			Name: "Rui",
			Age:  "20",
			Sex:  "Male",
		}
		agentType := agent.AgentType_PEDESTRIAN
		if agentInfo.Type == "car" {
			agentType = agent.AgentType_CAR
		}

		agentInfo := &agent.AgentInfo{
			Time:        uint32(1),
			AgentId:     uint32(agentInfo.Id),
			AgentStatus: agentStatus,
			AgentType:   agentType,
			Route:       route,
		}
		agentsInfo = append(agentsInfo, agentInfo)
	}
	return agentsInfo
}

func CalcDirectionAndDistance(sLat float32, sLon float32, gLat float32, gLon float32) (float32, float32) {

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

/*func CalcMovedLatLon(sLat float32, sLon float32, distance float32, direction float32) (float32, float32) {

	r := float64(6378137) // equatorial radius
	// 緯線、経線上の移動距離
	latDistance := float64(distance) * math.Cos(float64(direction)*math.Pi/180)
	lonDistance := float64(distance) * math.Sin(float64(direction)*math.Pi/180)

	// 1mあたりの緯度
	latEarthCircle := 2 * math.Pi * r
	latParMeter := 360 / latEarthCircle

	// 緯度の変化量
	dLat := latDistance * latParMeter
	newLat := sLat + float32(dLat)

	// 1mあたりの経度
	lonEarthRadius := r * math.Cos(float64(newLat)*math.Pi/100)
	lonEarthCircle := 2 * math.Pi * lonEarthRadius
	lonPerMeter := 360 / lonEarthCircle

	// 経度の変化量
	dLon := lonDistance * lonPerMeter
	newLon := sLon + float32(dLon)

	return newLat, newLon
}*/

// TODO: Why Calc Error ? newLat=nan and newLon = inf
func CalcMovedLatLon(sLat float32, sLon float32, gLat float32, gLon float32, distance float32, speed float32) (float32, float32) {

	//r := float64(6378137) // equatorial radius

	// 割合
	x := speed * 1000 / 3600 / distance

	newLat := sLat + (gLat-sLat)*x
	newLon := sLon + (gLon-sLon)*x

	return newLat, newLon
}

func CheckFinishSync(pspMap map[uint64]*pb.Supply, idlist []uint32) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			senderId := uint32(sp.SenderId)
			if id == senderId {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}

// Finish Fix
// if agent type and coord satisfy, return true
func IsAgentInArea(agentInfo *agent.AgentInfo, areaInfo *area.AreaInfo, agentType int32) bool {
	lat := agentInfo.Route.Coord.Lat
	lon := agentInfo.Route.Coord.Lon
	slat := areaInfo.AreaCoord.StartLat
	elat := areaInfo.AreaCoord.EndLat
	slon := areaInfo.AreaCoord.StartLon
	elon := areaInfo.AreaCoord.EndLon
	if agentInfo.AgentType.String() == agent.AgentType_name[agentType] && slat <= lat && lat < elat && slon <= lon && lon < elon {
		return true
	} else {
		//log.Printf("agent type and coord is not match...\n\n")
		return false
	}
}

// Fix now
// if agent type and coord satisfy, return true
func IsAgentInControlledArea(agentInfo *agent.AgentInfo, areaInfo *area.AreaInfo, agentType int32) bool {
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

func CreateIdListByChannel(participantsInfo []*participant.ParticipantInfo) *IdListByChannel {
	participantIdList := make([]uint32, 0)
	clockIdList := make([]uint32, 0)
	agentIdList := make([]uint32, 0)
	areaIdList := make([]uint32, 0)
	routeIdList := make([]uint32, 0)

	for _, participantInfo := range participantsInfo {
		channelId := participantInfo.ChannelId

		participantIdList = append(participantIdList, channelId.ParticipantChannelId)
		areaIdList = append(areaIdList, channelId.AreaChannelId)
		agentIdList = append(agentIdList, channelId.AgentChannelId)
		clockIdList = append(clockIdList, channelId.ClockChannelId)
		routeIdList = append(routeIdList, channelId.RouteChannelId)
	}

	i := &IdListByChannel{
		ParticipantIdList: participantIdList,
		AgentIdList:       agentIdList,
		AreaIdList:        areaIdList,
		ClockIdList:       clockIdList,
		RouteIdList:       routeIdList,
	}

	return i
}

func CreateParticipantsInfo(pspMap map[uint64]*pb.Supply) []*participant.ParticipantInfo {
	participantsInfo := make([]*participant.ParticipantInfo, 0)

	for _, psp := range pspMap {
		getParticipantSupply := psp.GetArg_GetParticipantSupply()
		participantInfo := getParticipantSupply.ParticipantInfo
		participantsInfo = append(participantsInfo, participantInfo)
	}

	return participantsInfo
}

// IsFinishSync is a helper function to check if synchronization finish or not
func IsFinishSync(pspMap map[uint64]*pb.Supply, idlist []uint32) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			senderId := uint32(sp.SenderId)
			if id == senderId {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}

// IsSupplyTarget is a helper function to check target
func IsSupplyTarget(sp *pb.Supply, idlist []uint64) bool {
	spid := sp.TargetId
	for _, id := range idlist {
		if id == spid {
			return true
		}
	}
	return false
}

func CheckDemandType(dm *pb.Demand) string {
	// clock
	if dm.GetArg_SetClockDemand() != nil {
		return "SET_CLOCK_DEMAND"
	}
	if dm.GetArg_ForwardClockDemand() != nil {
		return "FORWARD_CLOCK_DEMAND"
	}
	if dm.GetArg_BackClockDemand() != nil {
		return "BACK_CLOCK_DEMAND"
	}
	// area
	if dm.GetArg_GetAreaDemand() != nil {
		return "GET_AREA_DEMAND"
	}
	// agents
	if dm.GetArg_GetAgentsDemand() != nil {
		return "GET_AGENTS_DEMAND"
	}
	if dm.GetArg_SetAgentsDemand() != nil {
		return "SET_AGENTS_DEMAND"
	}
	// participant
	if dm.GetArg_GetParticipantDemand() != nil {
		return "GET_PARTICIPANT_DEMAND"
	}
	if dm.GetArg_SetParticipantDemand() != nil {
		return "SET_PARTICIPANT_DEMAND"
	}
	// route
	if dm.GetArg_GetAgentRouteDemand() != nil {
		return "GET_AGENT_ROUTE_DEMAND"
	}
	if dm.GetArg_GetAgentsRouteDemand() != nil {
		return "GET_AGENTS_ROUTE_DEMAND"
	}

	return "INVALID_TYPE"
}

func CheckSupplyType(sp *pb.Supply) string {
	// clock
	if sp.GetArg_SetClockSupply() != nil {
		return "SET_CLOCK_SUPPLY"
	}
	if sp.GetArg_ForwardClockSupply() != nil {
		return "FORWARD_CLOCK_SUPPLY"
	}
	if sp.GetArg_BackClockSupply() != nil {
		return "BACK_CLOCK_SUPPLY"
	}
	// area
	if sp.GetArg_GetAreaSupply() != nil {
		return "GET_AREA_SUPPLY"
	}
	// agents
	if sp.GetArg_GetAgentsSupply() != nil {
		return "GET_AGENTS_SUPPLY"
	}
	if sp.GetArg_SetAgentsSupply() != nil {
		return "SET_AGENTS_SUPPLY"
	}
	if sp.GetArg_ForwardAgentsSupply() != nil {
		return "FORWARD_AGENTS_SUPPLY"
	}
	// participant
	if sp.GetArg_GetParticipantSupply() != nil {
		return "GET_PARTICIPANT_SUPPLY"
	}
	if sp.GetArg_SetParticipantSupply() != nil {
		return "SET_PARTICIPANT_SUPPLY"
	}
	// route
	if sp.GetArg_GetAgentRouteSupply() != nil {
		return "GET_AGENT_ROUTE_SUPPLY"
	}
	if sp.GetArg_GetAgentsRouteSupply() != nil {
		return "GET_AGENTS_ROUTE_SUPPLY"
	}

	return "INVALID_TYPE"
}

func SendProposeSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts, spMap map[uint64]*sxutil.SupplyOpts, idlist []uint64) (map[uint64]*sxutil.SupplyOpts, []uint64) {
	mu.Lock()
	id := sclient.ProposeSupply(opts)
	idlist = append(idlist, id) // my demand list
	spMap[id] = opts            // my demand options
	mu.Unlock()
	//    log.Printf("Propose my supply as id %v, %v",id,idlist)
	return spMap, idlist
}

func SendSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts, spMap map[uint64]*sxutil.SupplyOpts, idlist []uint64) (map[uint64]*sxutil.SupplyOpts, []uint64) {
	mu.Lock()
	id := sclient.RegisterSupply(opts)
	idlist = append(idlist, id) // my demand list
	spMap[id] = opts            // my demand options
	mu.Unlock()
	//    log.Printf("Register my supply as id %v, %v",id,idlist)
	return spMap, idlist
}

func SendDemand(sclient *sxutil.SMServiceClient, opts *sxutil.DemandOpts, dmMap map[uint64]*sxutil.DemandOpts, idlist []uint64) (map[uint64]*sxutil.DemandOpts, []uint64) {
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id) // my demand list
	dmMap[id] = opts            // my demand options
	mu.Unlock()
	return dmMap, idlist
}

func SubscribeSupply(client *sxutil.SMServiceClient, supplyCallback func(*sxutil.SMServiceClient, *pb.Supply), wg *sync.WaitGroup) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback, wg)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func SubscribeDemand(client *sxutil.SMServiceClient, demandCallback func(*sxutil.SMServiceClient, *pb.Demand), wg *sync.WaitGroup) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback, wg)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}
