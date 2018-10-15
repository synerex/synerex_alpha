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

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart [provider name]",
	Short: "Restart the process running under Synerex Engine",
	Long: `Restart the process running under Synerex Engine.`,
	Run: func(cmd *cobra.Command, args []string) {

		s, err :=sioClient.Ack("restart",args, time.Second * 3)

		if err == nil {
			fmt.Printf("restart %s\n", s)
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)

}
