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

package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// InstanceInfo contains information about an instance
type InstanceInfo struct {
	IP           string
	CPU          int
	Memory       int
	InstanceName string
	InstanceType string
	CreationTime string
}

// QueryInstanceList will query info of instances with matched hostgroup tag
// The other tags are static: "Environment=production" and "Datacenter=alisg"
func (c *Client) QueryInstanceList(hostgroup string) ([]InstanceInfo, error) {
	req := ecs.CreateDescribeInstancesRequest()
	req.PageSize = requests.NewInteger(100)
	req.PageNumber = requests.NewInteger(1)
	req.Tag = &[]ecs.DescribeInstancesTag{
		ecs.DescribeInstancesTag{
			Value: "production",
			Key:   "Environment",
		},
		ecs.DescribeInstancesTag{
			Value: "alisg",
			Key:   "Datacenter",
		},
		ecs.DescribeInstancesTag{
			Value: hostgroup,
			Key:   "Hostgroup",
		},
	}

	var instanceList []InstanceInfo

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeInstances(req)
		if err != nil {
			return nil, err
		}

		for i := range resp.Instances.Instance {
			instance := resp.Instances.Instance[i]

			instanceList = append(instanceList,
				InstanceInfo{
					// IP: instance.NetworkInterfaces.
					// 	NetworkInterface[0].PrimaryIpAddress,
					IP: instance.VpcAttributes.
						PrivateIpAddress.IpAddress[0],
					CPU:          instance.Cpu,
					Memory:       instance.Memory,
					InstanceName: instance.InstanceName,
					CreationTime: instance.CreationTime,
				},
			)
		}
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.Instances.Instance))
	}

	return instanceList, nil
}
