// Copyright © 2017 Jake Teton-Landis <just.1.jake@gmail.com>
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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/justjake/encabulator/rekey"
)

var mapNameToRekey = map[string](func() *rekey.KeyEnsurer){
	"yubikey": rekey.Yubikey,
	"goldkey": rekey.Goldkey,
	"id":      rekey.DefaultIdentity,
}
var names = keys(mapNameToRekey)

// rekeyCmd represents the rekey command
var rekeyCmd = &cobra.Command{
	Use:   "rekey",
	Short: "Ensure that your SSH keys are loaded into ssh-agent",
	Long: `Checks ssh-agent for the given types of keys. If those keys are not
loaded into ssh-agent, rekey will ask you to add them.

Rekey is handy when working with hardware tokens with short caching lifetimes.
Add rekey to your aliases to ensure your commands always have the SSH keys that
they need.
`,
	Example: `  encabulator rekey --type=id && ssh foo@example.com
  encabulator rekey && git pull`,
	ValidArgs: names,
	// check to see if arguments are copacetic
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// no args
		if len(args) > 0 {
			return fmt.Errorf("Unknown arguments %v", args)
		}

		// type must be valid
		name := viper.GetString("rekey.Type")
		if _, ok := mapNameToRekey[name]; !ok {
			return fmt.Errorf("Unknown type %s.", name)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		kill := viper.GetBool("rekey.KillAgent")
		name := viper.GetString("rekey.Type")

		loader := mapNameToRekey[name]()
		if kill {
			loader = rekey.EnsureRestartingAgent(loader)
		}

		err := loader.EnsureLoaded()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	},
}

func keys(m map[string]func() *rekey.KeyEnsurer) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func init() {
	RootCmd.AddCommand(rekeyCmd)
	viper.SetDefault("rekey.Type", "yubikey")
	viper.SetDefault("rekey.Kill", true)

	// --type: type of flag
	rekeyCmd.Flags().StringP("type", "t", viper.GetString("rekey.Type"), fmt.Sprintf("Type of key to load. One of %v.", names))
	viper.BindPFlag("rekey.Type", rekeyCmd.Flags().Lookup("type"))

	// --kill: kill agent before attempting to load keys
	rekeyCmd.Flags().BoolP("kill", "k", viper.GetBool("rekey.Kill"), "Kill ssh-agent before loading any keys. This is useful on macOS.")
	viper.BindPFlag("rekey.KillAgent", rekeyCmd.Flags().Lookup("kill"))
}
