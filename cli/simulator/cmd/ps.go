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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var bp *bool

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List the process running under Synerex Engine",
	Long: `List the process running under Synerex Engine
`,
	Run: func(cmd *cobra.Command, args []string) {
		opt := "short"
		if *bp {
			opt = "long"
		}
		s, err :=sioClient.Ack("ps",opt, time.Second * 3)

//		fmt.Printf("%s, %v\n",s,err)
		if err != nil {
			fmt.Printf("ps: stop Error %s",err.Error())
		}

		// we need to convert s:
		results := make([]string,0)
		json.Unmarshal([]byte(s),&results)
		if len(results) == 0{
			fmt.Println("No provider is running.")
		}else {
			for _, v := range results {
				fmt.Print(v)
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)

	bp = psCmd.Flags().BoolP("long", "l", false, "Long output")
}
