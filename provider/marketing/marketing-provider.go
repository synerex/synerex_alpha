package main

// Draft code for Advertisement Service Provider

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist     []uint64
	dmMap      map[uint64]*sxutil.DemandOpts
	done       bool
	send       bool
	res        bool
	wg         sync.WaitGroup
	layout     = "2006-01-02 15:04:05"
	logFile    = "anslog.txt"
)

func init() {
	idlist = make([]uint64, 10)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
	wg = sync.WaitGroup{} // for syncing other goroutines
}

func msgCallback(clt *sxutil.SMServiceClient, msg *pb.MbusMsg) {
	log.Println("Got Mbus Msg callback")
	jsonStr := msg.ArgJson

	jsonBytes := ([]byte)(jsonStr)
	var data map[string]interface{}

	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		log.Fatalf("fail to unmarshal: %v", err)
		return
	}

	switch {
	case data["command"] == "RESULTS":
		// save data
		if data["results"] != nil {
			log.Println("Save Ans Data")
			// should get vehicle ID from onemile provider.
			device_id := ""
			if data["device_id"] != nil{
				device_id = data["device_id"].(string)
			}

			file, err := os.OpenFile("anslog.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			t := time.Now()
			s, _ := json.Marshal(data["results"])
			fmt.Fprintln(file, t.Format(layout)+" "+device_id+", "+string(s))
		}

		if done && !send {
			send = true
			sendEnqMsg(clt)
		}
		res = true

	case data["command"] == "Done":
		// Got Ad finish
		done = true

		if res && !send {
			send = true
			sendEnqMsg(clt)
		}
	}
}

func subscribeMBus(client *sxutil.SMServiceClient) {
	go func() {
		ctx := context.Background() // 必要？
		client.SubscribeMbus(ctx, msgCallback)
		// comes here if channel closed
		log.Printf("SubscribeMBus:%d", client.MbusID)

	}()
}

func sendMsg(client *sxutil.SMServiceClient, msg string) {
	log.Printf("SendMsg:%d", client.MbusID)

	done = false
	res = false
	m := new(pb.MbusMsg)
	m.ArgJson = msg
	ctx := context.Background() // 必要？
	client.SendMsg(ctx, m)
}

func sendAdMsg(client *sxutil.SMServiceClient) {
	var url = "img/20151216-092610.jpg"
	var period = 10

	request := map[string]interface{}{
		"command": "CONTENTS",
		"contents": []interface{}{
			map[string]interface{}{
				"type":   "AD",
				"data":   url,
				"period": period,
			},
			map[string]interface{}{
				"type":   "AD",
				"data":   "img/20151216-092920.jpg",
				"period": period,
			},
			map[string]interface{}{
				"type":   "AD",
				"data":   "img/20151222-114801.jpg",
				"period": period,
			},
			map[string]interface{}{
				"type":   "AD",
				"data":   "img/20160315-181937.jpg",
				"period": period,
			},
			map[string]interface{}{
				"type":   "AD",
				"data":   "img/20160315-182035.jpg",
				"period": period,
			},
			map[string]interface{}{
				"type":   "AD",
				"data":   "img/20160315-182129.jpg",
				"period": period,
			},
			map[string]interface{}{
				"type":   "AD",
				"data":   "img/20160315-182203.jpg",
				"period": period,
			},
			map[string]interface{}{
				"type":   "AD",
				"data":   "img/20160315-183639.jpg",
				"period": period,
			},
			map[string]interface{}{
				"type":   "AD",
				"data":   "img/20170531-132015.jpg",
				"period": period,
			},
		},
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("fail to marshal: %v", err)
		return
	}

	sendMsg(client, string(jsonBytes))
}

