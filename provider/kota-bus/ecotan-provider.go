package main

import (
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/synerex/synerex_alpha/sxutil"
	"github.com/synerex/synerex_alpha/api/ptransit"
	api "github.com/synerex/synerex_alpha/api"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
)



var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The synerex server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	mqsrv      = flag.String("mqtt_serv", "", "MQTT Server Address host:port")
	mquser      = flag.String("mqtt_user", "", "MQTT Username")
	mqpass      = flag.String("mqtt_pass", "", "MQTT UserID")
	mqtopic      = flag.String("mqtt_topic", "", "MQTT Topic")
	idlist     []uint64
	spMap      map[uint64]*sxutil.SupplyOpts
	mu         sync.Mutex
)

func messageHandler(client *MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}


func subscribe(client MQTT.Client, sub chan<- MQTT.Message) {
	subToken := client.Subscribe(
		*mqtopic,
		0,
		func(client MQTT.Client, msg MQTT.Message) {
			sub <- msg
		})
	if subToken.Wait() && subToken.Error() != nil {
		fmt.Println(subToken.Error())
		os.Exit(1) // todo: should handle MQTT error
	}
}

func handleMQMessage(sclient *sxutil.SMServiceClient, msg MQTT.Message){

	smo := sxutil.SupplyOpts{
		Name:  "Ecotan Bus Info",
		Fleet: &fleet,
	}

	sclient.RegisterSupply()

}



func main() {

	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "EcotanProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
//	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := pb.NewSMarketClient(conn)
	sclient := sxutil.NewSMServiceClient(client, api.MarketType_PT_SERVICE,"")

	// MQTT

	mopts := MQTT.NewClientOptions()
	mopts.AddBroker(*mqsrv)
	mopts.SetClientID("synerex-provider")
	mopts.SetUsername(*mquser)
	mopts.SetPassword(*mqpass)

	mclient := MQTT.NewClient(mopts)
	token := mclient.Connect()
	token.Wait()
	merr := token.Error()
	if merr != nil {
		panic(merr)
	}

	sub := make(chan MQTT.Message)
	go subscribe(mclient, sub)  // subscribe MQTT channel

	for {
		select {
			case s := <- sub:
				msg := string(s.Payload())
				fmt.Printf("\nmsg: %s\n", msg)
				handleMQMessage(sclient, msg)
		}
	}

	sxutil.CallDeferFunctions() // cleanup!

}
