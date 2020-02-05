package main

import (
	"flag"
	"log"
	"sync"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"sort"
	"bufio"
	"io"
	"strconv"
	"strings"
	"github.com/google/uuid"
	gosocketio "github.com/mtfelian/golang-socketio"
	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/daemon"
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	"github.com/synerex/synerex_alpha/provider/simulation/scenario/communicator"
	"github.com/synerex/synerex_alpha/provider/simulation/scenario/simulator"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"io/ioutil"
)

var (
	serverAddr       = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv          = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	version          = "0.01"
	isStart          bool
	mu               sync.Mutex
	com              *communicator.ScenarioCommunicator
	sim              *simulator.ScenarioSimulator
	providerManager *ProviderManager
)

const MAX_AGENTS_NUM = 1000

func init() {
	isStart = false
	providerManager = NewProviderManager()
}


var (
	fcs *geojson.FeatureCollection
	geofile string
	port = 9995
	assetsDir http.FileSystem
	server *gosocketio.Server = nil
	//providerMap map[string]*exec.Cmd
	runProviders map[uint64]*Provider
	providerMutex sync.RWMutex
	client daemon.SimDaemonClient
	providerArray []ProviderData
	orderArray []OrderData
)



func loadGeoJson(fname string) *geojson.FeatureCollection{

	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Print("Can't read file:", err)
		panic("load json")
	}
	fc, _ := geojson.UnmarshalFeatureCollection(bytes)

	return fc
}








type Coord struct{
	Longitude float64
	Latitude float64
}

// エリアの初期値
var initAreaNum = 4

// エリア分割の閾値(100人超えたら分割)
var areaAgentNum = 100

// 担当するAreaの範囲
var mockAreaData = []*common.Coord{
	{Latitude: 35.156431, Longitude: 136.97285,},
	{Latitude: 35.156431, Longitude: 136.981308,},
	{Latitude: 35.153578, Longitude: 136.981308,},
	{Latitude: 35.153578, Longitude: 136.97285,},
}


var providerStats = []*ProviderStats{
	{
		AgentNum: 50,
		Time: 0.0,
	},
	{
		AgentNum: 50,
		Time: 0.0,
	},
	{
		AgentNum: 50,
		Time: 0.0,
	},
}

/*var mockProviderStats = []*ProviderStats{
	{
		AgentNum: 50,
		Time: 0.0,
	},
}*/

func Sub(coord1 *common.Coord, coord2 *common.Coord)*common.Coord{
	return &common.Coord{Latitude: coord1.Latitude-coord2.Latitude, Longitude: coord1.Longitude-coord2.Longitude}
}

func Mul(coord1 *common.Coord, coord2 *common.Coord) float64 {
	return coord1.Latitude*coord2.Latitude + coord1.Longitude*coord2.Longitude
}

func Abs(coord1 *common.Coord) float64 {
	return math.Sqrt(Mul(coord1, coord1))
}

func Add(coord1 *common.Coord, coord2 *common.Coord) *common.Coord {
	return &common.Coord{Latitude: coord1.Latitude + coord2.Latitude, Longitude: coord1.Longitude + coord2.Longitude}
}

func Div(coord *common.Coord, s float64) *common.Coord {
	return &common.Coord{Latitude: coord.Latitude / s, Longitude: coord.Longitude / s}
}

type ByAbs struct {
	Coords []*common.Coord
}

func (b ByAbs) Less(i, j int) bool{
	return Abs(b.Coords[i]) < Abs(b.Coords[j])
}
func (b ByAbs) Len() int {
    return len(b.Coords)
}

func (b ByAbs) Swap(i, j int) {
    b.Coords[i], b.Coords[j] = b.Coords[j], b.Coords[i]
}

