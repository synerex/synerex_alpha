package main

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"log"
	"time"
)

// JR timetable for Naogya <-> Aimi
// as of 2018/12/8

type Trip struct {
	trip_id int32
}

type StopTime struct {
	trip_id int32
	stop_id int32
	arrival_min int32
	departure_min int32
	stop_sequence int32
}

type Stop struct {
	name  string
	id   int32
	points [2]float64
}

var stops = []Stop{
	{ "名古屋", 1141101, [2]float64{136.881638,35.170694}},
	{ "相見", 1150242, [2]float64{137.160248,34.887573}},
}

var stop_times = []StopTime{
/*	{ 1000, 1150242, 12*60+34, 12*60+34, 1},
	{ 1000, 1141101,  13*60+12, 13*60+12, 2},
	{ 1001, 1150242, 13*60+4, 13*60+4, 1},
	{ 1001, 1141101,  13*60+42, 13*60+42, 2},
*/	{ 1002, 1150242, 13*60+34, 13*60+34, 1},
	{ 1002, 1141101,  14*60+12, 14*60+12, 2},
/*	{ 1003, 1150242, 14*60+4, 14*60+4, 1},
	{ 1003, 1141101,  14*60+42, 14*60+42, 2},
	{ 1004, 1150242, 14*60+34, 14*60+34, 1},
	{ 1004, 1141101,  15*60+12, 15*60+12, 2},
	{ 1005, 1150242, 15*60+04, 15*60+04, 1},
	{ 1005, 1141101,  15*60+42, 15*60+42, 2},
	{ 1006, 1150242, 15*60+34, 15*60+34, 1},
	{ 1006, 1141101,  16*60+12, 16*60+12, 2},

	{ 2000, 1141101, 12*60+31, 12*60+31, 1},
	{ 2000, 1150242,  13*60+11, 13*60+11, 2},
	{ 2001, 1141101, 13*60+01, 13*60+01, 1},
	{ 2001, 1150242,  13*60+41, 13*60+41, 2},
	{ 2002, 1141101, 13*60+31, 12*60+31, 1},
	{ 2002, 1150242,  14*60+11, 14*60+11, 2},
*/
	{ 2003, 1141101, 14*60+01, 14*60+01, 1},
	{ 2003, 1150242,  14*60+41, 14*60+41, 2},
/*	{ 2004, 1141101, 14*60+31, 14*60+31, 1},
	{ 2004, 1150242,  15*60+11, 15*60+11, 2},
	{ 2005, 1141101, 15*60+01, 15*60+01, 1},
	{ 2005, 1150242,  15*60+41, 15*60+41, 2},
*/
}

var trips = []Trip{
	{1000},
	{1001},
	{1002},
	{1003},
	{1004},
	{1005},
	{1006},
	{2000},
	{2001},
	{2002},
	{2003},
	{2004},
}

func init (){
	fmt.Printf("stops: %v",stops)

}

