package main

// Draft code for Advertisement Service Provider

import (
	"context"
	"flag"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	smutil "github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"fmt"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist     []uint64
	dmMap      map[uint64]*smutil.DemandOpts
)

func init() {
	idlist = make([]uint64, 10)
	dmMap = make(map[uint64]*smutil.DemandOpts)
}

// callback for each Supply
func supplyCallback(clt *smutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got ad supply callback")
	// choice is supply for me? or not.
	if clt.IsSupplyTarget(sp, idlist) { //
		// always select Supply
		clt.SelectSupply(sp)
	}
}

func subscribeSupply(client *smutil.SMServiceClient) {
	// ここは goroutine!
	ctx := context.Background() // 必要？
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
}

func addDemand(sclient *smutil.SMServiceClient, nm string) {
	opts := &smutil.DemandOpts{Name: nm}
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id)
	dmMap[id] = opts
}

func main() {
	flag.Parse()
	smutil.RegisterNodeName(*nodesrv, "AdProvider", false)

	go smutil.HandleSigInt()
	smutil.RegisterDeferFunction(smutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure()) // only for draft version
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	smutil.RegisterDeferFunction(func() { conn.Close() })

	client := pb.NewSMarketClient(conn)
	argJson := fmt.Sprintf("{Client:Ad}")
	// create client wrapper
	sclient := smutil.NewSMServiceClient(client, pb.MarketType_AD_SERVICE, argJson)

	wg.Add(1)
	go subscribeSupply(sclient)
	addDemand(sclient, "Advertise video for 30s woman")
	wg.Wait()

	smutil.CallDeferFunctions() // cleanup!
}
