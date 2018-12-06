package main

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"log"
	"time"
)

// Ecotan for XXX <-> Aimi
// as of 2018/12/8


type BusStop struct {
	trip  int
	name  string
	id int
	time  string
	lat float64
	lon float64
}


var busStops = []BusStop{
//	{1,"役場","13:35",34.8643,137.1659},
//	{1,"横落郷前","13:37",34.86381,137.1724},
//	{1,"横落児童館","13:37",34.86573,137.1744},
//	{1,"町民会館・図書館","13:40",34.87212,137.1801},
	{1,"町民会館南",10,"13:41",34.87101,137.1774},
//	{1,"商工会館西","13:43",34.87124,137.1752},
//	{1,"大草松山（医療団地内）","13:44",34.87389,137.174},
//	{1,"菱池矢尻","13:46",34.87475,137.1722},
//	{1,"こうた眼科クリニック","13:47",34.87669,137.1712},
	{1,"幸田小学校",11,"13:47",34.87953,137.1719},
//	{1,"大草大正","13:48",34.88339,137.1723},
//	{1,"高力熊谷","13:49",34.88707,137.1691},
//	{1,"むらかみ整形外科","13:50",34.89168,137.1676},
//	{1,"カメリアガーデン東","13:52",34.88844,137.1641},
	{1,"ＪＲ相見駅",12,"13:53",34.88746,137.1606},

	{2,"ＪＲ相見駅",12,"14:10",34.88746,137.1606},
//	{2,"カメリアガーデン西","14:12",34.8885,137.164},
//	{2,"むらかみ整形外科","14:14",34.89181,137.1677},
//	{2,"高力熊谷","14:15",34.88729,137.1689},
//	{2,"大草大正","14:16",34.88337,137.1724},
	{2,"幸田小学校",11,"14:17",34.87925,137.1719},
//	{2,"こうた眼科クリニック","14:17",34.87702,137.1714},
//	{2,"菱池矢尻","14:18",34.87475,137.1722},
//	{2,"大草松山（医療団地内）","14:20",34.87389,137.174},
//	{2,"商工会館西","14:21",34.87134,137.1751},
	{2,"町民会館南",10,"14:23",34.87101,137.1774},
//	{2,"町民会館・図書館","14:24",34.87212,137.1801},
//	{2,"横落児童館","14:26",34.86574,137.1745},
//	{2,"横落郷前","14:26",34.86368,137.1723},
//	{2,"役場","14:28",34.8643,137.1659},

}
var busStopTime = make([]StopTime,0)

func init (){
	fmt.Printf("stops: %v",busStops)

	for i, bs := range busStops {
		s ,_ := time.Parse("15:04", bs.time)
		st := StopTime{
			int32(bs.trip),
			int32(bs.id),
			int32(s.Hour()*60+s.Minute()),
			int32(s.Hour()*60+s.Minute()),
			int32(i),
		}
		busStopTime = append(busStopTime, st)
	}

	fmt.Printf("stops: %v",busStopTime)

}

//
func getBusStopTimesByStation(fst int32, est int32) []StopTime{
	var sts []StopTime
	var sub []StopTime
	tripid := busStopTime[0].trip_id
	stflag := false

	for _, st := range busStopTime {
		if st.trip_id != tripid{
			tripid = st.trip_id
			stflag = false
		}
		if stflag == false && st.stop_id == fst {
			stflag = true
			sub = append(sub, st)
		}
		if stflag == true && st.stop_id == est {
			for _, nn := range sub {
				sts = append(sts, nn)
			}
			sts = append(sts, st)
			sub = make([]StopTime,0)
		}
	}
	return sts
}


func findClosestBusStop(pt0 *common.Point) int32{
	minDist := 100000.0
	minIx := -1

	for i, st := range busStops{
		pt := new(common.Point)
		pt.Latitude = st.lat
		pt.Longitude = st.lon
		dt, _ := pt0.Distance(pt)
		if dt < minDist {
			minDist = dt
			minIx = i
		}
	}
	if minIx != -1{
		return int32(busStops[minIx].id)
	}
	return 0
}

