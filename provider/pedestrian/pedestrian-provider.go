package main

// Pedestrian Provider to provide file-based simulation.

import (
	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"github.com/synerex/synerex_alpha/api/fleet"	
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	waitmsec    = flag.Int("waitmsec", 100, "Wait for millisecond")
	jsonFile    = flag.String("json", "data/sample.json", "Movement json file")
)

type Opt struct {
	Elv []int `json:"optElevation"`
	Color []int `json:"color"`
	Position []float32 `json:"position"`
	OptColor [][]int `json:"optColor"`
	Memo string `json:"memo"`
	ElapsedTime int64 `json:"elapsedtime"`
	Test bool `json:"test"`	
}

type Movement struct {
	Dtime int64 `json:"departuretime"`
	Atime int64 `json:"arrivaltime"`
	Id string `json:"uuid"`
	Operation []Opt `json:"operation"`
}

func init() {

}

func readJsonMovement(client *sxutil.SMServiceClient){


	raw, err := ioutil.ReadFile(*jsonFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1) //
	}

	var mv []Movement
	
	json.Unmarshal(raw, &mv)
	raw = nil

	fmt.Printf("Length: %d\n", len(mv))
	sz := len(mv)
//	startTime := time.Now()  // no future time?
	startTime := time.Date(2999,12,31,23,59,59,0, time.Local).Unix()
	endTime := time.Date(1000,12,31,23,59,59,0, time.Local).Unix()

	fmt.Println("Start:",time.Unix(startTime,0))	
	for  i := 0;  i < sz; i++ {
		mvt := mv[i]
		if mvt.Dtime < startTime {
			startTime = mvt.Dtime
		}
		if mvt.Atime > endTime {
			endTime = mvt.Atime
		}
		
//		fmt.Printf("%d:  optSize %d \n", i, len(mvt.Operation))
	}

	fmt.Println("First:",time.Unix(startTime,0))
	fmt.Println("End  :",time.Unix(endTime,0))	

	//	fbytes, _ := json.Marshal(mv[1])
	//	fmt.Printf(string(fbytes))


	for now := startTime ;  now < endTime ; now ++ {
		for  i := 0;  i < sz; i++ {
			mvt := mv[i]
			if mvt.Dtime <= now && mvt.Atime >= now { // object exists.
				// calc location.
				
				for j :=0 ; j < len(mvt.Operation); j++ {
					if mvt.Operation[j].ElapsedTime == now {
						fleet := fleet.Fleet{
							VehicleId: int32(i),
							Angle:     float32(0),
							Speed:     int32(0),
							Status:    int32(0),
							Coord: &fleet.Fleet_Coord{
								Lat: float32(mvt.Operation[j].Position[1]),
								Lon: float32(mvt.Operation[j].Position[0]),
							},
						}
					
						smo := sxutil.SupplyOpts{
							Name:  "Ped Supply",
							Fleet: &fleet,
						}
						client.RegisterSupply(&smo)
					}
				}
				
			}
		}

		time.Sleep(time.Duration(*waitmsec) * time.Millisecond )

	}

}

func loopJsonMovement(client *sxutil.SMServiceClient){
	for {
		readJsonMovement(client)
	}
}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "PedestrianProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Pedestrian}")
	sclient := sxutil.NewSMServiceClient(client, pb.ChannelType_RIDE_SHARE,argJson)

	wg.Add(1)

	// We add Pedestrian Provider to "RIDE_SHARE" Supply
	go loopJsonMovement(sclient)
	//	go subscribeDemand(sclient)
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
