package main

// Daemon code for Synergic Exchange
import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/gops/agent"
	"github.com/mtfelian/golang-socketio"

	"github.com/kardianos/service"
)

var version = "0.03"
var logger service.Logger
var errlog *log.Logger
var port = 9995
var isDaemon = false
var interrupt chan os.Signal

var assetsDir http.FileSystem
var server *gosocketio.Server

var providerMap map[string]*exec.Cmd
var providerMutex sync.RWMutex

type SynerexService struct {
}

type SubCommands struct {
	CmdName     string
	Description string
	SrcDir      string
	BinName     string
	GoFiles     []string
	RunFunc     func()
}

var cmdArray []SubCommands

func init() {
	//	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	providerMap = make(map[string]*exec.Cmd)
	providerMutex = sync.RWMutex{}

	cmdArray = []SubCommands{
		{
			CmdName: "All",
			RunFunc: runAllServ,
			GoFiles: nil,
		},
		{
			CmdName: "NodeIDServer",
			SrcDir:  "nodeserv",
			BinName: "nodeid-server",
			GoFiles: []string{"nodeid-server.go"},
		},
		{
			CmdName: "MonitorServer",
			SrcDir:  "monitor",
			BinName: "monitor-server",
			GoFiles: []string{"monitor-server.go"},
		},
		{
			CmdName: "SynerexServer",
			SrcDir:  "server",
			BinName: "synerex-server",
			GoFiles: []string{"synerex-server.go", "message-store.go"},
		},
		{
			CmdName: "Taxi",
			SrcDir:  "provider/taxi",
			BinName: "taxi-provider",
			GoFiles: []string{"taxi-provider.go"},
		},
		{
			CmdName: "Ad",
			SrcDir:  "provider/ad",
			BinName: "ad-provider",
			GoFiles: []string{"ad-provider.go"},
		},
		{
			CmdName: "Multi",
			SrcDir:  "provider/multi",
			BinName: "multi-provider",
			GoFiles: []string{"multi-provider.go"},
		},
		{
			CmdName: "User",
			SrcDir:  "provider/user",
			BinName: "user-provider",
			GoFiles: []string{"user-provider.go"},
		},
		{
			CmdName: "Fleet",
			SrcDir:  "provider/fleet",
			BinName: "fleet-provider",
			GoFiles: []string{"fleet-provider.go"},
		},
	}
}

func (sesrv *SynerexService) Start(s service.Service) error {
	go sesrv.run()
	return nil
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

func runMyCmd(cmd *exec.Cmd, cmdName string) {
	providerMutex.Lock()
	providerMap[cmdName] = cmd
	providerMutex.Unlock()

	pipe, err := cmd.StderrPipe()
	if err != nil {
		logger.Infof("Error for getting stdout pipe %s\n", cmd.Args[0])
		return
	}
	err = cmd.Start()
	if err != nil {
		logger.Infof("Error for executing %s %v\n", cmd.Args[0], err)
		return
	}
	logger.Infof("Starting %s..\n", cmd.Args[0])

	reader := bufio.NewReader(pipe)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			logger.Infof("Command [%s] EOF\n", cmdName)
			break
		} else if err != nil {
			logger.Infof("Err %v\n", err)
		}
		server.BroadcastToAll("log", "["+cmdName+"]"+string(line))
		logger.Infof("[%s]:%s", cmdName, string(line))
	}
	//	log.Printf("[%s]:Now ending...",cmdName)
	logger.Infof("[%s]:Now ending...", cmdName)

	cmd.Wait()
	providerMutex.Lock()
	delete(providerMap, cmdName)
	providerMutex.Unlock()

	logger.Infof("Command [%s] closed\n", cmdName)
}

