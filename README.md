# ali

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

Supports `dep ensure -v` to make life a bit easier.

## Usage

```sh
$ ali help
Using config file: /home/william/.ali.yaml
Personal CLI tool to interact with Aliyun or Alibaba Cloud console.

Usage:
  ali [command]

Available Commands:
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

```sh
$ go version
go version go1.11 linux/amd64
```
