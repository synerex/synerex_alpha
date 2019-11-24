package test

import (
	"testing"

	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
)

var (
	simpleRoute *simutil.SimpleRoute
)

func init(){
}

func TestIsAgentInControlledArea(t *testing.T) {
	t.Log("IsAgentInControlledAreaのテスト")

	t.Run("targetIdが自分のdmIdlistに含まれていればTrue", func(t *testing.T) {
		sp := &api.Supply{
			TargetId: uint64(1000),
		}
		idlist := []uint64{uint64(1000)}
		result := simutil.IsSupplyTarget(sp, idlist)
		expext := true
		if result != expext {
			t.Error("\n実際： ", result, "\n理想： ", expext)
		}
	})

	t.Log("IsAgentInControlledAreaのテスト終了")
}