module github.com/synerex/synerex_alpha/provider/simulation/area/communicator

require (
	github.com/RuiHirano/rvo2-go v0.0.0-20191123125933-81940413d701 // indirect
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/agent v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/area v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/clock v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/participant v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0

)

replace (
	github.com/synerex/synerex_alpha/api => ../../../../api
	github.com/synerex/synerex_alpha/api/simulation/agent => ../../../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ../../../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ../../../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/common => ../../../../api/simulation/common
	github.com/synerex/synerex_alpha/api/simulation/synerex => ../../../../api/simulation/synerex
	github.com/synerex/synerex_alpha/api/simulation/participant => ../../../../api/simulation/participant
	github.com/synerex/synerex_alpha/nodeapi => ../../../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../../../../sxutil
)