/*func updateProvider(newProviderStats []*ProviderStats) []*ProviderStats{
	// write area devide algorithm
	updatedProviderStats := make([]*ProviderStats, 0)
	for _, state := range newProviderStats{
		if state.AgentNum > areaAgentNum{
			// devide area
			//vectors := make([]*common.Coord, 0)
			
			point1, point2, point3, point4 := state.Area[0], state.Area[1], state.Area[2], state.Area[3]
			point1vecs := []*common.Coord{Sub(point1, point1), Sub(point2, point1), Sub(point3, point1), Sub(point4, point1)}
			// 昇順にする
			sort.Sort(ByAbs{point1vecs})
			devPoint1 := Div(point1vecs[2], 2)	//分割点1
			divPoint2 := Add(Div(point1vecs[2], 2), point1vecs[1]) //分割点2
			// 二つに分割
			coords1 := []*common.Coord{
				Add(point1vecs[0], point1), Add(point1vecs[1], point1), Add(devPoint1, point1), Add(divPoint2, point1),
			}
			coords2 := []*common.Coord{
				Add(point1vecs[2], point1), Add(point1vecs[3], point1), Add(devPoint1, point1), Add(divPoint2, point1),
			}
			// 追加
			state1 := &ProviderStats{
				Area: coords1,
				AgentNum: state.AgentNum/2,
				Time: state.Time,
			}
			state2 := &ProviderStats{
				Area: coords2,
				AgentNum: state.AgentNum/2,
				Time: state.Time,
			}

			updatedProviderStats = append(updatedProviderStats, state1)
			updatedProviderStats = append(updatedProviderStats, state2)
		}else{
			updatedProviderStats = append(updatedProviderStats, state)
		}
	}
	return updatedProviderStats
}

func areaTest(){
	go func(){
		for{
			log.Printf("time: ---\n")
			time.Sleep(1 * time.Second)
			// change provider stats
			newProviderStats := make([]*ProviderStatus, 0)
			for i, state := range mockProviderStats{
				if i == 0{
					log.Printf("agentNum: %v, providerNum: %v\n", state.AgentNum, len(mockProviderStats))
				}
				log.Printf("area: %v\n", state.Area)
				state.AgentNum += 30
				newProviderStats = append(newProviderStats, state)
			}
			// update provider
			mockProviderStats = updateProvider(newProviderStats)

			// startup provider
			//log.Printf("providerStats: %v\n", providerStats)
		}
	}()
}*/

type Option struct{
	Key string
	Value string
}

type ProviderData struct{
	CmdName     string
	Type ProviderType
	Description string
	SrcDir      string
	BinName     string
	GoFiles     []string
	Options     []Option
}

type OrderData struct{
	CmdName     string
	Type OrderType
	Options     []Option
}

