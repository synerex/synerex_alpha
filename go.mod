module github.com/synerex/synerex_alpha

require (
	github.com/bwmarrin/snowflake v0.0.0-20180412010544-68117e6bbede
	github.com/eclipse/paho.mqtt.golang v1.1.1
	github.com/golang/protobuf v1.2.0
	github.com/google/gops v0.3.5
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/kardianos/service v0.0.0-20180910224244-b1866cf76903
	github.com/mtfelian/golang-socketio v0.0.0-20181017124241-8d8ec6f9bb4c
	github.com/mtfelian/synced v0.0.0-20180626092057-b82cebd56589 // indirect
	github.com/sirupsen/logrus v1.1.1 // indirect
	golang.org/x/net v0.0.0-20181011144130-49bb7cea24b1
	google.golang.org/grpc v1.15.0
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/nodeapi v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	github.com/synerex/synerex_alpha/api/rideshare v0.0.0
	github.com/synerex/synerex_alpha/api/ptransit v0.0.0
	github.com/synerex/synerex_alpha/api/library v0.0.0
	github.com/synerex/synerex_alpha/api/fleet v0.0.0
	github.com/synerex/synerex_alpha/api/adservice v0.0.0
	github.com/synerex/synerex_alpha/api/clock v0.0.0
	github.com/synerex/synerex_alpha/api/area v0.0.0
	github.com/synerex/synerex_alpha/api/agent v0.0.0
)

replace (
	github.com/synerex/synerex_alpha/api => ./api
	github.com/synerex/synerex_alpha/nodeapi => ./nodeapi
	github.com/synerex/synerex_alpha/sxutil => ./sxutil
	github.com/synerex/synerex_alpha/api/rideshare => ./api/rideshare
	github.com/synerex/synerex_alpha/api/ptransit => ./api/ptransit
	github.com/synerex/synerex_alpha/api/library => ./api/library
	github.com/synerex/synerex_alpha/api/fleet => ./api/fleet
	github.com/synerex/synerex_alpha/api/adservice => ./api/adservice
	github.com/synerex/synerex_alpha/api/routing => ./api/routing
	github.com/synerex/synerex_alpha/api/marketing => ./api/marketing
)
