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
		totalCount = requests.NewInteger(resp.TotalCount)
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
	}

	return scalingConfigurationList, nil
}