func init() {
	//areaTest()
	geofile = "transit_points.geojson"
	fcs = loadGeoJson(geofile)
	runProviders = make(map[uint64]*Provider)
	//providerMap = make(map[string]*exec.Cmd)
	providerMutex = sync.RWMutex{}

	providerArray = []ProviderData{
		{
			CmdName: "NodeIDServer",
			Type: ProviderType_NODE_ID_SERVER,
			SrcDir:  "nodeserv",
			BinName: "nodeid-server",
			GoFiles: []string{"nodeid-server.go"},
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "MonitorServer",
			Type: ProviderType_MONITOR_SERVER,
			SrcDir:  "monitor",
			BinName: "monitor-server",
			GoFiles: []string{"monitor-server.go"},
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "SynerexServer",
			Type: ProviderType_SYNEREX_SERVER,
			SrcDir:  "server",
			BinName: "synerex-server",
			GoFiles: []string{"synerex-server.go", "message-store.go"},
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "Area",
			Type: ProviderType_AREA,
			SrcDir:  "provider/simulation/area",
			BinName: "area-provider",
			GoFiles: []string{"area-provider.go"},
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "Scenario",
			Type: ProviderType_SCENARIO,
			SrcDir:  "provider/simulation/scenario",
			BinName: "scenario-provider",
			GoFiles: []string{"scenario-provider.go"},
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "Pedestrian",
			Type: ProviderType_PEDESTRIAN,
			SrcDir:  "provider/simulation/pedestrian",
			BinName: "pedestrian-provider",
			GoFiles: []string{"pedestrian-provider.go"},
			Options: []Option{Option{
				Key: "areaId",
				Value: "1",
			}},
		},
		{
			CmdName: "Car",
			Type: ProviderType_CAR,
			SrcDir:  "provider/simulation/car",
			BinName: "car-provider",
			GoFiles: []string{"car-provider.go"},
			Options: []Option{Option{
				Key: "areaId",
				Value: "1",
			}},
		},
		{
			CmdName: "Visualization",
			Type: ProviderType_VISUALIZATION,
			SrcDir:  "provider/simulation/visualization",
			BinName: "visualization-provider",
			GoFiles: []string{"visualization-provider.go"},
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "Clock",
			Type: ProviderType_CLOCK,
			SrcDir:  "provider/simulation/clock",
			BinName: "clock-provider",
			GoFiles: []string{"clock-provider.go"},
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
	}

	orderArray = []OrderData{
		{
			CmdName: "SetAgents",
			Type: OrderType_SET_AGENTS,
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "SetArea",
			Type: OrderType_SET_AREA,
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "SetClock",
			Type: OrderType_SET_CLOCK,
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "StartClock",
			Type: OrderType_START_CLOCK,
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
		{
			CmdName: "StopClock",
			Type: OrderType_STOP_CLOCK,
			Options: []Option{Option{
				Key: "test",
				Value: "0",
			}},
		},
	}

}

////////////////////////////////////////////////////////////
//////////////////        Util          ///////////////////
///////////////////////////////////////////////////////////

// providerに変化があった場合にGUIに情報を送る
func sendRunnningProviders(){
	providerMutex.RLock()

	//fmt.Printf("providers---------- %v\n", len(runProviders))
	rpJsons := make([]string, 0)
	for _, rp := range runProviders {
		bytes, _ := json.Marshal(rp)
    	rpJson := string(bytes)
		fmt.Printf("provider----------\n")
		//fmt.Printf("Json: %v \n", rpJson)
		rpJsons = append(rpJsons, rpJson)
	}
	//c.Emit("providers", rpJsons)
	server.BroadcastToAll("providers", rpJsons)
	providerMutex.RUnlock()
}



// assetsFileHandler for static Data
func assetsFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}

	file := r.URL.Path

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

func createOption() *ProviderOption{
	po := &ProviderOption{
	}
	return po
}

func getGoPath() string{
	env := os.Environ()
	for _, ev := range env {
		if strings.Contains(ev,"GOPATH=") {
			return ev
		}
	}
	return ""
}

func getGoEnv() []string { // we need to get/set gopath
	d, _ := os.Getwd() // may obtain dir of se-daemon
	gopath := filepath.FromSlash(filepath.ToSlash(d) + "/../../../../")
	absGopath, _ := filepath.Abs(gopath)
	env := os.Environ()
	newenv := make([]string, 0, 1)
	foundPath := false
	for _, ev := range env {
		if strings.Contains(ev, "GOPATH=") {
			// this might depends on each OS
			newenv = append(newenv, ev+string(os.PathListSeparator)+filepath.FromSlash(filepath.ToSlash(absGopath)+"/"))
			foundPath = true
		} else {
			newenv = append(newenv, ev)
		}
	}
	if !foundPath { // this might happen at in the daemon..
		gp := getGoPath()
		newenv = append(newenv, gp)
	}
	return newenv
}



//////// To Order ////////////

func decideCoordInGeo(fcs *geojson.FeatureCollection) *common.Coord{

	geoCoords := fcs.Features[0].Geometry.(orb.MultiLineString)[0]
	transitNum := rand.Int63n(int64(len(geoCoords)-1))

	longitude := geoCoords[transitNum][0] + 0.0001 * rand.Float64()
	latitude := geoCoords[transitNum][1] + 0.0001 * rand.Float64()
	coord := &common.Coord{
		Longitude: longitude,
		Latitude: latitude,
	}
	//log.Printf("coord: ", coord)

	return coord
}

// Agentオブジェクトの変換
func calcRoute() *agent.PedRoute {

	var departure, destination *common.Coord
	
	departure = decideCoordInGeo(fcs)
	destination = decideCoordInGeo(fcs)

	transitPoints := make([]*common.Coord, 0)
	transitPoints = append(transitPoints, destination)

	route := &agent.PedRoute{
		Position:    departure,
		Direction:   100 * rand.Float64(),
		Speed:       100 * rand.Float64(),
		Departure:   departure,
		Destination: destination,
		TransitPoints: transitPoints,
		NextTransit: destination,
	}

	return route
}



func runDividedProvider(name string, providerType ProviderType){
	// 最初に分割して起動するプロバイダ数 2つに分割
	const INIT_PROVIDER_NUM = uint64(2)  // 1, 2, 4, 9, 16...
	isFirstRun := true
	for _, provider := range runProviders{
		if provider.Type == providerType{
			isFirstRun = false
		}
	}
	if isFirstRun {
		areaInfos := divideArea(mockAreaData, INIT_PROVIDER_NUM)
		for _, areaInfo := range areaInfos{
			fmt.Printf("areaInfo: %v\n", areaInfo)
			options := make([]*Option, 0)
			options = append(options, &Option{
				Key: "server_addr",
				Value: "127.0.0.1:10000",
			})
			options = append(options, &Option{
				Key: "nodesrv",
				Value: "127.0.0.1:9990",
			})
			bytes, _ := json.Marshal(areaInfo)
			options = append(options, &Option{
				Key: "area_json",
				Value: string(bytes),
			})
			provider, _ := NewProvider(name, options)
			provider.Run()
		}
	}
}


var mockProviderStats = &ProviderStats{
	AgentNum: 1500,
	Time: 0.0,
}

func divideArea(areaCoord []*common.Coord, num uint64) []*area.Area2{
	// エリアを分割する
	// 最初は単純にエリアを半分にする
	//providerStats := mockProviderStats
	//duplicateRate := 0.1	// areaCoordの10%の範囲
	// 二等分にするアルゴリズム
	point1, point2, point3, point4 := areaCoord[0], areaCoord[1], areaCoord[2], areaCoord[3]
	point1vecs := []*common.Coord{Sub(point1, point1), Sub(point2, point1), Sub(point3, point1), Sub(point4, point1)}
	// 昇順にする
	sort.Sort(ByAbs{point1vecs})
	divPoint1 := Div(point1vecs[2], 2)	//分割点1
	divPoint2 := Add(Div(point1vecs[2], 2), point1vecs[1]) //分割点2
	// 二つに分割
	coords1 := []*common.Coord{
		Add(point1vecs[0], point1), Add(point1vecs[1], point1), Add(divPoint1, point1), Add(divPoint2, point1),
	}
	coords2 := []*common.Coord{
		Add(point1vecs[2], point1), Add(point1vecs[3], point1), Add(divPoint1, point1), Add(divPoint2, point1),
	}
	areaInfos := []*area.Area2{&area.Area2{
		Id: 1,
		Name: "aaa",
		NeighborAreas: []uint64{1},
		DuplicateArea: coords1,
		ControlArea: coords1,
	}, &area.Area2{
		Id: 1,
		Name: "bbb",
		NeighborAreas: []uint64{1},
		DuplicateArea: coords2,
		ControlArea: coords2,
	}}

	for _, coord := range coords1{
		fmt.Printf("coord: %v\n", coord )
	}
	for _, coord := range coords2{
		fmt.Printf("coord: %v\n", coord )
	}

	return areaInfos
}

// ps commands from se cli
//////// To ps ////////////
func checkRunning(opt string) []string {
	isLong := false
	if opt == "long" {
		isLong = true
	}
	var procs []string
	i := 0
	providerMutex.RLock()
	if isLong {
		procs = make([]string, len(runProviders)+2)
		str := fmt.Sprintf("  pid: %-20s : \n", "process name")
		procs[i] = str
		procs[i+1] = "-----------------------------------------------------------------\n"
		i += 2
	} else {
		procs = make([]string, len(runProviders))
	}
	for _, provider := range runProviders {
		pid := provider.Cmd.Process.Pid
		name := provider.Name
		if isLong {
			str := fmt.Sprintf("%5d: %-20s : \n", pid, name)
			procs[i] = str
		} else {
			if i != 0 {
				procs[i] = ", " + name
			} else {
				procs[i] = name
			}
		}
		i += 1
	}
	providerMutex.RUnlock()
	return procs

}



////////////////////////////////////////////////////////////
//////////////         Order Class         /////////////////
///////////////////////////////////////////////////////////

// Order
type OrderType int
const (
	OrderType_SET_AGENTS  OrderType = 0
    OrderType_SET_AREA  OrderType = 1
    OrderType_SET_CLOCK  OrderType = 2
    OrderType_START_CLOCK  OrderType = 3
	OrderType_STOP_CLOCK OrderType = 4
)

type OrderOption struct{
	AgentNum string
	ClockTime string
}

type Order struct {
	Type   OrderType
	Name string
	Option *OrderOption
}

func NewOrder(name string, option *OrderOption) (*Order, error){
	for _, sc := range orderArray {
		if sc.CmdName == name {
			o := &Order{
				Type: sc.Type,
				Name: name,
				Option: option,
			}
			return o, nil
		}
	}
	msg := "invalid OrderName..."
	return nil, fmt.Errorf("Error: %s\n", msg)
}

func (o *Order)Send() string {
	target := o.Name
	fmt.Printf("Target is : %v\n", target)
	switch target {
	case "SetClock":
		fmt.Printf("SetClock\n")
		globalTime := float64(0)
		timeStep := float64(1)
		o.SetClock(globalTime, timeStep)
		return "ok"

	case "SetAgents":
		fmt.Printf("SetAgents\n")
		//agentNum, _ := strconv.Atoi(order.Option)
		agentNum := uint64(1)
		o.SetAgents(agentNum)
		return "ok"

	case "StartClock":
		fmt.Printf("StartClock\n")
		stepNum := uint64(1)
		o.StartClock(stepNum)
		return "ok"

	case "StopClock":
		fmt.Printf("StopClock\n")
		o.StopClock()
		return "ok"

	case "SetArea":
		fmt.Printf("SetArea\n")
		//o.StopClock()
		return "ok"

	default:
		err := "true"
		log.Printf("Can't find command %s", target)
		return err
	}
}

// startClock:
func (o *Order)StartClock(stepNum uint64) (bool, error) {

	// エージェントを設置するリクエスト
	com.StartClockRequest(stepNum)

	// 同期のため待機
	com.WaitStartClockResponse()
	return true, nil
}

// stopClock: Clockを停止する
func (o *Order)StopClock() (bool, error) {
	// エージェントを設置するリクエスト
	com.StopClockRequest()

	// 同期のため待機
	com.WaitStopClockResponse()
	return true, nil
}

// setAgents: agentをセットするDemandを出す関数
func (o *Order)SetAgents(agentNum uint64) (bool, error) {

	agents := make([]*agent.Agent, 0)

	for i := 0; i < int(agentNum); i++ {
		uuid, err := uuid.NewRandom()
		if err == nil {
			agent := &agent.Agent{
				Id:   uint64(uuid.ID()),
				Type: common.AgentType_PEDESTRIAN,
				Data: &agent.Agent_Pedestrian{
					Pedestrian: &agent.Pedestrian{
						Status: &agent.PedStatus{
							Age:  "20",
							Name: "rui",
						},
						Route: calcRoute(),
					},
				},
			}
			agents = append(agents, agent)
		}
	}

	// agentsに必要なプロバイダを起動
	runDividedProvider("Pedestrian", ProviderType_PEDESTRIAN)

	// プロバイダに参加者情報を配布
	//log.Printf("\x1b[30m\x1b[47m \n Finish: SetAgentRequest! \x1b[0m\n")
	//com.SetParticipantsRequest()

	// エージェントを設置するリクエスト
	com.SetAgentsRequest(agents)

	// 同期のため待機
	com.WaitSetAgentsResponse()

	log.Printf("\x1b[30m\x1b[47m \n Finish: Agents set \n Add: %v \x1b[0m\n", len(agents))
	return true, nil
}


// setClock : クロック情報をDaemonから受け取りセットする
func (o *Order)SetClock(globalTime float64, timeStep float64) (bool, error) {
	// クロックをセット
	sim.SetGlobalTime(globalTime)
	sim.SetTimeStep(timeStep)

	// クロック情報をプロバイダに送信
	clockInfo := sim.GetClock()
	com.SetClockRequest(clockInfo)
	// Responseを待機
	com.WaitSetClockResponse()
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.GlobalTime, sim.TimeStep)
	return true, nil
}


