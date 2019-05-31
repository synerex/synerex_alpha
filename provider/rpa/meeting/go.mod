module meeting

require (
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/fleet v0.0.0
	github.com/synerex/synerex_alpha/api/routing v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	github.com/tidwall/gjson v1.2.1
	github.com/tidwall/match v1.0.1 // indirect
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	google.golang.org/grpc v1.17.0
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
	github.com/synerex/synerex_alpha/nodeapi => ../../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../../../sxutil
)
