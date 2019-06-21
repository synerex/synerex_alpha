package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/synerex/synerex_alpha/provider/rpa/selenium"

	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idList     []uint64
	spMap      map[uint64]*sxutil.SupplyOpts
	mu         sync.RWMutex
	bj         *BookingJson
)

type BookingJson struct {
	cid    string
	status string
	year   string
	month  string
	day    string
	week   string
	start  string
	end    string
	people string
	title  string
}

func init() {
	idList = make([]uint64, 0)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
}

func checkMonth(month string) time.Month {
	var t time.Month
	switch month {
	case "1":
		t = time.January
	case "2":
		t = time.February
	case "3":
		t = time.March
	case "4":
		t = time.April
	case "5":
		t = time.May
	case "6":
		t = time.June
	case "7":
		t = time.July
	case "8":
		t = time.August
	case "9":
		t = time.September
	case "10":
		t = time.October
	case "11":
		t = time.November
	case "12":
		t = time.December
	}
	return t
}

func isPasted(year string, month string, day string) bool {
	flag := false
	y, _ := strconv.Atoi(year)
	m := checkMonth(month)
	d, _ := strconv.Atoi(day)

	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Println("Failed to get location of JST:", err)
	}

	now := time.Now().In(location)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	then := time.Date(y, m, d, 0, 0, 0, 0, location)
	subtract := then.Sub(now)

	// fmt.Println("now:", now)
	// fmt.Println("then:", then)
	// fmt.Println("subtract:", subtract)

	if subtract >= 0 {
		flag = true
	}
	return flag
}

func parseByGjson(json string) *BookingJson {
	cid := gjson.Get(json, "cid").String()
	status := gjson.Get(json, "status").String()
	year := gjson.Get(json, "date.Year").String()
	month := gjson.Get(json, "date.Month").String()
	day := gjson.Get(json, "date.Day").String()
	week := gjson.Get(json, "date.Week").String()
	start := gjson.Get(json, "date.Start").String()
	end := gjson.Get(json, "date.End").String()
	people := gjson.Get(json, "date.People").String()
	title := gjson.Get(json, "date.Title").String()

	bj := &BookingJson{
		cid:    cid,
		status: status,
		year:   year,
		month:  month,
		day:    day,
		week:   week,
		start:  start,
		end:    end,
		people: people,
		title:  title,
	}

	fmt.Println("parseByGjson is called:", bj)
	return bj
}

func demandCallback(clt *sxutil.SMServiceClient, dm *api.Demand) {
	log.Println("Got Meeting demand callback")
	fmt.Println(dm.ArgJson)
	bj = parseByGjson(dm.ArgJson)

	if dm.TargetId != 0 { // selected

		// ---------- This Error MUST fix ----------
		// if flag := selenium.Execute(bj.year, bj.month, bj.day, bj.week, bj.start, bj.end, bj.people, bj.title); flag == true {
		// 	log.Println("Select the room!")
		// 	clt.Confirm(sxutil.IDType(dm.Id))
		// } else {
		// 	log.Println("Failed to execute selenium")
		// }

		log.Println("Select the room!")
		clt.Confirm(sxutil.IDType(dm.Id))

	} else { // not selected

		switch bj.status {
		case "checking":
			if flag := isPasted(bj.year, bj.month, bj.day); flag == true {
				if flag := selenium.Schedules(bj.year, bj.month, bj.day, bj.start, bj.end, bj.people); flag == true {
					// bj.status = "OK"
					// b, err := json.Marshal(&bj)
					// if err != nil {
					// 	fmt.Println("Failed to json marshal:", err)
					// }
					json := `{"flag":"OK","data":` + dm.ArgJson + `}`
					sp := &sxutil.SupplyOpts{
						Target: dm.Id,
						Name:   "Valid schedules",
						JSON:   json,
					}

					mu.Lock()
					pid := clt.ProposeSupply(sp)
					idList = append(idList, pid)
					spMap[pid] = sp
					mu.Unlock()
				} else {
					// bj.status = "NG"
					// b, err := json.Marshal(&bj)
					// if err != nil {
					// 	fmt.Println("Failed to json marshal:", err)
					// }
					json := `{"flag":"NG","data":` + dm.ArgJson + `}`
					sp := &sxutil.SupplyOpts{
						Target: dm.Id,
						Name:   "Invalid schedules",
						JSON:   json,
					}

					mu.Lock()
					pid := clt.ProposeSupply(sp)
					idList = append(idList, pid)
					spMap[pid] = sp
					mu.Unlock()
				}
			} else {
				// bj.status = "NG"
				// b, err := json.Marshal(&bj)
				// if err != nil {
				// 	fmt.Println("Failed to json marshal:", err)
				// }
				json := `{"flag":"NG","data":` + dm.ArgJson + `}`
				sp := &sxutil.SupplyOpts{
					Target: dm.Id,
					Name:   "Invalid schedules",
					JSON:   json,
				}

				mu.Lock()
				pid := clt.ProposeSupply(sp)
				idList = append(idList, pid)
				spMap[pid] = sp
				mu.Unlock()
			}
		default:
			fmt.Printf("Switch case of default(%s) is called\n", bj.status)
		}

	}
}

func subscribeDemand(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
	log.Println("Server closed... on Meeting provider")
}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "RPAMeetingProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client: Meeting}")
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_MEETING_SERVICE, argJson)

	wg.Add(1)
	go subscribeDemand(sclient)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!
}
