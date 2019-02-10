package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// datastore provider provides Datastore Service.

type DataStore interface{
	store(dt *LocInfo)
}


var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	port      = flag.Int("port", 10080, "Map Provider Listening Port")
	mu         sync.Mutex
	version = "0.01"
	baseDir = "store"
	dataDir string
	ds DataStore
)

func init(){
	var err error
	dataDir, err =os.Getwd()
	if err != nil {
		fmt.Printf("Can't obtain current wd")
	}
	dataDir =filepath.ToSlash(dataDir) + "/" + baseDir
	ds = &FileSystemDataStore{
		storeDir:dataDir,
	}

}

type FileSystemDataStore struct{
	storeDir string
	storeFile *os.File
	todayStr string
}

// open file with today info
func (fs *FileSystemDataStore)store(li *LocInfo){
	const layout = "2006-01-02"
	day := time.Now()
	todayStr := day.Format(layout)+".csv"
	if fs.todayStr != "" && fs.todayStr != todayStr {
		fs.storeFile.Close()
		fs.storeFile = nil
	}
	if fs.storeFile == nil {
		_, er := os.Stat(fs.storeDir)
		if er != nil {// create dir
			er = os.MkdirAll(fs.storeDir, 0777)
			if er != nil {
				fmt.Printf("Can't make dir '%s'.",fs.storeDir)
				return
			}
		}
		fs.todayStr = todayStr
		file, err := os.OpenFile(filepath.FromSlash(fs.storeDir+"/"+todayStr),os.O_RDWR| os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Can't open file '%s'",todayStr)
			return
		}
		fs.storeFile =file
	}
	fs.storeFile.WriteString(li.argJson+"\n")
}

type LocInfo struct {
	mtype int32 `json:"mtype"`
	id int32 `json:"id"`
	lat float32 `json:"lat"`
	lon float32 `json:"lon"`
	angle float32 `json:"angle"`
	speed int32 `json:"speed"`
	argJson string `json:"json"`
}

func (m *LocInfo)GetJson() string {
	s := fmt.Sprintf("{\"id\":%d,\"lat\":%f,\"lon\":%f,\"angle\":%f,\"speed\":%d, %s}",
		m.id, m.lat, m.lon, m.angle, m.speed, m.argJson)
	return s
}


func supplyRideCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	flt := sp.GetArg_Fleet()
	if flt != nil{ // get Fleet supplu
		mm := &LocInfo{
			mtype:int32(api.ChannelType_RIDE_SHARE),
			id:flt.VehicleId,
			lat:flt.Coord.Lat,
			lon:flt.Coord.Lon,
			angle:flt.Angle,
			speed:flt.Speed,
		}
		ds.store( mm)

	}
}

func subscribeRideSupply(client *sxutil.SMServiceClient) {
	ctx := context.Background() //
	client.SubscribeSupply(ctx, supplyRideCallback)
}


func supplyPTCallback(clt *sxutil.SMServiceClient, sp *api.Supply) {
	pt := sp.GetArg_PTService()
	if pt != nil{ // get Fleet supplu
		mm := &LocInfo{
			mtype:int32(api.ChannelType_PT_SERVICE),
			id: pt.VehicleId,
			lat: float32(pt.CurrentLocation.GetPoint().Latitude),
			lon: float32(pt.CurrentLocation.GetPoint().Longitude),
			angle: pt.Angle,
			speed: pt.Speed,
			argJson: sp.ArgJson,
		}
		ds.store(mm)
	}
}

func subscribePTSupply(client *sxutil.SMServiceClient) {
	ctx := context.Background() //
	client.SubscribeSupply(ctx, supplyPTCallback)
}


func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "DataStoreProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Fail to Connect Synerex Server: %v", err)
	}

	client := api.NewSynerexClient(conn)
//	argJson := fmt.Sprintf("{Client:Map:RIDE}")
//	ride_client := sxutil.NewSMServiceClient(client, api.ChannelType_RIDE_SHARE,argJson)

	argJson2 := fmt.Sprintf("{Client:Datastore:PT}")
	pt_client := sxutil.NewSMServiceClient(client, api.ChannelType_PT_SERVICE,argJson2)

//	wg.Add(1)//
//	go subscribeRideSupply(ride_client)
	wg.Add(1)
	go subscribePTSupply(pt_client)

	wg.Wait()

}
