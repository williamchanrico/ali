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

	"github.com/spf13/cobra"
	"github.com/williamchanrico/ali/cmd/ess"
)

// etCmd represents the et command
var etCmd = &cobra.Command{
	Use:   "et [SCALING_GROUP_NAME]",
	Short: "Query Event-Trigger Task(s) from aliyun.",
	Long: `Query Event-Trigger Task(s) from aliyun. Requires SCALING_GROUP_NAME
because Event-Trigger Task can't be searched by name, and scaling group name is a better
compromise than using an Event-Trigger Task ID`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Querying Event-Trigger Task %v:\n", args[0])

		ess := ess.New()
		sgList, err := ess.QuerySGInfo(args[0])
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to query scaling group list: %s", err))
		}

		for i := range sgList {
			etList, err := ess.QueryETInfo(sgList[i].ScalingGroupID)
			if err != nil {
				fmt.Println(fmt.Errorf("Failed to query event-trigger task list: %s", err))
			}

			fmt.Printf("\n--------------- %v ---------------\n", i)
			for j := range etList {
				fmt.Println(etList[j].String())
			}
			fmt.Println(sgList[i].String())
			fmt.Printf("--- https://essnew.console.aliyun.com/?spm=5176.2020520101.203.4." +
				"278f7d33hepSMf#/task/alarm/region/ap-southeast-1 ---\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(etCmd)
}
