module github.com/synerex/synerex_alpha/provider/simulation/ped-area-provider

require (
	github.com/RuiHirano/rvo2-go v1.1.1 // indirect
	github.com/RuiHirano/rvo2-go/src/rvosimulator v0.0.0-20200118052731-21c801eb6c10 // indirect
	github.com/paulmach/orb v0.1.5
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/common v0.0.0-00010101000000-000000000000
	github.com/synerex/synerex_alpha/api/simulation/participant v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/synerex v0.0.0-00010101000000-000000000000
	github.com/synerex/synerex_alpha/provider/simulation/pedestrian/agent v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/pedestrian/communicator v0.0.0-00010101000000-000000000000
	github.com/synerex/synerex_alpha/provider/simulation/pedestrian/simulator v0.0.0-00010101000000-000000000000
	github.com/synerex/synerex_alpha/provider/simulation/simutil/communicator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	golang.org/x/sys v0.0.0-20190804053845-51ab0e2deafa // indirect
	google.golang.org/grpc v1.22.1
)

replace (
	github.com/synerex/synerex_alpha/api => ./../../../api
	github.com/synerex/synerex_alpha/api/common => ./../../../api/common
	github.com/synerex/synerex_alpha/api/simulation/agent => ./../../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ./../../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ./../../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/common => ./../../../api/simulation/common
	github.com/synerex/synerex_alpha/api/simulation/participant => ./../../../api/simulation/participant
	github.com/synerex/synerex_alpha/api/simulation/synerex => ./../../../api/simulation/synerex
	github.com/synerex/synerex_alpha/nodeapi => ./../../../nodeapi
	github.com/synerex/synerex_alpha/provider/simulation/pedestrian/agent => ../pedestrian/agent
	github.com/synerex/synerex_alpha/provider/simulation/pedestrian/communicator => ../pedestrian/communicator
	github.com/synerex/synerex_alpha/provider/simulation/pedestrian/simulator => ../pedestrian/simulator
	github.com/synerex/synerex_alpha/provider/simulation/simutil/communicator => ../simutil/communicator
	github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator => ../simutil/simulator
	github.com/synerex/synerex_alpha/sxutil => ./../../../sxutil
)
