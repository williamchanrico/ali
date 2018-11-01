# ali

[![Go Report Card](https://goreportcard.com/badge/github.com/williamchanrico/ali)](https://goreportcard.com/report/github.com/williamchanrico/ali)

Personal CLI tool to interact with Aliyun or Alibaba Cloud console.

This tool is designed with my workplace's cloud environment in mind.

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
Using config file: /home/william/.ali.yaml
Personal CLI tool to interact with Aliyun or Alibaba Cloud console.

Usage:
  ali [command]

Available Commands:
  downscale   Remove all upscaled instance down to minimum instance.
  et          Query Event-Trigger Task(s) from aliyun.
  help        Help about any command
  ip          Query active IP(s) of a service hostgroup
  sg          Query active ScalingGroup of a service by name

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
