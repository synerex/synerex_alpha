module nodeid-server

require (
	github.com/google/gops v0.3.6
	github.com/synerex/synerex_alpha/nodeapi v0.0.2
	google.golang.org/grpc v1.23.0
)

replace github.com/synerex/synerex_alpha/nodeapi => ../nodeapi
