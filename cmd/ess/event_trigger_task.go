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
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// ETInfo contains relevant info about an Event-Trigger Task
type ETInfo struct {
	AlarmTaskID        string
	Name               string
	MetricName         string
	EvaluationCount    int
	Threshold          float64
	ComparisonOperator string
	Enable             bool
}

// String pretty prints the struct ETInfo
func (e *ETInfo) String() string {
	return fmt.Sprintf(`> Alarm Name: %v (Enabled: %v)
> %v %v %v, %v times
`, e.Name,
		e.Enable,
		e.MetricName,
		e.ComparisonOperator,
		e.Threshold,
		e.EvaluationCount)
}

// QueryETInfo queries useful info about an event-trigger task
func (c *Client) QueryETInfo(scalingGroupID string) ([]ETInfo, error) {
	req := ess.CreateDescribeAlarmsRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	var eventTriggerTaskList []ETInfo

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeAlarms(req)
		if err != nil {
			return nil, err
		}

		for i := range resp.AlarmList.Alarm {
			et := resp.AlarmList.Alarm[i]
			eventTriggerTaskList = append(eventTriggerTaskList,
				ETInfo{
					AlarmTaskID:        et.AlarmTaskId,
					Name:               et.Name,
					MetricName:         et.MetricName,
					EvaluationCount:    et.EvaluationCount,
					Threshold:          et.Threshold,
					ComparisonOperator: et.ComparisonOperator,
					Enable:             et.Enable,
				},
			)
		}

		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.AlarmList.Alarm))
	}

	return eventTriggerTaskList, nil
}

// EnableEventTriggerTask will enable the event trigger task
func (c *Client) EnableEventTriggerTask(alarmTaskID string) error {
	req := ess.CreateEnableAlarmRequest()
	req.AlarmTaskId = alarmTaskID

	_, err := c.EnableAlarm(req)
	if err != nil {
		return err
	}

	return nil
}

// DisableEventTriggerTask will disable the event trigger task
func (c *Client) DisableEventTriggerTask(alarmTaskID string) error {
	req := ess.CreateDisableAlarmRequest()
	req.AlarmTaskId = alarmTaskID

	_, err := c.DisableAlarm(req)
	if err != nil {
		return err
	}

	return nil
}
