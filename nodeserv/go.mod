module nodeid-server

require (
	github.com/google/gops v0.3.6
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/synerex/synerex_alpha/nodeapi v0.0.2
	google.golang.org/grpc v1.23.0
)

replace github.com/synerex/synerex_alpha/nodeapi => ../nodeapi
