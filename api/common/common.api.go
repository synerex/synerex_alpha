package common

import (
	"errors"
	"fmt"
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
