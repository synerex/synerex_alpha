package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nkawa/gtfsparser"
	"github.com/nkawa/gtfsparser/gtfs"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/ptransit"
	"google.golang.org/grpc"
	"hash/crc32"
	"math"
	"strconv"
	"strings"

	//	"github.com/synerex/synerex_alpha/api/common"
	//	"github.com/synerex/synerex_alpha/api/ptransit"
	"github.com/synerex/synerex_alpha/sxutil"
	//	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The synerex server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	feedName      = flag.String("feed", "", "GTFS Feed Filename")
)

func demandPTCallback(clt *sxutil.SMServiceClient, sp *api.Demand) {
 sp.GetArg_PTService()
}


func subscribePTDemand(client *sxutil.SMServiceClient) {
	ctx := context.Background() //
	err := client.SubscribeDemand(ctx, demandPTCallback)
	// not finish.
	log.Printf("Error:Demand %s\n",err.Error())
}

// find closest index from shapepoints
func pointDistance(p1 *common.Point, lat2 float32, lon2 float32) float32{

	p2 := &common.Point{
		Latitude: float64(lat2),
		Longitude: float64(lon2),
	}
	d, _ := p1.Distance(p2)
	return float32(d)
}

func getClosestPoints(pts gtfs.ShapePoints, lat float32, lon float32,  indx int) (int, float32 ){
	dist := float32(math.MaxFloat32)
	di := -1
	p := &common.Point{
		Latitude:  float64(lat),
		Longitude: float64(lon),
	}
	dists := []float32{}
	idx  :=[]int{}
	for i, pt := range (pts[indx:]) {
		d := pointDistance(p, pt.Lat, pt.Lon)
		if d < dist {
			di = i+indx
			dist = d
		}
		if d < 15 { //less than 5m
			dists = append(dists,d)
			idx = append(idx,i+indx)
		}
	}
	if len(idx) ==0{
		dists = append(dists,dist)
		idx = append(idx,di)
		fmt.Printf("Distance stop!")
	}
//	fmt.Printf	("Got nearest %d %v %v\n",len(idx),idx,dists)
	return di, dist
}

func getLatLonFromRatio(shp *gtfs.Shape, from_stop *gtfs.Stop,to_stop *gtfs.Stop, ratio float32,shape_idx int) (float32, float32, int){
	// try to findout nearest point from shp.
	if ratio < 0.001 {
//		fmt.Printf("small ratio %f\n",ratio)
		return from_stop.Lat, from_stop.Lon, shape_idx
	}
	if ratio > 0.999 {
//		fmt.Printf("large ratio %f\n",ratio)
		return to_stop.Lat, to_stop.Lon, shape_idx
	}


	c1,d1 := getClosestPoints(shp.Points, from_stop.Lat, from_stop.Lon ,shape_idx)
	c2,d2 := getClosestPoints(shp.Points, to_stop.Lat, to_stop.Lon ,c1)

	if c1 == -1 || c2 ==- 1{
		fmt.Printf("Umm. failuer")
		fmt.Printf("Find closest index len:%d, %d, %.1f  -> %d, %.1f\n",len(shp.Points), c1,d1, c2, d2)
	}
	if  d2 > 15 {
		fmt.Printf("Distance %f, %s\n", d2, to_stop.Name)
	}

	var totalDist float32
	totalDist = 0


	if c1 < c2 {// ok order
//		fmt.Printf("C1:%d, C2:%d \n",c1, c2)
		for i := c1; i <c2; i++ {
			pt := shp.Points[i]
			pt2 := shp.Points[i+1]
			totalDist += pointDistance(&common.Point{Latitude:float64(pt.Lat),Longitude:float64(pt.Lon)},
							pt2.Lat, pt2.Lon)
		}
		ratioDist := totalDist * ratio // we need to step this.
		fmt.Printf("Dist:%f, Ratio: %f /  %f \n",totalDist,ratioDist,  ratio)
		var distance  float32
		distance = 0
		for i := c1; i <c2; i++ {
			pt := shp.Points[i]
			pt2 := shp.Points[i+1]
			diff := pointDistance(&common.Point{Latitude:float64(pt.Lat),Longitude:float64(pt.Lon)},
				pt2.Lat, pt2.Lon)
			if distance + diff >= ratioDist { // now We came here!
				restDist := ratioDist - distance
				ptRatio := restDist / totalDist // get ratio of pt -pt2
				fmt.Printf("lat,lon dist %.1f , %.1f, total %.1f ratio %.3f\n",restDist, distance, totalDist, ptRatio)
				return (pt.Lat*(1-ptRatio)+pt2.Lat*ptRatio), (pt.Lon *(1-ptRatio)+ pt2.Lon*ptRatio), c1
			}
			distance += diff
		}
		return	to_stop.Lat, to_stop.Lon, c2
	}else{
		fmt.Printf("Reverse order ")
		fmt.Printf("Find closest index len:%d, %d, %.1f  -> %d, %.1f\n",len(shp.Points), c1,d1, c2, d2)
	}

	return 0,0,0
}



