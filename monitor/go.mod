module monitor-server

require (
	github.com/google/gops v0.3.6
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/mtfelian/golang-socketio v0.0.0-20181017124241-8d8ec6f9bb4c
	github.com/mtfelian/synced v0.0.0-20180626092057-b82cebd56589 // indirect
	github.com/sirupsen/logrus v1.1.1 // indirect
	github.com/synerex/synerex_alpha/monitor/monitorapi v0.0.2
	github.com/synerex/synerex_alpha/nodeapi v0.0.2
	google.golang.org/grpc v1.23.0
)

replace github.com/synerex/synerex_alpha/monitor/monitorapi => ./monitorapi

replace github.com/synerex/synerex_alpha/nodeapi => ../nodeapi
