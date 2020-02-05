package main

// Daemon code for Synergic Exchange
import (
	"bufio"
	//"bytes"
	//"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	//"runtime/debug"
	//"strconv"
	"math"
	"strings"
	"sync"
	"time"
	
	gosocketio "github.com/mtfelian/golang-socketio"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	"github.com/synerex/synerex_alpha/api/simulation/daemon"
	"google.golang.org/grpc"
	"math/rand"
	//"math"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"io/ioutil"
)

var (
	fcs *geojson.FeatureCollection
	geofile string
	version = "0.04"
	port = 9995
	assetsDir http.FileSystem
	server *gosocketio.Server = nil
	//providerMap map[string]*exec.Cmd
	runProviders map[string]*Provider
	providerMutex sync.RWMutex
	client daemon.SimDaemonClient
	providerArray []ProviderData
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



type ProviderData struct{
	CmdName     string
	Type ProviderType
	Description string
	SrcDir      string
	BinName     string
	GoFiles     []string
	Options     []Option
}

// Log
type Log struct{
	ID int
	Description string
}



// Order Info


type Coord struct{
	Longitude float64
	Latitude float64
}

/*type OrderOption struct{
	AgentNum int
	Time float64
	AreaCoord []*Coord
}

type Order struct {
	ID string
	Type   OrderType
	Name string
	Option OrderOption
}*/

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

type ProviderState struct{
	Area []*common.Coord
	AgentNum int
	Time float64
}

var providerStats = []*ProviderState{
	{
		Area: []*common.Coord{
			{Latitude: 35.156431, Longitude: 136.97285,},
			{Latitude: 35.156431, Longitude: 136.981308,},
			{Latitude: 35.153578, Longitude: 136.981308,},
			{Latitude: 35.153578, Longitude: 136.97285,},
		},
		AgentNum: 50,
		Time: 0.0,
	},
	{
		Area: []*common.Coord{
			{Latitude: 35.156431, Longitude: 136.97285,},
			{Latitude: 35.156431, Longitude: 136.981308,},
			{Latitude: 35.153578, Longitude: 136.981308,},
			{Latitude: 35.153578, Longitude: 136.97285,},
		},
		AgentNum: 50,
		Time: 0.0,
	},
	{
		Area: []*common.Coord{
			{Latitude: 35.156431, Longitude: 136.97285,},
			{Latitude: 35.156431, Longitude: 136.981308,},
			{Latitude: 35.153578, Longitude: 136.981308,},
			{Latitude: 35.153578, Longitude: 136.97285,},
		},
		AgentNum: 50,
		Time: 0.0,
	},
}

var mockProviderStats = []*ProviderState{
	{
		Area: []*common.Coord{
			{Latitude: 0, Longitude: 100,},
			{Latitude: 100, Longitude: 100,},
			{Latitude: 100, Longitude: 0,},
			{Latitude: 0, Longitude: 0,},
		},
		AgentNum: 50,
		Time: 0.0,
	},
}

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

func updateProvider(newProviderStats []*ProviderState) []*ProviderState{
	// write area devide algorithm
	updatedProviderStats := make([]*ProviderState, 0)
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
				Add(point1vecs[0], point1vecs[0]), Add(point1vecs[1], point1vecs[0]), Add(devPoint1, point1vecs[0]), Add(divPoint2, point1vecs[0]),
			}
			coords2 := []*common.Coord{
				Add(point1vecs[2], point1vecs[0]), Add(point1vecs[3], point1vecs[0]), Add(devPoint1, point1vecs[0]), Add(divPoint2, point1vecs[0]),
			}
			// 追加
			state1 := &ProviderState{
				Area: coords1,
				AgentNum: state.AgentNum/2,
				Time: state.Time,
			}
			state2 := &ProviderState{
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
			newProviderStats := make([]*ProviderState, 0)
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
}

func init() {
	//areaTest()
	geofile = "transit_points.geojson"
	fcs = loadGeoJson(geofile)
	runProviders = make(map[string]*Provider)
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

}

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

// プロバイダの状態確認
/*func checkRunning2() {

	go func(){
		for{
			providerMutex.RLock()
			time.Sleep(1*time.Second)
			fmt.Printf("time----------\n")
			for key, cx := range providerMap {
				pid := cx.Process.Pid
				fmt.Printf("provider----------\n")
				fmt.Printf("%5d: %-20s : \n", pid, key)
			}
			providerMutex.RUnlock()
			//return 
		}
	}()
}*/


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

/////////////// Handle Command /////////////
/////////////////////////////////////////

func runMyCmd(provider *Provider) {
	cmd := provider.Cmd
	cmdName := provider.Name

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
	provider.ID = cmd.Process.Pid
	runProviders[cmdName] = provider
	providerMutex.Unlock()

	// logを送る
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
			ID: provider.ID,
			Description: string(line),
		}

		bytes, err  := json.Marshal(logInfo)
		logjson := string(bytes)
		//fmt.Printf("Log is %v \n", logjson)
 
		if server != nil {
			server.BroadcastToAll("log", logjson)
		}
		log.Printf("[%s]:%s", cmdName, string(line))
	}
	
	//	log.Printf("[%s]:Now ending...",cmdName)
	log.Printf("[%s]:Now ending...", cmdName)

	cmd.Wait()
	providerMutex.Lock()
	//delete(providerMap, cmdName)
	delete(runProviders, cmdName)
	providerMutex.Unlock()

	log.Printf("Command [%s] closed\n", cmdName)
}


