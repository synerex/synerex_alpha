package main

import (
	//"context"
	"flag"
	"log"
	"sync"
	//"math/rand"
	"runtime"
	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	//"github.com/synerex/synerex_alpha/api/fleet"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"google.golang.org/grpc"
	"time"
	//s"encoding/json"
	"fmt"
	//"reflect"
	//"os"
	//"os/exec"
	"github.com/mtfelian/golang-socketio/transport"
	"github.com/mtfelian/golang-socketio"
	"os"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist     []uint64
	dmMap      map[uint64]*sxutil.DemandOpts
	spMap		map[uint64]*pb.Supply
	pspMap		map[uint64]*pb.Supply
	ch 			chan *pb.Supply
	startSync 	bool
	startCollectId bool
	isStop bool
	isStart bool
	isSetClock bool
	isSetArea bool
	isSetAgent bool
	isGetParticipant bool
	participantIdList  []uint64
	Order string
	mu	sync.Mutex
	sclientArea *sxutil.SMServiceClient
	sclientAgent *sxutil.SMServiceClient
	sclientClock *sxutil.SMServiceClient
	sclientParticipant *sxutil.SMServiceClient
	sioClient *gosocketio.Client
	idListByChannel *simutil.IdListByChannel
	data *Data
)

func init() {
	Order = ""
	idlist = make([]uint64, 0)
	participantIdList = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*pb.Supply)
	pspMap = make(map[uint64]*pb.Supply)
	startSync = false
	isStop = false
	isSetClock = false
	isSetAgent = false
	isSetArea = false
	isStart = false
	isGetParticipant = false
	ch = make(chan *pb.Supply)
	data = new(Data)
}

type Data struct {
	ClockInfo *clock.ClockInfo
}

// this function waits
func collectParticipantId(clt *sxutil.SMServiceClient, d int){

	for i := 0; i < d; i++{
		time.Sleep(1 * time.Second)
		log.Printf("waiting... %v",i+1)
	}

	mu.Lock()
	idListByChannel = simutil.CreateIdListByChannel(pspMap)
	mu.Unlock()
	log.Printf("finish collecting Id %v", idListByChannel)
	log.Printf("clock Id %v", idListByChannel.ClockIdList)
	log.Printf("area Id %v", idListByChannel.AreaIdList)
	log.Printf("agent Id %v", idListByChannel.AgentIdList)
	startCollectId = false
	isGetParticipant = true
	pspMap = make(map[uint64]*pb.Supply)
}

func syncProposeSupply(sp *pb.Supply, idList []uint64, callback func()){
	go func() { 
		log.Println("Send Supply")
		ch <- sp
		return
	}()
	if !startSync {
		log.Println("Start Sync")
		startSync = true

		go func(){		
		for {
			select {
				case psp := <- ch:
					log.Println("recieve ProposeSupply")
					pspMap[psp.SenderId] = psp
					log.Printf("waitidList %v %v", pspMap, idList)
			
				if simutil.IsFinishSync(pspMap, idList){
					fmt.Printf("Finish Sync\n")
					log.Println(runtime.NumGoroutine())
					startSync = false
					pspMap = make(map[uint64]*pb.Supply)
					// if you need, return response
					callback()
					return
				}
			}
		}

		}()
	}	
}

func setClock(){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else{
		firstTime := uint32(1)

		clockInfo := clock.ClockInfo{
			Time: firstTime,
			StatusType: 0, // OK
			Meta: "",
		}
		data.ClockInfo = &clockInfo
	
		clockDemand := clock.ClockDemand{
			Time: firstTime,
			DemandType: 2, // SET
			CycleNum: uint32(1),
			CycleDuration: uint32(1),
			CycleInterval: uint32(1),
			StatusType: 2, // NONE
			Meta: "",
		}
		
		nm := "setClock order by scenario"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, ClockDemand: &clockDemand}
	
		dmMap, idlist = simutil.SendDemand(sclientClock, opts, dmMap, idlist)
	}
}