// filter sts by tripid
func filterBusTrip(sts []StopTime, tripid int32)[]StopTime{
	rst := make([]StopTime,0)
	for _, st := range sts {
		if st.trip_id == tripid {
			rst = append(rst,st)
		}
	}
	return rst
}

// Filter bus trip for departure time.
func findDepartureBusTrip(stFrom int32, stTo int32, dmin int32) []StopTime{
	var minTrip int32
	minTime := int32(100000000)

	sts := getBusStopTimesByStation(stFrom, stTo)

	log.Printf("Get Staions : %d  %d->%d  time:%d",len(sts), stFrom, stTo, dmin)
	for _, st := range sts {

		if st.stop_id == stFrom && dmin < st.departure_min {
			if st.departure_min-dmin < minTime {
				minTime = st.departure_min - dmin
				minTrip = st.trip_id
			}
		}
	}
	if minTrip != 0 { // found Trip!
		return filterBusTrip(sts, minTrip)
	}
	return make([]StopTime,0)
}

func findArrivalBusTrip(stFrom int32, stTo int32, amin int32) []StopTime{
	var minTrip int32
	minTime := int32(100000)

	sts := getBusStopTimesByStation(stFrom, stTo)
	for _, st := range sts {
		if st.stop_id == stTo && amin > st.arrival_min {
			if amin- st.arrival_min  < minTime {
				minTime = amin - st.arrival_min
				minTrip = st.trip_id
			}
		}
	}
	if minTrip != 0 { // found Trip!
		return filterBusTrip(sts, minTrip)
	}
	return make([]StopTime,0)
}

func getBusStopPoint(stopid int32) *common.Point {
	pt := new (common.Point)
	for _, st := range busStops {
		if int32(st.id) == stopid {
			pt.Longitude = st.lon
			pt.Latitude = st.lat
			break
		}
	}
	return pt
}


func getBusRouteFromTime(from *common.Point,to *common.Point, departMin int32, arriveMin int32 ) *rideshare.Route {
	var sts []StopTime
	stFrom := findClosestBusStop(from)
	stTo := findClosestBusStop(to)

	if stFrom == 0 || stTo == 0  || stFrom == stTo {
		log.Printf("Can't find proper station %d, %d", stFrom, stTo)
		return nil
	}

	if departMin > 0 { // find from departure
	log.Printf("Departure %d", departMin)
		sts = findDepartureBusTrip(stFrom, stTo, departMin)

	}else if arriveMin > 0 { // find from arrival
		log.Printf("Arrival")
		sts = findArrivalBusTrip(stFrom, stTo, arriveMin)
	}else {
		log.Printf("Time not specified.")
	}

	// we assume len=2 of sts

	if len(sts) >= 2 {
		fid := 0
		tid := len(sts)-1
		if sts[0].stop_id == stTo { // 順序が逆のこともありえる？
			fid = len(sts)-1
			tid = 0
		}
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		rt := new(rideshare.Route)
		rt.TrafficType = rideshare.TrafficType_TRAIN
		rt.TransportName = "えこたんバス"
		rt.TransportLine = "北ルート"
		rt.DepartPoint = common.NewPlace().WithPoint(getBusStopPoint(sts[fid].stop_id))
		rt.ArrivePoint = common.NewPlace().WithPoint(getBusStopPoint(sts[tid].stop_id))
		stTime := time.Date(2018,12,8,int(sts[fid].departure_min/60),int(sts[fid].departure_min%60),0,0,jst)
		edTime := time.Date(2018,12,8,int(sts[tid].arrival_min/60),int(sts[tid].arrival_min%60),0,0,jst)
		stTsp, _ := ptypes.TimestampProto(stTime)
		rt.DepartTime = common.NewTime().WithTimestamp(stTsp)
		edTsp, _ := ptypes.TimestampProto(edTime)
		rt.ArriveTime = common.NewTime().WithTimestamp(edTsp)
		rt.AmountTime = ptypes.DurationProto(edTime.Sub(stTime))
		rt.AmountPrice = 760 // yen
		rt.AmountSheets = 1//
		rt.AvailableSheets = 100//
		return rt
	}

	log.Printf("Can't make bus stations %d", len(sts))
	return nil
	// now we have single trip.
}