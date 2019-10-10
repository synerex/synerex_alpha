package simutil

import (
  "testing"
  "github.com/synerex/synerex_alpha/api"
)

// test IsSupplyTarget
func TestIsSupplyTarget(t *testing.T) {
	t.Log("IsSupplyTargetのテスト")

	t.Run("targetIdが自分のdmIdlistに含まれていればTrue", func(t *testing.T){
		sp := &api.Supply{
			TargetId:   uint64(1000),
		}
		idlist := []uint64{uint64(1000)}
		result := IsSupplyTarget(sp, idlist)
		expext := true
		if result != expext {
		  t.Error("\n実際： ", result, "\n理想： ", expext)
		}
	})	

	t.Run("targetIdが自分のdmIdlistに含まれていなければFalse", func(t *testing.T){
		sp := &api.Supply{
			TargetId:   uint64(1000),
		}
		idlist := []uint64{uint64(2000)}
		result := IsSupplyTarget(sp, idlist)
		expext := false
		if result != expext {
		  t.Error("\n実際： ", result, "\n理想： ", expext)
		}
	})	
  
	t.Log("IsSupplyTargetのテスト終了")
}

// test IsFinishSync
func TestIsFinishSync(t *testing.T) {
	t.Log("IsFinishSyncのテスト")

	t.Run("idlist内のIdがpspMapに含まれていればTrue", func(t *testing.T){
		pspMap := map[uint64]*api.Supply{
			1: &api.Supply{
				SenderId:   uint64(1000),
			},
			2: &api.Supply{
				SenderId:   uint64(2000),
			},
		}
		idlist := []uint64{uint64(1000), uint64(2000)}
		result := IsFinishSync(pspMap, idlist)
		expext := true
		if result != expext {
		  t.Error("\n実際： ", result, "\n理想： ", expext)
		}
	})	

	t.Run("idlist内のIdがpspMapに含まれていなければFalse", func(t *testing.T){
		pspMap := map[uint64]*api.Supply{
			1: &api.Supply{
				SenderId:   uint64(1000),
			},
		}
		idlist := []uint64{uint64(1000), uint64(2000)}
		result := IsFinishSync(pspMap, idlist)
		expext := false
		if result != expext {
		  t.Error("\n実際： ", result, "\n理想： ", expext)
		}
	})	
  
	t.Log("IsFinishSyncのテスト終了")
}