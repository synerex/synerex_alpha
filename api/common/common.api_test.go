package common

import (
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

const (
	ON  = true
	OFF = false
)

func makePoint(lat, lon float64) *Point {
	if math.IsNaN(lat) && math.IsNaN(lon) {
		return nil
	} else {
		return &Point{Latitude: lat, Longitude: lon}
	}
}

func makeTimestampPb(s string) *timestamp.Timestamp {
	tm, _ := time.Parse("2006-01-02", s)
	ts, _ := ptypes.TimestampProto(tm)
	return ts
}

func makeTimestampTm(s string) time.Time {
	tm, _ := time.Parse("2006-01-02", s)
	return tm
}

func TestPoint_ValidatePoint_Success(t *testing.T) {
	tests := []struct {
		lat float64
		lon float64
	}{
		{-90.000, 139.704},
		{90.000, 139.704},
		{35.689, -180.000},
		{35.689, 180.000},
		{35.689, 139.704},
	}

	for i, test := range tests {
		p := makePoint(test.lat, test.lon)

		got := ValidatePoint(p)

		if got != nil {
			t.Errorf("[%d] want=nil, got=%s", i, got.Error())
		}
	}
}

func TestPoint_ValidatePoint_Error(t *testing.T) {
	tests := []struct {
		lat    float64
		lon    float64
		regexp string
	}{
		{-90.001, 139.704, "Latitude .*"},
		{90.001, 139.704, "Latitude .*"},
		{35.689, -180.001, "Longitude .*"},
		{35.689, 180.001, "Longitude .*"},
		{math.NaN(), math.NaN(), "Point .*"},
	}

	for i, test := range tests {
		p := makePoint(test.lat, test.lon)
		r := regexp.MustCompile(test.regexp)

		got := ValidatePoint(p)

		if got == nil {
			t.Errorf("[%d] want=%s, got=nil\n", i, test.regexp)
		}

		if !r.MatchString(got.Error()) {
			t.Errorf("[%d] want=%s, got=%s\n", i, test.regexp, got.Error())
		}
	}
}

func TestPoint_IsSamePoint(t *testing.T) {
	tests := []struct {
		p1   *Point
		p2   *Point
		r    float64
		want bool
	}{
		{makePoint(35.689166, 139.704444), makePoint(35.689166, 139.704444), 0, true},    // 新宿-新宿
		{makePoint(35.689166, 139.704444), makePoint(35.654444, 139.706666), 0, false},   // 新宿-渋谷
		{makePoint(35.689166, 139.704444), makePoint(35.689904, 139.704163), 100, true},  // 新宿-新宿 (100m以内)
		{makePoint(35.689166, 139.704444), makePoint(35.654444, 139.706666), 100, false}, // 新宿-渋谷 (100m以内)
	}

	for i, test := range tests {
		want := test.want
		got := test.p1.IsSamePoint(test.p2, test.r)

		if want != got {
			t.Errorf("[%d], want=%t, got=%t", i, want, got)
		}
	}
}

func TestPoint_Distance(t *testing.T) {
	tests := []struct {
		Lat1 float64
		Lon1 float64
		Lat2 float64
		Lon2 float64
		want float64
	}{
		{35.689166, 139.704444, 35.654444, 139.706666, 3857.347000}, // 新宿-渋谷
		{35.654444, 139.706666, 35.632904, 139.715935, 2532.790000}, // 渋谷-目黒
		{35.647078, 139.710099, 35.632904, 139.715935, 1658.913000}, // 恵比寿-目黒
		{35.654444, 139.706666, 35.647078, 139.710099, 874.316000},  // 渋谷-恵比寿
	}

	for i, test := range tests {
		p1 := makePoint(test.Lat1, test.Lon1)
		p2 := makePoint(test.Lat2, test.Lon2)

		want := test.want
		got, _ := p1.Distance(p2)

		if math.Abs(want-got) > 100 {
			t.Errorf("[%d] want=%f, got=%f", i, want, got)
		}
	}
}

func TestNewPlace(t *testing.T) {
	p := NewPlace()

	assert.Nil(t, p.Value)
}

func TestWithPoint(t *testing.T) {
	p := NewPlace().WithPoint(&Point{Latitude: 35.689, Longitude: 139.704})

	assert.IsType(t, &Place_Point{}, p.Value, "p.Value shuld be *Place_Point")
	assert.Equal(t, 35.689, p.GetPoint().Latitude)
	assert.Equal(t, 139.704, p.GetPoint().Longitude)
}

func TestWithAreas(t *testing.T) {
	p := NewPlace().WithAreas([][]*Point{
		[]*Point{
			&Point{Latitude: 35.689, Longitude: 139.704},
			&Point{Latitude: 35.789, Longitude: 139.804},
			&Point{Latitude: 35.889, Longitude: 139.904},
		},
		[]*Point{
			&Point{Latitude: 36.589, Longitude: 139.604},
			&Point{Latitude: 36.489, Longitude: 139.504},
			&Point{Latitude: 36.389, Longitude: 139.404},
			&Point{Latitude: 36.289, Longitude: 139.304},
		},
	})

	assert.IsType(t, &Place_Areas{}, p.Value)
	assert.Len(t, p.GetAreas().Values, 2)
	assert.Len(t, p.GetAreas().Values[0].Points, 3)
	assert.Len(t, p.GetAreas().Values[1].Points, 4)
}

func TestNewTime(t *testing.T) {
	tm := NewTime()

	assert.Nil(t, tm.Value)
}

func TestWithTimestamp(t *testing.T) {
	tm := NewTime().WithTimestamp(&timestamp.Timestamp{Seconds: 1540542272, Nanos: 550878777})

	assert.IsType(t, &Time_Timestamp{}, tm.Value)
	assert.Equal(t, int64(1540542272), tm.GetTimestamp().Seconds)
	assert.Equal(t, int32(550878777), tm.GetTimestamp().Nanos)
}

func TestWithPeriods(t *testing.T) {
	tp := NewTime().WithPeriods([]*Period{
		// 10月 (第1-2週)
		&Period{
			From: makeTimestampPb("2018-10-01"),
			To:   makeTimestampPb("2018-10-31"),
			Options: []*RepeatOption{
				&RepeatOption{
					Weeks: []bool{ON, ON, OFF, OFF, OFF},
				},
			},
		},
		// 11月 (毎土・日曜日)
		&Period{
			From: makeTimestampPb("2018-11-01"),
			To:   makeTimestampPb("2018-11-30"),
			Options: []*RepeatOption{
				&RepeatOption{
					Weekdays: []bool{ON, OFF, OFF, OFF, OFF, OFF, ON},
				},
			},
		},
	})

	assert.IsType(t, &Time_Periods{}, tp.Value)
	assert.Len(t, tp.GetPeriods().Values, 2)
	assert.Equal(t, makeTimestampPb("2018-10-01"), tp.GetPeriods().Values[0].From)
}

func TestWithOtherTime(t *testing.T) {
	to := NewTime().WithOtherTime(OtherTime_AS_SOON_AS)

	assert.IsType(t, &Time_Other{}, to.Value)
	assert.Equal(t, OtherTime_AS_SOON_AS, to.GetOther())
}
