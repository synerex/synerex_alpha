package main

// Draft code for Advertisement Service Provider

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist     []uint64
	dmMap      map[uint64]*sxutil.DemandOpts
)

func init() {
	idlist = make([]uint64, 10)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
}

func msgCallback(clt *sxutil.SMServiceClient, msg *pb.MbusMsg) {
	log.Println("Got Mbus Msg callback")
}

func subscribeMBus(client *sxutil.SMServiceClient) {
	log.Printf("SubscribeMBus:%d", client.MbusID)

	// ここは goroutine!
	ctx := context.Background() // 必要？
	client.SubscribeMbus(ctx, msgCallback)
	// comes here if channel closed
}

func sendMsg(client *sxutil.SMServiceClient, msg string) {
	log.Printf("SendMsg:%d", client.MbusID)

	m := new(pb.MbusMsg)
	m.ArgJson = msg
	ctx := context.Background() // 必要？
	client.SendMsg(ctx, m)
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got marketing supply callback:" + sp.GetSupplyName())

	// choice is supply for me? or not.
	if clt.IsSupplyTarget(sp, idlist) {
		// always select Supply
		log.Println("before SelectSupply")
		clt.SelectSupply(sp)
		log.Println("after SelectSupply")

		go subscribeMBus(clt)

		//time.Sleep(time.Second * 5)
		// send Msg
		sendMsg(clt, "json")
	}

}

func subscribeSupply(client *sxutil.SMServiceClient) {
	// ここは goroutine!
	ctx := context.Background() // 必要？
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
}

func addDemand(sclient *sxutil.SMServiceClient, nm string) {
	opts := &sxutil.DemandOpts{Name: nm}
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id)
	dmMap[id] = opts
}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "MarketingProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure()) // only for draft version
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Marketing}")
	// create client wrapper
	sclient := sxutil.NewSMServiceClient(client, pb.ChannelType_MARKETING_SERVICE, argJson)

	wg.Add(1)
	go subscribeSupply(sclient)

	addDemand(sclient, "Kota-City citizen")

	wg.Wait()

	sxutil.CallDeferFunctions() // cleanup!
}
