module github.com/synerex/synerex_alpha/provider/simulation/visualization-provider

require (
	github.com/mtfelian/golang-socketio v1.5.2
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/agent v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/area v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/clock v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/common v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/api/simulation/synerex v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/simutil/communicator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/visualization/communicator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/visualization/simulator v0.0.0-00010101000000-000000000000 // indirect
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
	github.com/synerex/synerex_alpha/api/simulation/participant => ./../../../api/simulation/participant
	github.com/synerex/synerex_alpha/api/simulation/synerex => ./../../../api/simulation/synerex
	github.com/synerex/synerex_alpha/nodeapi => ./../../../nodeapi
	github.com/synerex/synerex_alpha/provider/simulation/simutil/communicator => ../simutil/communicator
	github.com/synerex/synerex_alpha/provider/simulation/simutil/simulator => ../simutil/simulator
	github.com/synerex/synerex_alpha/provider/simulation/visualization/communicator => ../visualization/communicator
	github.com/synerex/synerex_alpha/provider/simulation/visualization/simulator => ../visualization/simulator
	github.com/synerex/synerex_alpha/sxutil => ./../../../sxutil
)
