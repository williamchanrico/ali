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

// upscaleCmd represents the upscale command
var upscaleCmd = &cobra.Command{
	Use:   "upscale SCALING_GROUP_NAME NUM_OF_INSTANCES_TO_ADD",
	Short: "Upscale a scaling group to add specified number of instances.",
	Long: `Upscale a scaling group to add specified number of instances.
Will temporarily change the upscale scaling rule to match NUM_OF_INSTANCES_TO_ADD number,
execute it, and revert the number.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Upscale instance(s) in scaling group ")
		color.Green("%v\n", args[0])

		ess := ess.New()

		numToAdd, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		color.White("Adding %v instance(s)\n", numToAdd)
		ok, err := ess.UpscaleInstances(args[0], numToAdd)
		if err != nil {
			fmt.Println(err)
		}

		if ok {
			color.Green("Successfully upscaled %v instance(s)", numToAdd)
		} else {
			color.Red("Failed to upscale instance(s)")
		}
	},
}

func init() {
	rootCmd.AddCommand(upscaleCmd)
}
