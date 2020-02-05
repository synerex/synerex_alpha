module github.com/synerex/synerex_alpha/provider/simulation/scenario-provider

require (
	github.com/google/uuid v1.1.1 // indirect
	github.com/mtfelian/golang-socketio v1.5.2
	github.com/paulmach/orb v0.1.5 // indirect
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/agent v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/area v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/clock v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/common v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/api/simulation/daemon v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/api/simulation/synerex v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/scenario/communicator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/scenario/simulator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/simutil/communicator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	google.golang.org/grpc v1.22.1
)

replace (
	github.com/synerex/synerex_alpha/api => ./../../../api
	github.com/synerex/synerex_alpha/api/common => ./../../../api/common
	github.com/synerex/synerex_alpha/api/simulation/agent => ./../../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ./../../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ./../../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/common => ./../../../api/simulation/common
	github.com/synerex/synerex_alpha/api/simulation/daemon => ../../../api/simulation/daemon
	github.com/synerex/synerex_alpha/api/simulation/participant => ./../../../api/simulation/participant
	github.com/synerex/synerex_alpha/api/simulation/synerex => ./../../../api/simulation/synerex
	github.com/synerex/synerex_alpha/nodeapi => ./../../../nodeapi
	github.com/synerex/synerex_alpha/provider/simulation/scenario/communicator => ../scenario2/communicator
	github.com/synerex/synerex_alpha/provider/simulation/scenario/simulator => ../scenario2/simulator
	github.com/synerex/synerex_alpha/provider/simulation/simutil/communicator => ../simutil/communicator
	github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator => ../simutil/simulator
	github.com/synerex/synerex_alpha/sxutil => ./../../../sxutil
)