func startClock(){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else if isSetAgent == false{
		fmt.Printf("Error... please order setAgent")
	}else if isSetArea == false{
		fmt.Printf("Error... please order setArea")
	}else if isSetClock == false{
		fmt.Printf("Error... please order setClock")
	}else{
		isStart = true
		// forward clock
		clockDemand := clock.ClockDemand{
			Time: data.ClockInfo.Time,
			DemandType: 0, // Forward
			CycleNum: uint32(1),
			CycleDuration: uint32(1),
			CycleInterval: uint32(1),
			StatusType: 2, // NONE
			Meta: "",
		}
		
		nm := "startClock order by scenario"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, ClockDemand: &clockDemand}
		
		dmMap, idlist = simutil.SendDemand(sclientClock, opts, dmMap, idlist)
		
		// calc next time
		nextTime := data.ClockInfo.Time + 1
		clockInfo := clock.ClockInfo{
			Time: uint32(nextTime),
			StatusType: 0, // OK
			Meta: "",
		}
		data.ClockInfo = &clockInfo
	}
}

func stopClock(){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else{
		if isStart == true{
			isStop = true
		}else{
			fmt.Printf("Clock has not been started. ")
		}
	}
}

func setArea(){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else{
		areaDemand := area.AreaDemand{
			Time: uint32(1),
			AreaId: 0, // A
			DemandType: 0, // SET
			StatusType: 2, // NONE
			Meta: "",
		}
	
		nm := "setArea order by scenario"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, AreaDemand: &areaDemand}

		dmMap, idlist = simutil.SendDemand(sclientArea, opts, dmMap, idlist)
	}
}

func setAgent(){
	if isGetParticipant == false{
		fmt.Printf("Error... please order getParticipant")
	}else{
		route := agent.Route{
			Coord: &agent.Route_Coord{
				Lat: float32(10),
				Lon: float32(10), 
			},
			Direction: float32(0),
			Speed: float32(10),
			Destination: float32(10),
			Departure: float32(100),
		}
		rule := agent.Rule{
			RuleInfo: "nothing",
		}
		agentStatus := agent.AgentStatus{
			Age: "20",
			Sex: "Male",
		}
	
		agentDemand := agent.AgentDemand{
			Time: uint32(1),
			AgentId: 1,
			AgentName: "Agent1",
			AgentType: 0,
			AgentStatus: &agentStatus,
			Route: &route,
			Rule: &rule,
			DemandType: 0, // SET
			StatusType: 2, // NONE
			Meta: "",
		}
		
		nm := "setAgent order by scenario"
		js := ""
		opts := &sxutil.DemandOpts{Name: nm, JSON: js, AgentDemand: &agentDemand}
	
		dmMap, idlist = simutil.SendDemand(sclientAgent, opts, dmMap, idlist)
	}
}

func getParticipant(){
	participantDemand := participant.ParticipantDemand {
		ClientId: uint64(sclientParticipant.ClientID),
		DemandType: 0, // GET
		StatusType: 2, // NONE
		Meta: "",
	}
	
	nm := "getParticipant order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, ParticipantDemand: &participantDemand}

	dmMap, idlist = simutil.SendDemand(sclientParticipant, opts, dmMap, idlist)
}


func callbackForSetAgent(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for set_agent callback")
	if clt.IsSupplyTarget(sp, idlist) { 

		opts :=	dmMap[sp.TargetId]
		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp
		
		callback := func (){
			fmt.Printf("Callback SetAgent! return OK to Simulation Synerex Engine")
			isSetAgent = true
		}
		agentIdList := idListByChannel.AgentIdList
		syncProposeSupply(sp, agentIdList, callback)	

	}else{
		log.Printf("This is not propose supply \n")
	}
}

func callbackForStartClock(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for start_clock callback")
	if clt.IsSupplyTarget(sp, idlist) { 

		opts :=	dmMap[sp.TargetId]
		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp
		
		callback := func (){
			fmt.Printf("Callback StartClock! move next clock")
			time.Sleep(3* time.Second)
			if isStop == false{
				startClock()
			}else{
				fmt.Printf("Stop clock!")
				isStop = false
				isStart = false
			}
		}
	
		clockIdList := idListByChannel.ClockIdList
		syncProposeSupply(sp, clockIdList, callback)

	}else{
		log.Printf("This is not propose supply \n")
	}
}

