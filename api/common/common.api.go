package common

import (
	"errors"
	"fmt"
	"math"

	"github.com/golang/protobuf/ptypes/timestamp"
)

const (
	MIN_LATITUDE  = float64(-90)
	MAX_LATITUDE  = float64(90)
	MIN_LONGITUDE = float64(-180)
	MAX_LONGITUDE = float64(180)
)

// ValidatePoint determins whether a Point is valid.
// Latitude is in range [-90, 90] and Longitude is in rage [-180, 180].
func ValidatePoint(p *Point) error {
	if p == nil {
		return errors.New("Point is nil")
	} else if p.GetLatitude() < MIN_LATITUDE || MAX_LATITUDE < p.GetLatitude() {
		return fmt.Errorf("Latitude is out of range. (%f)", p.GetLatitude())
	} else if p.GetLongitude() < MIN_LONGITUDE || MAX_LONGITUDE < p.GetLongitude() {
		return fmt.Errorf("Longitude is out of range. (%f)", p.GetLongitude())
	}
	return nil
}

// IsSamePoint determins whether two points are same.
// If r is 0, matching is strict.
// Otherwise, matching is lenient. (consider as same if they are within r meters)
func (p1 *Point) IsSamePoint(p2 *Point, r float64) bool {
	if r <= 0.0 {
		return p1.GetLatitude() == p2.GetLatitude() && p1.GetLatitude() == p2.GetLatitude()
	} else {
		d, _ := p1.Distance(p2)
		return d <= r
	}
}

// convert degree to radian
func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// calculate distance using Hubeny formula.
func (p1 *Point) Distance(p2 *Point) (float64, error) {
	a := 6378137.000
	b := 6356752.314
	e := math.Sqrt((math.Pow(a, 2) - math.Pow(b, 2)) / math.Pow(a, 2))

	x1 := deg2rad(p1.GetLongitude())
	y1 := deg2rad(p1.GetLatitude())
	x2 := deg2rad(p2.GetLongitude())
	y2 := deg2rad(p2.GetLatitude())

	dy := y1 - y2
	dx := x1 - x2
	uy := (y1 + y2) / 2.0

	W := math.Sqrt(1 - math.Pow(e, 2)*math.Pow(math.Sin(uy), 2))
	M := a * (1 - math.Pow(e, 2)) / math.Pow(W, 3)
	N := a / W

	d := math.Sqrt(math.Pow(dy*M, 2) + math.Pow(dx*N*math.Cos(uy), 2))

	return d, nil
}

// NewPlace returns empty Place.
func NewPlace() *Place {
	return &Place{}
}

// WithPoint set a given Point to Place.Place_Point.Point.
func (p *Place) WithPoint(point *Point) *Place {
	p.Value = &Place_Point{point}
	return p
}

// WithAreas makes []*Area from a given points and set to Place.Place_Areas.Areas.
func (p *Place) WithAreas(areas [][]*Point) *Place {
	values := make([]*Area, len(areas))
	for i, points := range areas {
		values[i] = &Area{Points: points}
	}
	p.Value = &Place_Areas{&Areas{Values: values}}
	return p
}

// NewTime returns empty Time.
func NewTime() *Time {
	return &Time{}
}

// WithTimestamp set a given timestamp to Time.Time_Timestamp.Timestamp.
func (t *Time) WithTimestamp(ts *timestamp.Timestamp) *Time {
	t.Value = &Time_Timestamp{ts}
	return t
}

// WithPeriods set given periods to Time.Time_Periods.Periods.
func (t *Time) WithPeriods(periods []*Period) *Time {
	t.Value = &Time_Periods{&Periods{Values: periods}}
	return t
}

// WithOtherTime set a given time to Time.Time_Other.Other.
func (t *Time) WithOtherTime(other OtherTime) *Time {
	t.Value = &Time_Other{other}
	return t
}
