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
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// IPInfo contains IP along with its associated consul_tags if available
type IPInfo struct {
	IP        string
	ConsulTag string
	IsRunning bool
}

// QueryIPList will query IP of instances with matched hostgroup tag
// The other tags are static: "Environment=production" and "Datacenter=alisg"
func (c *Client) QueryIPList(hostgroup string) ([]IPInfo, error) {
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

	var ipList []IPInfo

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeInstances(req)
		if err != nil {
			return nil, err
		}

		for i := range resp.Instances.Instance {
			instance := resp.Instances.Instance[i]
			consulTag := ""
			for tagIdx := range instance.Tags.Tag {
				if instance.Tags.Tag[tagIdx].TagKey == "consul_tags" {
					consulTag = instance.Tags.Tag[tagIdx].TagValue
				}
			}

			isRunning := false
			if instance.Status == "Running" {
				isRunning = true
			}

			ipList = append(ipList,
				IPInfo{
					IP: instance.NetworkInterfaces.
						NetworkInterface[0].PrimaryIpAddress,
					ConsulTag: consulTag,
					IsRunning: isRunning,
				},
			)
		}
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.Instances.Instance))
	}

	return ipList, nil
}

// QueryAllHostgroupIP will query all IP of existing running instances
// Returns map["tkp-hostgroup"] = []string{"IP_1", "IP_2", ...}
func (c *Client) QueryAllHostgroupIP() (map[string][]string, error) {
	req := ecs.CreateDescribeInstancesRequest()
	req.PageSize = requests.NewInteger(100)
	req.PageNumber = requests.NewInteger(1)

	hostgroupList := make(map[string][]string)

	processedInstance := 0
	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeInstances(req)
		if err != nil {
			return nil, err
		}

		processedInstance += len(resp.Instances.Instance)
		fmt.Printf("Processing %-4v/%-4v instances [%-4v hostgroup found so far]\n",
			processedInstance, resp.TotalCount, len(hostgroupList))

		for i := range resp.Instances.Instance {
			instance := resp.Instances.Instance[i]
			if instance.Status != "Running" {
				continue
			}

			hostgroup := ""
			for tagIdx := range instance.Tags.Tag {
				if instance.Tags.Tag[tagIdx].TagKey == "Hostgroup" {
					hostgroup = instance.Tags.Tag[tagIdx].TagValue
				}
			}
			if hostgroup == "" {
				continue
			}

			ipAddress := instance.NetworkInterfaces.NetworkInterface[0].PrimaryIpAddress
			hostgroupList[hostgroup] = append(hostgroupList[hostgroup], ipAddress)
		}

		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.Instances.Instance))
	}

	return hostgroupList, nil
}
