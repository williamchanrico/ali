package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// IPInfo contains IP along with its associated consul_tags if available
type IPInfo struct {
	IP        string
	ConsulTag string
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

			ipList = append(ipList,
				IPInfo{
					IP: instance.NetworkInterfaces.
						NetworkInterface[0].PrimaryIpAddress,
					ConsulTag: consulTag,
				},
			)
		}
		totalCount = requests.NewInteger(resp.TotalCount)
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
	}

	return ipList, nil
}