//   st[ix-1].Departure_time __ bus(t) __   st[ix].Arrival_time
// return bus location with index and time.
func locWithTime(feed *gtfsparser.Feed, trip_id string,tp *gtfs.Trip, ix int,  t gtfs.Time, shape_idx int) (nexit_index int, lat float32, lon float32, shape_index int ){
	if ix < 0 {return ix, 0, 0, shape_idx} // just for error check
	st := tp.StopTimes
	if ix == 0 && t.Minus(st[0].Arrival_time) <= 0 { // before start
		return 0, st[0].Stop.Lat, st[0].Stop.Lon, shape_idx
	}

	if ix >= len(st) {
		fmt.Printf("Index Error %d at %s\n",ix, trip_id)
		return -1, 0, 0, 0
	}
	for ; t.Minus(st[ix].Arrival_time) >= 0;  { // arrived current dest.
		ix ++
		if len(st)==ix	{ // it was final station.
			return -1, st[ix-1].Stop.Lat, st[ix-1].Stop.Lon, shape_idx
		}
	}
	//
	duration := st[ix].Arrival_time.Minus(st[ix-1].Departure_time) // total time for next station
	dt := t.Minus(st[ix-1].Departure_time)  // time already past from last stoptime

//	shape_id := tp.Route.Id // shape_id shoud taken from tip!
//	trip_id
	shape_id := tp.Shape.Id

	shapes , ok := feed.Shapes[shape_id]
	if !ok {
		fmt.Printf("Can't find shape from shape_id %s\n", shape_id)
	}
	ratio := float32(dt)/float32(duration)
	if trip_id == "1020" {
		fmt.Printf("Get trip %s, Lonlat %d  time: %d, duration:%d  ratio:%f\n",trip_id, ix, dt, duration, ratio)
	}
	lat, lon ,shape_idx = getLatLonFromRatio(shapes,st[ix-1].Stop, st[ix].Stop, ratio, shape_idx)

	// we should use shape index. but just use interpolated..

	return ix, lat, lon, shape_idx

}
func deg2rad(deg float32) float64 {
	return float64(math.Pi *deg / 180.0)
}
func calcAngle(lat1 float32, lon1 float32, lat2 float32, lon2 float32) float32{
	a := 6378137.000
	b := 6356752.314
	e := math.Sqrt((math.Pow(a, 2) - math.Pow(b, 2)) / math.Pow(a, 2))

	lat0 :=(deg2rad(lat1)+deg2rad(lat2))/2.0
	dlat := deg2rad(lat1)-deg2rad(lat2)
	dlon := deg2rad(lon1)-deg2rad(lon2)

	W := math.Sqrt(1 - math.Pow(e, 2)*math.Pow(math.Sin(lat0), 2))
	M := a * (1 - math.Pow(e, 2)) / math.Pow(W, 3)
	N := a / W

	ddi := dlat * M
	ddk := dlon * N * math.Cos(lat0)

//	ret := float32(180.0*math.Atan2(ddi,ddk)/math.Pi)
//	fmt.Printf("Ret %f\n",ret)
	ret := float32(180 + 180.0*math.Atan2(ddk, ddi)/math.Pi)
	if ret >= 360 {
		ret -= 360
	}

	return ret

}


