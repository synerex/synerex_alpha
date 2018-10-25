package common

import (
	"errors"
	"fmt"
	"math"
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
func IsSamePoint(p1, p2 *Point) bool {
	return p1.GetLatitude() == p2.GetLatitude() && p1.GetLatitude() == p2.GetLatitude()
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
