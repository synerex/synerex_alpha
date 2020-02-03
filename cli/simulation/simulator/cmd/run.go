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
)

// cmdInfo represents the run command aliases
type cmdInfo struct {
	Aliases []string
	CmdName string
}

var cmds = [...]cmdInfo{
	{
		Aliases: []string{"All", "all"},
		CmdName: "All",
	},
	{
		Aliases: []string{"nodeserv", "nodesrv", "ndsrv", "NodeIDServer"},
		CmdName: "NodeIDServer",
	},
	{
		Aliases: []string{"Synerex", "smarket", "server", "synerex", "SynerexServer"},
		CmdName: "SynerexServer",
	},
	{
		Aliases: []string{"monitor", "MonitorServer", "mon"},
		CmdName: "MonitorServer",
	},
	{
		Aliases: []string{"Area", "area"},
		CmdName: "Area",
	},
	{
		Aliases: []string{"Scenario", "scenario"},
		CmdName: "Scenario",
	},
	{
		Aliases: []string{"Pedestrian", "pedestrian", "ped"},
		CmdName: "Pedestrian",
	},
	{
		Aliases: []string{"Car", "car"},
		CmdName: "Car",
	},
	{
		Aliases: []string{"simulator", "sim", "Simulator"},
		CmdName: "Simulator",
	},
}

func getCmdName(alias string) string {
	for _, ci := range cmds {
		for _, str := range ci.Aliases {
			if alias == str {
				return ci.CmdName
			}
		}
	}
	return "" // can'f find alias
}

func handleProvider(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		for n := range args {
			findflag := false
			fmt.Printf("cmd is:'%s'\n", cmds)
			for _, ci := range cmds {
				for _, str := range ci.Aliases {
					if args[n] == str {
						fmt.Printf("simulator: Starting '%s'\n", ci.CmdName)

						//todo: we should use ack for this. but its not working....
						res, err := sioClient.Ack("run", ci.CmdName, 20*time.Second)
						//					err := sioClient.Emit("run",ci.CmdName) //, 20*time.Second)
						time.Sleep(3 * time.Second)

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
				fmt.Printf("simulator: Can't find command run '%s'.\n", args[n])
				fmt.Printf("cmd is:'%s'\n", cmds)
				break
			}
		}
	}
}

var runCmd = &cobra.Command{
	Use:   "run [provider name] [options..]",
	Short: "Start a provider",
	Long: `Start a provider with options 
For example:
    se run nodeserv   // start a node server
	se run server     // start a synergic exchange server
    se run taxi       // start a taxi provider

    se run all        // start all basic providers

`,
	Run: handleProvider,
}

func init() {
	rootCmd.AddCommand(runCmd)
}
