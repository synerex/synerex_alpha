package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

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
)

// vehicle
type vehicle struct {
	VehicleId   string     `json:"vehicle_id"`   // unique id
	VehicleType string     `json:"vehicle_type"` // [onemile | bus | train | ...]
	Status      string     `json:"status"`       // [pickup | free | ride]
	Coord       [2]float64 `json:"coord"`        // current position (lon/lat)
	sockId      string     `json:"-"`            // Socket.IO socket id
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

// register OnemileProvider to NodeServer
func registerOnemileProvider() {
	sxutil.RegisterNodeName(*nodesrv, "OnemileProvider", false)
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)
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

// TODO: 乗車シーケンス
// subscribe rideshare channel
func subscribeRideShare(rdClient, rtClient *sxutil.SMServiceClient) {
	ctx := context.Background()
	rdClient.SubscribeDemand(ctx, func(clt *sxutil.SMServiceClient, dm *api.Demand) {
		if dm.GetDemandName() == "" {
			// Confirm
			// TODO: 迎車処理 (routing-providerからのSelectSupply受信〜乗車まで)
		} else {
			// ProposeSupply
			// TODO: 経路取得 (routing-providerからのRegisterDemand受信〜ProposeSupplyまで)
		}
	})
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
				// emit start event for each display
				for taxi := range dispMap {
					dispMap[taxi].wg.Add(1)
					go func(taxi, name string, payload interface{}) {
						// wait unti a taxi will depart
						dispMap[taxi].wg.Wait()
						// emit event
						dispMap[taxi].socket.Emit(name, payload)
						log.Printf("Emit [taxi: %s, name: %s, payload: %s]\n", taxi, name, payload)
					}(taxi, "disp_start", msg.ArgJson)
				}
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
				v.sockId = so.Id()
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

			for k, v := range vehicleMap {
				if v.sockId == so.Id() {
					v.Coord[0] = data.(map[string]interface{})["latlng"].([]interface{})[0].(float64)
					v.Coord[1] = data.(map[string]interface{})["latlng"].([]interface{})[1].(float64)
					log.Printf("Update position [taxi: %s, coord:[%f, %f]\n", k, v.Coord[0], v.Coord[1])
				}
			}
		})

		// TODO: 車位置アップデート
		so.On("xxxxx", func(data interface{}) interface{} {
			return nil
		})

		// TODO: 移動処理 (乗車〜降車まで)
		so.On("xxxxx", func(data interface{}) interface{} {
			return nil
		})
		so.On("xxxxx", func(data interface{}) interface{} {
			return nil
		})

		// [Display] register taxi and display mapping
		so.On("disp_register", func(data interface{}) {
			log.Printf("Receive disp_register from %s [%v]\n", so.Id(), data)

			taxi := data.(map[string]interface{})["taxi"].(string)
			disp := data.(map[string]interface{})["disp"].(string)

			if _, ok := dispMap[taxi]; !ok {
				dispMap[taxi] = &display{dispId: disp, socket: so, wg: sync.WaitGroup{}}
				log.Printf("Register display [taxi: %s => display: %v]\n", taxi, dispMap[taxi])
				dispWg.Done()
			}
		})

		// [Display] complete ad and enquate
		so.On("disp_complete", func(data interface{}) {
			log.Printf("Receive disp_complete from %s [%v]\n", so.Id(), data)

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

			taxi := data.(map[string]interface{})["taxi"].(string)

			dispMap[taxi].wg.Done()
		})
		so.On("arrive", func(data interface{}) {
			log.Printf("Receive arrive from %s [%v]\n", so.Id(), data)
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

func main() {
	flag.Parse()

	// init vehicles
	for i := 0; i < *n; i++ {
		var id = fmt.Sprintf("%02d", i+1)
		vehicleMap[id] = &vehicle{"vehicle" + id, "onemile", "free", [2]float64{0.0, 0.0}, ""}
	}

	// set number of display
	dispWg.Add(*n)

	// register onemile-provider
	registerOnemileProvider()

	var wg sync.WaitGroup

	wg.Add(1)
	// subscribe rideshare channel
	rdClient := createSMServiceClient(api.ChannelType_RIDE_SHARE, "")
	rtClient := createSMServiceClient(api.ChannelType_ROUTING_SERVICE, "")
	go subscribeRideShare(rdClient, rtClient)

	wg.Add(1)
	// subscribe marketing channel
	mktClient := createSMServiceClient(api.ChannelType_MARKETING_SERVICE, "")
	go subscribeMarketing(mktClient)

	wg.Add(1)
	// start Websocket Server
	go runSocketIOServer(rdClient, mktClient)

	wg.Wait()
}
