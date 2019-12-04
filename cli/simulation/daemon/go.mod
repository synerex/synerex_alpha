module sim-daemon

require (
	github.com/google/gops v0.3.5
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/kardianos/service v0.0.0-20180910224244-b1866cf76903
	github.com/mtfelian/golang-socketio v0.0.0-20181017124241-8d8ec6f9bb4c
	github.com/mtfelian/synced v0.0.0-20180626092057-b82cebd56589 // indirect
	github.com/sirupsen/logrus v1.1.1 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/agent v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a
)

replace (
	github.com/synerex/synerex_alpha/api => ../../../api
	github.com/synerex/synerex_alpha/api/adservice => ../../../api/adservice
	github.com/synerex/synerex_alpha/api/common => ../../../api/common
	github.com/synerex/synerex_alpha/api/fleet => ../../../api/fleet
	github.com/synerex/synerex_alpha/api/library => ../../../api/library
	github.com/synerex/synerex_alpha/api/ptransit => ../../../api/ptransit
	github.com/synerex/synerex_alpha/api/rideshare => ../../../api/rideshare
	github.com/synerex/synerex_alpha/api/routing => ../../../api/routing
	github.com/synerex/synerex_alpha/api/simulation/agent => ../../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ../../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ../../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/participant => ../../../api/simulation/participant
	github.com/synerex/synerex_alpha/api/simulation/route => ../../../api/simulation/route
	github.com/synerex/synerex_alpha/nodeapi => ../../../nodeapi
	github.com/synerex/synerex_alpha/provider/simulation/simutil => ../../../provider/simulation/simutil
	github.com/synerex/synerex_alpha/provider/simulation/simutil/objects/agent => ../../../provider/simulation/simutil/objects/agent
	github.com/synerex/synerex_alpha/sxutil => ./../../../sxutil
)
