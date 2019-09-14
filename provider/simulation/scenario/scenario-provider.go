package main

import (
	"context"
	"flag"
	"log"
	"sync"
	//"math/rand"

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
	participantIdList  []uint64
	Order string
	mu	sync.Mutex
	sclientArea *sxutil.SMServiceClient
	sclientAgent *sxutil.SMServiceClient
	sclientClock *sxutil.SMServiceClient
	sclientParticipant *sxutil.SMServiceClient
	sioClient *gosocketio.Client
	idListByChannel *simutil.IdListByChannel
)

func init() {
	Order = ""
	idlist = make([]uint64, 0)
	participantIdList = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*pb.Supply)
	pspMap = make(map[uint64]*pb.Supply)
	startSync = false
	ch = make(chan *pb.Supply)
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
	startCollectId = false
	pspMap = make(map[uint64]*pb.Supply)
}

func callback(){
	fmt.Printf("Callback! return OK to Simulation Synerex Engine")
}

func syncProposeSupply(sp *pb.Supply, idList []uint64){
	go func() { 
		log.Println("Send Supply")
		ch <- sp
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
			
				if simutil.IsFinishSync(pspMap, idList){
					fmt.Printf("Finish Sync")
					startSync = false
					pspMap = make(map[uint64]*pb.Supply)
					// if you need, return response
					callback()
					break
				}
			}
		}
		
		}()
	}	
}


func startUpAreaAgentProvider(clt *sxutil.SMServiceClient, sp *pb.Supply){
	// start up area-agent-provider with area property
	// but now, run case provider is already running
	//sendDemand(sclientArea, "START_UP_A", "{Area:{Latitude:36.5, Longitude:135.6}}")
	//sendDemand(sclientArea, "START_UP_B", "{Area:{Latitude:40.5, Longitude:140.6}}")

}

func startUpOK(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("start up and set area ok")
	log.Printf("supply is %v",sp)
}

func setAgentOK(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("set agent ok")
	log.Printf("supply is %v",sp)
}

func setClockOK(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("set clock ok")
	log.Printf("supply is %v",sp)
}

func forwardClockOK(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("forwardClock OK")
}

func forwardAgentOK(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("forwardAgent OK")
}

func forwardAreaOK(clt *sxutil.SMServiceClient,sp *pb.Supply){
	log.Println("forwardArea OK")
}



func setClock(){
	clockDemand := clock.ClockDemand{
		Time: uint32(1),
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

	sendDemand(sclientClock, opts)
}

func setArea(){
	areaDemand := area.AreaDemand{
		Time: uint32(1),
		AreaId: 0, // B
		DemandType: 0, // SET
		StatusType: 2, // NONE
		Meta: "",
	}
	
	nm := "setArea order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, AreaDemand: &areaDemand}

	sendDemand(sclientArea, opts)
}

func setAgent(){
	route := agent.Route{
		Coord: &agent.Route_Coord{
			Lat: float32(0),
			Lon: float32(0), 
		},
		Direction: float32(0),
		Speed: float32(10),
		Destination: float32(10),
		Departure: float32(100),
	}
	rule := agent.Rule{
		RuleInfo: "nothing",
	}

	agentDemand := agent.AgentDemand{
		Time: uint32(1),
		AgentId: 1,
		AgentName: "Agent1",
		AgentType: 0,
		Route: &route,
		Rule: &rule,
		DemandType: 0, // SET
		StatusType: 2, // NONE
		Meta: "",
	}
	
	nm := "setAgent order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, AgentDemand: &agentDemand}

	sendDemand(sclientAgent, opts)
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

	sendDemand(sclientParticipant, opts)
}


func callbackForSetClock(clt *sxutil.SMServiceClient, sp *pb.Supply){
	log.Println("Got for set_clock callback")
	if clt.IsSupplyTarget(sp, idlist) { 

		opts :=	dmMap[sp.TargetId]
		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.SenderId] = sp
		

		if idListByChannel != nil{
			clockIdList := idListByChannel.ClockIdList
			syncProposeSupply(sp, clockIdList)
		}else{
			log.Printf("Error... please order getParticipantId")
		}

	}else{

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
			//log.Printf("This is not my supply id %v, %v",sp,idlist)
			// do not need to say.
		}
		mu.Unlock()
		log.Println("finish participant")
	}
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got supply callback")
	log.Printf("supply is %v",sp)
	switch Order {
		case "GetParticipant":
			callbackForGetParticipant(clt, sp)
			fmt.Println("getParticipant")
		case "SetTime":
			callbackForSetClock(clt, sp)
			fmt.Println("setClock")
		/*case "SetArea":
			setArea()
			fmt.Println("setArea")
		case "SetAgent":
			setAgent()
			fmt.Println("skip clock")
		case "Start":
			fmt.Println("set clock")
			//sendDemand(sclientClock, order, "{Date: '2019-7-29T22:32:13.234252Z'")
		case "Stop":
			fmt.Println("start clock")
			//sendDemand(sclientClock, order, "Start Clock")
		case "Forward":
			fmt.Println("set agent")
			//sendDemand(sclientAgent, order, "{Principle: {}, Position, {36.5, 138.5}, Agent: 'Pedestrian'}")
		case "Back":
			fmt.Println("set area")
			//sendDemand(sclientArea, order, "{Area:{Latitude:36.5, Longitude:135.6}}")*/
		default:
			fmt.Println("error")
	}

	
}

func subscribeSupply(client *sxutil.SMServiceClient) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func sendDemand(sclient *sxutil.SMServiceClient, opts *sxutil.DemandOpts) {
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id) // my demand list
	dmMap[id] = opts            // my demand options
	mu.Unlock()
	log.Printf("Register my demand as id %v, %v",id,idlist)
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
	go subscribeSupply(sclientAgent)
	go subscribeSupply(sclientClock)
	go subscribeSupply(sclientArea)
	go subscribeSupply(sclientParticipant)

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

	sioClient.On("scenario", func(c *gosocketio.Channel, order string) {
		fmt.Printf("get order is: %s\n", order)
		Order = order
		switch order {
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
			fmt.Println("skip clock")
		case "Start":
			fmt.Println("set clock")
			//sendDemand(sclientClock, order, "{Date: '2019-7-29T22:32:13.234252Z'")
		case "Stop":
			fmt.Println("start clock")
			//sendDemand(sclientClock, order, "Start Clock")
		case "Forward":
			fmt.Println("set agent")
			//sendDemand(sclientAgent, order, "{Principle: {}, Position, {36.5, 138.5}, Agent: 'Pedestrian'}")
		case "Back":
			fmt.Println("set area")
			//sendDemand(sclientArea, order, "{Area:{Latitude:36.5, Longitude:135.6}}")
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
