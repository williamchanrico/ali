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

package ess

import (
	"encoding/base64"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// SCInfo contains useful info about a scaling configuration
type SCInfo struct {
	UserData                 string
	ScalingConfigurationName string
}

// QuerySCInfo query relatively useful info about a scaling configuration
func (c *Client) QuerySCInfo(id string) ([]SCInfo, error) {
	req := ess.CreateDescribeScalingConfigurationsRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingConfigurationId1 = id

	var scalingConfigurationList []SCInfo

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeScalingConfigurations(req)
		if err != nil {
			return nil, err
		}

		for i := range resp.ScalingConfigurations.ScalingConfiguration {
			sc := resp.ScalingConfigurations.ScalingConfiguration[i]
			data, err := base64.StdEncoding.DecodeString(sc.UserData)
			if err != nil {
				return nil, err
			}
			userData := string(data[:])
			scalingConfigurationList = append(scalingConfigurationList,
				SCInfo{
					UserData:                 userData,
					ScalingConfigurationName: sc.ScalingConfigurationName,
				},
			)
		}
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.ScalingConfigurations.ScalingConfiguration))
	}

	return scalingConfigurationList, nil
}
