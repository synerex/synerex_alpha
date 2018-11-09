module github.com/synerex/synerex_alpha/sxutil

require (
	github.com/bwmarrin/snowflake v0.0.0-20180412010544-68117e6bbede
	github.com/golang/protobuf v1.2.0
	github.com/synerex/synerex_alpha v0.0.1
	github.com/synerex/synerex_alpha/api v0.0.1
	github.com/synerex/synerex_alpha/api/adservice v0.0.1
	github.com/synerex/synerex_alpha/api/fleet v0.0.1
	github.com/synerex/synerex_alpha/api/library v0.0.1
	github.com/synerex/synerex_alpha/api/ptransit v0.0.1
	github.com/synerex/synerex_alpha/api/rideshare v0.0.1
	google.golang.org/grpc v1.16.0
)

replace github.com/synerex/synerex_alph/api/rideshare => ../api/rideshare

replace github.com/synerex/synerex_alph/api/ptransit => ../api/ptransit

replace github.com/synerex/synerex_alph/api/library => ../api/library

replace github.com/synerex/synerex_alph/api/fleet => ../api/fleet

replace github.com/synerex/synerex_alph/api/adservice => ../api/adservice
