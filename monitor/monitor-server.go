package main

//go:generate protoc -I monitorapi --go_out=plugins=grpc:monitorapi monitorapi/monitor.proto

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/gops/agent"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	monitorpb "github.com/synerex/synerex_alpha/monitor/monitorapi"
	//	nodeapi "../nodeapi"
	"github.com/mtfelian/golang-socketio"
	"google.golang.org/grpc"
)

type monitorInfo struct {
	serv *gosocketio.Server
}

var (
	port      = flag.Int("port", 9999, "Monitor Server Listening Port")
	mesPort   = flag.Int("mesPort", 9998, "Monitor gRPC Port")
	mInfo     monitorInfo
	assetsDir http.FileSystem
)

// assetsFileHandler for static Data
func assetsFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}

	file := r.URL.Path
	//	log.Printf("Open File '%s'",file)
	if file == "/" {
		file = "/index.html"
	}
	f, err := assetsDir.Open(file)
	if err != nil {
		log.Printf("can't open file %s: %v\n", file, err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Printf("can't open file %s: %v\n", file, err)
		return
	}
	http.ServeContent(w, r, file, fi.ModTime(), f)
}

func message(serv *gosocketio.Server) {

	for {
		time.Sleep(10 * time.Second)

		log.Printf("Send from Server!")
		serv.BroadcastToAll("event", "Hello from Server")
	}

}

// SendReport send a monitor Message to each Socket.IO clients
func (s *monitorInfo) SendReport(ctx context.Context, m *monitorpb.Mes) (*monitorpb.Response, error) {
	s.serv.BroadcastToAll("event", m.GetJson())
	// broadcast always success (?)
	return &monitorpb.Response{Ok: true}, nil
}

func prepareGrpcServer(opts ...grpc.ServerOption) *grpc.Server {
	monitorServer := grpc.NewServer(opts...)
	monitorpb.RegisterMonitorServer(monitorServer, &mInfo)
	return monitorServer
}

func main() {
	// we need to moniter messages from smarket-server
	if gerr := agent.Listen(agent.Options{}); gerr != nil{
		log.Fatal(gerr)
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *mesPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	monitorServer := prepareGrpcServer(opts...)
	log.Printf("Start waiting Monitor Server at port :%d ...", *mesPort)
	go monitorServer.Serve(lis)

	currentRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	d := filepath.Join(currentRoot, "client", "build")

	assetsDir = http.Dir(d)
	log.Println("AssetDir:", assetsDir)

	server := gosocketio.NewServer()

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected %s", c.Id())
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected %s", c.Id())
	})
	server.On("node", func(c *gosocketio.Channel, param interface{}) {
		nid := param.(string)
		log.Printf("Get node query %s", nid)
	})

	serveMux := http.NewServeMux()

	serveMux.Handle("/socket.io/", server)
	serveMux.HandleFunc("/", assetsFileHandler)

	log.Printf("Starting Server at Web %d", *port)
	mInfo.serv = server
	//	go message(server)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), serveMux); err != nil {
		log.Panic(err)
	}

}
