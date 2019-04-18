package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	spMap      map[uint64]*sxutil.SupplyOpts
	mu         sync.RWMutex
)

func init() {
	spMap = make(map[uint64]*sxutil.SupplyOpts)
}

// wait for routing demand.
func subscribeDemand(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeDemand(ctx, routingDemandCallback)
	// comes here if channel closed
	log.Printf("Server closed... on Routing provider")
}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "RPAMeetingProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client: Meeting}")
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_ROUTING_SERVICE, argJson)

	wg.Add(1)
	go subscribeDemand(sclient)
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!
}
