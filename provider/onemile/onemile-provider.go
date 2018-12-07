package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/synerex/synerex_alpha/api/fleet"

	socketio "github.com/googollee/go-socket.io"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	version = "0.01"

	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	client     api.SynerexClient

	port   = flag.Int("port", 7777, "Onemile Provider Listening Port")
	ioserv *socketio.Server

	n      = flag.Int("n", 1, "Number of taxi (or display)")
	dispWg sync.WaitGroup

	pitch = flag.Int("pitch", 2, "Pitch of vehicle status update (sec)")
)

// vehicle
type vehicle struct {
	VehicleId   string          `json:"vehicle_id"`   // unique id
	VehicleType string          `json:"vehicle_type"` // [onemile | bus | train | ...]
	Status      string          `json:"status"`       // [pickup | free | ride]
	Coord       [2]float64      `json:"coord"`        // current position (lat/lng)
	Angle       float32         `json:"angle"`        // current angle
	socket      socketio.Socket `json:"-"`            // Socket.IO socket
	Mission     *mission        `json:"mission"`      // assigned mission
	mu          sync.RWMutex    `json:"-"`            // mutex lock for vehicle read/write
}

// mission
type mission struct {
	MissionId string  `json:"mission_id"` // mission id
	Title     string  `json:"title"`      // mission title (option)
	Detail    string  `json:"detail"`     // mission detail (option)
	Events    []event `json:"events"`     // events
	Accepted  bool    `json:"accepted"`   // mission accepted by client?
}

// event
type event struct {
	EventId     string       `json:"event_id"`    // event id
	EventType   string       `json:"event_type"`  // [pickup | ride]
	StartTime   int64        `json:"start_time"`  // event start time (msec)
	EndTime     int64        `json:"end_time"`    // event end time (msec)
	Destination string       `json:"destination"` // destination (option)
	Route       [][2]float64 `json:"route"`       // routing (lat/lng)
	Status      string       `json:"status"`      // [none | start | end]
}

// managed vehicles by onemile-provider
var vehicleMap = make(map[string]*vehicle)

// display
type display struct {
	dispId string          // display id
	socket socketio.Socket // Socket.IO socket
	wg     sync.WaitGroup  // for synchronization to display ad and enquate
}

// taxi/display mapping
var dispMap = make(map[string]*display)

// convert mission to map
func (m *mission) toMap() map[string]interface{} {
	return map[string]interface{}{
		"mission_id": m.MissionId,
		"title":      m.Title,
		"detail":     m.Detail,
		"events": func(events []event) []map[string]interface{} {
			ret := make([]map[string]interface{}, len(events))
			for i, evt := range events {
				ret[i] = evt.toMap()
			}
			return ret
		}(m.Events),
	}
}

// convert event to map
func (e event) toMap() map[string]interface{} {
	return map[string]interface{}{
		"event_id":    e.EventId,
		"event_type":  e.EventType,
		"start_time":  e.StartTime,
		"end_time":    e.EndTime,
		"destination": e.Destination,
		"route":       e.Route,
	}
}

// utility function for converting time in milliseconds
func toMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// register OnemileProvider to NodeServer
func registerOnemileProvider() {
	sxutil.RegisterNodeName(*nodesrv, "OnemileProvider", false)
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)
	sxutil.RegisterDeferFunction(saveVehicles)
	go sxutil.HandleSigInt()
}

// create SMServiceClient for a given ChannelType
func createSMServiceClient(ch api.ChannelType, arg string) *sxutil.SMServiceClient {
	// create grpc client (at onece)
	if client == nil {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())

		conn, err := grpc.Dial(*serverAddr, opts...)
		if err != nil {
			log.Fatalf("Fail to Connect Synerex Server: %v", err)
		}

		client = api.NewSynerexClient(conn)
	}

	// create SMServiceClient
	return sxutil.NewSMServiceClient(client, ch, arg)
}

