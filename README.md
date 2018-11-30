# ali

[![Go Report Card](https://goreportcard.com/badge/github.com/williamchanrico/ali)](https://goreportcard.com/report/github.com/williamchanrico/ali)

Personal CLI tool to interact with Aliyun or Alibaba Cloud console.

[![asciicast](https://asciinema.org/a/gee4XkKWpvENAuBOaHbSMFIIN.png)](https://asciinema.org/a/gee4XkKWpvENAuBOaHbSMFIIN)

## Setup

```sh
> /home/william/.ali.yaml
---
ALICLOUD_ACCESS_KEY: <YOUR_ACCESS_KEY>
ALICLOUD_SECRET_KEY: <YOUR_SECRET_KEY>
ALICLOUD_REGION_ID: <REGION_ID>
```

Uses `viper.AutomaticEnv()`, means it can also read those config from Env variable.

Supports `dep ensure -v` to make our life a bit easier.

## Usage

```txt
$ ali help
Personal CLI tool to interact with Aliyun or Alibaba Cloud console.

Usage:
  ali [command]

Available Commands:
  change      Change attributes of cloud objects
  downscale   Remove all upscaled instance down to minimum instance.
  et          Query Event-Trigger Task(s) from aliyun.
  help        Help about any command
  ip          Query active IP(s) of a service hostgroup
  memoryUsage Get and calculate the memory usage of instance(s) by hostgroups
  price       Show real-time price per hour for the instance type in USD (default region: ap-southeast-1).
  sg          Query active ScalingGroup of a service by name
  upscale     Upscale a scaling group to add specified number of instances.

Flags:
      --config string   config file (default is $HOME/.ali.yaml)
  -h, --help            help for ali

Use "ali [command] --help" for more information about a command.
```

## Go Version

```txt
$ go version
go version go1.11 linux/amd64
```
