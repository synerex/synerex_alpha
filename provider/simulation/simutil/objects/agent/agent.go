package agent

type Coord struct {
	Latitude  float64
	Longitude float64
}

type Route struct {
	Position      *Coord
	Direction     float64
	Speed         float64
	Departure     *Coord
	Destination   *Coord
	TransitPoints []*Coord
	NextTransit   *Coord
	TotalDistance float64
	RequiredTime  float64
}

type Type int64

const (
	PEDESTRIAN Type = iota
	CAR
)

type Agent struct {
	ID   uint64
	Type Type
}