func callbackForSetClock(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for set_clock callback")
	if clt.IsSupplyTarget(sp, idlist) { 

		opts :=	dmMap[sp.TargetId]
		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp
		
		callback := func (){
			fmt.Printf("Callback SetClock! return OK to Simulation Synerex Engine")
			isSetClock = true
		}
		clockIdList := idListByChannel.ClockIdList
		syncProposeSupply(sp, clockIdList, callback)

	}else{
		log.Printf("This is not propose supply \n")
	}
}

func callbackForSetArea(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for set_area callback")
	if clt.IsSupplyTarget(sp, idlist) { 

		opts :=	dmMap[sp.TargetId]
		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp
		

		callback := func (){
			fmt.Printf("Callback SetArea! return OK to Simulation Synerex Engine")
			isSetArea = true
		}
	
		areaIdList := idListByChannel.AreaIdList
		syncProposeSupply(sp, areaIdList, callback)

	}else{
		log.Printf("This is not propose supply \n")
	}
}

func callbackForGetParticipant(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for get_participant callback")
	
	if idListByChannel != nil{
		log.Println("already corrected IdList!")
	}else{
		mu.Lock()
		if clt.IsSupplyTarget(sp, idlist) { 

			opts :=	dmMap[sp.TargetId]
			log.Printf("Got Supply for %v as '%v'",opts, sp )
			spMap[sp.SenderId] = sp
			pspMap[sp.SenderId] = sp

			if !startCollectId {
				log.Println("start selection")
				startCollectId = true
				go collectParticipantId(clt, 3)
			}
		}else{
			log.Printf("This is not propose supply \n")
		}
		mu.Unlock()
	}
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got supply callback")
	//log.Printf("supply is %v",sp)
	switch Order {
		case "GetParticipant":
			fmt.Println("getParticipant")
			callbackForGetParticipant(clt, sp)
		case "SetTime":
			fmt.Println("setClock")
			callbackForSetClock(clt, sp)
		case "SetArea":
			fmt.Println("setArea")
			callbackForSetArea(clt, sp)
		case "SetAgent":
			fmt.Println("setAgent")
			callbackForSetAgent(clt, sp)
		case "Start":
			fmt.Println("startClock")
			callbackForStartClock(clt, sp)
		default:
			fmt.Println("order is invalid")
	}
}


func main() {

	flag.Parse()

	sxutil.RegisterNodeName(*nodesrv, "ScenarioProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Scenario}")
	sclientAgent = sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE,argJson)
	sclientClock = sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE,argJson)
	sclientArea = sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE,argJson)
	sclientParticipant = sxutil.NewSMServiceClient(client, pb.ChannelType_PARTICIPANT_SERVICE,argJson)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go simutil.SubscribeSupply(sclientAgent, supplyCallback)
	go simutil.SubscribeSupply(sclientClock, supplyCallback)
	go simutil.SubscribeSupply(sclientArea, supplyCallback)
	go simutil.SubscribeSupply(sclientParticipant, supplyCallback)

	var sioErr error
	sioClient, sioErr = gosocketio.Dial("ws://localhost:9995/socket.io/?EIO=3&transport=websocket", transport.DefaultWebsocketTransport())
	if sioErr != nil {
		fmt.Println("se: Error to connect with se-daemon. You have to start se-daemon first.") //,err)
		os.Exit(1)
	}else{
		fmt.Println("se: connect OK")
	}
	sioClient.On(gosocketio.OnConnection, func(c *gosocketio.Channel,param interface{}) {
		fmt.Println("Go socket.io connected ")
		c.Emit("setCh", "Scenario")
	})

	sioClient.On("scenario", func(c *gosocketio.Channel, test *simutil.Test) {
		fmt.Printf("get order is: %v\n", test)
		Order = test.Order
		switch Order {
		case "GetParticipant":
			getParticipant()
			fmt.Println("getParticipant")
		case "SetTime":
			setClock()
			fmt.Println("setClock")
		case "SetArea":
			setArea()
			fmt.Println("setArea")
		case "SetAgent":
			setAgent()
			fmt.Println("set agent")
		case "Start":
			startClock()
			fmt.Println("start clock")
		case "Stop":
			stopClock()
			fmt.Println("stop clock")
		default:
			fmt.Println("error")
		}
		
	})

	sioClient.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel,param interface{}) {
		fmt.Println("Go socket.io disconnected ",c)
	})

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
