// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/williamchanrico/ali/cmd/ecs"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list TARGET_FILE",
	Short: "Generate a JSON file containing all hostgroups with IP(s).",
	Long: `Generate a JSON file containing all hostgroups along with it's
associated IP addresses.

Example of the generated JSON:
	{
	  "tkp-asdf": [
		"172.21.42.220"
	  ],
	  "tkp-qwerty": [
		"172.21.45.11"
	  ],
	  "tkp-hjkl": [
		"172.21.45.20"
	  ]
	}
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		color.White("Querying all hostgroups")
		color.White("---")
		ecs := ecs.New()
		hostgroupList, err := ecs.QueryAllHostgroupIP()
		if err != nil {
			fmt.Println("Error querying all hostgroup IP(s):", err)
			return
		}

		var invFile *os.File
		invFile, err = os.OpenFile(args[0], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer invFile.Close()

		hostgroupListJSON, err := json.Marshal(hostgroupList)
		if err != nil {
			fmt.Println("Error converting hostgroup list to JSON:", err)
			return
		}

		_, err = invFile.WriteString(string(hostgroupListJSON))
		if err != nil {
			fmt.Println("Error writing hostgroup list to file:", err)
			return
		}

		color.Green("Finished writing to %v", args[0])
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
