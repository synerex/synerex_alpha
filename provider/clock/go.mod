module github.com/synerex/synerex_alpha/provider/user

require (
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	google.golang.org/grpc v1.17.0
)

replace (
	github.com/synerex/synerex_alpha/api => ./../../api
	github.com/synerex/synerex_alpha/api/agent => ./../../api/agent
	github.com/synerex/synerex_alpha/api/area => ./../../api/area
	github.com/synerex/synerex_alpha/api/clock => ./../../api/clock
	github.com/synerex/synerex_alpha/api/common => ./../../api/common
	github.com/synerex/synerex_alpha/nodeapi => ./../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ./../../sxutil
)