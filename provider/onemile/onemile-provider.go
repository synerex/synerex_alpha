package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mtfelian/golang-socketio"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	version = "0.01"

	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")

	client api.SynerexClient
	port   = flag.Int("port", 7777, "OneMile Provider Listening Port")
	ioserv *gosocketio.Server
)

// display
type display struct {
	dispId string // display id
	chanId string // channel id
}

// taxi/display mapping
var dispMap = make(map[string]*display)

// register OneMileProvider to NodeServer
func registerOneMileProvider() {
	sxutil.RegisterNodeName(*nodesrv, "OneMileProvider", false)
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)
	go sxutil.HandleSigInt()
}

// create SMServiceClient for a given ChannelType
func createSMServiceClient(ch api.ChannelType, arg string) *sxutil.SMServiceClient {
	// create grpc client (at onece)
	if client == nil {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())

		conn, err := grpc.Dial(*serverAddr, opts...)
		if err != nil {
			log.Fatalf("Fail to Connect Synerex Server: %v", err)
		}

		client = api.NewSynerexClient(conn)
	}

	// create SMServiceClient
	return sxutil.NewSMServiceClient(client, ch, arg)
}

// subscribe marketing channel
func subscribeMarketing(client *sxutil.SMServiceClient) {
	ctx := context.Background()
	client.SubscribeDemand(ctx, func(clt *sxutil.SMServiceClient, dm *api.Demand) {
		if dm.GetDemandName() == "" {
			log.Printf("Receive SelectSupply [id: %s, name: %s]\n", dm.GetId(), dm.GetDemandName())
			clt.Confirm(sxutil.IDType(dm.GetId()))
		} else {
			log.Printf("Receive RegisterDemand [id: %s, name: %s]\n", dm.GetId(), dm.GetDemandName())
			sp := &sxutil.SupplyOpts{
				Target: dm.GetId(),
				Name:   "onemile-provider has a display for advertising and enqueting",
			}
			clt.ProposeSupply(sp)
		}
	})
}

// run Socket.IO server for OneMile-Display-Client
func runSocketIOServer() {
	ioserv := gosocketio.NewServer()

	ioserv.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected from %s as %s\n", c.IP(), c.Id())
	})

	ioserv.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected from %s as %s\n", c.IP(), c.Id())
	})

	// register taxi and display mapping
	ioserv.On("register", func(c *gosocketio.Channel, data interface{}) {
		log.Printf("Receive register from %s [%v]\n", c.Id(), data)

		taxi := data.(map[string]interface{})["taxi"].(string)
		disp := data.(map[string]interface{})["disp"].(string)

		_, ok := dispMap[taxi]
		if !ok {
			dispMap[taxi] = &display{dispId: disp, chanId: c.Id()}
			log.Printf("Register display [taxi: %s => display: %v]\n", taxi, dispMap[taxi])
		}
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", ioserv)
	serveMux.Handle("/", http.FileServer(http.Dir("./display-client")))

	log.Printf("Starting OneMile Provider %s on port %d", version, *port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), serveMux)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	// register onemile-provider
	registerOneMileProvider()

	// subscribe marketing channel
	mktClient := createSMServiceClient(api.ChannelType_MARKETING_SERVICE, "")
	subscribeMarketing(mktClient)

	// start Websocket Server
	runSocketIOServer()
}
