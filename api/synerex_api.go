package api

import (
	"github.com/synerex/synerex_alpha/api/adservice"
	"github.com/synerex/synerex_alpha/api/fleet"
	"github.com/synerex/synerex_alpha/api/library"
	"github.com/synerex/synerex_alpha/api/marketing"
	"github.com/synerex/synerex_alpha/api/ptransit"
	"github.com/synerex/synerex_alpha/api/rideshare"
	"github.com/synerex/synerex_alpha/api/routing"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
)

// Demand
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

// WithClockInfo set a given ClockInfo to Demand.Demand_Arg_ClockInfo.Arg_ClockInfo .
//func (dm *Demand) WithClockInfo(c *clock.ClockInfo) *Demand {
//	dm.ArgOneof = &Demand_Arg_ClockInfo{c}
//	return dm
//}

// WithClockDemand set a given ClockDemand to Demand.Demand_Arg_ClockDemand.Arg_ClockDemand .
func (dm *Demand) WithClockDemand(c *clock.ClockDemand) *Demand {
	dm.ArgOneof = &Demand_Arg_ClockDemand{c}
	return dm
}

// WithAreaInfo set a given AreaInfo to Demand.Demand_Arg_AreaInfo.Arg_AreaInfo.
//func (dm *Demand) WithAreaInfo(a *area.AreaInfo) *Demand {
//	dm.ArgOneof = &Demand_Arg_AreaInfo{a}
//	return dm
//}

// WithAreaDemand set a given AreaDemand to Demand.Demand_Arg_AreaDemand.Arg_AreaDemand.
func (dm *Demand) WithAreaDemand(a *area.AreaDemand) *Demand {
	dm.ArgOneof = &Demand_Arg_AreaDemand{a}
	return dm
}

// WithAgentInfo set a given AgentInfo to Demand.Demand_Arg_AgentInfo.Arg_AgentInfo.
//func (dm *Demand) WithAgentInfo(a *agent.AgentInfo) *Demand {
//	dm.ArgOneof = &Demand_Arg_AgentInfo{a}
//	return dm
//}

// WithAgentsInfo set a given AgentsInfo to Demand.Demand_Arg_AgentsInfo.Arg_AgentsInfo.
//func (dm *Demand) WithAgentsInfo(a *agent.AgentsInfo) *Demand {
//	dm.ArgOneof = &Demand_Arg_AgentsInfo{a}
//	return dm
//}

// WithAgentDemand set a given AgentDemand to Demand.Demand_Arg_AgentDemand.Arg_AgentDemand.
func (dm *Demand) WithAgentDemand(a *agent.AgentDemand) *Demand {
	dm.ArgOneof = &Demand_Arg_AgentDemand{a}
	return dm
}

// Supply
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

// WithClockInfo set a given ClockInfo to Supply.Supply_Arg_ClockInfo.Arg_ClockInfo .
func (sp *Supply) WithClockInfo(c *clock.ClockInfo) *Supply {
	sp.ArgOneof = &Supply_Arg_ClockInfo{c}
	return sp
}

// WithClockDemand set a given ClockDemand to Supply.Supply_Arg_ClockDemand.Arg_ClockDemand .
//func (sp *Supply) WithClockDemand(c *clock.ClockDemand) *Supply {
//	sp.ArgOneof = &Supply_Arg_ClockDemand{c}
//	return sp
//}

// WithAreaInfo set a given AreaInfo to Supply.Supply_Arg_AreaInfo.Arg_AreaInfo.
func (sp *Supply) WithAreaInfo(a *area.AreaInfo) *Supply {
	sp.ArgOneof = &Supply_Arg_AreaInfo{a}
	return sp
}

// WithAreaDemand set a given AreaDemand to Supply.Supply_Arg_AreaDemand.Arg_AreaDemand.
//func (sp *Supply) WithAreaDemand(a *area.AreaDemand) *Supply {
//	sp.ArgOneof = &Supply_Arg_AreaDemand{a}
//	return sp
//}

// WithAgentInfo set a given AgentInfo to Supply.Supply_Arg_AgentInfo.Arg_AgentInfo.
func (sp *Supply) WithAgentInfo(a *agent.AgentInfo) *Supply {
	sp.ArgOneof = &Supply_Arg_AgentInfo{a}
	return sp
}

// WithAgentsInfo set a given AgentsInfo to Supply.Supply_Arg_AgentsInfo.Arg_AgentsInfo.
func (sp *Supply) WithAgentsInfo(a *agent.AgentsInfo) *Supply {
	sp.ArgOneof = &Supply_Arg_AgentsInfo{a}
	return sp
}

// WithAgentDemand set a given AgentDemand to Supply.Supply_Arg_AgentDemand.Arg_AgentDemand.
//func (sp *Supply) WithAgentDemand(a *agent.AgentDemand) *Supply {
//	sp.ArgOneof = &Supply_Arg_AgentDemand{a}
//	return sp
//}
