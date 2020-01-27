package main

import (
	"flag"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
	"github.com/synerex/synerex_alpha/provider/simulation/area/communicator"
	"github.com/synerex/synerex_alpha/sxutil"
	"github.com/paulmach/orb/geojson"
	//"github.com/paulmach/orb"

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
	areafile string
	fcs *geojson.FeatureCollection
)

func init() {
	areaData = readAreaData()
	areafile = "area.geojson"
	//fcs = loadGeoJson(areafile)
}

/*func loadGeoJson(fname string) *geojson.FeatureCollection{

	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Print("Can't read file:", err)
		panic("load json")
	}
	fc, _ := geojson.UnmarshalFeatureCollection(bytes)
	log.Printf("fc: %v\n", fc.Features[0].Geometry.(orb.MultiLineString)[0])

	return fc
}*/

func readAreaData() []*area.Area {
	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile("area3.json")
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	var areaData []*area.Area

	if err := json.Unmarshal(bytes, &areaData); err != nil {
		log.Fatal(err)
	}
	log.Printf("\x1b[30m\x1b[47m \n Finish: Area data got. \n AreaData lenth : %v \x1b[0m\n", len(areaData))
	return areaData
}


// notifyStartUp : 起動時に、他プロバイダの参加者情報を集める
func notifyStartUp() {

	// 情報をプロバイダに送信
	com.NotifyStartUpRequest(participant.ProviderType_AREA)
}

func callbackGetAreaRequest(dm *pb.Demand) {
	areaId := dm.GetSimDemand().GetGetAreaRequest().GetId()
	targetId := dm.GetId()
	// AreaDataからエリア情報を取得
	for _, data := range areaData {
		if data.Id == areaId {
			// area情報を送信
			com.GetAreaResponse(targetId, data)
			log.Printf("\x1b[30m\x1b[47m \n Finish: Area information sent. \n AreaData : %v \x1b[0m\n", data.GetDuplicateArea())
			break
		}
	}
}

/*func callbackGetAreaRequest2(dm *pb.Demand) {
	areaId := dm.GetSimDemand().GetGetAreaRequest().GetId()
	targetId := dm.GetId()
	// AreaDataからエリア情報を取得
	for i, feature := range fcs.Features {
		multiPosition := feature.Geometry.(orb.MultiLineString)[0]
		if data.Id == areaId {
			// area情報を送信
			com.GetAreaResponse(targetId, data)
			log.Printf("\x1b[30m\x1b[47m \n Finish: Area information sent. \n AreaData : %v \x1b[0m\n", data)
			break
		}
	}
}*/

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {

	switch dm.GetSimDemand().DemandType {

	case synerex.DemandType_GET_AREA_REQUEST:
		// エリアを取得する要求
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

	// 起動したことを通知
	notifyStartUp()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
