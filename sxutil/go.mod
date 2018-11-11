module github.com/synerex/synerex_alpha/sxutil

require (
	cloud.google.com/go v0.32.0 // indirect
	github.com/bwmarrin/snowflake v0.0.0-20180412010544-68117e6bbede
	github.com/golang/lint v0.0.0-20181026193005-c67002cb31c3 // indirect
	github.com/golang/protobuf v1.2.0
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/adservice v0.0.0
	github.com/synerex/synerex_alpha/api/common v0.0.0
	github.com/synerex/synerex_alpha/api/fleet v0.0.0
	github.com/synerex/synerex_alpha/api/library v0.0.0
	github.com/synerex/synerex_alpha/api/ptransit v0.0.0
	github.com/synerex/synerex_alpha/api/rideshare v0.0.0
	github.com/synerex/synerex_alpha/nodeapi v0.0.0
	golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3 // indirect
	golang.org/x/net v0.0.0-20181108082009-03003ca0c849 // indirect
	golang.org/x/oauth2 v0.0.0-20181106182150-f42d05182288 // indirect
	golang.org/x/sync v0.0.0-20181108010431-42b317875d0f // indirect
	golang.org/x/sys v0.0.0-20181107165924-66b7b1311ac8 // indirect
	golang.org/x/tools v0.0.0-20181108221941-77439c55185e // indirect
	google.golang.org/appengine v1.3.0 // indirect
	google.golang.org/genproto v0.0.0-20181107211654-5fc9ac540362 // indirect
	google.golang.org/grpc v1.16.0
	honnef.co/go/tools v0.0.0-20180920025451-e3ad64cb4ed3 // indirect
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
