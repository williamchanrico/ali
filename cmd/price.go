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

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/williamchanrico/ali/cmd/ecs"
)

// priceCmd represents the price command
var priceCmd = &cobra.Command{
	Use:   "price INSTANCE_TYPE [REGION_ID]",
	Short: "Show real-time price per hour for the instance type in USD (default region: ap-southeast-1).",
	Long: `Will query real-time price per hour in USD and output the specified instance type
pricing (default region is Singapore ap-southeast-1)`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			args = append(args, "ap-southeast-1")
		}

		ecs := ecs.New()
		priceList, err := ecs.QueryPriceList()
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to query price list: %s", err))
		}

		fmt.Print("Price for instance type: ")
		color.Green(args[0])

		fmt.Print("Region ID: ")
		color.Yellow(args[1])

		for _, price := range priceList {
			if price.RegionID == args[1] && price.InstanceType == args[0] {
				fmt.Printf("$%v/hour (%v)\n", price.PricePerHour, price.OSType)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(priceCmd)
}