////////////////////////////////////////////////////////////
//////////////       Provider Manager Class      //////////
///////////////////////////////////////////////////////////

type ProviderManager struct{
	Providers []*Provider
}

func NewProviderManager() *ProviderManager{
	pm := &ProviderManager{
		Providers: make([]*Provider, 0),
	}
	return pm
}

func (pm *ProviderManager)AddProvider(provider *Provider){
	pm.Providers = append(pm.Providers, provider)
	//log.Printf("Providers: %v\n", pm.Providers)
}

func (pm *ProviderManager)SetProvider(index int, provider *Provider){
	log.Printf("\x1b[31m\x1b[47m \n Provider Registed!!!: %v \x1b[0m\n", provider)
	pm.Providers[index] = provider
}

func (pm *ProviderManager)DeleteProvider(id uint64){
	newProviders := make([]*Provider, 0)
    for _, provider := range pm.Providers {
        if provider.ID == id {
            continue
        }
        newProviders = append(newProviders, provider)
    }
	pm.Providers = newProviders
}

////////////////////////////////////////////////////////////
////////////         Provider Class         ////////////////
///////////////////////////////////////////////////////////

// Provider
type ProviderType int
const (
	ProviderType_SCENARIO  ProviderType = 0
    ProviderType_AREA  ProviderType = 1
    ProviderType_CAR  ProviderType = 2
    ProviderType_PEDESTRIAN  ProviderType = 3
    ProviderType_ROUTE ProviderType = 4
	ProviderType_VISUALIZATION  ProviderType = 5
	ProviderType_NODE_ID_SERVER ProviderType = 6
	ProviderType_SYNEREX_SERVER ProviderType = 7
	ProviderType_MONITOR_SERVER ProviderType = 8
	ProviderType_CLOCK ProviderType = 9
)

