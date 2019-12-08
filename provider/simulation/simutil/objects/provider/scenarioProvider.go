package provider

import (
	"github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/agent"
)

const (
	SET_AGENTS      = "SET_AGENTS"
	GET_PARTICIPANT = "GET_PARTICIPANT"
	START_CLOCK     = "START_CLOCK"
	STOP_CLOCK      = "STOP_CLOCK"
)

type SetAgentsOrder struct {
	Type string
	Peds []*agent.Pedestrian
	Cars []*agent.Car
}

type GetParticipantOrder struct {
	Type string
}

type StartClockOrder struct {
	Type string
}

type StopClockOrder struct {
	Type string
}

// ScenarioProvider :
type ScenarioProvider struct {
	*SynerexProvider //埋め込み
}

// NewScenarioProvider:
func NewScenarioProvider() *ScenarioProvider {

	provider := &ScenarioProvider{
		SynerexProvider: NewSynerexProvider(),
	}

	return provider
}
