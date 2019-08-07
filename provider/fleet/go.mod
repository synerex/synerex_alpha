module fleet-provider

require (
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/mtfelian/golang-socketio v0.0.0-20181017124241-8d8ec6f9bb4c
	github.com/mtfelian/synced v0.0.0-20181026093311-f1dd911faaa7 // indirect
	github.com/sirupsen/logrus v1.2.0 // indirect
	github.com/synerex/synerex_alpha/api v0.0.1
	github.com/synerex/synerex_alpha/api/fleet v0.0.1
	github.com/synerex/synerex_alpha/sxutil v0.0.1
	google.golang.org/grpc v1.22.1
)

replace (
	github.com/synerex/synerex_alpha/api => ../../api
	github.com/synerex/synerex_alpha/api/adservice => ../../api/adservice
	github.com/synerex/synerex_alpha/api/common => ../../api/common
	github.com/synerex/synerex_alpha/api/fleet => ../../api/fleet
	github.com/synerex/synerex_alpha/api/library => ../../api/library
	github.com/synerex/synerex_alpha/api/ptransit => ../../api/ptransit
	github.com/synerex/synerex_alpha/api/rideshare => ../../api/rideshare
	github.com/synerex/synerex_alpha/api/routing => ../../api/routing
	github.com/synerex/synerex_alpha/api/simulation/agent => ../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ../../api/simulation/clock
	github.com/synerex/synerex_alpha/nodeapi => ../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../../sxutil
)