// subscribe marketing channel
func subscribeMarketing(mktClient *sxutil.SMServiceClient) {
	// wait until completing display registration
	dispWg.Wait()

	ctx := context.Background()
	seen := make(map[string]struct{})

	mktClient.SubscribeDemand(ctx, func(clt *sxutil.SMServiceClient, dm *api.Demand) {
		if dm.GetDemandName() == "" {
			// Confirm
			log.Printf("Receive SelectSupply [id: %d, name: %s]\n", dm.GetId(), dm.GetDemandName())
			clt.Confirm(sxutil.IDType(dm.GetId()))

			// SubscribeMbus
			clt.SubscribeMbus(context.Background(), func(clt *sxutil.SMServiceClient, msg *api.MbusMsg) {
				// wait until emitting all display
				emit := sync.WaitGroup{}
				emit.Add(*n)

				// emit start event for each display
				for taxi := range dispMap {
					dispMap[taxi].wg.Add(1)
					go func(taxi, name string, payload interface{}) {
						// wait unti a taxi will depart
						dispMap[taxi].wg.Wait()
						// emit event
						dispMap[taxi].socket.Emit(name, payload)
						log.Printf("Emit [taxi: %s, name: %s, payload: %s]\n", taxi, name, payload)
						// count down emit
						emit.Done()
					}(taxi, "disp_start", msg.ArgJson)
				}

				// send Done msg (and receive next msg)
				emit.Wait()
				clt.SendMsg(context.Background(), &api.MbusMsg{ArgJson: `{"command": "Done"}`})
			})
		} else {
			// ProposeSupply
			if _, ok := seen[dm.GetDemandName()]; !ok {
				seen[dm.GetDemandName()] = struct{}{}
				log.Printf("Receive RegisterDemand [id: %d, name: %s]\n", dm.GetId(), dm.GetDemandName())
				sp := &sxutil.SupplyOpts{
					Target: dm.GetId(),
					Name:   "a display for advertising and enqueting",
				}
				clt.ProposeSupply(sp)
			}
		}
	})
}

func printStackTrace(skip int) {
	for depth := skip; ; depth++ {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		log.Printf("====> %d: %v:%d\n", depth, file, line)
	}
}

