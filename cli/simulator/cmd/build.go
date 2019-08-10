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

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build all supported providers",
	Long: `Build all supported providers`,
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
		fmt.Println("Try to build ", cmdList)
		s, err :=sioClient.Ack("build",cmdList, time.Second * 100)

		if err == nil {
			fmt.Printf("build results[%s]\n", s)
		}else{
			fmt.Printf("error %s\n",err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
