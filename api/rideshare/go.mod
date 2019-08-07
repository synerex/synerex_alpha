module github.com/synerex/synerex_alpha/api/rideshare

require (
	github.com/golang/protobuf v1.3.2
	github.com/synerex/synerex_alpha/api/common v0.0.0
)

replace (
	github.com/synerex/synerex_alpha/api => ../../api
	github.com/synerex/synerex_alpha/api/adservice => ../../api/adservice
	github.com/synerex/synerex_alpha/api/common => ../../api/common
	github.com/synerex/synerex_alpha/api/fleet => ../../api/fleet
	github.com/synerex/synerex_alpha/api/library => ../../api/library
	github.com/synerex/synerex_alpha/api/ptransit => ../../api/ptransit
	github.com/synerex/synerex_alpha/api/rideshare => ../../api/rideshare
	github.com/synerex/synerex_alpha/nodeapi => ../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../../sxutil
)