func sendEnqMsg(client *sxutil.SMServiceClient) {
	var enq_json = []byte(`{"questions":[{"label":"年齢","type":"select","name":"age","option":{"multiple":"false","options":[{"value":"0","text":"10代未満"},{"value":"10","text":"10代"},{"value":"20","text":"20代"},{"value":"30","text":"30代"},{"value":"40","text":"40代"},{"value":"50","text":"50代"},{"value":"60","text":"60代"},{"value":"70","text":"70代"},{"value":"80","text":"80代"},{"value":"90","text":"上記以外"}]}},{"label":"性別","type":"select","name":"sex","option":{"multiple":"false","options":[{"value":"0","text":"女性"},{"value":"1","text":"男性"}]}},{"label":"職業","type":"select","name":"occupation","option":{"multiple":"false","options":[{"value":"0","text":"会社員"},{"value":"1","text":"公務員"},{"value":"2","text":"自営業"},{"value":"3","text":"学生"},{"value":"4","text":"主婦"},{"value":"5","text":"アルバイト/パート"},{"value":"6","text":"無職"},{"value":"7","text":"その他"}]}},{"label":"ご自宅から駅までの交通手段を教えてください（複数回答可能）","type":"checkbox","name":"transportation","option":{"multiple":"true","options":[{"value":"0","text":"自動車（自分で運転する）"},{"value":"1","text":"自動車（送迎してもらう）"},{"value":"2","text":"オートバイ/原付"},{"value":"3","text":"JR"},{"value":"4","text":"えこたんバス"},{"value":"5","text":"自転車"},{"value":"6","text":"徒歩"},{"value":"7","text":"タクシー"}]}},{"label":"コミュニティバスへの乗り換えがスマートなモビリティが実現した場合に利用しますか？","type":"range","name":"fiveGrade1","option":{"max":"5","maxText":"利用する","min":"1","minText":"利用しない","options":[{"value":"2","text":"2"},{"value":"3","text":"3"},{"value":"4","text":"4"}]}},{"label":"移動中に広告を見たりアンケートに答えることで利用料が安くなれば","type":"range","name":"fiveGrade2","option":{"max":"5","maxText":"積極的に使う","min":"1","minText":"使わない","options":[{"value":"2","text":"2"},{"value":"3","text":"3"},{"value":"4","text":"4"}]}},{"label":"宅配などのモノの移動と一緒に移動することによって利用料が安くなれば","type":"range","name":"fiveGrade3","option":{"max":"5","maxText":"時間を要しても使う","min":"1","minText":"使わない","options":[{"value":"2","text":"2"},{"value":"3","text":"3"},{"value":"4","text":"4"}]}},{"label":"アプリの操作について","type":"range","name":"fiveGrade4","option":{"max":"5","maxText":"使いやすい","min":"1","minText":"使いにくい","options":[{"value":"2","text":"2"},{"value":"3","text":"3"},{"value":"4","text":"4"}]}},{"label":"お買い物や通院など、自宅からJR駅までスマートなモビリティ・サービスを利用する場合、片道が何円であれば利用しますか？","type":"select","name":"cost","option":{"multiple":"false","options":[{"value":"100","text":"100円"},{"value":"200","text":"200円"},{"value":"300","text":"300円"},{"value":"400","text":"400円"},{"value":"500","text":"500円"},{"value":"600","text":"600円"},{"value":"700","text":"700円"},{"value":"800","text":"800円"},{"value":"900","text":"900円"},{"value":"1000","text":"1000円"}]}},{"label":"ご感想","type":"textarea","name":"thoughts","option":{"placeholder":"研究のためにご感想をお寄せください。","options":[]}}]}`)
	var enq interface{}
	err := json.Unmarshal(enq_json, &enq)

	request := map[string]interface{}{
		"command": "CONTENTS",
		"contents": []interface{}{
			map[string]interface{}{
				"type":   "ENQ",
				"data":   enq,
				"period": 0,
			},
		},
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("fail to marshal: %v", err)
		return
	}

	sendMsg(client, string(jsonBytes))
}

func processMBus(clt *sxutil.SMServiceClient) {
	go subscribeMBus(clt)

	sendAdMsg(clt)
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got marketing supply callback:" + sp.GetSupplyName())

	// choice is supply for me? or not.
	if clt.IsSupplyTarget(sp, idlist) {
		// always select Supply
		clt.SelectSupply(sp)

		processMBus(clt)
	}

}

func subscribeSupply(client *sxutil.SMServiceClient) {
	// ここは goroutine!
	ctx := context.Background() // 必要？
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
}

func addDemand(sclient *sxutil.SMServiceClient, nm string) {
	opts := &sxutil.DemandOpts{Name: nm}
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id)
	dmMap[id] = opts
	if len(idlist)>10{ // remove
		rid := idlist[0]
		idlist = idlist[1:]
		delete(dmMap,rid)
	}
}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "MarketingProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure()) // only for draft version
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Marketing}")
	// create client wrapper
	sclient := sxutil.NewSMServiceClient(client, pb.ChannelType_MARKETING_SERVICE, argJson)

	wg.Add(1)
	go subscribeSupply(sclient)

	for {
		addDemand(sclient, "Kota-City citizen")
		time.Sleep(time.Second * time.Duration(10+rand.Int()%10))
	}

	wg.Wait()

	sxutil.CallDeferFunctions() // cleanup!
}
