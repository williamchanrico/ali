# ali

Personal CLI tool to interact with Aliyun or Alibaba Cloud console.

## Setup

```
> /home/william/.ali.yaml
---
ALICLOUD_ACCESS_KEY: <YOUR_ACCESS_KEY>
ALICLOUD_SECRET_KEY: <YOUR_SECRET_KEY>
ALICLOUD_REGION_ID: ap-southeast-1
```

Uses `viper.AutomaticEnv()`, means it can also read those config from Env variable.

Supports `dep ensure -v` to make life a bit easier.

## Usage

```
$ ali -h
Personal CLI tool to interact with Aliyun or Alibaba Cloud console.

Usage:
  ali [command]

Available Commands:
  help        Help about any command
  ip          Query active IP(s) of a service hostgroup
  sg          Query active ScalingGroup of a service by name

Flags:
      --config string   config file (default is $HOME/.ali.yaml)
  -h, --help            help for ali

Use "ali [command] --help" for more information about a command.
```

## Go Version

```
$ go version
go version go1.11 linux/amd64
```
