package simutil

import (
	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/sxutil"

	//	"github.com/synerex/synerex_alpha/api/simulation/route"
	"fmt"
	//"time"
	"context"
	"log"
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
	Departure   string  `json:"departure"`
	Destination string  `json:"destination"`
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
			Coord: &agent.Route_Coord{
				Lat: float32(agentInfo.Route.Coord.Lat),
				Lon: float32(agentInfo.Route.Coord.Lon),
			},
			Direction:   float32(agentInfo.Route.Direction),
			Speed:       float32(agentInfo.Route.Speed),
			Destination: float32(10),
			Departure:   float32(100),
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

func SyncProposeSupply(sp *pb.Supply, syncIdList []uint32, pspMap map[uint64]*pb.Supply, callback func(pspMap map[uint64]*pb.Supply)) {
	go func() {
		log.Println("Send Supply")
		ch <- sp
		return
	}()
	log.Printf("StartSync : %v", startSync)
	if !startSync {
		log.Println("Start Sync")
		startSync = true

		go func() {
			for {
				select {
				case psp := <-ch:
					log.Println("recieve ProposeSupply")
					pspMap[psp.SenderId] = psp
					//					log.Printf("waitidList %v %v", pspMap, idList)

					if IsFinishSync(pspMap, syncIdList) {
						fmt.Printf("Finish Sync\n")
						// init pspMap
						pspMap = make(map[uint64]*pb.Supply)
						startSync = false
						fmt.Printf("startSync to false: %v\n", startSync)

						// if you need, return response
						callback(pspMap)

						return
					}
				}
			}

		}()
	}
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
	if dm.GetArg_GetRouteDemand() != nil {
		return "GET_ROUTE_DEMAND"
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
	if sp.GetArg_GetRouteSupply() != nil {
		return "GET_ROUTE_SUPPLY"
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

func SubscribeSupply(client *sxutil.SMServiceClient, supplyCallback func(*sxutil.SMServiceClient, *pb.Supply)) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func SubscribeDemand(client *sxutil.SMServiceClient, demandCallback func(*sxutil.SMServiceClient, *pb.Demand)) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}