// delete all bin files
func cleanCmd(sc SubCommands) string { // build local node server
	logger.Infof("clean '%s'\n", sc.CmdName)

	d, err := getRegisteredDir()
	if err != nil {
		logger.Errorf("%s", err.Error())
		return "cannot get dir: " + err.Error()
	}

	binpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../" + sc.SrcDir + "/" + binName(sc.BinName))
	_, err = os.Stat(binpath)

	if err != nil {
		return "none"
	}
	os.Remove(binpath)
	return "ok"
}

// bulid From SubCommand
func buildCmd(sc SubCommands) string { // build local node server
	logger.Infof("build '%s'\n", sc.CmdName)
	providerMutex.RLock()
	_, ok := providerMap[sc.CmdName]
	providerMutex.RUnlock()
	if ok {
		logger.Warningf("%s is already running\n", sc.CmdName)
		return sc.CmdName + " is already running" // return to se command
	}

	d, err := getRegisteredDir()
	if err != nil {
		logger.Errorf("%s", err.Error())
		return "cannot get dir: " + err.Error()
	}

	// get src dir
	srcpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../" + sc.SrcDir)
	binpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../" + sc.SrcDir + "/" + binName(sc.BinName))
	fi, err := os.Stat(binpath)
	sfi, _ := os.Stat(srcpath)

	// obtain most latest source file time.
	modTime := sfi.ModTime()
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
		// if binary is newer than sources
		return "ok"
	} else {
		runArgs := append([]string{"build"}, sc.GoFiles...)
		cmd = exec.Command("go", runArgs...) // run go build with srcfile
	}

	cmd.Dir = srcpath
	cmd.Env = getGoEnv()

	go runMyCmd(cmd, sc.CmdName)
	// no way to check the command result...
	return "ok"
}

// run From SubCommand
func runProp(sc SubCommands) string { // start local node server
	logger.Infof("run '%s'\n", sc.CmdName)
	providerMutex.RLock()
	_, ok := providerMap[sc.CmdName]
	providerMutex.RUnlock()
	if ok {
		logger.Warningf("%s is already running\n", sc.CmdName)
		return sc.CmdName + " is already running" // return to se command
	}

	d, err := getRegisteredDir()
	if err != nil {
		logger.Errorf("%s", err.Error())
		return "cannot get dir: " + err.Error()
	}

	// get src dir
	srcpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../" + sc.SrcDir)
	binpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../" + sc.SrcDir + "/" + binName(sc.BinName))
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
		cmd = exec.Command("./" + binName(sc.BinName)) // run binary
	} else {
		runArgs := append([]string{"run"}, sc.GoFiles...)
		logger.Infof("runArgs: [%s] %d, %s %s", sc.CmdName, len(sc.GoFiles), runArgs[0], runArgs[1])
		cmd = exec.Command("go", runArgs...) // run go with srcfile
	}

	cmd.Dir = srcpath
	cmd.Env = getGoEnv()

	go runMyCmd(cmd, sc.CmdName)
	// no way to check the command result...
	return "ok"
}

func getGoEnv() []string { // we need to get/set gopath
	d, _ := getRegisteredDir() // may obtain dir of se-daemon
	gopath := filepath.FromSlash(filepath.ToSlash(d) + "/../../../")
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
		gp, err := getRegisteredGoPath()
		if err == nil {
			newenv = append(newenv, gp)
		}
	}
	return newenv
}

func runAllServ() {
	runSubCmd("NodeIDServer")
	runSubCmd("MonitorServer")
	time.Sleep(1 * time.Second)
	runSubCmd("SynerexServer")
	time.Sleep(1 * time.Second)
	runSubCmd("Taxi")
	runSubCmd("Ad")
	runSubCmd("Multi")
	time.Sleep(1 * time.Second)
	runSubCmd("User")
}

func killAll() {
	killCmd("User")
	killCmd("Multi")
	killCmd("Ad")
	killCmd("Taxi")
	killCmd("SynerexServer")
	killCmd("MonitorServer")
	killCmd("NodeIDServer")
}