// run Socket.IO server for Onemile-Client and Onemile-Display-Client
func runSocketIOServer(rdClient, mktClient *sxutil.SMServiceClient) {
	ioserv, e := socketio.NewServer(nil)
	if e != nil {
		log.Fatal(e)
	}

	ioserv.On("connection", func(so socketio.Socket) {
		log.Printf("Connected from %s as %s\n", so.Request().RemoteAddr, so.Id())


		so.On("clt_profile",func(data interface{}) {
			// get profile
			log.Println("Got Client Profile:", data)

			defer func() { // for error recovery
				if err := recover(); err != nil {
					log.Printf("Can't convert json data: %s\n", err)
				}
			}()

			taxi := data.(map[string]interface{})["device_id"].(string)

			if v, ok := vehicleMap[taxi]; ok {
				v.socket = so
				log.Println("OneMile:Sock",so.Id(),"<-> VID:",v.VehicleId,":Taxi=", taxi)
			}else{
				log.Println("No Taxi for :",taxi)
			}
		})

		qerr :=	so.Emit("clt_who_are_you","")
		if qerr != nil{
			log.Printf("Error on who are you %v",qerr)
		}

		// [Client] login
		so.On("clt_login", func(data interface{}) (ret interface{}) {
			log.Printf("Receive clt_login from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic clt_login: %s\n", err)
					printStackTrace(2)
					ret = map[string]interface{}{"panic": err.(error).Error()}
				}
			}()

			taxi := data.(map[string]interface{})["device_id"].(string)

			vehicleId := "unknown"
			if v, ok := vehicleMap[taxi]; ok {
				vehicleId = v.VehicleId
				v.socket = so
				log.Println("OneMile:Sock",so.Id(),"<-> VID:",v.VehicleId,":Taxi=", taxi)
			}

			ret = map[string]interface{}{
				"act":  "clt_login",
				"code": 0,
				"results": map[string]interface{}{
					"provider_id": "onemile-provider",
					"vehicle_id":  vehicleId,
					"token":       "1234567890",
				},
			}

			log.Printf("Ack [taxi: %s, name: %s, payload: %v]\n", taxi, "clt_login", ret)
			return ret
		})

		// [Client] update position
		so.On("clt_update_position", func(data interface{}) {
			log.Printf("Receive clt_update_position from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic clt_update_position: %s\n", err)
					printStackTrace(2)
				}
			}()

			taxi := data.(map[string]interface{})["device_id"].(string)

			if v, ok := vehicleMap[taxi]; ok {
				// lock vehicle
				v.mu.Lock()
				defer v.mu.Unlock()

				// update vehicle coord
				log.Println("Got Event:", data)

				v.Coord[0] = data.(map[string]interface{})["latlng"].([]interface{})[0].(float64)
				v.Coord[1] = data.(map[string]interface{})["latlng"].([]interface{})[1].(float64)
				log.Printf("Update position [taxi: %s, coord:[%f, %f]\n", taxi, v.Coord[0], v.Coord[1])

				// we should send this to server.
				f_coord := &fleet.Fleet_Coord{
					Lat: float32(v.Coord[0]),
					Lon: float32(v.Coord[1]),
				}
				vid, _ := strconv.Atoi(taxi)
				flt := &fleet.Fleet{
					VehicleId: int32(vid + 1000),
					Status:    0,
					Coord:     f_coord,
				}
				spo := &sxutil.SupplyOpts{
					Name:  "SupplyFromOnemile",
					JSON:  "{OnemileFleet}",
					Fleet: flt,
				}
				rdClient.RegisterSupply(spo)

				return
			}

			log.Println("Update ignored: [taxi: %s]\n", taxi)
		})

		// [Client] accept mission
		so.On("clt_accept_mission", func(data interface{}) (ret interface{}) {
			log.Printf("Receive clt_accept_mission from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic clt_accept_mission: %s\n", err)
					printStackTrace(2)
					ret = map[string]interface{}{"panic": err.(error).Error()}
				}
			}()

			taxi := data.(map[string]interface{})["device_id"].(string)
			missionId := data.(map[string]interface{})["mission_id"].(string)

			if v, ok := vehicleMap[taxi]; ok {
				if v.Mission != nil && v.Mission.MissionId == missionId {
					// lock vehicle
					v.mu.Lock()
					defer v.mu.Unlock()

					v.Mission.Accepted = true
					log.Printf("Mission accepted: [taxi: %s, missionId: %s]\n", taxi, missionId)

					return map[string]interface{}{"code": 0}
				}
			}

			log.Printf("Mission ignored: [taxi: %s, missionId: %s]\n", taxi, missionId)
			return map[string]interface{}{"code": 1}
		})

		// [Client] start mission event
		so.On("clt_start_mission_event", func(data interface{}) (ret interface{}) {
			log.Printf("Receive clt_start_mission_event from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic clt_start_mission_event: %s\n", err)
					printStackTrace(2)
					ret = map[string]interface{}{"panic": err.(error).Error()}
				}
			}()

			taxi := data.(map[string]interface{})["device_id"].(string)
			missionId := data.(map[string]interface{})["mission_id"].(string)
			eventId := data.(map[string]interface{})["event_id"].(string)

			if v, ok := vehicleMap[taxi]; ok {
				if v.Mission != nil && v.Mission.MissionId == missionId {
					for _, evt := range v.Mission.Events {
						if evt.EventId == eventId {
							// lock vehicle
							v.mu.Lock()
							defer v.mu.Unlock()

							// update status
							evt.Status = "start"
							log.Printf("Event start: [tax: %s, missionId: %s, eventId: %s]\n", taxi, missionId, eventId)

							// start display for marketing
							if evt.EventType == "ride" {
								dispClt := dispMap[taxi]
								if dispClt != nil {
									dispClt.wg.Done()
								}
							}

							// update vehicle status
							v.Status = evt.EventType

							return map[string]interface{}{"code": 0}
						}
					}
				}
			}

			log.Printf("Event ignored: [taxi: %s, missionId: %s, event_id: %s]\n", taxi, missionId, eventId)
			return map[string]interface{}{"code": 1}
		})

		// [Client] end mission event
		so.On("clt_end_mission_event", func(data interface{}) (ret interface{}) {
			log.Printf("Receive clt_end_mission_event from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic clt_end_mission_event: %s\n", err)
					printStackTrace(2)
					ret = map[string]interface{}{"panic": err.(error).Error()}
				}
			}()

			taxi := data.(map[string]interface{})["device_id"].(string)
			missionId := data.(map[string]interface{})["mission_id"].(string)
			eventId := data.(map[string]interface{})["event_id"].(string)

			if v, ok := vehicleMap[taxi]; ok {
				if v.Mission != nil && v.Mission.MissionId == missionId {
					for i, evt := range v.Mission.Events {
						if evt.EventId == eventId {
							// lock vehicle
							v.mu.Lock()
							defer v.mu.Unlock()

							// update status
							evt.Status = "end"
							log.Printf("Event end: [tax: %s, missionId: %s, eventId: %s]\n", taxi, missionId, eventId)

							if i != len(v.Mission.Events)-1 {
								// emit next event if any
								m := v.Mission.Events[i+1].toMap()
								m["mission_id"] = v.Mission.MissionId
								emitToClient(taxi, "clt_mission_event", m)
							} else {
								// all event done
								v.Status = "free"
							}

							return map[string]interface{}{"code": 0}
						}
					}
				}
			}

			log.Printf("Event ignored: [taxi: %s, missionId: %s, event_id: %s]\n", taxi, missionId, eventId)
			return map[string]interface{}{"code": 1}
		})

		// [Display] register taxi and display mapping
		so.On("disp_register", func(data interface{}) {
			log.Printf("Receive disp_register from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic disp_register: %s\n", err)
					printStackTrace(2)
				}
			}()

			taxi := data.(map[string]interface{})["taxi"].(string)
			disp := data.(map[string]interface{})["disp"].(string)

			if _, ok := dispMap[taxi]; !ok {
				dispMap[taxi] = &display{dispId: disp, socket: so, wg: sync.WaitGroup{}}
				log.Printf("Register display [taxi: %s => display: %v]\n", taxi, dispMap[taxi])
				dispWg.Done()
			} else {
				dispMap[taxi].socket = so
				log.Printf("Update display [taxi: %s, display: %v]\n", taxi, dispMap[taxi])
			}
		})

		// [Display] complete ad and enquate
		so.On("disp_complete", func(data interface{}) {
			log.Printf("Receive disp_complete from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic disp_complete: %s\n", err)
					printStackTrace(2)
				}
			}()

			// marshal json
			bytes, err := json.Marshal(data)
			if err != nil {
				log.Printf("Marshal error: %s\n", err)
			}

			// send results via Mbus
			mktClient.SendMsg(context.Background(), &api.MbusMsg{ArgJson: string(bytes)})
		})

		// [DEBUG] (simulate departure or arrive of taxi in disp-test.html)
		so.On("depart", func(data interface{}) {
			log.Printf("Receive depart from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic depart: %s\n", err)
					printStackTrace(2)
				}
			}()

			taxi := data.(map[string]interface{})["taxi"].(string)

			dispMap[taxi].wg.Done()
		})
		so.On("arrive", func(data interface{}) {
			log.Printf("Receive arrive from %s [%v]\n", so.Id(), data)
		})

		// [DEBUG] (register mission in clt-test.html)
		so.On("clt_register_mission", func(data interface{}) (ret interface{}) {
			log.Printf("Receive clt_register_mission from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic clt_register_mission: %s\n", err)
					printStackTrace(2)
					ret = map[string]interface{}{"panic": err.(error).Error()}
				}
			}()

			taxi := data.(map[string]interface{})["device_id"].(string)

			if v, ok := vehicleMap[taxi]; ok {
				// convert to string
				bytes, err := json.Marshal(data)
				if err != nil {
					log.Printf("Marshal failed: %s\n", err)
					return map[string]interface{}{"code": 1}
				}

				// lock vehicle
				v.mu.Lock()
				defer v.mu.Unlock()

				// convert to mission
				v.Mission = &mission{}
				err = json.Unmarshal(bytes, v.Mission)
				if err != nil {
					log.Printf("Unmarshal failed: %s\n", err)
					return map[string]interface{}{"code": 1}
				}

				emitToClient(taxi, "clt_request_mission", v.Mission.toMap())

				log.Printf("Mission registerd: [taxi: %s, mission: %#v]\n", taxi, v.Mission)
				return map[string]interface{}{"code": 0}
			}

			log.Printf("Mission ignored: [taxi: %s, mission: %#v]\n", taxi, data)
			return map[string]interface{}{"code": 1}
		})

		// [DEBUG] (order first event in mission)
		so.On("clt_mission_event", func(data interface{}) (ret interface{}) {
			log.Printf("Receive clt_mission_event from %s [%v]\n", so.Id(), data)

			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic clt_mission_event: %s\n", err)
					printStackTrace(2)
					ret = map[string]interface{}{"panic": err.(error).Error()}
				}
			}()

			taxi := data.(map[string]interface{})["device_id"].(string)
			missionId := data.(map[string]interface{})["mission_id"].(string)

			if v, ok := vehicleMap[taxi]; ok {
				if v.Mission.MissionId == missionId {
					// lock vehicle
					v.mu.RLock()
					defer v.mu.RUnlock()

					m := v.Mission.Events[0].toMap()
					m["mission_id"] = v.Mission.MissionId
					emitToClient(taxi, "clt_mission_event", m)

					log.Printf("Event ordered: [taxi: %s, missionId: %s, eventId: %s]\n", taxi, missionId, v.Mission.Events[0].EventId)
					return map[string]interface{}{"code": 0}
				}
			}

			log.Printf("Mission ignored: [taxi: %s, missionId: %s]\n", taxi, missionId)
			return map[string]interface{}{"code": 1}
		})

		// by Kawaguchi to set next event
		so.On("clt_cancel_all", func(data interface{}) (ret interface{}) {
			log.Printf("Receive clt_cancel_all from %s [%v]\n", so.Id(), data)

			// reset vehicle status:
			for _, v := range vehicleMap {
				v.mu.Lock()
				v.Mission = nil
				v.Status = "free"
				v.mu.Unlock()
			}

			return map[string]interface{}{"code": 0}
		})

		so.On("clt_dump_vehicles", func(data interface{}) {
			log.Printf("Receive clt_dump_vehicles from %s [%v]\n", so.Id(), data)
			bytes, _ := json.Marshal(vehicleMap)
			log.Printf("vehiceMap: %s\n", string(bytes))
		})
	})

	ioserv.On("disconnection", func(so socketio.Socket) {
		log.Printf("Disconnected from %s as %s\n", so.Request().RemoteAddr, so.Id())
	})

	ioserv.On("error", func(so socketio.Socket, err error) {
		log.Printf("Websocket error: %s\n", err)
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", ioserv)
	serveMux.Handle("/", http.FileServer(http.Dir("./display-client")))

	log.Printf("Starting Socket.IO Server %s on port %d", version, *port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), serveMux)
	if err != nil {
		log.Fatal(err)
	}
}

