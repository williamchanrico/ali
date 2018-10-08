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

		totalCount = requests.NewInteger(resp.TotalCount)
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
	}

	return eventTriggerTaskList, nil
}
