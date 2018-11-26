package main

// Draft code for Advertisement Service Provider

import (
	"context"
	"encoding/json"
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

type Request struct {
	Command  string    `json:"command"`
	Contents []Content `json:"contents"`
}

type Content struct {
	Type   string `json:"type"`
	Data   string `json:"data"`
	Period int    `json:"period"`
}

type Result struct {
	Command string `json:"command"`
	Results []struct {
		Data string `json:"data"`
	}
}

func init() {
	idlist = make([]uint64, 10)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
}

func msgCallback(clt *sxutil.SMServiceClient, msg *pb.MbusMsg) {
	log.Println("Got Mbus Msg callback")
	jsonStr := msg.ArgJson
	log.Println("JSON:" + jsonStr)

	jsonBytes := ([]byte)(jsonStr)
	data := new(Result)

	if err := json.Unmarshal(jsonBytes, data); err != nil {
		log.Fatalf("fail to unmarshal: %v", err)
		return
	}

	// save data

	if data.Command == "RESULTS" {
		sendEnqMsg(clt)
	}
}

func subscribeMBus(client *sxutil.SMServiceClient) {
	go func() {
		ctx := context.Background() // 必要？
		client.SubscribeMbus(ctx, msgCallback)
		// comes here if channel closed
		log.Printf("SubscribeMBus:%d", client.MbusID)

	}()
}

func sendMsg(client *sxutil.SMServiceClient, msg string) {
	log.Printf("SendMsg:%d", client.MbusID)

	m := new(pb.MbusMsg)
	m.ArgJson = msg
	ctx := context.Background() // 必要？
	client.SendMsg(ctx, m)
}

func sendAdMsg(client *sxutil.SMServiceClient) {
	var url = "url"

	content := Content{Type: "AD", Data: url, Period: 0}
	request := Request{Command: "CONTENTS", Contents: []Content{content}}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("fail to marshal: %v", err)
		return
	}

	sendMsg(client, string(jsonBytes))
}

func sendEnqMsg(client *sxutil.SMServiceClient) {
	var enq = "json:enq"

	content := Content{Type: "ENQ", Data: enq, Period: 0}
	request := Request{Command: "CONTENTS", Contents: []Content{content}}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("fail to marshal: %v", err)
		return
	}

	sendMsg(client, string(jsonBytes))
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got marketing supply callback:" + sp.GetSupplyName())

	// choice is supply for me? or not.
	if clt.IsSupplyTarget(sp, idlist) {
		// always select Supply
		clt.SelectSupply(sp)

		go subscribeMBus(clt)

		sendAdMsg(clt)
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
