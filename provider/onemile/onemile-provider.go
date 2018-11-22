package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mtfelian/golang-socketio"
)

var (
	version = "0.01"
	port    = flag.Int("port", 7777, "OneMile Provider Listening Port")
	ioserv  *gosocketio.Server
)

// display
type display struct {
	dispId string // display id
	chanId string // channel id
}

// taxi/display mapping
var dispMap = make(map[string]*display)

// run Socket.IO server for OneMile-Display-Client
func runSocketIOServer() {
	ioserv := gosocketio.NewServer()

	ioserv.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected from %s as %s\n", c.IP(), c.Id())
	})

	ioserv.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected from %s as %s\n", c.IP(), c.Id())
	})

	// register taxi and display mapping
	ioserv.On("register", func(c *gosocketio.Channel, data interface{}) {
		log.Printf("Receive register from %s [%v]\n", c.Id(), data)

		taxi := data.(map[string]interface{})["taxi"].(string)
		disp := data.(map[string]interface{})["disp"].(string)

		_, ok := dispMap[taxi]
		if !ok {
			dispMap[taxi] = &display{dispId: disp, chanId: c.Id()}
		}

		log.Printf("Register display [taxi: %s => display: %v]\n", taxi, dispMap[taxi])
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", ioserv)
	serveMux.Handle("/", http.FileServer(http.Dir("./display-client")))

	log.Printf("Starting OneMile Provider %s on port %d", version, *port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), serveMux)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	runSocketIOServer()
}
