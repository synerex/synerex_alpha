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
}

type Order struct {
	Type string
	Arg  string
}

type Options struct {
	optJsonName string
}

var o = &Options{}

var orderCmds = [...]orderCmdInfo{
	{
		Aliases: []string{"SetAll", "setAll"},
		CmdName: "SetAll",
	},
	{
		Aliases: []string{"SetTime", "setTime", "settime", "set-time"},
		CmdName: "SetTime",
	},
	{
		Aliases: []string{"SetArea", "setArea", "setarea", "set-area"},
		CmdName: "SetArea",
	},
	{
		Aliases: []string{"SetAgent", "setAgent", "setagent", "set-agent"},
		CmdName: "SetAgent",
	},
	{
		Aliases: []string{"Start", "start"},
		CmdName: "Start",
	},
	{
		Aliases: []string{"Stop", "stop"},
		CmdName: "Stop",
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
	return "" // can'f find alias
}

func handleOrder(cmd *cobra.Command, args []string) {

	//simData := handleUserDialogue()
	fmt.Printf("Dialogue Result: %v\n", args)
	if len(args) > 0 {
		for n := range args {
			findflag := false
			for _, ci := range orderCmds {
				for _, str := range ci.Aliases {
					if args[n] == str {
						fmt.Printf("simulator: Starting '%s'\n", ci.CmdName)

						//todo: we should use ack for this. but its not working....
						res, err := sioClient.Ack("order", &Order{Type: ci.CmdName}, 20*time.Second)
						//					err := sioClient.Emit("run",ci.CmdName) //, 20*time.Second)
						time.Sleep(1 * time.Second)

						if err != nil || res != "\"ok\"" {
							fmt.Printf("simulator: Got error on reply:'%s',%v\n", res, err)
							return
						} else {
							fmt.Printf("simulator: Reply [%s]\n", res)
							fmt.Printf("simulator: Run '%s' succeeded.\n", ci.CmdName)
							findflag = true
						}
						break
					}
				}

			}
			if !findflag {
				fmt.Printf("simulation: Can't find command run '%s'.\n", args[n])
				fmt.Printf("cmd is:'%s'\n", orderCmds)
				break
			}
		}
	}
}

func handleSetAll(cmd *cobra.Command, args []string) {
	cmdName := "SetAll"
	//simData := handleUserDialogue()
	fmt.Printf("Dialogue Result: %v\n", o.optJsonName)
	fmt.Printf("simulator: Starting SetAll Order\n")

	res, err := sioClient.Ack("order", &Order{Type: cmdName, Arg: o.optJsonName}, 20*time.Second)
	time.Sleep(1 * time.Second)

	if err != nil || res != "\"ok\"" {
		fmt.Printf("simulator: Got error on reply:'%s',%v\n", res, err)
		return
	} else {
		fmt.Printf("simulator: Reply [%s]\n", res)
		fmt.Printf("simulator: Run '%s' succeeded.\n", cmdName)
	}
}

func handleStart(cmd *cobra.Command, args []string) {
	cmdName := "Start"
	fmt.Printf("simulator: Starting Start Order\n")

	res, err := sioClient.Ack("order", &Order{Type: cmdName}, 20*time.Second)
	time.Sleep(1 * time.Second)

	if err != nil || res != "\"ok\"" {
		fmt.Printf("simulator: Got error on reply:'%s',%v\n", res, err)
		return
	} else {
		fmt.Printf("simulator: Reply [%s]\n", res)
		fmt.Printf("simulator: Run '%s' succeeded.\n", cmdName)
	}
}

func handleStop(cmd *cobra.Command, args []string) {
	cmdName := "Stop"
	fmt.Printf("simulator: Starting Stop Order\n")

	res, err := sioClient.Ack("order", &Order{Type: cmdName}, 20*time.Second)
	time.Sleep(1 * time.Second)

	if err != nil || res != "\"ok\"" {
		fmt.Printf("simulator: Got error on reply:'%s',%v\n", res, err)
		return
	} else {
		fmt.Printf("simulator: Reply [%s]\n", res)
		fmt.Printf("simulator: Run '%s' succeeded.\n", cmdName)
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

var setAllCmd = &cobra.Command{
	Use:   "setAll",
	Short: "set time, area and agent info",
	Long: `-j options 
	For example:
    ./simulator setAll -j (json-name)  `,
	Run: handleSetAll,
}

var startClockCmd = &cobra.Command{
	Use:   "start",
	Short: "start simulation",
	Long: ` 
	For example:
    ./simulator start`,
	Run: handleStart,
}

var stopClockCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop simulation",
	Long: ` 
	For example:
    ./simulator stop`,
	Run: handleStop,
}

func init() {
	rootCmd.AddCommand(orderCmd)
	rootCmd.AddCommand(setAllCmd)
	rootCmd.AddCommand(startClockCmd)
	rootCmd.AddCommand(stopClockCmd)
	setAllCmd.Flags().StringVarP(&o.optJsonName, "json", "j", "sample.json", "string option")
}
