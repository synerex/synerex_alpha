package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	spMap      map[uint64]*sxutil.SupplyOpts
	mu         sync.RWMutex
	port       = flag.Int("port", 8888, "RPA User Provider Listening Port")
)

func init() {
	spMap = make(map[uint64]*sxutil.SupplyOpts)
}

func routingDemandCallback(clt *sxutil.SMServiceClient, dm *api.Demand) {
	log.Println("Got routing demand callback on SRouting")
}

// wait for routing demand.
func subscribeDemand(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeDemand(ctx, routingDemandCallback)
	// comes here if channel closed
	log.Printf("Server closed... on Routing provider")
}

// run Socket.IO Server for User Interface
func runSocketIOServer() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("OnConnect:", s.ID())
		return nil
	})

	server.OnError("/", func(e error) {
		fmt.Println("OnError:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("OnDisconnect", msg)
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	serveMux.Handle("/", http.FileServer(http.Dir("./client")))

	log.Printf("Starting Socket.IO Server at localhost:%d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), serveMux); err != nil {
		log.Fatal(err)
	}
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
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_ROUTING_SERVICE, argJson)

	wg.Add(1)
	go subscribeDemand(sclient)

	wg.Add(1)
	go runSocketIOServer()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!
}