type ProviderStatusType int
const (
	ProviderStatusType_NONE  ProviderStatusType = 0
    ProviderStatusType_REGISTED  ProviderStatusType = 1
    ProviderStatusType_RUNNING  ProviderStatusType = 2
)


type ProviderStats struct{
	AgentNum int
	Time float64
}

type ProviderOption struct{
	NodeServAddr string
	SynerexServAddr string
	AreaCoord string
}

type Provider struct{
	ID uint64
	Name string
	Type ProviderType
	Options []*Option
	Status ProviderStatusType
	Cmd *exec.Cmd
	Sc ProviderData
	Stats []*ProviderStats
	ParticipantInfo *participant.Participant
}

func  NewProvider(name string, options []*Option) (*Provider, error){
	for _, sc := range providerArray {
		if sc.CmdName == name {
			uuid, _ := uuid.NewRandom()
			pid := uint64(uuid.ID())
			// server以外にoptionにidを追加
			if sc.Type != ProviderType_MONITOR_SERVER && sc.Type != ProviderType_NODE_ID_SERVER && sc.Type != ProviderType_SYNEREX_SERVER{
				log.Printf("\x1b[31m\x1b[47m \n AppendPID %v \x1b[0m\n", strconv.FormatUint(pid,10))
				options = append(options, &Option{
					Key: "pid",
					Value: strconv.FormatUint(pid,10),
				})
			}
			p := &Provider{
				ID: pid,
				Name: name,
				Type: sc.Type,
				Sc: sc,
				Options: options,
				Status: ProviderStatusType_NONE,
			}
			return p, nil
		}
	}
	msg := "invalid CmdName..."
	return nil, fmt.Errorf("Error: %s\n", msg)
}

