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