// utility function for emitting event to client
func emitToClient(taxi, name string, payload interface{}) {
	if v, ok := vehicleMap[taxi]; ok {
		if v.socket != nil {
			v.socket.Emit(name, payload)
			log.Printf("emit %s: [taxi: %s, payload: %#v]\n", name, taxi, payload)
		}
	}
}

// send each vehicle status to all vehicles
func sendVehicleStatus(pitch int) {
	for {
		time.Sleep(time.Second * time.Duration(pitch))

		// all vehicle status
		m := map[string]interface{}{"vehicles": []interface{}{}}

		// add each vehicle status
		for _, v := range vehicleMap {
			v.mu.RLock()
			stat := map[string]interface{}{
				"vehicle_id":   v.VehicleId,
				"provider_id":  "onemile-provider",
				"vehicle_type": v.VehicleType,
				"status":       v.Status,
				"coord":        v.Coord,
			}
			v.mu.RUnlock()

			m["vehicles"] = append(m["vehicles"].([]interface{}), stat)
		}

		//		log.Printf("clt_vehicle_status: %v\n", m)

		// broadcast to all vehicles :
		for k, v := range vehicleMap {

			// also check the status of the mission/events.
			if v.Mission != nil && v.Mission.Accepted {
				evs := v.Mission.Events
				for i, ev := range evs {
					if ev.Status == "none" { // not started?
						tm := time.Now()
						evtm := ev.StartTime / 1000
						evns := (ev.StartTime % 1000) * 1000
						//						log.Println(tm, "is after? ", time.Unix(evtm, evns))
						if tm.After(time.Unix(evtm, evns)) {
							ms := ev.toMap()
							ms["mission_id"] = v.Mission.MissionId
							// we should start
							emitToClient(k, "clt_mission_event", ms)

							v.Mission.Events[i].Status = "start"
						}
					}
				}

			}

			emitToClient(k, "clt_vehicle_status", m)

		}

		// check for mission start.

	}
}

