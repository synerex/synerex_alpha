package main

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"time"
)

// Ecotan for XXX <-> Aimi
// as of 2018/12/8


type BusStop struct {
	trip  int
	name  string
	time  string
	lat float64
	lon float64
}
var busStops = []BusStop{
	{10,"役場","13:35",34.8643,137.1659},
	{11,"横落郷前","13:37",34.86381,137.1724},
	{12,"横落児童館","13:37",34.86573,137.1744},
	{13,"町民会館・図書館","13:40",34.87212,137.1801},
	{14,"町民会館南","13:41",34.87101,137.1774},
	{15,"商工会館西","13:43",34.87124,137.1752},
	{16,"大草松山（医療団地内）","13:44",34.87389,137.174},
	{17,"菱池矢尻","13:46",34.87475,137.1722},
	{18,"こうた眼科クリニック","13:47",34.87669,137.1712},
	{19,"幸田小学校","13:47",34.87953,137.1719},
	{20,"大草大正","13:48",34.88339,137.1723},
	{21,"高力熊谷","13:49",34.88707,137.1691},
	{22,"むらかみ整形外科","13:50",34.89168,137.1676},
	{23,"カメリアガーデン東","13:52",34.88844,137.1641},
	{24,"ＪＲ相見駅","13:53",34.88746,137.1606},

	{50,"ＪＲ相見駅","14:10",34.88746,137.1606},
	{51,"カメリアガーデン西","14:12",34.8885,137.164},
	{52,"むらかみ整形外科","14:14",34.89181,137.1677},
	{53,"高力熊谷","14:15",34.88729,137.1689},
	{54,"大草大正","14:16",34.88337,137.1724},
	{55,"幸田小学校","14:17",34.87925,137.1719},
	{56,"こうた眼科クリニック","14:17",34.87702,137.1714},
	{57,"菱池矢尻","14:18",34.87475,137.1722},
	{58,"大草松山（医療団地内）","14:20",34.87389,137.174},
	{59,"商工会館西","14:21",34.87134,137.1751},
	{60,"町民会館南","14:23",34.87101,137.1774},
	{61,"町民会館・図書館","14:24",34.87212,137.1801},
	{62,"横落児童館","14:26",34.86574,137.1745},
	{63,"横落郷前","14:26",34.86368,137.1723},
	{64,"役場","14:28",34.8643,137.1659},
}

type ShapeLoc struct{
	lat, lon float64
	num int
}

var lineShape = []ShapeLoc{
	{34.872117,137.180105,21},
	{34.8720456960156,137.179007126852,22},
	{34.8708841460199,137.178405148664,23},
	{34.8710123,137.177364,24},
	{34.8713379,137.175073,25},
	{34.8715336044152,137.17362754679,26},
	{34.8730356569777,137.174030231966,27},
	{34.8733844674077,137.174081193608,28},
	{34.8736473074207,137.173975187656,29},
	{34.8738884,137.174035,30},
	{34.8745164201022,137.174099386728,31},
	{34.8747598,137.172045,32},
	{34.8748381478923,137.170789537506,33},
	{34.876694,137.171227,34},
	{34.8795303,137.1718544,35},
	{34.8811303664493,137.172317466741,36},
	{34.8818060908628,137.172441780897,37},
	{34.8833874,137.1722943,38},
	{34.8838199361432,137.172277933671,39},
	{34.8847872266681,137.171622657874,40},
	{34.886690544001,137.169841202839,41},
	{34.8871124,137.1690354,42},
	{34.8874790749114,137.168425696215,43},
	{34.8911186322117,137.167659518296,44},
	{34.8916821,137.1676004,45},
	{34.892206841538,137.167633863531,46},
	{34.8922712508996,137.167467799378,47},
	{34.8917432482337,137.166924826491,48},
	{34.8884391,137.1641135,49},
	{34.8869394947474,137.162742330426,50},
	{34.887411248764,137.161805422352,51},
	{34.8874896445564,137.161356443514,52},
	{34.8874578604852,137.16094315052,53},
	{34.8871278604852,137.160934404488,54},
	{34.8871410463596,137.160620468641,55},
	{34.8874579,137.160573,56},
}


var stopOrder = [][2]int{
{3,14 },{3,14 },{9,14},
{15,26},{15,26},{15,30},
}

var shapeOrder =[][2]int{
{0, 35}, {0, 35}, {14, 35},
{35, 0}, {35, 0}, {35, 14},
}



type BusStopTime struct{
	id int
	name string
	hour, min int
	lat,lon float64
}
var busStopTimes = make([]BusStopTime,0)

func init (){
	fmt.Printf("stops: %v",len(busStops))

	for _, bs := range busStops {
		var b BusStopTime
		s ,_ := time.Parse("15:04", bs.time)
		b.hour = s.Hour()
		b.min = s.Minute()
		b.id = bs.trip
		b.lat = bs.lat
		b.lon = bs.lon
		busStopTimes = append(busStopTimes, b)
	}
//	fmt.Printf("stops: %v",len(busStopTime)

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
		return int32(busStops[minIx].trip)
	}
	return 0
}

/*
func NewPoint(lon, lat float64) *common.Point {
	pt := new(common.Point)
	pt.Latitude = lat
	pt.Longitude = lat
	return pt
}*/

func fromToPoint(spid [2]int) []*common.Point{
	newLs := make([]*common.Point, 0)
	for i, lis := range lineShape{
		if i>=spid[0] {
			newLs = append(newLs,NewPoint(lis.lat, lis.lon))
		}
		if i >spid[1] {
			break
		}
	}
	return newLs
}

func getBusStopPoint(bst BusStopTime) *common.Point{
	tp := new(common.Point)
	tp.Longitude = bst.lon
	tp.Latitude = bst.lat
	return tp
}


// Experiment Specific Bus route
//
func getBusRouteForExp(ftid [2]int , spid  [2]int ) *rideshare.Route {
	fid := ftid[0]
	tid := ftid[1]
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	rt := new(rideshare.Route)
	rt.TrafficType = rideshare.TrafficType_BUS
	rt.TransportName = "えこたんバス"
	rt.TransportLine = "北ルート"
	rt.DepartPoint = common.NewPlace().WithPoint(getBusStopPoint(busStopTimes[fid]))
	rt.ArrivePoint = common.NewPlace().WithPoint(getBusStopPoint(busStopTimes[tid]))
	stTime := time.Date(2018,12,8,busStopTimes[fid].hour, busStopTimes[fid].min,0,0,jst)
	edTime := time.Date(2018,12,8,busStopTimes[fid].hour, busStopTimes[fid].min,0,0,jst)
	stTsp, _ := ptypes.TimestampProto(stTime)
	rt.DepartTime = common.NewTime().WithTimestamp(stTsp)
	edTsp, _ := ptypes.TimestampProto(edTime)
	rt.ArriveTime = common.NewTime().WithTimestamp(edTsp)
	rt.AmountTime = ptypes.DurationProto(edTime.Sub(stTime))
	rt.Points = fromToPoint(spid)
	rt.AmountPrice = 0 // yen
	rt.AmountSheets = 1//
	rt.AvailableSheets = 30//
	return rt

}


func getBusRoute(toStation bool, subj int) *rideshare.Route {
	if !toStation {
		subj+=3
	}
	return getBusRouteForExp(stopOrder[subj], shapeOrder[subj])
}
