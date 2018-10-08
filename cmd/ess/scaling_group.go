package ess

import (
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
		totalCount = requests.NewInteger(resp.TotalCount)
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
	}

	return scalingGroupList, nil
}