func supplyPTransitFeed(clt *sxutil.SMServiceClient, feed *gtfsparser.Feed){

//	go subscribePTDemand(clt) // wait for demand to give "TimeTable Info"

	tripStatus := make(map[string]int)

//	for { // run for every 2 secs.
//		time.Sleep(time.Second * 2)

//  we start think from current time
//	now := time.Now()

	now := time.Date(2018,11,10,8,40,0,0,time.Local)
	ct := 0
//		year, month, date := now.Date()
	lastLat := make(map[string]float32)
	lastLon := make(map[string]float32)
	lastAng := make(map[string]float32)
	shape_idx := make(map[string]int)
	for {
		t := gtfs.Time{
			Hour: int8(now.Hour()),
			Minute: int8(now.Minute()),
			Second: int8(now.Second()),
		}
		for k,v := range(feed.Trips){
			// we should filter calendar

			if strings.Contains(feed.Trips[k].Id, "土日") {
				continue
			}


			if tripStatus[k] < 0 {
				continue
			}

			rid , _ := strconv.ParseInt(feed.Trips[k].Route.Id,10,32)
			trip_id , nerr := strconv.ParseInt(feed.Trips[k].Id,10,32)

			if nerr != nil {

				trip_id = int64(crc32.ChecksumIEEE([]byte(feed.Trips[k].Id)))
				fmt.Printf("Convert:"+feed.Trips[k].Id+" -> %d \n", trip_id)
			}
			if trip_id == 0{
				trip_id = rid
			}


			st ,lat, lon, current_shape_idx:=	locWithTime(feed, k, v, tripStatus[k], t, shape_idx[k])
			if st < 0 {
				shape_idx[k] = 0
			}
			shape_idx[k] = current_shape_idx
			tripStatus[k] = st
//			fmt.Printf("%d:st %d: %s \n",rid, st, feed.Trips[k].StopTimes[st].Stop.Name)
			if st > 0 {
				var angle float32
				if math.Abs(float64(lastLat[k]-lat)) < 0.00001 && math.Abs(float64(lastLon[k]- lat)) > 0.00001 {
					angle = lastAng[k]
//					fmt.Printf("SameAngle : %.2f  ==",angle)
				} else {
					angle = calcAngle(lastLat[k], lastLon[k], lat, lon)
				}
				fmt.Printf("%s: %s, %d: %d, (%.4f,%.4f)-(%.4f,%.4f) angle:%.2f\n",now.Format("15:04:05"), k,trip_id, st, lastLat[k], lastLon[k],lat, lon, angle)


				place := common.NewPlace().WithPoint(&common.Point{
					Latitude: float64(lat),
					Longitude: float64(lon),
				})


				pts := &ptransit.PTService{
					VehicleId: int32(trip_id),
					Angle: float32(angle),
					Speed: int32(0.0),
					CurrentLocation: place,
					VehicleType: int32( feed.Trips[k].Route.Type),
				}
				lastLat[k]=lat
				lastLon[k]=lon
				lastAng[k] = angle

				smo := sxutil.SupplyOpts{
					Name:  "GTFS Supply",
					PTService: pts,
					JSON: "",
				}
				clt.RegisterSupply(&smo)

			}else{
				if lat != 0.0 {
					lastLat[k] = lat
					lastLon[k] = lon
//					lastAng[k] = angle
//					fmt.Printf("%s: %s, %d: %d, (%.4f,%.4f)-\n",now.Format("15:04:05"), k,rid, st, lastLat[k], lastLon[k])
				}
			}
		}
		time.Sleep(time.Millisecond * 500)
		now = now.Add(time.Second * 10)
		ct++
		if ct > 8000 {
			break
		}
	}

}


// start reading gtfs
func main(){
	flag.Parse()


	if len(*feedName) ==0 {
		fmt.Printf("Please speficy GTFS Feed name.")
		os.Exit(0)
	}

	feed := gtfsparser.NewFeed()

	err := feed.Parse(*feedName)
	if err != nil{
		fmt.Printf("Error %s\n",err.Error())
	}
	fmt.Printf("Done, parsed %d agencies, %d stops, %d routes, %d trips, %d fare attributes\n",
		len(feed.Agencies),len(feed.Stops), len(feed.Routes), len(feed.Trips), len(feed.FareAttributes))

//	for k, v := range feed.Stops{
//		fmt.Printf("[%s] %s (@ %f, %f)\n", k, v.Name, v.Lat, v.Lon)
//	}



// here for usual provider
	sxutil.RegisterNodeName(*nodesrv, "PTransit:"+*feedName, false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	//	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to connect synerex server: %v", err)
	}

	client := api.NewSynerexClient(conn)
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_PT_SERVICE,"")

	supplyPTransitFeed(sclient, feed)

	sxutil.CallDeferFunctions() // cleanup!

}
