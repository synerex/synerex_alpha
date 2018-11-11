module synerex-server

require (
	github.com/sirupsen/logrus v1.2.0
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/monitor/monitorapi v0.0.0
	github.com/synerex/synerex_alpha/nodeapi v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	golang.org/x/crypto v0.0.0-20181106171534-e4dc69e5b2fd // indirect
	google.golang.org/grpc v1.16.0
)

replace (
	github.com/synerex/synerex_alpha/api => ../api
	github.com/synerex/synerex_alpha/api/adservice => ../api/adservice
	github.com/synerex/synerex_alpha/api/common => ../api/common
	github.com/synerex/synerex_alpha/api/fleet => ../api/fleet
	github.com/synerex/synerex_alpha/api/library => ../api/library
	github.com/synerex/synerex_alpha/api/ptransit => ../api/ptransit
	github.com/synerex/synerex_alpha/api/rideshare => ../api/rideshare
	github.com/synerex/synerex_alpha/monitor/monitorapi => ../monitor/monitorapi
	github.com/synerex/synerex_alpha/nodeapi => ../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../sxutil
)
