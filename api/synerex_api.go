package api

import (
	"github.com/synerex/synerex_alpha/api/simulation/synerex"
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

func (dm *Demand) WithSimDemand(r *synerex.SimDemand) *Demand {
	dm.ArgOneof = &Demand_SimDemand{r}
	return dm
}

func (sp *Supply) WithSimSupply(c *synerex.SimSupply) *Supply {
	sp.ArgOneof = &Supply_SimSupply{c}
	return sp
}
