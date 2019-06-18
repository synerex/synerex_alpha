package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	gosocketio "github.com/mtfelian/golang-socketio"
	"github.com/tidwall/gjson"

	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	spMap      map[uint64]*api.Supply
	mu         sync.RWMutex
	port       = flag.Int("port", 8888, "RPA User Provider Listening Port")
	server     = gosocketio.NewServer()
)

func init() {
	spMap = make(map[uint64]*api.Supply)
}

func confirmBooking(cid string, date string, clt *sxutil.SMServiceClient, sp *api.Supply) {
	var msg string

	// emit to client
	channel, err := server.GetChannel(cid)
	if err != nil {
		fmt.Println("Failed to get socket channel:", err)
	}
	msg = "Are you sure to booking? " + date
	channel.Emit("check_booking", msg)

	server.On("confirm_booking", func(c *gosocketio.Channel, data interface{}) {
		if data == "yes" {
			clt.SelectSupply(sp)
			msg = "Success: " + date
		} else {
			msg = "Stop: " + date
		}
		channel.Emit("server_to_client", msg)
	})
}

func supplyCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	log.Println("Got RPA User supply callback")

	// parse JSON by gjson
	flag := gjson.Get(sp.ArgJson, "flag").String()
	cid := gjson.Get(sp.ArgJson, "data.cid").String()
	year := gjson.Get(sp.ArgJson, "data.date.Year").String()
	month := gjson.Get(sp.ArgJson, "data.date.Month").String()
	day := gjson.Get(sp.ArgJson, "data.date.Day").String()
	week := gjson.Get(sp.ArgJson, "data.date.Week").String()
	start := gjson.Get(sp.ArgJson, "data.date.Start").String()
	end := gjson.Get(sp.ArgJson, "data.date.End").String()
	people := gjson.Get(sp.ArgJson, "data.date.People").String()
	title := gjson.Get(sp.ArgJson, "data.date.Title").String()

	if people == "" {
		people = "0"
	}

	if flag == "true" {
		date := year + "/" + month + "/" + day + " " + week + " " + start + "~" + end + " " + title + " (" + people + " people)"
		confirmBooking(cid, date, clt, sp)
	} else {
		// emit to client
		channel, err := server.GetChannel(cid)
		if err != nil {
			fmt.Println("Failed to get socket channel:", err)
		}
		channel.Emit("server_to_client", "Invalid date")
	}
}

func subscribeSupply(client *sxutil.SMServiceClient) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func runSocketIOServer(sclient *sxutil.SMServiceClient) {
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected %s", c.Id())
	})

	server.On("client_to_server", func(c *gosocketio.Channel, data interface{}) {
		log.Println("client_to_server:", data)
		byte, _ := json.Marshal(data)
		json := `{"cid":"` + c.Id() + `","date":` + string(byte) + `}`
		sendDemand(sclient, "Booking meeting room", json)
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected %s", c.Id())
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	serveMux.Handle("/", http.FileServer(http.Dir("./client")))

	log.Printf("Starting Server at localhost:%d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), serveMux); err != nil {
		log.Fatal(err)
	}
}

func sendDemand(sclient *sxutil.SMServiceClient, nm string, js string) {
	opts := &sxutil.DemandOpts{Name: nm, JSON: js}
	mu.Lock()
	id := sclient.RegisterDemand(opts)
	mu.Unlock()
	log.Printf("Register meeting demand as id:%v\n", id)
}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "RPAUserProvider", false)

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
	argJson := fmt.Sprintf("{Client: RPAUser}")
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_MEETING_SERVICE, argJson)

	wg.Add(1)
	go subscribeSupply(sclient)

	wg.Add(1)
	go runSocketIOServer(sclient)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!
}
