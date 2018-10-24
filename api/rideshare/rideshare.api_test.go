package rideshare

import (
	"github.com/synerex/synerex_alpha/api/common"
	"math"
	"regexp"

	//"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	durpb "github.com/golang/protobuf/ptypes/duration"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

// invalid Timestamp (Seconds < minValidSeconds)
var invalidTimestamp = &tspb.Timestamp{Seconds: -62135596800 - 1, Nanos: 0}

// makeTimestamp makes protobuf Timestamp from given string.
// If s is empty and valid is true, returns nil.
// If s is empty and valid is false, returns invalid Timestamp.
// Otherwise, returns protobuf Timestamp
func makeTimestamp(s string, valid bool) (*tspb.Timestamp, error) {
	if s == "" {
		if valid {
			return nil, nil
		} else {
			return invalidTimestamp, nil
		}
	} else {
		ts, _ := time.Parse("2006-01-02 15:04", s)
		return ptypes.TimestampProto(ts)
	}
}

// makeDuration makes protobuf Duration from given string.
// If s is empty, returns nil.
// Otherwise, returns protobuf Duration.
func makeDuration(s string) (*durpb.Duration, error) {
	if s == "" {
		return nil, nil
	} else {
		d, _ := time.ParseDuration(s)
		return ptypes.DurationProto(d), nil
	}
}

// makeRouteTm makes Route with DepartTime, ArriveTime from given string.
// If given string is empty and valid is true, it is converted to nil.
// If given string is empty and valid is false, it is converted to invalid Timestamp.
// Otherwise, given string is converted valid Timestamp.
func makeRouteTm(deptTime, arrvTime string, valid bool) *Route {
	deptTimePb, _ := makeTimestamp(deptTime, valid)
	arrvTimePb, _ := makeTimestamp(arrvTime, valid)

	return &Route{
		DepartTime: &common.Time{Value: &common.Time_Timestamp{deptTimePb}},
		ArriveTime: &common.Time{Value: &common.Time_Timestamp{arrvTimePb}},
	}
}

// makeRoutePt makes Route with DepartPoint, ArrivePoint from given points.
func makeRoutePt(deptLat, deptLon, arrvLat, arrvLon float64) *Route {
	deptPoint := &common.Point{Latitude: deptLat, Longitude: deptLon}
	arrvPoint := &common.Point{Latitude: arrvLat, Longitude: arrvLon}

	return &Route{
		DepartPoint: &common.Place{Value: &common.Place_Point{deptPoint}},
		ArrivePoint: &common.Place{Value: &common.Place_Point{arrvPoint}},
	}
}

// makeRoutePr makes Route with AmountPrice from given price.
func makeRoutePr(price uint32) *Route {
	return &Route{AmountPrice: price}
}

func makeRideShare(routes []*Route) *RideShare {
	return &RideShare{Routes: routes}
}

// test CalcAmountTime (Success)
func TestRoute_CalcAmountTime_Success(t *testing.T) {
	tests := []struct {
		deptTime string
		arrvTime string
		want     string
	}{
		{"", "", ""},
		{"2018-10-05 11:00", "", ""},
		{"", "2018-10-05 12:30", ""},
		{"2018-10-05 11:00", "2018-10-05 12:30", "1h30m"},
	}

	for i, test := range tests {
		route := makeRouteTm(test.deptTime, test.arrvTime, true)

		want, _ := makeDuration(test.want)
		got, _ := route.CalcAmountTime()

		if want.String() != got.String() {
			t.Errorf("[%d] want=%v, got=%v\n", i, want, got)
		}
	}
}

// test CalcAmountTime (Error)
func TestRoute_CalcAmountTime_Error(t *testing.T) {
	tests := []struct {
		deptTime string
		arrvTime string
		regexp   string
	}{
		{"", "2018-10-05 12:30", "depart_time=timestamp: .*, arrive_time=<nil>"},
		{"2018-10-05 11:00", "", "depart_time=<nil>, arrive_time=timestamp: .*"},
		{"", "", "depart_time=timestamp: .*, arrive_time=timestamp: .*"},
	}

	for i, test := range tests {
		route := makeRouteTm(test.deptTime, test.arrvTime, false)

		r := regexp.MustCompile(test.regexp)
		_, got := route.CalcAmountTime()

		if got == nil {
			t.Errorf("[%d] got is nil", i)
		}

		if !r.MatchString(got.Error()) {
			t.Errorf("[%d] want=%s, got=%s\n", i, test.regexp, got.Error())
		}
	}
}

func TestRoute_CalcAmountDistance_Success(t *testing.T) {
	tests := []struct {
		deptLat float64
		deptLon float64
		arrvLat float64
		arrvLon float64
		want    float64
	}{
		{35.689166, 139.704444, 35.654444, 139.706666, 3857.347000}, // 新宿-渋谷
		{35.654444, 139.706666, 35.632904, 139.715935, 2532.790000}, // 渋谷-目黒
		{35.647078, 139.710099, 35.632904, 139.715935, 1658.913000}, // 恵比寿-目黒
		{35.654444, 139.706666, 35.647078, 139.710099, 874.316000},  // 渋谷-恵比寿
	}

	for i, test := range tests {
		route := makeRoutePt(test.deptLat, test.deptLon, test.arrvLat, test.arrvLon)
		//fmt.Printf("[%d] route=%#v\n", i, route)

		want := test.want
		got, _ := route.CalcAmountDistance()

		// TODO: 誤差をどの程度許容するか?
		if math.Abs(want-got) > 100 {
			t.Errorf("[%d] want=%f, got=%f", i, want, got)
		}
	}
}

func TestRoute_CalcAmountDistance_Error(t *testing.T) {
	tests := []struct {
		deptLat float64
		deptLon float64
		arrvLat float64
		arrvLon float64
		regexp  string
	}{
		{-90.001000, +139.704444, +35.654444, +139.706666, "DepartPoint=Latitude .*, ArrivePoint=<nil>"},
		{+35.689166, +139.704444, -90.001000, +139.706666, "DepartPoint=<nil>, ArrivePoint=Latitude .*"},
	}

	for i, test := range tests {
		route := makeRoutePt(test.deptLat, test.deptLon, test.arrvLat, test.arrvLon)

		r := regexp.MustCompile(test.regexp)
		_, got := route.CalcAmountDistance()

		if got == nil {
			t.Errorf("[%d] got is nil", i)
		}

		//fmt.Println(got.Error())
		if !r.MatchString(got.Error()) {
			t.Errorf("[%d] want=%s, got=%s\n", i, test.regexp, got.Error())
		}
	}
}

func TestRideSahre_CalcAmountTime_Success(t *testing.T) {
	tests := []struct {
		routes []*Route
		want   string
	}{
		{
			routes: []*Route{
				makeRouteTm("2018-10-11 13:20", "2018-10-11 13:40", true),
				makeRouteTm("2018-10-11 13:45", "2018-10-11 14:15", true),
				makeRouteTm("2018-10-11 14:25", "2018-10-11 15:10", true),
			},
			want: "1h50m",
		},
	}

	for i, test := range tests {
		rideShare := makeRideShare(test.routes)

		want, _ := time.ParseDuration(test.want)
		got, err := rideShare.CalcAmountTime()

		if err != nil {
			t.Errorf("[%d] want=<nil>, got=%s", i, err.Error())
		}

		if d, _ := ptypes.Duration(got); want != d {
			t.Errorf("[%d] want=%s, got=%s", i, want, d)
		}
	}
}

func TestRideShare_CalcAmountTime_Error(t *testing.T) {
	tests := []struct {
		routes []*Route
		regexp string
	}{
		{
			routes: nil,
			regexp: "Routes is nil",
		},
		{
			routes: []*Route{
				makeRouteTm("2018-10-11 13:20", "2018-10-11 13:40", true),
				makeRouteTm("", "", false),
				makeRouteTm("2018-10-11 14:25", "2018-10-11 15:10", true),
			},
			regexp: `Routes\[1\]: .*`,
		},
	}

	for i, test := range tests {
		rideShare := makeRideShare(test.routes)

		r := regexp.MustCompile(test.regexp)
		_, got := rideShare.CalcAmountTime()

		if got == nil {
			t.Errorf("[%d] want=%s, got=<nil>", i, test.regexp)
		}

		if !r.MatchString(got.Error()) {
			t.Errorf("[%d] want=%s, got=%s", i, test.regexp, got.Error())
		}
	}
}

func TestRideShare_CalcAmountDistance_Success(t *testing.T) {
	tests := []struct {
		routes []*Route
		want   float64
	}{
		{
			// 乗り換え距離＝0
			routes: []*Route{
				makeRoutePt(35.689166, 139.704444, 35.654444, 139.706666), // 新宿-渋谷
				makeRoutePt(35.654444, 139.706666, 35.632904, 139.715935), // 渋谷-目黒
			},
			want: 6390.137000,
		},
		{
			// 乗り換え距離≠0
			routes: []*Route{
				makeRoutePt(35.689166, 139.704444, 35.654444, 139.706666), // 新宿-渋谷
				makeRoutePt(35.647078, 139.710099, 35.632904, 139.715935), // 恵比寿-目黒
			},
			want: 6390.57600,
		},
	}

	for i, test := range tests {
		rideShare := makeRideShare(test.routes)

		want := test.want
		got, err := rideShare.CalcAmountDistance()

		if err != nil {
			t.Errorf("[%d] err: %s", i, err.Error())
		}

		if math.Abs(want-got) > 100 {
			t.Errorf("[%d] want=%f, got=%f", i, want, got)
		}
	}
}

func TestRideShare_CalcAmountPrice_Success(t *testing.T) {
	tests := []struct {
		routes []*Route
		want   uint32
	}{
		{
			routes: []*Route{
				makeRoutePr(100),
				makeRoutePr(250),
				makeRoutePr(180),
			},
			want: 530,
		},
	}

	for i, test := range tests {
		rideShare := makeRideShare(test.routes)

		want := test.want
		got, err := rideShare.CalcAmountPrice()

		if err != nil {
			t.Errorf("[%d] want=<nil>, got=%s", i, err.Error())
		}

		if want != got {
			t.Errorf("[%d] want=%d, got=%d", i, want, got)
		}
	}
}

func TestRideShare_CalcAmountPrice_Error(t *testing.T) {
	tests := []struct {
		routes []*Route
		regexp string
	}{
		{
			routes: nil,
			regexp: "Routes is nil",
		},
	}

	for i, test := range tests {
		rideShare := makeRideShare(test.routes)

		r := regexp.MustCompile(test.regexp)
		_, got := rideShare.CalcAmountPrice()

		if got == nil {
			t.Errorf("[%d] want=%s, got=<nil>", i, test.regexp)
		}

		if !r.MatchString(got.Error()) {
			t.Errorf("[%d] want=%s, got=%s", i, test.regexp, got.Error())
		}
	}
}
