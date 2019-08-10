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
	"github.com/spf13/cobra"
	"time"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [provider names]",
	Short: "Stop the process running under Synerex Engine",
	Long: `Stop the process running under Synerex Engine`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdList := make([]string,len(args))
		n := 0
		for _, alias := range args {
			cmd := getCmdName(alias)
			if len(cmd) != 0 {
				cmdList[n] =  cmd
				n++
			}
		}
		fmt.Println("Try to stop ", cmdList)
		s, err :=sioClient.Ack("stop",cmdList, time.Second * 3)

		if err == nil {
			fmt.Printf("stop %s\n", s)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)


}
