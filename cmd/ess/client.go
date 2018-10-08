package ess

import (
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