func (p *Provider) Run() error{
	log.Printf("Run '%s'\n", p.Name)
	sc := p.Sc

	d, err := os.Getwd()
	if err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("cannot get dir: %s", err.Error())
	}

	// get src dir
	srcpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../../" + sc.SrcDir)
	binpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../../" + sc.SrcDir + "/" + sc.BinName)
	//fi, err := os.Stat(binpath)
	_, err = os.Stat(binpath)

	// バイナリが最新かどうか
	modTime := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	for _, fn := range sc.GoFiles {
		sp := filepath.FromSlash(filepath.ToSlash(srcpath) + "/" + fn)
		ss, _ := os.Stat(sp)
		if ss.ModTime().After(modTime) {
			modTime = ss.ModTime()
		}
	}

	// 最新でない場合、run
	var cmd *exec.Cmd
	if err == nil{ //&& fi.ModTime().After(modTime) { // check binary time
		cmdArgs := make([]string, 0)
		for _, option := range p.Options{
			cmdArgs =  append(cmdArgs, "-" + option.Key)
			cmdArgs =  append(cmdArgs, option.Value)
		}
		cmd = exec.Command("./" + sc.BinName, cmdArgs...) // run binary
	} else {
		log.Printf("Error: [provider].go file isn't done build command\n")
		return fmt.Errorf("Error: [provider].go file isn't done build command")
		//runArgs := append([]string{"run"}, sc.GoFiles...)
		//log.Printf("runArgs: [%s] %d, %s %s", sc.CmdName, len(sc.GoFiles), runArgs[0], runArgs[1])
		//cmd = exec.Command("go", runArgs...) // run go with srcfile
	}

	cmd.Dir = srcpath
	cmd.Env = getGoEnv()
	p.Cmd = cmd

	go p.RunMyCmd()

	return nil
}

