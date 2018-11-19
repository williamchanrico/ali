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
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/williamchanrico/ali/cmd/ess"
)

// maxSizeCmd represents the maxSize command
var maxSizeCmd = &cobra.Command{
	Use:   "maxSize NEW_MAX_SIZE [SCALING_GROUP_NAME]",
	Short: "Change max size of scaling group.",
	Long:  `Change max size of scaling group.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var newMaxSize int

		newMaxSize, err = strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}

		ess := ess.New()

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			panic(err)
		}

		color.Yellow("---------------------------------\n")
		if all {
			fmt.Println("Changing max size of all scaling group")
			ess.ChangeAllMaxSize(newMaxSize)
			if err != nil {
				fmt.Println("")
			}
		} else {
			if len(args) < 2 {
				fmt.Println("Please specify SCALING_GROUP_NAME")
				return
			}

			sgList, err := ess.QuerySGInfo(args[1])
			if err != nil {
				fmt.Println(fmt.Errorf("Failed to query scaling group list: %s", err))
			}

			sgIdx := -1
			for i := range sgList {
				if sgList[i].ScalingGroupName == args[1] {
					sgIdx = i
					break
				}
			}

			if sgIdx != -1 {
				sg := sgList[sgIdx]

				fmt.Print("ScalingGroupName: ")
				color.Green(sg.ScalingGroupName)

				fmt.Printf("Changing max size from %v to ", sg.MaxSize)
				color.Green("%v", newMaxSize)

				err = ess.ChangeMaxSize(sg.ScalingGroupID, newMaxSize)
				if err != nil {
					color.Red("Failed to change max size:", err)
				} else {
					color.Green("Successfully changed the max size")
				}
			} else {
				color.Red("Scaling group not found: %v", args[1])
			}
		}
		color.Yellow("\n---------------------------------\n")
	},
}

func init() {
	changeCmd.AddCommand(maxSizeCmd)

	maxSizeCmd.Flags().BoolP("all", "a", false, "Apply changes to all scaling groups")
}
