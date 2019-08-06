package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"math/rand"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"time"
	"encoding/json"
	"fmt"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist     []uint64
	dmMap      map[uint64]*sxutil.DemandOpts
	spMap		map[uint64]*pb.Supply
	selection 	bool
	mu	sync.Mutex
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

// this function waits
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
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got Ride_Share supply callback")
	// choice is supply for me? or not.
	mu.Lock()
	if clt.IsSupplyTarget(sp, idlist) { //
		// always select Supply
		// this is not good..
//		clt.SelectSupply(sp)
		// just show the supply Information
		opts :=	dmMap[sp.TargetId]
		log.Printf("Got Supply for %v as '%v'",opts, sp )
		spMap[sp.TargetId] = sp
		// should wait several seconds to find the best proposal.
		// if there is no selection .. lets start
		if !selection {
			selection = true
			go startSelection(clt, time.Second*5)
		}
	}else{
//		log.Printf("This is not my supply id %v, %v",sp,idlist)
		// do not need to say.
	}
	mu.Unlock()
}

func subscribeSupply(client *sxutil.SMServiceClient) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func sendDemand(sclient *sxutil.SMServiceClient, nm string, js string) {
	opts := &sxutil.DemandOpts{Name: nm, JSON: js}
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id) // my demand list
	dmMap[id] = opts            // my demand options
	mu.Unlock()
	log.Printf("Register my demand as id %v, %v",id,idlist)
}

func main() {
	flag.Parse()

	sxutil.RegisterNodeName(*nodesrv, "UserProvider", false)

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
	argJson := fmt.Sprintf("{Client:User}")
	sclient := sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE,argJson)
	sclient2 := sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE,argJson)
	sclient3 := sxutil.NewSMServiceClient(client, pb.ChannelType_AREA_SERVICE,argJson)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go subscribeSupply(sclient)
	go subscribeSupply(sclient2)
	go subscribeSupply(sclient3)

	for {
		sendDemand(sclient, "Share Ride to Home", "{Destination:{Latitude:36.5, Longitude:135.6}, Duration: 1200}")
		time.Sleep(time.Second * time.Duration(10 + rand.Int()%10))
	}
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
