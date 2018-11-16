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
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/williamchanrico/ali/cmd/ecs"
	"github.com/williamchanrico/ali/cmd/util"
)

const timeFormat = "2006-01-02 15:04:05"

// memoryUsageCmd represents the memoryUsage command
var memoryUsageCmd = &cobra.Command{
	Use:   "memoryUsage HOSTGROUP_A,HOSTGROUP_B,...",
	Short: "Get and calculate the memory usage of instance(s) by hostgroups",
	Long: `Get and calculate the memory usage of instance(s) by hostgroups.
HOSTGROUP_LIST is separated by comma (default, use -d flag to specify delimiter).

Will output the following information: <ip_address> <creation_time_in_utc> <memory_mb>,
along with total memory usage.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inv := strings.Split(args[0], ",")

		ecs := ecs.New()

		totalMem := 0
		summary := map[string]int{}
		for _, hostgroup := range inv {
			hgTotalMem := 0

			instanceList, err := ecs.QueryInstanceList(hostgroup)
			if err != nil {
				color.Red("Failed to query instance list:", err)
			}

			if len(instanceList) == 0 {
				color.Red("!!! Hostgroup %v has no instance(s) !!!\n", hostgroup)
				continue
			}

			color.Green("[%s]", hostgroup)
			for _, i := range instanceList {
				creationTime, err := util.ParseTimeStr(i.CreationTime)
				if err != nil {
					creationTime = time.Time{}
				}
				// if creationTime.Before(time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC)) {
				// }

				color.Yellow("%s\t%s\t%dMB", i.IP, creationTime.Format(timeFormat), i.Memory)
				hgTotalMem += i.Memory
			}
			totalMem += hgTotalMem
			summary[hostgroup] = hgTotalMem

			color.White("------------------\n")
			color.White("Mem usage:\t\t\t\t%dMB\n\n", hgTotalMem)
		}

		color.Green("Summary")
		color.White("==============================================")
		for k, v := range summary {
			fmt.Printf("%dMB\t\t%s\n", v, k)
		}

		color.White("==============================================")
		color.Green("Grand total: %dMB\n\n", totalMem)
	},
}

func init() {
	rootCmd.AddCommand(memoryUsageCmd)
}
