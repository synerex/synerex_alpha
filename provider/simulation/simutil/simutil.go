package simutil

import (  
    pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
    //"fmt"
    //"time"
    "log"
    "sync"
    "context"
)

var (
    ch 			chan *pb.Supply
    startSync 	bool
    mu	sync.Mutex
)

func init() {
	startSync = false
	ch = make(chan *pb.Supply)
}

type IdListByChannel struct {
	ParticipantIdList []uint64
	ClockIdList    []uint64
	AgentIdList    []uint64
	AreaIdList    []uint64
}

type Order struct{
	Type string 
	ClockInfo ClockInfo
	AreaInfo AreaInfo
	AgentInfo AgentInfo
}

type Coord struct{
	Lat float32 	`json:"lat"`
	Lon float32		`json:"lon"`
}

type Route struct{
	Coord Coord	`json:"coord"`
	Direction float32	`json:"direction"`
	Speed float32	`json:"speed"`
	Departure string	`json:"departure"`
	Destination string	`json:"destination"`
}

type Status struct{
	Name string	`json:"name"`
	Age string	`json:"age"`
	Sex string	`json:"sex"`
}

type Rule struct{

}

type ClockInfo struct{
	Time string `json:"time"`
}

type AreaInfo struct{
	Id uint32	`json:"id"`
	Name string 	`json:"name"`
}

type AgentInfo struct{
	Id uint32	`json:"id"`
	Type string		`json:"type"`
	Status Status	`json:"status"`
	Route Route	`json:"route"`
	Rule Rule		`json:"rule"`
}

type SimData struct{
	Time string 	`json:"time"`
	Area []AreaInfo	`json:"area"`
	Agent []AgentInfo	`json:"agent"`
}

func CreateIdListByChannel(pspMap map[uint64]*pb.Supply) *IdListByChannel {
    participantIdList := make([]uint64, 0)
    clockIdList := make([]uint64, 0)
    agentIdList := make([]uint64, 0)
    areaIdList := make([]uint64, 0)
    
    for _, psp := range pspMap {
        argOneof := psp.GetArg_ParticipantInfo()

        participantIdList = append(participantIdList, argOneof.ClientParticipantId)
		areaIdList = append(areaIdList, argOneof.ClientAreaId)
		agentIdList = append(agentIdList, argOneof.ClientAgentId)
        clockIdList = append(clockIdList, argOneof.ClientClockId)
    }

    i := &IdListByChannel{
		ParticipantIdList: participantIdList,
	    AgentIdList:    agentIdList,
        AreaIdList:    areaIdList,
	    ClockIdList:    clockIdList,
    }
    
	return i
}

// IsFinishSync is a helper function to check if synchronization finish or not 
func IsFinishSync(pspMap map[uint64]*pb.Supply, idlist []uint64) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			senderId := sp.SenderId
			if id == senderId{
				log.Printf("match! %v %v",id, senderId)
				isMatch = true
			}
		}
		if isMatch == false {
			log.Printf("false")
			return false
		} 
	}
	return true
}

func CheckDemandArgOneOf(dm *pb.Demand) string {
	if(dm.GetArg_ClockDemand() != nil){
		argOneof := dm.GetArg_ClockDemand()
		switch(argOneof.DemandType.String()){
			case "SET": return "SET_CLOCK"
			case "FORWARD": return "FORWARD_CLOCK"
			case "STOP": return "STOP_CLOCK"
			case "BACK": return "BACK_CLOCK"
			case "START": return "START_CLOCK"
		}
	}
	if(dm.GetArg_AreaDemand() != nil){
		argOneof := dm.GetArg_AreaDemand()
		switch(argOneof.DemandType.String()){
			case "SET": return "SET_AREA"
			case "GET": return "GET_AREA"
		}
	}
	if(dm.GetArg_AgentDemand() != nil){
		argOneof := dm.GetArg_AgentDemand()
		switch(argOneof.DemandType.String()){
			case "SET": return "SET_AGENT"
		}
	}
	if(dm.GetArg_AgentsDemand() != nil){
		argOneof := dm.GetArg_AgentsDemand()
		switch(argOneof.DemandType.String()){
			case "GET": return "GET_AGENTS"
		}
    }
	if(dm.GetArg_ParticipantDemand() != nil){
		argOneof := dm.GetArg_ParticipantDemand()
		switch(argOneof.DemandType.String()){
			case "GET": return "GET_PARTICIPANT"
		}
	}
	return "INVALID_TYPE"
}

func CheckSupplyArgOneOf(sp *pb.Supply) string {
	if(sp.GetArg_ClockInfo() != nil){
		argOneof := sp.GetArg_ClockInfo()
		switch(argOneof.SupplyType.String()){
			case "RES_SET": return "RES_SET_CLOCK"
			case "RES_FORWARD": return "RES_FORWARD_CLOCK"
			case "RES_STOP": return "RES_STOP_CLOCK"
			case "RES_BACK": return "RES_BACK_CLOCK"
			case "RES_START": return "RES_START_CLOCK"
		}
	}
	if(sp.GetArg_AreaInfo() != nil){
		argOneof := sp.GetArg_AreaInfo()
		switch(argOneof.SupplyType.String()){
			case "RES_SET": return "RES_SET_AREA"
			case "RES_GET": return "RES_GET_AREA"
		}
	}
	if(sp.GetArg_AgentInfo() != nil){
		argOneof := sp.GetArg_AgentInfo()
		switch(argOneof.SupplyType.String()){
			case "RES_SET": return "RES_SET_AGENT"
		}
    }
    if(sp.GetArg_AgentsInfo() != nil){
		argOneof := sp.GetArg_AgentsInfo()
		switch(argOneof.SupplyType.String()){
			case "RES_SET": return "RES_SET_AGENTS"
		}
	}
	if(sp.GetArg_ParticipantInfo() != nil){
		argOneof := sp.GetArg_ParticipantInfo()
		switch(argOneof.SupplyType.String()){
			case "RES_GET": return "RES_GET_PARTICIPANT"
		}
	}
	return "INVALID_TYPE"
}

func SendProposeSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts, spMap map[uint64]*sxutil.SupplyOpts, idlist []uint64) (map[uint64]*sxutil.SupplyOpts, []uint64){
	mu.Lock()
	id := sclient.ProposeSupply(opts)
	idlist = append(idlist, id) // my demand list
	spMap[id] = opts            // my demand options
	mu.Unlock()
//    log.Printf("Propose my supply as id %v, %v",id,idlist)
    return spMap, idlist
}

func SendSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts, spMap map[uint64]*sxutil.SupplyOpts, idlist []uint64) (map[uint64]*sxutil.SupplyOpts, []uint64){
	mu.Lock()
	id := sclient.RegisterSupply(opts)
	idlist = append(idlist, id) // my demand list
	spMap[id] = opts            // my demand options
	mu.Unlock()
//    log.Printf("Register my supply as id %v, %v",id,idlist)
    return spMap, idlist
}

func SendDemand(sclient *sxutil.SMServiceClient, opts *sxutil.DemandOpts, dmMap map[uint64]*sxutil.DemandOpts, idlist []uint64) (map[uint64]*sxutil.DemandOpts, []uint64){
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