// run From SubCommand
func runProp(sc ProviderData) string { // start local node server
	log.Printf("run '%s'\n", sc.CmdName)
	providerMutex.RLock()
	//_, ok := providerMap[sc.CmdName]
	_, ok := runProviders[sc.CmdName]
	providerMutex.RUnlock()
	if ok {
		log.Printf("%s is already running\n", sc.CmdName)
		return sc.CmdName + " is already running" // return to se command
	}

	d, err := os.Getwd()
	if err != nil {
		log.Printf("%s", err.Error())
		return "cannot get dir: " + err.Error()
	}

	// get src dir
	srcpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../../" + sc.SrcDir)
	binpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../../" + sc.SrcDir + "/" + sc.BinName)
	fi, err := os.Stat(binpath)

	// obtain most latest source file time.
	modTime := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	for _, fn := range sc.GoFiles {
		sp := filepath.FromSlash(filepath.ToSlash(srcpath) + "/" + fn)
		ss, _ := os.Stat(sp)
		if ss.ModTime().After(modTime) {
			modTime = ss.ModTime()
		}
	}

	var cmd *exec.Cmd

	// check mod time
	if err == nil && fi.ModTime().After(modTime) { // check binary time
		cmd = exec.Command("./" + sc.BinName) // run binary
	} else {
		runArgs := append([]string{"run"}, sc.GoFiles...)
		log.Printf("runArgs: [%s] %d, %s %s", sc.CmdName, len(sc.GoFiles), runArgs[0], runArgs[1])
		cmd = exec.Command("go", runArgs...) // run go with srcfile
	}

	cmd.Dir = srcpath
	cmd.Env = getGoEnv()

	provider := &Provider{
		Name: sc.CmdName,
		Type: sc.Type,
		Option: createOption(),
		Cmd: cmd,
	}
	go runMyCmd(provider)
	// no way to check the command result...
	return "ok"
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


/////////////// Handle SubCommand /////////////
/////////////////////////////////////////

func handleRun(target string) string {
	var res string
		for _, sc := range providerArray {
			if sc.CmdName == target {
				res = runProp(sc)
				
				return res
			}
	}
	
	log.Printf("Can't find command %s", target)
	return "Can't find command " + target
}

func decideCoordInGeo(fcs *geojson.FeatureCollection) *common.Coord{

	geoCoords := fcs.Features[0].Geometry.(orb.MultiLineString)[0]
	transitNum := rand.Int63n(int64(len(geoCoords)-1))

	longitude := geoCoords[transitNum][0] + 0.0001 * rand.Float64()
	latitude := geoCoords[transitNum][1] + 0.0001 * rand.Float64()
	coord := &common.Coord{
		Longitude: longitude,
		Latitude: latitude,
	}
	log.Printf("coord: ", coord)

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

func createRandomAgentType() string {
	num := rand.Float32()
	if num > 0.5 {
		return "car"
	} else {
		return "pedestrian"
	}
}


func handleOrder(order *UIOrder) string {
	target := order.Name
	fmt.Printf("Target is : %v\n", target)
	switch target {
	case "SetClock":
		fmt.Printf("SetClock\n")
		return "ok"
	case "ClearClock":
		fmt.Printf("ClearClock\n")
		return "ok"
	case "SetAgents":
		fmt.Printf("SetAgents\n")
		return "ok"
	case "ClearAgents":
		fmt.Printf("ClearAgents\n")
		return "ok"
	case "StartClock":
		fmt.Printf("StartClock\n")
		return "ok"
	case "StopClock":
		fmt.Printf("StopClock\n")
		return "ok"
	default:
		err := "true"
		log.Printf("Can't find command %s", target)
		return err
	}
	/*for _, sc := range orderArray {
		if sc.Name == target {
			var res string
			switch target{
			case "SetClock":
				fmt.Printf("SetClock\n")
			case "ClearClock":
				fmt.Printf("ClearClock\n")
			case "SetAgents":
				fmt.Printf("SetAgents\n")
			case "ClearAgents":
				fmt.Printf("ClearAgents\n")
			case "StartClock":
				fmt.Printf("StartClock\n")
			case "StopClock":
				fmt.Printf("StopClock\n")
			}*/
			///if target == "Clock" {
				// JSONファイル読み込み
				/*jsonName := order.Option
				fmt.Printf("jsonName is : %v\n", order.Option)
				bytes, err := ioutil.ReadFile(jsonName)
				if err != nil {
					log.Fatal(err)
				}
				// JSONデコード
				var simData SimData

				if err := json.Unmarshal(bytes, &simData); err != nil {
					log.Fatal(err)
				}
				*/
			//}
			
			
			/*else if target == "SetClock" {
				message := &daemon.SetClockMessage{
					GlobalTime: float64(0),
					TimeStep:   float64(1),
				}
				r, err := client.SetClockOrder(ctx, message)
				if err != nil {
					log.Fatalf("could not order: %v", err)
				}
				log.Printf("Response: %s", r.Ok)

			} else if target == "SetAgents" {
				agentNum, _ := strconv.Atoi(order.Option)
				agents := make([]*agent.Agent, 0)

				for i := 0; i < agentNum; i++ {
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

				message := &daemon.SetAgentsMessage{
					Agents: agents,
				}
				r, err := client.SetAgentsOrder(ctx, message)
				if err != nil {
					log.Fatalf("could not order: %v", err)
				}
				log.Printf("Response: %s", r.Ok)

			} else if target == "StartClock" {

				message := &daemon.StartClockMessage{
					StepNum: uint64(1),
				}
				r, err := client.StartClockOrder(ctx, message)
				if err != nil {
					log.Fatalf("could not order: %v", err)
				}
				log.Printf("Response: %s", r.Ok)

			} else if target == "StopClock" {
				message := &daemon.StopClockMessage{}
				r, err := client.StopClockOrder(ctx, message)
				if err != nil {
					log.Fatalf("could not order: %v", err)
				}
				log.Printf("Response: %s", r.Ok)
			} else if target == "ClearAgents" {
				message := &daemon.ClearAgentsMessage{}
				r, err := client.ClearAgentsOrder(ctx, message)
				if err != nil {
					log.Fatalf("could not order: %v", err)
				}
				log.Printf("Response: %s", r.Ok)
			}*/

		/*	res = "ok"
			log.Printf("your order is %s", target)
			return res
		}
	}
	log.Printf("Can't find command %s", target)
	return "Can't find command " + target*/
}


// ps commands from se cli
/*func checkRunning(opt string) []string {
	isLong := false
	if opt == "long" {
		isLong = true
	}
	var procs []string
	i := 0
	providerMutex.RLock()
	if isLong {
		procs = make([]string, len(providerMap)+2)
		str := fmt.Sprintf("  pid: %-20s : \n", "process name")
		procs[i] = str
		procs[i+1] = "-----------------------------------------------------------------\n"
		i += 2
	} else {
		procs = make([]string, len(providerMap))
	}
	for key, cx := range providerMap {
		pid := cx.Process.Pid
		if isLong {
			str := fmt.Sprintf("%5d: %-20s : \n", pid, key)
			procs[i] = str
		} else {
			if i != 0 {
				procs[i] = ", " + key
			} else {
				procs[i] = key
			}
		}
		i += 1
	}
	providerMutex.RUnlock()
	return procs

}*/

type UIOrder struct{
	Name string
	Options []*Option
}

// Order
type OrderType int
const (
	OrderType_SET_AGENTS  OrderType = 0
    OrderType_SET_AREA  OrderType = 1
    OrderType_SET_CLOCK  OrderType = 2
    OrderType_START_CLOCK  OrderType = 3
	OrderType_STOP_CLOCK OrderType = 4
)

type Option struct{
	Key string
	Value string
}

type Order struct {
	Type   OrderType
	Name string
	Options []*Option
}

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

type ProviderOption struct{
	NodeServAddr string
	SynerexServAddr string
	AreaCoord string
}

type Provider struct{
	ID int
	Name string
	Type ProviderType
	Option *ProviderOption
	Cmd *exec.Cmd
}

// simulator cliからの通信
func runSimulatorServer() error {
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
		//opt := param.(string)

		//return checkRunning(opt)
		return []string{"ok"}
	})
	

	server.On("run", func(c *gosocketio.Channel, param interface{}) string {
		targetName := param.(string)
		log.Printf("Get run command %s", targetName)

		ok := handleRun(targetName)
		sendRunnningProviders()
		return ok
	})

	server.On("order", func(c *gosocketio.Channel, param *Order) string {
		targetName := param.Type
		log.Printf("Get order command %s", targetName)

		//ok := handleRun(targetName)
		//sendRunnningProviders()
		return "ok"
	})

	server.On("command", func(c *gosocketio.Channel, param *UIOrder) string {
		targetName := param.Name
		log.Printf("Get order command %s", targetName)

		return handleOrder(param)
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	serveMux.HandleFunc("/", assetsFileHandler)
	log.Println("Serving at localhost:9995...")
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), serveMux); err != nil {
		log.Panic(err)
	}

	return nil
}


func runScenarioServer(){
	// Connect Daemon Server：シナリオとの通信
	conn, err := grpc.Dial("127.0.0.1:9996", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client = daemon.NewSimDaemonClient(conn)
	fmt.Println("Succsess Connection with Daemon Server ")
}


func main() {

	// scenarioへの送信クライアントサーバ
	go runScenarioServer()

	// cli, monitorの受信サーバ
	go runSimulatorServer()
	
	// run server
	handleRun("NodeIDServer")
	time.Sleep(500 * time.Millisecond)
	handleRun("MonitorServer")
	time.Sleep(500 * time.Millisecond)
	handleRun("SynerexServer")
	time.Sleep(500 * time.Millisecond)
	handleRun("Clock")
	//checkRunning2()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