//
func getStopTimesByStation(fst int32, est int32) []StopTime{
	var sts []StopTime
	var sub []StopTime
	tripid := stop_times[0].trip_id
	stflag := false

	for _, st := range stop_times {
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


func findClosestStation(pt0 *common.Point) int32{
	minDist := 100000.0
	minIx := -1

	for i, st := range stops{
		pt := new(common.Point)
		pt.Latitude = st.points[1]
		pt.Longitude = st.points[0]
		dt, _ := pt0.Distance(pt)
		if dt < minDist {
			minDist = dt
			minIx = i
		}
//		log.Printf("%f, min %f",dt, minDist)
	}
	if minIx != -1{
		log.Printf("find closest station from %v to %d",*pt0, stops[minIx].id)
		return stops[minIx].id
	}
	return 0
}

// filter sts by tripid
func filterTrip(sts []StopTime, tripid int32)[]StopTime{
	rst := make([]StopTime,0)
	for _, st := range sts {
		if st.trip_id == tripid {
			rst = append(rst,st)
		}
	}
	return rst
}

func findDepartureTrip(stFrom int32, stTo int32, dmin int32) []StopTime{
	var minTrip int32
	minTime := int32(100000)

	sts := getStopTimesByStation(stFrom, stTo)
	for _, st := range sts {
		if st.stop_id == stFrom && dmin < st.departure_min {
			if st.departure_min-dmin < minTime {
				minTime = st.departure_min - dmin
				minTrip = st.trip_id
			}
		}
	}
	if minTrip != 0 { // found Trip!
		return filterTrip(sts, minTrip)
	}
	return make([]StopTime,0)
}

func findArrivalTrip(stFrom int32, stTo int32, amin int32) []StopTime{
	var minTrip int32
	minTime := int32(100000)

	sts := getStopTimesByStation(stFrom, stTo)
	for _, st := range sts {
		if st.stop_id == stTo && amin > st.arrival_min {
			if amin- st.arrival_min  < minTime {
				minTime = amin - st.arrival_min
				minTrip = st.trip_id
			}
		}
	}
	if minTrip != 0 { // found Trip!
		return filterTrip(sts, minTrip)
	}
	return make([]StopTime,0)
}

func getStopPoint(stopid int32) *common.Point {
	pt := new (common.Point)
	for _, st := range stops {
		if st.id == stopid {
			pt.Longitude = st.points[0]
			pt.Latitude = st.points[1]
			break
		}
	}
	return pt
}



func getTrainRouteFromTime(from *common.Point,to *common.Point, departMin int32, arriveMin int32 ) *rideshare.Route {
	var sts []StopTime
	stFrom := findClosestStation(from)
	stTo := findClosestStation(to)

	if stFrom == 0 || stTo == 0 {
		log.Printf("Can't find proper station")
		return nil
	}

	if departMin > 0 { // find from departure
		sts = findDepartureTrip(stFrom, stTo, departMin)

	}else if arriveMin > 0 { // find from arrival
		sts = findArrivalTrip(stFrom, stTo, arriveMin)
	}else {
		log.Printf("Time not specified.")
	}
//	if len(sts) == 0 {
//		return sts
//	}

	// we assume len=2 of sts

	if len(sts) == 2 {
		fid := 0
		tid := 1
		if sts[1].stop_id == stFrom {
			fid = 1
			tid = 0
		}
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		rt := new(rideshare.Route)
		rt.TrafficType = rideshare.TrafficType_TRAIN
		rt.TransportName = "JR東海"
		rt.TransportLine = "JR東海道線"
		rt.DepartPoint = common.NewPlace().WithPoint(getStopPoint(sts[fid].stop_id))
		rt.ArrivePoint = common.NewPlace().WithPoint(getStopPoint(sts[tid].stop_id))
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


	return nil
	// now we have single trip.
}




func getTrainRouteFromTimeExp(toStation bool) *rideshare.Route {
	var stFrom, stTo int32
	var sl int
	if toStation {
		stFrom =  stops[1].id
		stTo = stops[0].id
		sl = 2
	}else{
		stFrom = stops[0].id
		stTo = stops[1].id
		sl = 0
	}

		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		rt := new(rideshare.Route)
		rt.TrafficType = rideshare.TrafficType_TRAIN
		rt.TransportName = "JR東海"
		rt.TransportLine = "JR東海道線"
		rt.DepartPoint = common.NewPlace().WithPoint(getStopPoint(stFrom))
		rt.ArrivePoint = common.NewPlace().WithPoint(getStopPoint(stTo))
		stTime := time.Date(2018,12,8,int(stop_times[sl].departure_min/60),int(stop_times[sl].departure_min%60),0,0,jst)
		edTime := time.Date(2018,12,8,int(stop_times[sl+1].arrival_min/60),int(stop_times[sl+1].arrival_min%60),0,0,jst)
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