func buildAll() []string {
	resp := make([]string, len(cmdArray))
	i := 0
	for _, sc := range cmdArray {
		if sc.CmdName != "All" {
			resp[i] = buildSubCmd(sc.CmdName)
			i++
		}
	}
	return resp
}

func cleanAll() []string {
	resp := make([]string, len(cmdArray))
	i := 0
	for _, sc := range cmdArray {
		if sc.CmdName != "All" {
			resp[i] = cleanSubCmd(sc.CmdName)
			i++
		}
	}
	return resp
}

// obtain key of the cmdMap array
func getCmdArray() []string {
	keys := make([]string, len(cmdArray))
	for i, sc := range cmdArray {
		keys[i] = sc.CmdName
	}
	return keys
}

func runSubCmd(cmd string) {
	for _, sc := range cmdArray {
		if sc.CmdName == cmd {
			runProp(sc)
			return
		}
	}
	logger.Errorf("Can't find subcommand %s", cmd)
}

func buildSubCmd(cmd string) string {
	for _, sc := range cmdArray {
		if sc.CmdName == cmd {
			return buildCmd(sc)
		}
	}
	str := fmt.Sprintf("Can't find subcommand %s", cmd)
	return str
}
func cleanSubCmd(cmd string) string {
	for _, sc := range cmdArray {
		if sc.CmdName == cmd {
			return cleanCmd(sc)
		}
	}
	str := fmt.Sprintf("Can't find subcommand %s", cmd)
	return str
}

func handleBuild(target []string) []string {
	resp := make([]string, len(target))
	for i, proc := range target {
		if proc == "All" || proc == "all" {
			return buildAll()
			// you don't need to build others
		} else {
			resp[i] = buildSubCmd(proc)
		}
	}
	return resp

}
func handleClean(target []string) []string {
	resp := make([]string, len(target))
	for i, proc := range target {
		if proc == "All" || proc == "all" {
			return cleanAll()
			// you don't need to build others
		} else {
			resp[i] = cleanSubCmd(proc)
		}
	}
	return resp
}

func handleRun(target string) string {
	for _, sc := range cmdArray {
		if sc.CmdName == target {
			var res string
			if sc.RunFunc == nil {
				res = runProp(sc)
			} else {
				//			res = sc.RunFunc()
				res = "ok"
				sc.RunFunc()
			}
			return res
		}
	}
	logger.Infof("Can't find command %s", target)
	return "Can't find command " + target
}

func killCmd(target string) string {
	res := "no"
	providerMutex.RLock()
	cmd, ok := providerMap[target]
	providerMutex.RUnlock()
	if ok {
		logger.Infof("Try to stop %s\n", target)
		err := cmd.Process.Signal(os.Kill)
		if err != nil {
			res = err.Error()
		} else {
			logger.Infof("OK to kill %s\n", target)
			res = "ok"
			/*			err = cmd.Process.Release()
						if err != nil {
							logger.Info("Release failed on %s\n",target)
						}
			*/
		}
	}
	return res
}

func handleStop(target []string) []string {
	resp := make([]string, len(target))
	for i, proc := range target {
		if proc == "All" || proc == "all" {
			killAll()
			break // you don't need to kill others
		} else {
			resp[i] = killCmd(proc)
		}
	}
	return resp
}

func handleRestart(target []string) []string {
	resp := handleStop(target)
	for _, tg := range target {
		handleRun(tg)
	}
	return resp
}

// ps commands from se cli
func checkRunning(opt string) []string {
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

}

func interfaceToString(target interface{}) []string {
	procs := target.([]interface{})
	resp := make([]string, len(procs))
	for i, pp := range procs {
		resp[i] = pp.(string)
	}
	return resp
}

