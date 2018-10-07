package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/spf13/viper"
)

// Client struct
type Client struct {
	*ecs.Client
}

// New creates a new ECS client
func New() *Client {
	ecsClient, err := ecs.NewClientWithAccessKey(
		viper.GetString("ALICLOUD_REGION_ID"),
		viper.GetString("ALICLOUD_ACCESS_KEY"),
		viper.GetString("ALICLOUD_SECRET_KEY"),
	)
	if err != nil {
		panic(err)
	}

	return &Client{
		ecsClient,
	}
}

// QueryIPList will query IP of instances with matched hostgroup tag
func (c *Client) QueryIPList(hostgroup string) ([]string, error) {
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

	var ipList []string

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeInstances(req)
		if err != nil {
			return nil, err
		}

		for i := range resp.Instances.Instance {
			ipList = append(ipList, resp.Instances.Instance[i].NetworkInterfaces.
				NetworkInterface[0].PrimaryIpAddress)
		}
		totalCount = requests.NewInteger(resp.TotalCount)
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
	}

	return ipList, nil
}
