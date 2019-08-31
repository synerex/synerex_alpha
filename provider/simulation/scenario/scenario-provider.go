package main

import (
	"context"
	"flag"
	"log"
	"sync"
	//"math/rand"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	//"github.com/synerex/synerex_alpha/api/fleet"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	//"github.com/synerex/synerex_alpha/api/simulation/agent"
	//"github.com/synerex/synerex_alpha/api/simulation/area"
	"google.golang.org/grpc"
	//"time"
	//"encoding/json"
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
	selection 	bool
	mu	sync.Mutex
	sclientArea *sxutil.SMServiceClient
	sclientAgent *sxutil.SMServiceClient
	sclientClock *sxutil.SMServiceClient
	sioClient *gosocketio.Client
)

func init() {
	idlist = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	spMap = make(map[uint64]*pb.Supply)
	selection = false
}

type LonLat struct{
	Latitude	float32
	Longitude	float32
}

type TaxiDemand struct {
	Price	int
	Distance	int
	Arrival int
	Destination int
	Position LonLat
}

/*// this function waits
func startSelection(clt *sxutil.SMServiceClient,d time.Duration){
	var sid uint64

	for i := 0; i < 5; i++{
		time.Sleep(d / 5)
		log.Printf("waiting... %v",i)
	}
	mu.Lock()
	log.Printf("From now, let's find best taxi..")
	lowest := 99999
	// find most valuable proposal
	for k, sp := range spMap {
		dat := TaxiDemand{}
		js := sp.ArgJson
		err := json.Unmarshal([]byte(js),&dat)
		if err != nil {
			log.Printf("Err JSON %v",err)
		}
		if dat.Price < lowest { // we should check some ...
			sid = k
		}
	}
	mu.Unlock()
	log.Printf("Select supply %v", spMap[sid])
	clt.SelectSupply(spMap[sid])
	// we have to cleanup all info.
}*/

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

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got supply callback")
	log.Printf("supply is %v",sp)
	switch sp.SupplyName{
		case "SEND_AREA": startUpAreaAgentProvider(clt, sp)
		case "SET_AGENT_OK": setAgentOK(clt, sp)
		case "START_UP_OK": startUpOK(clt, sp)
		case "SET_CLOCK_ALL_OK": setClockOK(clt, sp)
		case "FORWARD_CLOCK_OK": forwardClockOK(clt, sp)
		case "FORWARD_AGENT_OK": forwardAgentOK(clt, sp)
		case "FORWARD_AREA_OK": forwardAreaOK(clt, sp)
		default: log.Println("demand callback is valid.")

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

func userSelect() string{
	Options := []string{"SET_CLOCK", "START_CLOCK", "FORWARD_CLOCK", "BACK_CLOCK", "SKIP_CLOCK",
						"SET_AGENT", "SET_AREA"}
	for i, s := range Options {
		fmt.Printf("%d. %s\n", i, s)
	} 
	fmt.Print("\n命令を指定してください[int]\n")
	var selectNum int
	fmt.Scan(&selectNum)
	fmt.Printf("命令: %d, %s\n", selectNum, Options[selectNum])
	return Options[selectNum]
}

func setClock(){
	clockDemand := clock.ClockDemand{
		Time: uint32(1),
		DemandType: 2, // SET
		NumCycle: uint32(1),
		CycleDuration: uint32(1),
		CycleTime: uint32(1),
		StatusType: 2, // NONE
		Meta: "",
	}
	
	nm := "setClock order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, ClockDemand: &clockDemand}

	sendDemand(sclientClock, opts)
}

func setArea(){
	
	/*areaRequest := area.AreaService_AreaRequest{
		Time: uint32(1),
		AreaId: uint32(1),	// area a: 1, b: 2
	}
	
	areaService := area.AreaService{
		AreaRequest: &areaRequest,
	}*/
	nm := "setArea order by scenario"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js}

	sendDemand(sclientArea, opts)
}

func setAgent(){

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

	wg := sync.WaitGroup{}

	wg.Add(1)
	go subscribeSupply(sclientAgent)
	go subscribeSupply(sclientClock)
	go subscribeSupply(sclientArea)

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
		switch order {
		case "SetTime":
			setClock()
			fmt.Println("setClock")
		case "SetArea":
			setArea()
			fmt.Println("setArea")
		case "SetAgent":
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
