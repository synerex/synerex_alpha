package main

// Draft code for Taxi Display Service Provider

import (
	"context"
	"flag"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	sxutil "github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"fmt"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist     []uint64
	dmMap      map[uint64]*sxutil.SupplyOpts
)

func init() {
	idlist = make([]uint64, 10)
	dmMap = make(map[uint64]*sxutil.SupplyOpts)
}

// callback for each Demand
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	log.Println("Got ad demand callback")
	// choice is supply for me? or not.
	if clt.IsDemandTarget(dm, idlist) { //
		// always select Supply
		clt.SelectDemand(dm)
	}
}

func subscribeDemand(client *sxutil.SMServiceClient) {
	// ここは goroutine!
	ctx := context.Background() // 必要？
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
}

func addSupply(sclient *sxutil.SMServiceClient, nm string) {
	opts := &sxutil.SupplyOpts{Name: nm}
//	log.Printf("addSuply %v",*opts)
	id := sclient.RegisterSupply(opts)
	idlist = append(idlist, id)
	dmMap[id] = opts
}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "TaxiDisplayProvider", false)

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

	client := pb.NewSMarketClient(conn)
	argJson := fmt.Sprintf("{Client:TaxiDisplay}")
	// create client wrapper
	sclient := sxutil.NewSMServiceClient(client, pb.MarketType_AD_SERVICE, argJson)

	wg.Add(1)
	go subscribeDemand(sclient)
	addSupply(sclient, "Display for Ad/Entertainment")
	wg.Wait()

	sxutil.CallDeferFunctions() // cleanup!
}
