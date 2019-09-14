module github.com/synerex/synerex_alpha/provider/simulation/simutil

require (
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0
)

replace (
	github.com/synerex/synerex_alpha/provider/simulation/simutil => ../simutil
	github.com/synerex/synerex_alpha/api => ./../../../api
	github.com/synerex/synerex_alpha/sxutil => ./../../../sxutil
)