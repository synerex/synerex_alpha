module selenium

require (
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/sclevine/agouti v3.0.0+incompatible
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
