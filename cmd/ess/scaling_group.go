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
	"errors"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// SGInfo contains useful info about a scaling group
type SGInfo struct {
	ScalingGroupName         string
	ScalingGroupID           string
	TotalCapacity            int
	MinSize                  int
	MaxSize                  int
	ProtectedCapacity        int
	ScalingConfigurationName string
	UserData                 string
}

// String pretty prints the struct SGInfo
func (s *SGInfo) String() string {
	return fmt.Sprintf(`ScalingGroupName: %v
ScalingGroupID: %v
TotalCapacity: %v (Protected: %v)
MinSize: %v - MaxSize: %v
ScalingConfigurationName: %v
UserData:
>>> BEGIN - USERDATA
%v
<<< END - USERDATA
`, s.ScalingGroupName,
		s.ScalingGroupID,
		s.TotalCapacity,
		s.ProtectedCapacity,
		s.MinSize,
		s.MaxSize,
		s.ScalingConfigurationName,
		s.UserData)
}

// QuerySGInfo will query relatively useful info about a scaling group
func (c *Client) QuerySGInfo(name string) ([]SGInfo, error) {
	return c.queryScalingGroupsAndScalingConf(name)
}

func (c *Client) queryScalingGroupsAndScalingConf(name string) ([]SGInfo, error) {
	req := ess.CreateDescribeScalingGroupsRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupName = name

	var scalingGroupList []SGInfo

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeScalingGroups(req)
		if err != nil {
			return nil, err
		}

		for i := range resp.ScalingGroups.ScalingGroup {
			sg := resp.ScalingGroups.ScalingGroup[i]
			sc, err := c.QuerySCInfo(sg.ActiveScalingConfigurationId)
			if err != nil {
				return nil, err
			}

			scalingGroupList = append(scalingGroupList,
				SGInfo{
					ScalingGroupName:         sg.ScalingGroupName,
					ScalingGroupID:           sg.ScalingGroupId,
					TotalCapacity:            sg.TotalCapacity,
					MinSize:                  sg.MinSize,
					MaxSize:                  sg.MaxSize,
					ProtectedCapacity:        sg.ProtectedCapacity,
					UserData:                 sc[0].UserData,
					ScalingConfigurationName: sc[0].ScalingConfigurationName,
				},
			)
		}
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.ScalingGroups.ScalingGroup))
	}

	return scalingGroupList, nil
}

func (c *Client) queryScalingGroups(name string) ([]SGInfo, error) {
	req := ess.CreateDescribeScalingGroupsRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupName = name

	var scalingGroupList []SGInfo

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeScalingGroups(req)
		if err != nil {
			return nil, err
		}

		for i := range resp.ScalingGroups.ScalingGroup {
			sg := resp.ScalingGroups.ScalingGroup[i]
			scalingGroupList = append(scalingGroupList,
				SGInfo{
					ScalingGroupName:  sg.ScalingGroupName,
					ScalingGroupID:    sg.ScalingGroupId,
					TotalCapacity:     sg.TotalCapacity,
					MinSize:           sg.MinSize,
					MaxSize:           sg.MaxSize,
					ProtectedCapacity: sg.ProtectedCapacity,
				},
			)
		}
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.ScalingGroups.ScalingGroup))
	}

	return scalingGroupList, nil
}

func (c *Client) getScalingGroupID(scalingGroupName string) (string, error) {
	sgInfo, err := c.queryScalingGroupsAndScalingConf(scalingGroupName)
	if err != nil {
		return "", err
	}

	for _, sg := range sgInfo {
		if sg.ScalingGroupName == scalingGroupName {
			return sg.ScalingGroupID, nil
		}
	}

	return "", errors.New("ScalingGroup not found")
}

// ChangeAllMaxSize will change the max size of all scaling groups
func (c *Client) ChangeAllMaxSize(newMaxSize int) {
	sgList, err := c.queryScalingGroups("")
	if err != nil {
		fmt.Printf("Failed to query all scaling groups: %v", err)
	}

	for _, sg := range sgList {
		err = c.changeMaxSize(sg.ScalingGroupID, newMaxSize)
		if err != nil {
			fmt.Printf("Error on scaling group %v: %v", sg.ScalingConfigurationName, err)
		}

		fmt.Println("Changed max size of scaling group:", sg.ScalingGroupName)
	}
}

// ChangeMaxSize will change the max size of the specified scaling group
func (c *Client) ChangeMaxSize(scalingGroupID string, newMaxSize int) error {
	err := c.changeMaxSize(scalingGroupID, newMaxSize)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) changeMaxSize(scalingGroupID string, newMaxSize int) error {
	req := ess.CreateModifyScalingGroupRequest()
	req.ScalingGroupId = scalingGroupID
	req.MaxSize = requests.NewInteger(newMaxSize)

	_, err := c.ModifyScalingGroup(req)
	if err != nil {
		return err
	}

	return nil
}

// ChangeAllMinSize will change the min size of all scaling groups
func (c *Client) ChangeAllMinSize(newMinSize int) {
	sgList, err := c.queryScalingGroups("")
	if err != nil {
		fmt.Printf("Failed to query all scaling groups: %v", err)
	}

	for _, sg := range sgList {
		err = c.changeMinSize(sg.ScalingGroupID, newMinSize)
		if err != nil {
			fmt.Printf("Error on scaling group %v: %v", sg.ScalingConfigurationName, err)
		}

		fmt.Println("Changed min size of scaling group:", sg.ScalingGroupName)
	}
}

// ChangeMinSize will change the min size of the specified scaling group
func (c *Client) ChangeMinSize(scalingGroupID string, newMinSize int) error {
	err := c.changeMinSize(scalingGroupID, newMinSize)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) changeMinSize(scalingGroupID string, newMinSize int) error {
	req := ess.CreateModifyScalingGroupRequest()
	req.ScalingGroupId = scalingGroupID
	req.MinSize = requests.NewInteger(newMinSize)

	_, err := c.ModifyScalingGroup(req)
	if err != nil {
		return err
	}

	return nil
}
