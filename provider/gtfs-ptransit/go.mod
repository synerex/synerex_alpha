module ptransit-provider

require (
	github.com/nkawa/gtfsparser v0.0.0-20181118053010-0f0c06281a7e
	github.com/synerex/synerex_alpha/api v0.0.2
	github.com/synerex/synerex_alpha/api/common v0.0.2
	github.com/synerex/synerex_alpha/api/ptransit v0.0.2
	github.com/synerex/synerex_alpha/sxutil v0.0.2
	google.golang.org/grpc v1.23.0
)

replace (
	github.com/synerex/synerex_alpha/api => ../../api
	github.com/synerex/synerex_alpha/api/adservice => ../../api/adservice
	github.com/synerex/synerex_alpha/api/common => ../../api/common
	github.com/synerex/synerex_alpha/api/fleet => ../../api/fleet
	github.com/synerex/synerex_alpha/api/library => ../../api/library
	github.com/synerex/synerex_alpha/api/ptransit => ../../api/ptransit
	github.com/synerex/synerex_alpha/api/rideshare => ../../api/rideshare
	github.com/synerex/synerex_alpha/api/routing => ../../api/routing
	github.com/synerex/synerex_alpha/nodeapi => ../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../../sxutil
)
