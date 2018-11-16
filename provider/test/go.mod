module test-provider

require (
	github.com/synerex/synerex_alpha v0.0.1 // indirect
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/common v0.0.0
	github.com/synerex/synerex_alpha/api/ptransit v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	golang.org/x/sys v0.0.0-20181107165924-66b7b1311ac8
	google.golang.org/grpc v1.16.0
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
