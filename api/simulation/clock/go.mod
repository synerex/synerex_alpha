module github.com/synerex/synerex_alpha/api/simulation/clock

require (
	github.com/golang/protobuf v1.2.0
	golang.org/x/net v0.0.0-20181102091132-c10e9556a7bc // indirect
	golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f // indirect
)

replace (
	github.com/synerex/synerex_alpha/api => ../../../api
	github.com/synerex/synerex_alpha/api/simulation/agent => ../../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ../../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ../../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/common => ../../../api/common
	github.com/synerex/synerex_alpha/nodeapi => ../../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../../../sxutil
)