func saveVehicles() {
	// save vehicleMap json to file
	if bytes, err := json.Marshal(vehicleMap); err == nil {
		ioutil.WriteFile("vehicles.json", bytes, 0600)
		log.Printf("Save vehicles to vehicles.json")
	} else {
		log.Printf("Marshal failed: %s\n", err)
	}
}

func initVehicles() {
	// restore json if any
	if _, err := os.Stat("vehicles.json"); !os.IsNotExist(err) {
		if bytes, err := ioutil.ReadFile("vehicles.json"); err == nil {
			err = json.Unmarshal(bytes, &vehicleMap)
		}

		if err != nil {
			log.Printf("Restore failed: %s\n", err)
		}
	}

	// fix map size
	sz := *n
	if sz < len(vehicleMap) {
		// TODO: to be shrinked?
		sz = len(vehicleMap)
	}

	// init vehicles
	for i := 0; i < sz; i++ {
		var id = fmt.Sprintf("%02d", i+1)

		if _, ok := vehicleMap[id]; ok {
			// restored vehicle
			vehicleMap[id].mu = sync.RWMutex{}
		} else {
			// new vehicle
			vehicleMap[id] = &vehicle{"vehicle" + id, "onemile", "free", [2]float64{34.87101, 137.1774}, 0.0, nil, nil, sync.RWMutex{}}
		}
	}

	// print vehicles
	bytes, _ := json.Marshal(vehicleMap)
	log.Printf("VehicleMap: %s\n", string(bytes))
}

func main() {
	flag.Parse()

	initVehicles()

	// set number of display
	dispWg.Add(*n)

	// register onemile-provider
	registerOnemileProvider()

	var wg sync.WaitGroup

	wg.Add(1)
	// subscribe rideshare channel
	rdClient := createSMServiceClient(api.ChannelType_RIDE_SHARE, "{Onemile:RideShareDM}")
	rtClient := createSMServiceClient(api.ChannelType_ROUTING_SERVICE, "{Onemile:Routing}")
	go subscribeRideShare(rdClient, rtClient)

	wg.Add(1)
	// subscribe marketing channel
	mktClient := createSMServiceClient(api.ChannelType_MARKETING_SERVICE, "")
	go subscribeMarketing(mktClient)

	wg.Add(1)
	// start Websocket Server
	go runSocketIOServer(rdClient, mktClient)

	wg.Add(1)
	// send each vehicle status to all vehicles
	go sendVehicleStatus(*pitch)

	wg.Wait()
}