type Log struct{
	ID uint64
	Description string
}

func (p *Provider)RunMyCmd() {
	cmd := p.Cmd
	cmdName := p.Name

	pipe, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error for getting stdout pipe %s\n", cmd.Args[0])
		return
	}
	err = cmd.Start()
	if err != nil {
		log.Printf("Error for executing %s %v\n", cmd.Args[0], err)
		return
	}
	log.Printf("Starting %s..\n", cmd.Args[0])

	// プロバイダーをリストに追加
	providerMutex.Lock()
	providerManager.AddProvider(p)
	runProviders[p.ID] = p
	providerMutex.Unlock()

	// logを読み取る
	reader := bufio.NewReader(pipe)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			log.Printf("Command [%s] EOF\n", cmdName)
			break
		} else if err != nil {
			log.Printf("Err %v\n", err)
		}

		logInfo := &Log{
			ID: p.ID,
			Description: string(line),
		}

		bytes, err  := json.Marshal(logInfo)
		logjson := string(bytes)
 
		if server != nil {
			server.BroadcastToAll("log", logjson)
		}
		log.Printf("[%s]:%s", cmdName, string(line))
	}
	
	log.Printf("[%s]:Now ending...", cmdName)

	cmd.Wait()
	providerMutex.Lock()
	providerManager.DeleteProvider(p.ID)
	delete(runProviders, p.ID)
	providerMutex.Unlock()

	log.Printf("Command [%s] closed\n", cmdName)
}


////////////////////////////////////////////////////////////
////////////     Simulator CLI GUI Server    //////////////
//////////////////////////////////////////////////////////

type SimulatorServer struct{}

func NewSimulatorServer() *SimulatorServer{
	ss := &SimulatorServer{}
	return ss
}

func (ss *SimulatorServer)Run() error {
	go func(){
		log.Printf("Starting.. Synergic Engine:" + version)
		currentRoot, err := os.Getwd()
		if err != nil {
			log.Printf("se-daemon: Can' get registered directory: %s", err.Error())
		}
		d := filepath.Join(currentRoot, "monitor", "build")
	
		assetsDir = http.Dir(d)
		server = gosocketio.NewServer()
	
		server.On(gosocketio.OnConnection, func(c *gosocketio.Channel, param interface{}) {
			log.Printf("Connected from %s as %s", c.IP(), c.Id())
			// we need to send providers array
			time.Sleep(1000 * time.Millisecond)
			sendRunnningProviders()
	
		})
		server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
			log.Printf("Disconnected from %s as %s", c.IP(), c.Id())
		})
	
	
		server.On("ps", func(c *gosocketio.Channel, param interface{}) []string {
			// need to check param short or long
			opt := param.(string)
	
			return checkRunning(opt)
		})
		
	
		server.On("run", func(c *gosocketio.Channel, param interface{}) string {
			targetName := param.(string)
			log.Printf("Get run command %s", targetName)
	
			provider, _ := NewProvider("targetName", nil)
			provider.Run()
			//ok := handleRun(targetName)
			sendRunnningProviders()
			return "ok"
		})
	
		server.On("order", func(c *gosocketio.Channel, param *Order) string {
			name := param.Name
			log.Printf("Get order command %s\n", name)
			log.Printf("Get order command %v\n", param)
			log.Printf("Get order command %v\n", param.Option)
			order, _ := NewOrder(name, nil)
			order.Send()
			return "ok"
		})
	
	
		serveMux := http.NewServeMux()
		serveMux.Handle("/socket.io/", server)
		serveMux.HandleFunc("/", assetsFileHandler)
		log.Println("Serving at localhost:9995...")
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), serveMux); err != nil {
			log.Panic(err)
		}
	
		return

	}()
	return nil
}


