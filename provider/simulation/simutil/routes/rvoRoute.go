package routes

import (
	"fmt"
	"log"

	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/api/simulation/simutil"
)

var (
	isRVO2 bool
)

func init() {
	isRVO2 = false
}
