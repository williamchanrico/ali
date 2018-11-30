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
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/williamchanrico/ali/cmd/ecs"
)

var hostGroup string
var noConsulTags bool
var includeStoppedInstances bool
var generateAnsibleInventory bool

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip HOSTGROUP",
	Short: "Query active IP(s) of a service hostgroup",
	Long: `Query IP(s) of a service hostgroup with these tags:
- Env: production
- Datacenter: alisg
- Hostgroup: {HOSTGROUP}

And by default will only show running instance(s),
Use --all flag to include stopped instance(s).`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Querying IP(s) of hostgroup ")
		color.Green("%v\n", args[0])

		var invFile *os.File
		var err error
		if generateAnsibleInventory {
			invFile, err = os.OpenFile("inventory."+args[0]+".ini", os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}

			_, err = invFile.WriteString("[" + args[0] + "]\n")
			if err != nil {
				fmt.Println("Error writing hostgroup to file:", err)
				return
			}
		}
		defer invFile.Close()

		ecs := ecs.New()
		ipList, err := ecs.QueryIPList(args[0])
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to query IP list: %s", err))
		}

		for _, ip := range ipList {
			if !includeStoppedInstances && !ip.IsRunning {
				continue
			}

			if generateAnsibleInventory {
				if ip.ConsulTag == "" {
					_, err = invFile.WriteString(ip.IP + "\n")
				} else {
					_, err = invFile.WriteString(ip.IP + "\tconsul_tags=" + ip.ConsulTag + "\n")
				}
				if err != nil {
					fmt.Println("Error writing IP(s) to file:", err)
				}
			}

			color.Set(color.FgWhite)
			fmt.Printf("%v", ip.IP)

			if !noConsulTags {
				color.Set(color.FgYellow)
				fmt.Printf("\t%v\n", ip.ConsulTag)
			} else {
				fmt.Printf("\n")
			}

			color.Unset()
		}
		color.Yellow("---")
	},
}

func init() {
	ipCmd.Flags().BoolVarP(&noConsulTags, "no-tag", "n", false, "Hide consul_tags")
	ipCmd.Flags().BoolVarP(&includeStoppedInstances, "all", "a", false, "Include stopped instance(s)")
	ipCmd.Flags().BoolVarP(&generateAnsibleInventory, "generate-inv", "g", false, "Generate dynamic ansible inventory.")
	rootCmd.AddCommand(ipCmd)
}
