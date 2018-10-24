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

// etCmd represents the et command
var downscaleCmd = &cobra.Command{
	Use:   "downscale [SCALING_GROUP_NAME] [RETRY_COUNT] [RETRY_INTERVAL]",
	Short: "Remove all upscaled instance down to minimum instance.",
	Long:  `Will remove all upscaled instance down to minimum instance`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Removing upscaled instances in scaling group ")
		color.Green("%v\n", args[0])

		ess := ess.New()

		retryCount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
		}

		retryInterval, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
		}

		ok, err := ess.RemoveUpscaledInstances(args[0], retryCount, retryInterval)
		if err != nil {
			fmt.Println(err)
		}

		if ok {
			color.Green("Successfully removed upscaled instances")
		} else {
			color.Red("Failed to remove upscaled instances")
		}
	},
}

func init() {
	rootCmd.AddCommand(downscaleCmd)
}
