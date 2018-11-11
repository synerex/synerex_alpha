package main

import (
	"flag"
	"fmt"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/ptransit"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"log"
	"math"
	"math/rand"
	"time"
)

// test provider provides Testing Performance of provider.

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	throttle   = flag.Int("throttle", 100, "Throttring duration")
)


func main() {
	flag.Parse()
	err := sxutil.RegisterNodeName(*nodesrv, "TestProvider", false)
	if err != nil {
		log.Printf("Can't register Node: %s\n",err.Error())
		return
	}

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Printf("Fail to Connect Synerex Server: %v\n", err.Error())
	}

	client := api.NewSMarketClient(conn)
	//	argJson := fmt.Sprintf("{Client:Map:RIDE}")
	//	ride_client := sxutil.NewSMServiceClient(client, api.MarketType_RIDE_SHARE,argJson)

	argJson2 := fmt.Sprintf("{Client:TestProvider}")
	ptClient := sxutil.NewSMServiceClient(client, api.MarketType_PT_SERVICE,argJson2)

	rand.Seed(time.Now().UnixNano())
	count := 0
	angle := 0.0
	lat , lng := 34.8644,  137.16559
	for {
		count ++
		pts :=  &ptransit.PTService{
			VehicleId: int32(99),
			Angle: float32(angle),
			Speed: int32(10.0),
			CurrentLocation: common.NewPlace().WithPoint(&common.Point{
				Latitude: lat,
				Longitude: lng,
			}),
		}
		smo := &sxutil.SupplyOpts{
			Name:  "Kotacho Info",
			PTService: pts,
			JSON: fmt.Sprintf("count %d",count),
		}
		id :=	ptClient.RegisterSupply(smo)
		if id == 0 {
			log.Printf("Error for Register Supply\n")
			break
		}

		time.Sleep(time.Duration(*throttle)* time.Millisecond)
		angle = angle + 20-float64(rand.Int31n(40))
		lat = lat+math.Cos(angle/180*math.Pi)*0.0001
		lng = lng+math.Sin(angle/180*math.Pi)*0.0001
	}

}
