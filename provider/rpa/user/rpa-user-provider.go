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

type BookingJson struct {
	cid    string
	status string
	year   string
	month  string
	day    string
	week   string
	start  string
	end    string
	people string
	title  string
}

func init() {
	spMap = make(map[uint64]*api.Supply)
}

func parseByGjson(json string) BookingJson {
	cid := gjson.Get(json, "data.cid").String()
	status := gjson.Get(json, "flag").String()
	year := gjson.Get(json, "data.date.Year").String()
	month := gjson.Get(json, "data.date.Month").String()
	day := gjson.Get(json, "data.date.Day").String()
	week := gjson.Get(json, "data.date.Week").String()
	start := gjson.Get(json, "data.date.Start").String()
	end := gjson.Get(json, "data.date.End").String()
	people := gjson.Get(json, "data.date.People").String()
	title := gjson.Get(json, "data.date.Title").String()

	bj := BookingJson{
		cid:    cid,
		status: status,
		year:   year,
		month:  month,
		day:    day,
		week:   week,
		start:  start,
		end:    end,
		people: people,
		title:  title,
	}

	fmt.Println("parseByGjson is called:", bj)
	return bj
}

func confirmBooking(bj BookingJson, clt *sxutil.SMServiceClient, sp *api.Supply) {
	// emit to client
	channel, err := server.GetChannel(bj.cid)
	if err != nil {
		fmt.Println("Failed to get socket channel:", err)
	}
	channel.Emit("check_booking", "Are you sure to booking?")
	fmt.Printf("check_booking: Are you sure to booking? %v\n", bj)

	server.On("confirm_booking", func(c *gosocketio.Channel, data interface{}) {
		msg := ""
		if data == "yes" {
			clt.SelectSupply(sp)
			msg = "Success: " + bj.year + "/" + bj.month + "/" + bj.day + " " + bj.start + "~" + bj.end + " " + bj.title + " (" + bj.people + " 人)"
		} else {
			msg = "Stop: " + bj.year + "/" + bj.month + "/" + bj.day + " " + bj.start + "~" + bj.end + " " + bj.title + " (" + bj.people + " 人)"
		}
		channel.Emit("server_to_client", msg)
	})
}

func supplyCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	log.Println("Got RPA User supply callback and sp.ArgJson:", sp.ArgJson)

	bj := parseByGjson(sp.ArgJson)

	if bj.people == "" {
		bj.people = "0"
	}

	switch bj.status {
	case "OK":
		confirmBooking(bj, clt, sp)
	case "NG":
		// emit to client
		channel, err := server.GetChannel(bj.cid)
		if err != nil {
			fmt.Println("Failed to get socket channel:", err)
		}
		channel.Emit("server_to_client", "Invalid schedules")
	default:
		fmt.Printf("Switch case of default(%s) is called\n", bj.status)
	}
}

func subscribeSupply(client *sxutil.SMServiceClient) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
	log.Println("SMarket Server Closed?")
}

func runSocketIOServer(sclient *sxutil.SMServiceClient) {
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected %s\n", c.Id())
	})

	server.On("client_to_server", func(c *gosocketio.Channel, data interface{}) {
		log.Println("client_to_server:", data)
		byte, _ := json.Marshal(data)
		json := `{"cid":"` + c.Id() + `","status":"checking","date":` + string(byte) + `}`
		sendDemand(sclient, "Booking meeting room", json)
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected %s\n", c.Id())
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
