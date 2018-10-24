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
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// SRInfo contains scaling rule info
type SRInfo struct {
	ScalingRuleID   string
	ScalingRuleAri  string
	ScalingRuleName string
	AdjustmentValue int
	AdjustmentType  string
}

// ModifyDownScalingRule will modify downscale scaling rule (scaling rule with "-downscale" suffix)
func (c *Client) ModifyDownScalingRule(scalingGroupName string, adjValue int) error {
	if adjValue > 0 {
		return errors.New("Downscale rule must have a negative adjustment value")
	}

	downScalingRule, err := c.getDownScalingRule(scalingGroupName)
	if err != nil {
		return err
	}

	if err = c.modifyScalingRuleAdjValue(downScalingRule.ScalingRuleID, adjValue); err != nil {
		return err
	}

	return nil
}

// ModifyUpScalingRule will modify upscale scaling rule (scaling rule with "-upscale" suffix)
func (c *Client) ModifyUpScalingRule(scalingGroupName string, adjValue int) error {
	if adjValue < 0 {
		return errors.New("Upscale rule must have a positive adjustment value")
	}

	upScalingRule, err := c.getUpScalingRule(scalingGroupName)
	if err != nil {
		return err
	}

	if err = c.modifyScalingRuleAdjValue(upScalingRule.ScalingRuleID, adjValue); err != nil {
		return err
	}

	return nil
}

// RemoveUpscaledInstances will remove currently upscaled instances down to minimum instances number.
// Removal will happen by executing downscale rule, so instances number will not go below
// minimum instances number. But we'll remove 80% instead of 100% of total capacity, in case
// the minimum instance is zero
func (c *Client) RemoveUpscaledInstances(scalingGroupName string, retryCount int, retryInterval int) (bool, error) {
	ok := false
	downscaleRule, err := c.getDownScalingRule(scalingGroupName)
	if err != nil {
		return ok, err
	}

	fmt.Println("Temporarily modifying downscale rule to -80")
	err = c.modifyScalingRuleAdjValue(downscaleRule.ScalingRuleID, -80)
	if err != nil {
		return ok, err
	}

	fmt.Printf("Executing downscale rule (will retry %v times per %vs)\n", retryCount, retryInterval)
	err = c.executeScalingRule(downscaleRule.ScalingRuleAri, retryCount, retryInterval)
	if err != nil {
		fmt.Println(err)
	} else {
		ok = true
	}

	fmt.Println("Reverting downscale rule to", downscaleRule.AdjustmentValue)
	err = c.modifyScalingRuleAdjValue(downscaleRule.ScalingRuleID, downscaleRule.AdjustmentValue)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

// executeScalingRule will execute specified scaling rule, will retry every 'retryInterval' for 'retryCount' times
// 'retryInterval' is in seconds
func (c *Client) executeScalingRule(scalingRuleAri string, retryCount int, retryInterval int) error {
	req := ess.CreateExecuteScalingRuleRequest()
	req.ScalingRuleAri = scalingRuleAri

	var err error

	for i := 1; i <= retryCount; i++ {
		if _, err = c.ExecuteScalingRule(req); err == nil || i == retryCount {
			return err
		}

		fmt.Printf("Failed to execute, retrying in %vs\n", retryInterval)
		time.Sleep(time.Duration(retryInterval) * time.Second)
	}

	return err
}

// getUpScalingRule will return scaling rule with "-upscale" suffix
func (c *Client) getUpScalingRule(scalingGroupName string) (*SRInfo, error) {
	scalingRuleList, err := c.describeScalingRules(scalingGroupName)
	if err != nil {
		return nil, err
	}

	for i := range scalingRuleList {
		if strings.HasSuffix(scalingRuleList[i].ScalingRuleName, "-upscale") {
			return &scalingRuleList[i], nil
		}
	}

	return nil, errors.New("Upscale rule not found")
}

// getDownScalingRule will return scaling rule with "-downscale" suffix
func (c *Client) getDownScalingRule(scalingGroupName string) (*SRInfo, error) {
	scalingRuleList, err := c.describeScalingRules(scalingGroupName)
	if err != nil {
		return nil, err
	}

	for i := range scalingRuleList {
		if strings.HasSuffix(scalingRuleList[i].ScalingRuleName, "-downscale") {
			return &scalingRuleList[i], nil
		}
	}

	return nil, errors.New("Downscale rule not found")
}

// modifyScalingRuleAdjValue will modify adjustment value of the scaling rule
func (c *Client) modifyScalingRuleAdjValue(scalingRuleID string, adjValue int) error {
	req := ess.CreateModifyScalingRuleRequest()
	req.AdjustmentValue = requests.NewInteger(adjValue)
	req.ScalingRuleId = scalingRuleID

	if _, err := c.ModifyScalingRule(req); err != nil {
		return err
	}

	return nil
}

// describeScalingRules will return scaling rules of a scaling group
func (c *Client) describeScalingRules(scalingGroupName string) ([]SRInfo, error) {
	scalingGroupID, err := c.getScalingGroupID(scalingGroupName)
	if err != nil {
		return nil, err
	}

	req := ess.CreateDescribeScalingRulesRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	var scalingRuleList []SRInfo

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeScalingRules(req)
		if err != nil {
			return nil, err
		}

		for i := range resp.ScalingRules.ScalingRule {
			sr := resp.ScalingRules.ScalingRule[i]
			scalingRuleList = append(scalingRuleList,
				SRInfo{
					ScalingRuleName: sr.ScalingRuleName,
					ScalingRuleAri:  sr.ScalingRuleAri,
					ScalingRuleID:   sr.ScalingRuleId,
					AdjustmentValue: sr.AdjustmentValue,
					AdjustmentType:  sr.AdjustmentType,
				},
			)
		}
		totalCount = requests.NewInteger(resp.TotalCount)
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
	}

	return scalingRuleList, nil
}