////////////////////////////////////////////////////////////
////////////     Demand Supply Callback     ////////////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	switch sp.GetSimSupply().SupplyType {
	case synerex.SupplyType_SET_AGENTS_RESPONSE:
		com.SendToSetAgentsResponse(sp)
	case synerex.SupplyType_SET_CLOCK_RESPONSE:
		com.SendToSetClockResponse(sp)
	case synerex.SupplyType_FORWARD_CLOCK_RESPONSE:
		com.SendToForwardClockResponse(sp)
	default:
		//fmt.Println("order is invalid")
	}
}

// Demandのコールバック関数
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	switch dm.GetSimDemand().DemandType {
	case synerex.DemandType_DIVIDE_AREA_REQUEST:

		//divideArea(dm)
	case synerex.DemandType_REGIST_PARTICIPANT_REQUEST:
		participantInfo := dm.GetSimDemand().GetRegistParticipantRequest().GetParticipant()
		log.Printf("\x1b[31m\x1b[47m \n Provider Request: %v %v\x1b[0m\n",participantInfo.Id, participantInfo.ProviderType )
		for i, provider := range providerManager.Providers{
			log.Printf("\x1b[31m\x1b[47m \n Provider Request: %v %v\x1b[0m\n",participantInfo.Id, provider.ID )
			if provider.ID == participantInfo.Id{
				log.Printf("\x1b[31m\x1b[47m \n Provider Request \x1b[0m\n")
				
				provider.ParticipantInfo = participantInfo
				provider.Status = ProviderStatusType_REGISTED
				providerManager.SetProvider(i, provider)
				log.Printf("\x1b[30m\x1b[47m \n Start: SetAgentRequest! \x1b[0m\n")
				com.SetParticipantsRequest()
			}
		}
		// 返信
		com.RegistParticipantResponse(dm)
	case synerex.DemandType_DELETE_PARTICIPANT_REQUEST:
		participantInfo := dm.GetSimDemand().GetRegistParticipantRequest().GetParticipant()
		providerManager.DeleteProvider(participantInfo.Id)

	default:
		//fmt.Println("order is invalid")
	}

}


func main() {

	flag.Parse()

	// CLI, GUIの受信サーバ
	simulatorServer := NewSimulatorServer()
	simulatorServer.Run()
	
	// Run Server and Provider
	nodeServer, _ := NewProvider("NodeIDServer", nil)
	nodeServer.Run()
	//handleRun("NodeIDServer")
	time.Sleep(500 * time.Millisecond)
	monitorServer, _ := NewProvider("MonitorServer", nil)
	monitorServer.Run()
	//handleRun("MonitorServer")
	time.Sleep(500 * time.Millisecond)
	synerexServer, _ := NewProvider("SynerexServer", nil)
	synerexServer.Run()
	//handleRun("SynerexServer")
	time.Sleep(500 * time.Millisecond)
	clockProvider, _ := NewProvider("Clock", nil)
	clockProvider.Run()
	//handleRun("Clock")
	time.Sleep(500 * time.Millisecond)
	visProvider, _ := NewProvider("Visualization", nil)
	visProvider.Run()

	// Connect to Node Server
	sxutil.RegisterNodeName(*nodesrv, "Scenario2Provider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })


	// Connect to Synerex Server
	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Scenario2}")

	// Simulator
	timeStep := float64(1)
	globalTime := float64(0)
	sim = simulator.NewScenarioSimulator(timeStep, globalTime)

	// Communicator
	com = communicator.NewScenarioCommunicator()

	wg := sync.WaitGroup{}
	wg.Add(1)
	com.RegistClients(client, argJson)	// channelごとのClientを作成
	
	com.SubscribeAll(demandCallback, supplyCallback, &wg)	// ChannelにSubscribe
	wg.Wait()

	wg.Add(1)
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
