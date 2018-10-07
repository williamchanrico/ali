package ess

import (
	"encoding/base64"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/spf13/viper"
)

// Client struct
type Client struct {
	*ess.Client
}

// New creates a new ECS client
func New() *Client {
	essClient, err := ess.NewClientWithAccessKey(
		viper.GetString("ALICLOUD_REGION_ID"),
		viper.GetString("ALICLOUD_ACCESS_KEY"),
		viper.GetString("ALICLOUD_SECRET_KEY"),
	)
	if err != nil {
		panic(err)
	}

	return &Client{
		essClient,
	}
}

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
TotalCapacity: %v
MinSize: %v
MaxSize: %v
ProtectedCapacity: %v
ScalingConfigurationName: %v
UserData:
%v
`, s.ScalingGroupName,
		s.ScalingGroupID,
		s.TotalCapacity,
		s.MinSize,
		s.MaxSize,
		s.ProtectedCapacity,
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
		totalCount = requests.NewInteger(resp.TotalCount)
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
	}

	return scalingConfigurationList, nil
}
