package main

import (
  "testing"
  "github.com/synerex/synerex_alpha/api/simulation/agent"
  "github.com/synerex/synerex_alpha/api/simulation/area"
)

// test isAgentInArea
func TestIsAgentInArea(t *testing.T) {
	t.Log("IsAgentInAreaのテスト")

	t.Run("Ped, 位置がエリア内の場合にTrueが返る", func(t *testing.T){
		agentInfo := &agent.AgentInfo{
			AgentType: 0, //Ped
			Route: &agent.Route{
				Coord: &agent.Route_Coord{
					Lat: float32(35.156678),
					Lon: float32(136.977031), 
				},
			},
		}
		data := &Data{
			AreaInfo: &area.AreaInfo{
				Map: &area.Map{
					Coord: &area.Map_Coord{
						StartLat: float32(35.152476),
						StartLon: float32(136.973172),
						EndLat: float32(35.160678),
						EndLon: float32(136.984031), 
					},
				},
			},
		}
		agentType := int(0)	//Ped
		result := isAgentInArea(agentInfo, data, agentType)
		expext := true
		if result != expext {
		  t.Error("\n実際： ", result, "\n理想： ", expext)
		}
	})

	t.Run("Ped, 位置がエリア外の場合にFalseが返る", func(t *testing.T){
		agentInfo := &agent.AgentInfo{
			AgentType: 0, //Ped
			Route: &agent.Route{
				Coord: &agent.Route_Coord{
					Lat: float32(35.176678),
					Lon: float32(136.977031), 
				},
			},
		}
		data := &Data{
			AreaInfo: &area.AreaInfo{
				Map: &area.Map{
					Coord: &area.Map_Coord{
						StartLat: float32(35.152476),
						StartLon: float32(136.973172),
						EndLat: float32(35.160678),
						EndLon: float32(136.984031), 
					},
				},
			},
		}
		agentType := int(0)	//Ped
		result := isAgentInArea(agentInfo, data, agentType)
		expext := false
		if result != expext {
		  t.Error("\n実際： ", result, "\n理想： ", expext)
		}
	})

	t.Run("AgentTypeが違う場合にFalseが返る", func(t *testing.T){
		agentInfo := &agent.AgentInfo{
			AgentType: 1, //Car
			Route: &agent.Route{
				Coord: &agent.Route_Coord{
					Lat: float32(35.176678),

					Lon: float32(136.977031), 
				},
			},
		}
		data := &Data{
			AreaInfo: &area.AreaInfo{
				Map: &area.Map{
					Coord: &area.Map_Coord{
						StartLat: float32(35.152476),
						StartLon: float32(136.973172),
						EndLat: float32(35.160678),
						EndLon: float32(136.984031), 
					},
				},
			},
		}
		agentType := int(0)	//Ped
		result := isAgentInArea(agentInfo, data, agentType)
		expext := false
		if result != expext {
		  t.Error("\n実際： ", result, "\n理想： ", expext)
		}
	})
	
  
	t.Log("TestIsAgentInAreaのテスト終了")
}
