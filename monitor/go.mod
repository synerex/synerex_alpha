module monitor-server

require (
	github.com/google/gops v0.3.5
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/mtfelian/golang-socketio v0.0.0-20181017124241-8d8ec6f9bb4c
	github.com/mtfelian/synced v0.0.0-20180626092057-b82cebd56589 // indirect
	github.com/sirupsen/logrus v1.1.1 // indirect
	github.com/synerex/synerex_alpha/monitor/monitorapi v0.0.1
	github.com/synerex/synerex_alpha/nodeapi v0.0.1
	google.golang.org/grpc v1.16.0
)

replace github.com/synerex/synerex_alpha/monitor/monitorapi => ./monitorapi

replace github.com/synerex/synerex_alpha/nodeapi => ../nodeapi
