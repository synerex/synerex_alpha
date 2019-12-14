package main

import (
	"flag"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
	"github.com/synerex/synerex_alpha/provider/simulation/area/communicator"
	"github.com/synerex/synerex_alpha/sxutil"

	//	"github.com/synerex/synerex_alpha/api/simulation/area"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/synerex/synerex_alpha/provider/simulation/area/simulator"
	"google.golang.org/grpc"
	//	"os/exec"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	com        *communicator.AreaCommunicator
	sim        *simulator.AreaSimulator
	areaData   []*area.Area
)

func init() {
	areaData = readAreaData()
}

func readAreaData() []*area.Area {
	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile("area.json")
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	var areaData []*area.Area

	if err := json.Unmarshal(bytes, &areaData); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("areaData is : %v\n", areaData)
	return areaData
}

type Coord struct {
	StartLat float32 `json:"slat"`
	EndLat   float32 `json:"elat"`
	StartLon float32 `json:"slon"`
	EndLon   float32 `json:"elon"`
}

type Map struct {
	Id         uint32   `json:"id"`
	Name       string   `json:"name"`
	Coord      Coord    `json:"coord"`
	Controlled Coord    `json:"controlled"`
	Neighbor   []uint32 `json:"neighbor"`
}

func callbackGetAreaRequest(dm *pb.Demand) {
	areaId := dm.GetSimDemand().GetGetAreaRequest().GetId()
	targetId := dm.GetId()
	fmt.Printf("getArea2: %v\n", areaId)
	// AreaDataからエリア情報を取得
	for _, data := range areaData {
		if data.Id == areaId {
			fmt.Printf("getArea3: %v\n")
			// area情報を送信
			com.GetAreaResponse(targetId, data)
			break
		}
	}
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {

	fmt.Printf("registParticipant: %v\n", dm)
	switch dm.GetSimDemand().DemandType {

	case synerex.DemandType_GET_AREA_REQUEST:
		// エリアを取得する要求
		fmt.Printf("getArea: %v\n")
		callbackGetAreaRequest(dm)
	default:
		//log.Println("demand callback is invalid.")
	}
}

func main() {
	flag.Parse()

	sxutil.RegisterNodeName(*nodesrv, "AreaProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })

	// synerex simulator
	sim = simulator.NewAreaSimulator(1.0, 0.0)

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Area}")

	// Clientとして登録
	com = communicator.NewAreaCommunicator()

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	// channelごとのClientを作成
	com.RegistClients(client, argJson)
	// ChannelにSubscribe
	com.SubscribeAll(demandCallback, func(clt *sxutil.SMServiceClient, sp *pb.Supply) {}, &wg)
	wg.Wait()

	// start up(setArea)
	wg.Add(1)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
