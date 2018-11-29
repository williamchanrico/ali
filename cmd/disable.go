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
	"github.com/williamchanrico/ali/cmd/ess"
)

// disableCmd represents the enable command
var disableCmd = &cobra.Command{
	Use:   "disable SCALING_GROUP_NAME",
	Short: "Disables upscale/downscale alarm task",
	Long:  `Will disable upscale or downscale event-trigger task`,
	Run: func(cmd *cobra.Command, args []string) {
		target := ""
		fmt.Print("Disabling Event-Trigger Task ")
		if len(args) > 1 {
			target = args[0]
			color.Green("%v\n", target)
		} else {
			color.Green("for all scaling group(s)\n")
		}

		ess := ess.New()
		sgList, err := ess.QuerySGInfo(target)
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to query scaling group list: %s", err))
		}

		for i := range sgList {
			// Only processes exact scaling group name
			if !all && sgList[i].ScalingGroupName != target {
				continue
			}

			etList, err := ess.QueryETInfo(sgList[i].ScalingGroupID)
			if err != nil {
				fmt.Println(fmt.Errorf("Failed to query event-trigger task list: %s", err))
			}

			color.Yellow("\n--------------- %v ---------------\n", i)
			for j := range etList {
				if (upscale || both) && strings.Contains(etList[j].Name, "upscale") {
					disableAlarm(ess, etList[j])
				} else if (downscale || both) && strings.Contains(etList[j].Name, "downscale") {
					disableAlarm(ess, etList[j])
				}

			}

			color.Yellow("--- https://essnew.console.aliyun.com/?spm=5176.2020520101.203.4." +
				"278f7d33hepSMf#/task/alarm/region/ap-southeast-1 ---\n")
		}
	},
}

func disableAlarm(c *ess.Client, et ess.ETInfo) {
	fmt.Print("> Alarm Name: ")
	if et.Enable == false {
		color.Red("%v (Already Disabled)", et.Name)
	} else {
		err := c.DisableEventTriggerTask(et.AlarmTaskID)
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to disable event-trigger task: %s", err))
		} else {
			color.Green("%v (Successfully Disabled)", et.Name)
		}
	}

	fmt.Printf("> %v %v ", et.MetricName, et.ComparisonOperator)
	color.White("%v (%v times)", et.Threshold, et.EvaluationCount)
	fmt.Println()
}

func init() {
	etCmd.AddCommand(disableCmd)

	disableCmd.Flags().BoolVarP(&upscale, "upscale", "u", false, "Disable upscale event-trigger task")
	disableCmd.Flags().BoolVarP(&downscale, "downscale", "d", false, "Disable downscale event-trigger task")
	disableCmd.Flags().BoolVarP(&both, "both", "b", false, "Disable upscale and downscale event-trigger task")
	disableCmd.Flags().BoolVarP(&all, "all", "a", false, "Apply to all scaling groups")
}
