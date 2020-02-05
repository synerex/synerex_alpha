// Copyright © 2018 Synergic Mobility Project (https://synergic.mobi)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	//"strings"
	//"strconv"
	//"os/exec"
)

// cmdInfo represents the run command aliases
type orderCmdInfo struct {
	Aliases []string
	CmdName string
	Type OrderType
}

/*type Order struct {
	Type   string
	Option string
}*/

type Options struct {
	Json string
	AgentNum string
	ClockTime string
}

// Order
type OrderType int
const (
	OrderType_SET_AGENTS  OrderType = 0
    OrderType_SET_AREA  OrderType = 1
    OrderType_SET_CLOCK  OrderType = 2
    OrderType_START_CLOCK  OrderType = 3
	OrderType_STOP_CLOCK OrderType = 4
)

type Option struct{
	Key string
	Value string
}

type Order struct {
	Type   OrderType
	Name string
	Options []*Option
}

var o = Options{}

var orderCmds = [...]orderCmdInfo{
	{
		Aliases: []string{"SetClock", "setClock", "setclock", "set-clock"},
		CmdName: "SetClock",
		Type: OrderType_SET_CLOCK,
	},
	{
		Aliases: []string{"SetArea", "setArea", "setarea", "set-area"},
		CmdName: "SetArea",
		Type: OrderType_SET_AREA,
	},
	{
		Aliases: []string{"SetAgents", "setAgents", "setagents", "set-agents", "SetAgent", "setAgent", "setagent", "set-agent"},
		CmdName: "SetAgents",
		Type: OrderType_SET_AGENTS,
	},
	{
		Aliases: []string{"StartClock", "startClock", "start"},
		CmdName: "StartClock",
		Type: OrderType_START_CLOCK,
	},
	{
		Aliases: []string{"StopClock", "stopClock", "stop"},
		CmdName: "StopClock",
		Type: OrderType_STOP_CLOCK,
	},
}

func getOrderCmdName(alias string) string {
	for _, ci := range orderCmds {
		for _, str := range ci.Aliases {
			if alias == str {
				return ci.CmdName
			}
		}
	}
	return "" // can't find alias
}

func sendOrder(cmdName string, order *Order) bool {
	//todo: we should use ack for this. but its not working....
	fmt.Printf("simulator order [%v]\n", order)
	res, err := sioClient.Ack("order", order, 20*time.Second)
	//					err := sioClient.Emit("run",ci.CmdName) //, 20*time.Second)
	time.Sleep(1 * time.Second)

	if err != nil || res != "\"ok\"" {
		fmt.Printf("simulator: Got error on reply:'%s',%v\n", res, err)
		return false
	} else {
		fmt.Printf("simulator: Reply [%s]\n", res)
		fmt.Printf("simulator: Run '%s' succeeded.\n", cmdName)
		return true
	}
}

func handleOrder(cmd *cobra.Command, args []string) {

	//simData := handleUserDialogue()
	fmt.Printf("Dialogue Result: %v\n", o)
	if len(args) > 0 {
		findflag := false
		order := new(Order)
		//order.Option = "&o"
		for _, ci := range orderCmds {
			for _, str := range ci.Aliases {
				if args[0] == str {
					switch ci.CmdName {
					case "SetAgents":
						//order.Option = o.optAgentNum
					}

					fmt.Printf("simulator: Starting '%s'\n", ci.CmdName)
					order.Type = ci.Type
					order.Name = ci.CmdName
					findflag = sendOrder(ci.CmdName, order)
					break
				}
			}

		}
		if !findflag {
			fmt.Printf("simulation: Can't find command run '%s'.\n", args[0])
			fmt.Printf("cmd is:'%s'\n", orderCmds)

		}
	}
}

var orderCmd = &cobra.Command{
	Use:   "order [order name] [options..]",
	Short: "Start a provider",
	Long: `Start a provider with options 
For example:
    simulation order start   
	simulation order set-time   
	simulation order set-area   
`,
	Run: handleOrder,
}

func init() {
	rootCmd.AddCommand(orderCmd)
	orderCmd.Flags().StringVarP(&o.Json, "json", "j", "sample.json", "json name")
	orderCmd.Flags().StringVarP(&o.AgentNum, "num", "n", "1", "agent num")
	orderCmd.Flags().StringVarP(&o.ClockTime, "time", "t", "0", "clock time")
}
