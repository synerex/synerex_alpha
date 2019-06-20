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
)

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

func exeSelenium(date string) bool {
	// parse by gjson
	year := gjson.Get(date, "date.Year").String()
	month := gjson.Get(date, "date.Month").String()
	day := gjson.Get(date, "date.Day").String()
	week := gjson.Get(date, "date.Week").String()
	start := gjson.Get(date, "date.Start").String()
	end := gjson.Get(date, "date.End").String()
	people := gjson.Get(date, "date.People").String()
	title := gjson.Get(date, "date.Title").String()

	flag := isPasted(year, month, day)
	if flag == true {
		selenium.Execute(year, month, day, week, start, end, people, title)

		// // check schedules
		// if s := selenium.Schedules(year, month, day, start, end, people); s == true {
		// 	selenium.Execute(year, month, day, week, start, end, people, title)
		// }
	}
	return flag
}

func demandCallback(clt *sxutil.SMServiceClient, dm *api.Demand) {
	log.Println("Got Meeting demand callback")

	if dm.TargetId != 0 { // selected

		log.Println("Select the room!")
		clt.Confirm(sxutil.IDType(dm.Id))

	} else { // not selected

		if flag := exeSelenium(dm.ArgJson); flag == true {
			json := `{"flag":"true","data":` + dm.ArgJson + `}`

			sp := &sxutil.SupplyOpts{
				Target: dm.Id,
				Name:   "Option of meeting room",
				JSON:   json,
			}

			mu.Lock()
			pid := clt.ProposeSupply(sp)
			idList = append(idList, pid)
			spMap[pid] = sp
			mu.Unlock()
		} else {
			json := `{"flag":"false","data":` + dm.ArgJson + `}`

			sp := &sxutil.SupplyOpts{
				Target: dm.Id,
				Name:   "Invalid booking",
				JSON:   json,
			}

			mu.Lock()
			pid := clt.ProposeSupply(sp)
			idList = append(idList, pid)
			spMap[pid] = sp
			mu.Unlock()
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
