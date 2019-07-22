package api

import (
	"github.com/synerex/synerex_alpha/api/adservice"
	"github.com/synerex/synerex_alpha/api/fleet"
	"github.com/synerex/synerex_alpha/api/library"
	"github.com/synerex/synerex_alpha/api/marketing"
	"github.com/synerex/synerex_alpha/api/ptransit"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"github.com/synerex/synerex_alpha/api/routing"
	"github.com/synerex/synerex_alpha/api/clock"
	"github.com/synerex/synerex_alpha/api/area"
	"github.com/synerex/synerex_alpha/api/agent"
)

// NewDemand returns empty Demand.
func NewDemand() *Demand {
	return &Demand{}
}

// NewSupply returns empty Supply.
func NewSupply() *Supply {
	return &Supply{}
}

// WithFleet set a given Fleet to Demand.Demand_Arg_Fleet.Arg_Fleet.
func (dm *Demand) WithFleet(f *fleet.Fleet) *Demand {
	dm.ArgOneof = &Demand_Arg_Fleet{f}
	return dm
}

// WithRideShare set a given RideShare to Demand.Demand_Arg_RideShare.Arg_RideShare.
func (dm *Demand) WithRideShare(r *rideshare.RideShare) *Demand {
	dm.ArgOneof = &Demand_Arg_RideShare{r}
	return dm
}

// WithAdShare set a given AdShare to Demand.Demand_Arg_AdShare.Arg_AdShare.
func (dm *Demand) WithAdService(a *adservice.AdService) *Demand {
	dm.ArgOneof = &Demand_Arg_AdService{a}
	return dm
}

// WithLibService set a given LibService to Demand.Demand_Arg_LibService.Arg_LibService.
func (dm *Demand) WithLibService(l *library.LibService) *Demand {
	dm.ArgOneof = &Demand_Arg_LibService{l}
	return dm
}

// WithPTService set a given PTService to Demand.Demand_Arg_PTService.Arg_PTService.
func (dm *Demand) WithPTService(p *ptransit.PTService) *Demand {
	dm.ArgOneof = &Demand_Arg_PTService{p}
	return dm
}

// WithRoutingService set a given RoutingService to Demand.Demand_Arg_RoutingService.Arg_RoutingService.
func (dm *Demand) WithRoutingService(r *routing.RoutingService) *Demand {
	dm.ArgOneof = &Demand_Arg_RoutingService{r}
	return dm
}

// WithMarketingService set a given MarketingService to Demand.Demand_Arg_MarketingService.Arg_MarketingService.
func (dm *Demand) WithMarketingService(m *marketing.MarketingService) *Demand {
	dm.ArgOneof = &Demand_Arg_MarketingService{m}
	return dm
}

// WithFleet set a given Fleet to Supply.Supply_Arg_Fleet.Arg_Fleet.
func (sp *Supply) WithFleet(f *fleet.Fleet) *Supply {
	sp.ArgOneof = &Supply_Arg_Fleet{f}
	return sp
}

// WithRideShare set a given RideShare to Supply.Supply_Arg_RideShare.Arg_RideShare.
func (sp *Supply) WithRideShare(r *rideshare.RideShare) *Supply {
	sp.ArgOneof = &Supply_Arg_RideShare{r}
	return sp
}

// WithAdService set a given AdService to Supply.Supply_Arg_AdService.Arg_AdService.
func (sp *Supply) WithAdService(a *adservice.AdService) *Supply {
	sp.ArgOneof = &Supply_Arg_AdService{a}
	return sp
}

// WithLibService set a given LibService to Supply.Supply_Arg_LibService.Arg_LibService.
func (sp *Supply) WithLibService(l *library.LibService) *Supply {
	sp.ArgOneof = &Supply_Arg_LibService{l}
	return sp
}

// WithPTService set a given PTService to Supply.Supply_Arg_PTService.Arg_PTService.
func (sp *Supply) WithPTService(p *ptransit.PTService) *Supply {
	sp.ArgOneof = &Supply_Arg_PTService{p}
	return sp
}

// WithRoutingService set a given RoutingService to Supply.Supply_Arg_RoutingService.Arg_RoutingService.
func (sp *Supply) WithRoutingService(r *routing.RoutingService) *Supply {
	sp.ArgOneof = &Supply_Arg_RoutingService{r}
	return sp
}

// WithMarketingService set a given MarketingService to Supply.Supply_Arg_MarketingService.Arg_MarketingService.
func (sp *Supply) WithMarketingService(m *marketing.MarketingService) *Supply {
	sp.ArgOneof = &Supply_Arg_MarketingService{m}
	return sp
}

// WithClockService set a given ClockService to Supply.Supply_Arg_ClockService.Arg_ClockService .
func (sp *Supply) WithClockService(c *clock.ClockService) *Supply {
	sp.ArgOneof = &Supply_Arg_ClockService{c}
	return sp
}

// WithAreaService set a given AreaService to Supply.Supply_Arg_AreaService.Arg_AreaService.
func (sp *Supply) WithAreaService(a *area.AreaService) *Supply {
	sp.ArgOneof = &Supply_Arg_AreaService{a}
	return sp
}

// WithAgentService set a given AgentService to Supply.Supply_Arg_AgentService.Arg_AgentService.
func (sp *Supply) WithAgentService(a *agent.AgentService) *Supply {
	sp.ArgOneof = &Supply_Arg_AgentService{a}
	return sp
}
