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

func demandCallback(clt *sxutil.SMServiceClient, dm *api.Demand) {
	log.Printf("Meeting demand callback is called\nId:%d, SenderId:%d, TargetId:%d, Type:%v, DemandName:%s, TimeStamp:%v, ArgJson:%v, MbusId:%d, ArgMeeting:%v\n", dm.Id, dm.SenderId, dm.TargetId, dm.Type, dm.DemandName, dm.Ts, dm.ArgJson, dm.MbusId, dm.GetArg_MeetingService())
}

func subscribeDemand(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
	log.Println("Server closed... on Meeting provider")
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
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_MEETING_SERVICE, argJson)

	wg.Add(1)
	go subscribeDemand(sclient)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!
}
