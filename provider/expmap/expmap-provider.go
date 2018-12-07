package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mtfelian/golang-socketio"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// map provider provides Kota-map information to Web Service through socket.io.

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	port       = flag.Int("port", 10080, "Map Provider Listening Port")
	mu         sync.Mutex
	version    = "0.01"
	assetsDir  http.FileSystem
	ioserv     *gosocketio.Server
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

func run_server() *gosocketio.Server {

	currentRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	d := filepath.Join(currentRoot, "mclient", "build")

	assetsDir = http.Dir(d)
	log.Println("AssetDir:", assetsDir)

	assetsDir = http.Dir(d)
	server := gosocketio.NewServer()

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected from %s as %s", c.IP(), c.Id())
		// do something.

		// draw specific for Kota.


	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected from %s as %s", c.IP(), c.Id())
	})

	return server
}

type MapMarker struct {
	mtype int32   `json:"mtype"`
	id    int32   `json:"id"`
	lat   float32 `json:"lat"`
	lon   float32 `json:"lon"`
	angle float32 `json:"angle"`
	speed int32   `json:"speed"`
	text  string  `json:"text"`
}

func (m *MapMarker) GetJson() string {
	s := fmt.Sprintf("{\"mtype\":%d,\"id\":%d,\"lat\":%f,\"lon\":%f,\"angle\":%f,\"speed\":%d}",
		m.mtype, m.id, m.lat, m.lon, m.angle, m.speed)
	return s
}

func supplyRideCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	flt := sp.GetArg_Fleet()
	if flt != nil { // get Fleet supplu
		vtype := 0
		if sp.SupplyName == nil {
			vtype =1
		}
		str :=":"
		if flt.VehicleId > 10000{
			str = fmt.Sprintf("%2d", flt.VehicleId%100)
		}

		mm := &MapMarker{
			mtype: int32(vtype),
			id:    flt.VehicleId,
			lat:   flt.Coord.Lat,
			lon:   flt.Coord.Lon,
			angle: flt.Angle,
			speed: flt.Speed,
			text: "No"+str,
		}
		mu.Lock()
		ioserv.BroadcastToAll("event", mm.GetJson())
		mu.Unlock()
	}
}

func subscribeRideSupply(client *sxutil.SMServiceClient) {
	ctx := context.Background() //
	err := client.SubscribeSupply(ctx, supplyRideCallback)
	log.Printf("Error:Supply %s\n",err.Error())
}

func supplyPTCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	pt := sp.GetArg_PTService()
	if pt != nil { // get Fleet supplu
		mm := &MapMarker{
			mtype: pt.VehicleType, // depends on type of GTFS: 1 for Subway, 2, for Rail, 3 for bus
			id:    pt.VehicleId,
			lat:   float32(pt.CurrentLocation.GetPoint().Latitude),
			lon:   float32(pt.CurrentLocation.GetPoint().Longitude),
			angle: pt.Angle,
			speed: pt.Speed,
			text: fmt.Sprintf("%d",VehicleId),
		}
		mu.Lock()
		ioserv.BroadcastToAll("event", mm.GetJson())
		mu.Unlock()
	}
}

func subscribePTSupply(client *sxutil.SMServiceClient) {
	ctx := context.Background() //
	err := client.SubscribeSupply(ctx, supplyPTCallback)
	log.Printf("Error:Supply %s\n",err.Error())
}

func monitorStatus(){
	for{
		sxutil.SetNodeStatus(int32(runtime.NumGoroutine()),"ExpMapGoroute")
		time.Sleep(time.Second * 10)
	}
}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "ExpMapProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Fail to Connect Synerex Server: %v", err)
	}
	ioserv = run_server()
	fmt.Printf("Running ExpMap Server..\n")
	if ioserv == nil {
		os.Exit(1)
	}

	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:ExMap:RIDE}")
	ride_client := sxutil.NewSMServiceClient(client, api.ChannelType_RIDE_SHARE, argJson)

	argJson2 := fmt.Sprintf("{Client:ExMap:PT}")
	pt_client := sxutil.NewSMServiceClient(client, api.ChannelType_PT_SERVICE, argJson2)

	wg.Add(1)
	go subscribeRideSupply(ride_client)
	wg.Add(1)
	go subscribePTSupply(pt_client)

	go monitorStatus() // keep status

	serveMux := http.NewServeMux()

	serveMux.Handle("/socket.io/", ioserv)
	serveMux.HandleFunc("/", assetsFileHandler)

	log.Printf("Starting KotaExpMap Provider %s  on port %d", version, *port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), serveMux)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

}
