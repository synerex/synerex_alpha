module github.com/synerex/synerex_alpha/api/agent

require (
	github.com/golang/protobuf v1.2.0
	golang.org/x/net v0.0.0-20181102091132-c10e9556a7bc // indirect
	golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f // indirect
)

replace (
	github.com/synerex/synerex_alpha/api => ../../api
	github.com/synerex/synerex_alpha/api/agent => ../../api/agent
	github.com/synerex/synerex_alpha/api/area => ../../api/area
	github.com/synerex/synerex_alpha/api/clock => ../../api/clock
	github.com/synerex/synerex_alpha/api/common => ../../api/common
	github.com/synerex/synerex_alpha/nodeapi => ../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../../sxutil
)
