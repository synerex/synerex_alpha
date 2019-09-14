package simutil

import (  
    pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
    //"fmt"
    //"time"
    "log"
    "sync"
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
	//demandType := ""
	log.Printf("demandType1.5 is %v", dm)
	if(dm.GetArg_ClockDemand() != nil){
		argOneof := dm.GetArg_ClockDemand()
		log.Printf("demandType2 is %v", argOneof.DemandType.String())
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
		argOneof := dm.GetArg_AreaDemand()
		switch(argOneof.DemandType.String()){
			case "SET": return "SET_AGENT"
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

func SendProposeSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts, spMap map[uint64]*sxutil.SupplyOpts, idlist []uint64) {
	mu.Lock()
	id := sclient.ProposeSupply(opts)
	idlist = append(idlist, id) // my demand list
	spMap[id] = opts            // my demand options
	mu.Unlock()
	log.Printf("Propose my supply as id %v, %v",id,idlist)
}

func SendSupply(sclient *sxutil.SMServiceClient, opts *sxutil.SupplyOpts, spMap map[uint64]*sxutil.SupplyOpts, idlist []uint64) {
	mu.Lock()
	id := sclient.RegisterSupply(opts)
	idlist = append(idlist, id) // my demand list
	spMap[id] = opts            // my demand options
	mu.Unlock()
	log.Printf("Register my supply as id %v, %v",id,idlist)
}

func SendDemand(sclient *sxutil.SMServiceClient, opts *sxutil.DemandOpts, dmMap map[uint64]*sxutil.DemandOpts, idlist []uint64) {
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id) // my demand list
	dmMap[id] = opts            // my demand options
	mu.Unlock()
	log.Printf("Register my demand as id %v, %v",id,idlist)
}