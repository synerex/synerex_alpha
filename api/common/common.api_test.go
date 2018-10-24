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
		{makePoint(35.689, 139.704), makePoint(35.689, 139.704), true},
		{makePoint(35.689, 139.704), makePoint(35.654, 139.706), false},
	}

	for i, test := range tests {
		want := test.want
		got := IsSamePoint(test.p1, test.p2)

		if want != got {
			t.Errorf("[%d], want=%t, got=%t", i, want, got)
		}
	}
}
