package api

import (
	"testing"

	"github.com/synerex/synerex_alpha/api/adservice"
	"github.com/synerex/synerex_alpha/api/fleet"
	"github.com/synerex/synerex_alpha/api/library"
	"github.com/synerex/synerex_alpha/api/marketing"
	"github.com/synerex/synerex_alpha/api/ptransit"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"github.com/synerex/synerex_alpha/api/routing"

	"github.com/stretchr/testify/assert"
)

func TestNewDemand(t *testing.T) {
	dm := NewDemand()

	assert.IsType(t, &Demand{}, dm)
	assert.Nil(t, dm.ArgOneof)
}

func TestNewSupply(t *testing.T) {
	sp := NewSupply()

	assert.IsType(t, &Supply{}, sp)
	assert.Nil(t, sp.ArgOneof)
}

//
// Demand
//
func TestDemandWithFleet(t *testing.T) {
	dm := NewDemand().WithFleet(&fleet.Fleet{})

	assert.IsType(t, &fleet.Fleet{}, dm.GetArg_Fleet())
}

func TestDemandWithRideShare(t *testing.T) {
	dm := NewDemand().WithRideShare(&rideshare.RideShare{})

	assert.IsType(t, &rideshare.RideShare{}, dm.GetArg_RideShare())
}

func TestDemandWithAdService(t *testing.T) {
	dm := NewDemand().WithAdService(&adservice.AdService{})

	assert.IsType(t, &adservice.AdService{}, dm.GetArg_AdService())
}

func TestDemandWithLibService(t *testing.T) {
	dm := NewDemand().WithLibService(&library.LibService{})

	assert.IsType(t, &library.LibService{}, dm.GetArg_LibService())
}

func TestDemandWithPTService(t *testing.T) {
	dm := NewDemand().WithPTService(&ptransit.PTService{})

	assert.IsType(t, &ptransit.PTService{}, dm.GetArg_PTService())
}

func TestDemandWithRoutingService(t *testing.T) {
	dm := NewDemand().WithRoutingService(&routing.RoutingService{})

	assert.IsType(t, &routing.RoutingService{}, dm.GetArg_RoutingService())
}

func TestDemandWithMarketingService(t *testing.T) {
	dm := NewDemand().WithMarketingService(&marketing.MarketingService{})

	assert.IsType(t, &marketing.MarketingService{}, dm.GetArg_MarketingService())
}

//
// Supply
//
func TestSupplyWithFleet(t *testing.T) {
	sp := NewSupply().WithFleet(&fleet.Fleet{})

	assert.IsType(t, &fleet.Fleet{}, sp.GetArg_Fleet())
}

func TestSupplyWithRideShare(t *testing.T) {
	sp := NewSupply().WithRideShare(&rideshare.RideShare{})

	assert.IsType(t, &rideshare.RideShare{}, sp.GetArg_RideShare())
}

func TestSupplyWithAdService(t *testing.T) {
	sp := NewSupply().WithAdService(&adservice.AdService{})

	assert.IsType(t, &adservice.AdService{}, sp.GetArg_AdService())
}

func TestSupplyWithLibService(t *testing.T) {
	sp := NewSupply().WithLibService(&library.LibService{})

	assert.IsType(t, &library.LibService{}, sp.GetArg_LibService())
}

func TestSupplyWithPTService(t *testing.T) {
	sp := NewSupply().WithPTService(&ptransit.PTService{})

	assert.IsType(t, &ptransit.PTService{}, sp.GetArg_PTService())
}

func TestSupplyWithRoutingService(t *testing.T) {
	sp := NewSupply().WithRoutingService(&routing.RoutingService{})

	assert.IsType(t, &routing.RoutingService{}, sp.GetArg_RoutingService())
}

func TestSupplyWithMarketingService(t *testing.T) {
	sp := NewSupply().WithMarketingService(&marketing.MarketingService{})

	assert.IsType(t, &marketing.MarketingService{}, sp.GetArg_MarketingService())
}
