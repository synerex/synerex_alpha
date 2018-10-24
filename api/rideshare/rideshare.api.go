package rideshare

import (
	"github.com/synerex/synerex_alpha/api/common"
	"errors"
	"fmt"
	"math"

	"github.com/golang/protobuf/ptypes"
	durpb "github.com/golang/protobuf/ptypes/duration"
)

// CalcAmountDistance returns distance between DepartPoint and ArrivalPoint.
// If both point has valid Point, returns the distance(m) between both point.
// Otherwise, returns -1.
func (m *Route) CalcAmountDistance() (float64, error) {
	// get protobuf Point
	deptPoint := m.GetDepartPoint().GetPoint()
	arrvPoint := m.GetArrivePoint().GetPoint()

	if deptPoint != nil && arrvPoint != nil {
		// check whether a Point is valid
		err1 := common.ValidatePoint(deptPoint)
		err2 := common.ValidatePoint(arrvPoint)

		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("DepartPoint=%v, ArrivePoint=%v", err1, err2)
		}

		// get distance
		return distance(deptPoint, arrvPoint)
	} else {
		return -1, nil
	}
}

// convert degree to radian
func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// calculate distance using Hubeny formula.
func distance(p1, p2 *common.Point) (float64, error) {
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

// CalcAmountTime returns difference between ArriveTime and DepartTime.
// If both time has timestamp, CalcAmountTime returns protobuf Duration.
// Otherwise, CalcAmountTime returns nil.
func (m *Route) CalcAmountTime() (*durpb.Duration, error) {
	// get protobuf Timestamp
	deptTimePb := m.GetDepartTime().GetTimestamp()
	arrvTimePb := m.GetArriveTime().GetTimestamp()

	if deptTimePb != nil && arrvTimePb != nil {
		// convert to time.Time
		deptTimeTm, err1 := ptypes.Timestamp(deptTimePb)
		arrvTimeTm, err2 := ptypes.Timestamp(arrvTimePb)

		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("depart_time=%v, arrive_time=%v", err1, err2)
		}

		// get time.Duration
		durTm := arrvTimeTm.Sub(deptTimeTm)

		// convert to protobuf Duration
		return ptypes.DurationProto(durTm), nil
	} else {
		return nil, nil
	}
}

// CalcAmountTime returns sum of each Route's amount time.
// If each Route has valid timestamp, returns differens between first Route.DepartTime and last Route.ArriveTime.
// Otherwise, returns nil.
func (r *RideShare) CalcAmountTime() (*durpb.Duration, error) {
	if r.GetRoutes() == nil {
		return nil, errors.New("Routes is nil")
	}

	// validate each route has valid timestamp
	for i, route := range r.GetRoutes() {
		_, err := route.CalcAmountTime()

		if err != nil {
			return nil, fmt.Errorf("Routes[%d]: %s", i, err.Error())
		}
	}

	// get pb.Timestamp
	deptTimePb := r.GetRoutes()[0].GetDepartTime().GetTimestamp()
	arrvTimePb := r.GetRoutes()[len(r.GetRoutes())-1].GetArriveTime().GetTimestamp()

	// convert to time.Time
	deptTimeTm, _ := ptypes.Timestamp(deptTimePb)
	arrvTimeTm, _ := ptypes.Timestamp(arrvTimePb)

	// get duration
	dur := arrvTimeTm.Sub(deptTimeTm)

	// convert to pb.Duration
	return ptypes.DurationProto(dur), nil
}

// CalcAmountDistance returns sum of each Route's amount distance.
// If each Route has valid points, returns distance.
// Otherwise, returns -1.
func (r *RideShare) CalcAmountDistance() (float64, error) {
	if r.GetRoutes() == nil {
		return -1, errors.New("Routes is nil")
	}

	amntDist := float64(0.0)

	for i, route := range r.GetRoutes() {
		dist, err := route.CalcAmountDistance()

		if err != nil {
			return -1, fmt.Errorf("Routes[%d]: %s", i, err.Error())
		}

		amntDist += dist

		// add exchange distance
		if 0 < i {
			prevArrvPoint := r.GetRoutes()[i-1].GetArrivePoint().GetPoint()
			currDeptPoint := route.GetDepartPoint().GetPoint()

			if !common.IsSamePoint(prevArrvPoint, currDeptPoint) {
				dist, _ := distance(prevArrvPoint, currDeptPoint)
				amntDist += dist
			}
		}
	}

	return amntDist, nil
}

// CalcAmountPrice returns sum of each Route's amount price.
// If each Route has valid price, returns amount of price.
// Otherwise, returns 0.
func (r *RideShare) CalcAmountPrice() (uint32, error) {
	if r.GetRoutes() == nil {
		return 0, errors.New("Routes is nil")
	}

	amntPrice := uint32(0)

	for _, route := range r.GetRoutes() {
		amntPrice += route.GetAmountPrice()
	}

	return amntPrice, nil
}
