package ecs

import (
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
