// Copyright © 2018 William Chanrico
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
	"strings"

	"github.com/fatih/color"

	"github.com/spf13/cobra"
	"github.com/williamchanrico/ali/cmd/ecs"
)

var hostGroup string

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip [HOSTGROUP]",
	Short: "Query active IP(s) of a service hostgroup",
	Long: `Query IP(s) of a service hostgroup with these tags:
- Env: production
- Datacenter: alisg
- Hostgroup: [HOSTGROUP]`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Querying IP(s) of hostgroup ")
		color.Green("%v\n", args[0])

		ecs := ecs.New()
		ipList, err := ecs.QueryIPList(args[0])
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to query IP list: %s", err))
		}

		color.Yellow("\n---")
		color.White(strings.Join(ipList, "\n"))
		color.Yellow("---")
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}
