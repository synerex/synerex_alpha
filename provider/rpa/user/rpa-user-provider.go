package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/synerex/synerex_alpha/api/rpa"

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
	rm         *rpa.MeetingService
)

func init() {
	spMap = make(map[uint64]*api.Supply)
}

func confirmBooking(clt *sxutil.SMServiceClient, sp *api.Supply) {
	spMap[sp.Id] = sp

	// emit to client
	channel, err := server.GetChannel(rm.Cid)
	if err != nil {
		fmt.Println("Failed to get socket channel:", err)
	}
	js := `{"id":"` + strconv.FormatUint(sp.Id, 10) + `","room":"` + rm.Room + `"}`
	channel.Emit("check_booking", js)

	server.On("confirm_booking", func(c *gosocketio.Channel, data interface{}) {
		switch d := data.(type) {
		case string:
			uintData, err := strconv.ParseUint(d, 0, 64)
			if err != nil {
				fmt.Println("Failed to parse uint64:", err)
			}
			clt.SelectSupply(spMap[uintData])
			channel.Emit("server_to_client", spMap[uintData])
		}
	})
}

func setMeetingService(json string) {
	cid := gjson.Get(json, "cid").String()
	status := gjson.Get(json, "status").String()
	year := gjson.Get(json, "year").String()
	month := gjson.Get(json, "month").String()
	day := gjson.Get(json, "day").String()
	week := gjson.Get(json, "week").String()
	start := gjson.Get(json, "start").String()
	end := gjson.Get(json, "end").String()
	people := gjson.Get(json, "people").String()
	title := gjson.Get(json, "title").String()
	room := gjson.Get(json, "room").String()

	rm = &rpa.MeetingService{
		Cid:    cid,
		Status: status,
		Year:   year,
		Month:  month,
		Day:    day,
		Week:   week,
		Start:  start,
		End:    end,
		People: people,
		Title:  title,
		Room:   room,
	}
}

func supplyCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	log.Println("Got RPA User supply callback")
	setMeetingService(sp.ArgJson)

	if rm.People == "" {
		rm.People = "0"
	}

	switch rm.Status {
	case "OK":
		confirmBooking(clt, sp)
	case "NG":
		// emit to client
		channel, err := server.GetChannel(rm.Cid)
		if err != nil {
			fmt.Println("Failed to get socket channel:", err)
		}
		msg := "NG from id:" + strconv.FormatUint(sp.Id, 10)
		channel.Emit("server_to_client", msg)
	default:
		fmt.Printf("Switch case of default(%s) is called\n", rm.Status)
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
		st := string(byte)
		year := gjson.Get(st, "Year").String()
		month := gjson.Get(st, "Month").String()
		day := gjson.Get(st, "Day").String()
		week := gjson.Get(st, "Week").String()
		start := gjson.Get(st, "Start").String()
		end := gjson.Get(st, "End").String()
		people := gjson.Get(st, "People").String()
		title := gjson.Get(st, "Title").String()
		rm := rpa.MeetingService{
			Cid:    c.Id(),
			Status: "checking",
			Year:   year,
			Month:  month,
			Day:    day,
			Week:   week,
			Start:  start,
			End:    end,
			People: people,
			Title:  title,
		}
		b, _ := json.Marshal(rm)
		sendDemand(sclient, "Booking meeting room", string(b))
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
