package agent

type Map struct {
	SLatitude  float64
	SLongitude float64
	ELatitude  float64
	ELongitude float64
}

type Area struct {
	ID  uint64
	Map *Map
}

func NewArea() *Area {
	a := &Area{}
	return a
}
