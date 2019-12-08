package agent

import (
	"github.com/synerex/synerex_alpha/api/simulation/area"
)

type MapCoord struct {
	SLatitude  float64
	SLongitude float64
	ELatitude  float64
	ELongitude float64
}

type Map struct {
	ID        uint64
	Control   *MapCoord
	Duplicate *MapCoord
	Neighbors []*MapCoord
}

func NewMap() *Map {
	a := &Map{}
	return a
}

func (p *Map) SetControlMap(mapInfo *area.AreaCoord) {
	p.Control = &MapCoord{
		SLatitude:  float64(mapInfo.StartLat),
		SLongitude: float64(mapInfo.StartLon),
		ELatitude:  float64(mapInfo.EndLat),
		ELongitude: float64(mapInfo.EndLon),
	}
}

func (p *Map) SetDuplicateMap(mapInfo *area.AreaCoord) {
	p.Duplicate = &MapCoord{
		SLatitude:  float64(mapInfo.StartLat),
		SLongitude: float64(mapInfo.StartLon),
		ELatitude:  float64(mapInfo.EndLat),
		ELongitude: float64(mapInfo.EndLon),
	}
}

func (p *Map) SetNeighborMaps(mapsInfo []*area.AreaCoord) {
	neighbors := make([]*MapCoord, 0)
	for _, mapInfo := range mapsInfo {
		neighbors = append(neighbors, &MapCoord{
			SLatitude:  float64(mapInfo.StartLat),
			SLongitude: float64(mapInfo.StartLon),
			ELatitude:  float64(mapInfo.EndLat),
			ELongitude: float64(mapInfo.EndLon),
		})
	}
	p.Neighbors = neighbors
}
