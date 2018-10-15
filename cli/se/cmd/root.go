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
	"github.com/mtfelian/golang-socketio/transport"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mtfelian/golang-socketio"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var sioClient *gosocketio.Client
var Providers []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sm [OPTIONS] COMMAND [ARG...]",
	Short: "Synergic Exchange command launcher",
	Long: `Synergic Exchange command launcher
For example:

se is a CLI launcher for Synergic Exchange.

   se run all
   se status   // show status of provider/servers
`,
	//	Run: func(cmd *cobra.Command, args[]string){},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
//		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.se.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var st ="01234567890012345678900123456789001234567890012345678900123456789001234567890"
	var lp = &st

	rootCmd.Flags().StringVar(lp, "se-addr", "ws://localhost:9995/socket.io/?EIO=3&transport=websocket", "Default address for se-daemon")
	// need socket.io connection with daemon.
	var err error
	sioClient, err = gosocketio.Dial(*lp, transport.DefaultWebsocketTransport())
	if err != nil {
		fmt.Println("se: Error to connect with se-daemon. You have to start se-daemon first.") //,err)
		os.Exit(1)
	}
	sioClient.On(gosocketio.OnConnection, func(c *gosocketio.Channel,param interface{}) {
//		fmt.Println("Go socket.io connected ",c)
	})
	sioClient.On("providers", func(c *gosocketio.Channel,param interface{}) {
		//fmt.Println("Get Providers ",param)
		// we have to keep this to check parameters
		procs := param.([]interface{})
		Providers = make([]string, len(procs))
		for i, pp := range procs {
			Providers[i] = pp.(string)
		}
	})

	sioClient.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel,param interface{}) {
		fmt.Println("Go socket.io disconnected ",c)
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".se" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".se")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