func (sesrv *SynerexService) run() error {
	logger.Info("Starting.. Synergic Engine:" + version)
	currentRoot, err := getRegisteredDir()
	if err != nil {
		logger.Errorf("se-daemon: Can' get registered directory: %s", err.Error())
	}
	d := filepath.Join(currentRoot, "dclient", "build")

	assetsDir = http.Dir(d)
	server = gosocketio.NewServer()

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		logger.Infof("Connected from %s as %s", c.IP(), c.Id())
		// we need to send providers array
		// send Provider info to the web client
		c.Emit("providers", getCmdArray())
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		logger.Infof("Disconnected from %s as %s", c.IP(), c.Id())
	})

	server.On("ps", func(c *gosocketio.Channel, param interface{}) []string {
		// need to check param short or long
		opt := param.(string)

		return checkRunning(opt)
	})

	server.On("stop", func(c *gosocketio.Channel, param interface{}) []string {
		procs := interfaceToString(param)
		return handleStop(procs)
	})

	server.On("restart", func(c *gosocketio.Channel, param interface{}) []string {
		procs := interfaceToString(param)
		return handleRestart(procs)
	})

	server.On("build", func(c *gosocketio.Channel, param interface{}) []string {
		procs := interfaceToString(param)
		return handleBuild(procs)
	})

	server.On("clean", func(c *gosocketio.Channel, param interface{}) []string {
		procs := interfaceToString(param)
		return handleClean(procs)
	})

	server.On("run", func(c *gosocketio.Channel, param interface{}) string {
		nid := param.(string)
		//		fmt.Printf("Get Run Command %s\n", nid)
		logger.Infof("Get run command %s", nid)
		return handleRun(nid)
	})

	serveMux := http.NewServeMux()

	serveMux.Handle("/socket.io/", server)
	serveMux.HandleFunc("/", assetsFileHandler)

	logger.Info("Starting Synerex Engine daemon on port ", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), serveMux)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (sesrv *SynerexService) Stop(s service.Service) error {
	// how to stop serv.
	logger.Info("Stopping Synerex Engine.")
	// stop all running sub commands
	killAll()

	return nil
}

func (sesrv *SynerexService) Manage(s service.Service) (string, error) {
	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	usage := "Usage: se-daemon install | uninstall | start | stop | status"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			// Install se-daemon as a service
			//
			err := registerCurrentDir() // for each Operating System.
			if err != nil {
				logger.Error(err)
			}

			return "Synerex Engine Installed.", s.Install()
		case "uninstall":
			err := removeRegisteredDir() // uninstall
			if err != nil {
				logger.Error(err)
			}
			return "Synerex Engine Uninstalled.", s.Uninstall()
		case "start":
			return "Start Synerex Engine", s.Start()
		case "stop":

			return "Stop Synerex Engine", s.Stop()
		case "status":
			st, err := s.Status()
			var stst string
			switch st {
			case service.StatusRunning:
				stst = "Running"
			case service.StatusStopped:
				stst = "Stopped"
			default:
				stst = "Unknown"
			}
			return "Status :" + stst, err
		default:
			return usage, nil
		}
	}

	err := s.Run()
	if err != nil {
		logger.Error(err)
	}
	return "", nil
}

func main() {
	// add gops agent.
	//	fmt.Println("Start gops agent")
	if gerr := agent.Listen(agent.Options{}); gerr != nil {
		log.Fatal(gerr)
	}

	serv := &SynerexService{}

	svcConfig := &service.Config{
		Name:        "SynerexEngineDaemon",
		DisplayName: "Synerex Engine Daemon",
		Description: "Synerex Engine Daemon for controlling agents",
	}

	svc, err := service.New(serv, svcConfig)

	if err != nil {
		errlog.Println("Error:", err)
		os.Exit(1)
	}

	errs := make(chan error, 5)
	logger, err = svc.Logger(errs)
	st, err := svc.Status()
	if st == service.StatusRunning {
		isDaemon = true
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	//	logger.Info("Starting Synerex Engine "+version)
	status, err := serv.Manage(svc)
	if err != nil {
		errlog.Println(status, "\nError:", err)
		os.Exit(1)
	}

	fmt.Println(status)
}
