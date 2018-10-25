package common

import (
	"math"
	"regexp"
	"testing"
)

func makePoint(lat, lon float64) *Point {
	if math.IsNaN(lat) && math.IsNaN(lon) {
		return nil
	} else {
		return &Point{Latitude: lat, Longitude: lon}
	}
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
		want bool
	}{
		{makePoint(35.689166, 139.704444), makePoint(35.689166, 139.704444), true},  // 新宿-新宿
		{makePoint(35.689166, 139.704444), makePoint(35.654444, 139.706666), false}, // 新宿-渋谷
	}

	for i, test := range tests {
		want := test.want
		got := IsSamePoint(test.p1, test.p2)

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
