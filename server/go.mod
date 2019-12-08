module synerex-server

require (
	cloud.google.com/go v0.49.0 // indirect
	cloud.google.com/go/bigquery v1.3.0 // indirect
	cloud.google.com/go/pubsub v1.1.0 // indirect
	cloud.google.com/go/storage v1.4.0 // indirect
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20191203043605-d42048ed14fd // indirect
	github.com/creack/pty v1.1.9 // indirect
	github.com/envoyproxy/go-control-plane v0.9.1 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/google/pprof v0.0.0-20191205061153-f9b734f9ee64 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/prometheus/client_model v0.0.0-20191202183732-d1d2010b5bee // indirect
	github.com/rogpeppe/go-internal v1.5.0 // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/monitor/monitorapi v0.0.0
	github.com/synerex/synerex_alpha/nodeapi v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	go.opencensus.io v0.22.2 // indirect
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413 // indirect
	golang.org/x/exp v0.0.0-20191129062945-2f5052295587 // indirect
	golang.org/x/image v0.0.0-20191206065243-da761ea9ff43 // indirect
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f // indirect
	golang.org/x/mobile v0.0.0-20191130191448-5c0e7e404af8 // indirect
	golang.org/x/net v0.0.0-20191207000613-e7e4b65ae663 // indirect
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191206220618-eeba5f6aabab // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20191206204035-259af5ff87bd // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/genproto v0.0.0-20191206224255-0243a4be9c8f // indirect
	google.golang.org/grpc v1.25.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace (
	github.com/synerex/synerex_alpha/api => ../api
	github.com/synerex/synerex_alpha/api/adservice => ../api/adservice
	github.com/synerex/synerex_alpha/api/common => ../api/common
	github.com/synerex/synerex_alpha/api/fleet => ../api/fleet
	github.com/synerex/synerex_alpha/api/library => ../api/library
	github.com/synerex/synerex_alpha/api/marketing => ../api/marketing
	github.com/synerex/synerex_alpha/api/ptransit => ../api/ptransit
	github.com/synerex/synerex_alpha/api/rideshare => ../api/rideshare
	github.com/synerex/synerex_alpha/api/routing => ../api/routing
	github.com/synerex/synerex_alpha/api/simulation/agent => ../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/participant => ../api/simulation/participant
	github.com/synerex/synerex_alpha/monitor/monitorapi => ../monitor/monitorapi
	github.com/synerex/synerex_alpha/nodeapi => ../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../sxutil
)
