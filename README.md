<h1 align="center">Ali... that's it.</h1>

<div align="center">
  :house_with_garden:
</div>
<div align="center">
  <strong>CLI tool to ease the interaction with Alibaba Cloud</strong>
</div>
<div align="center">
  A <code>small</code> tool to simplify everyday tasks.
</div>

<br />

<div align="center">
  <!-- GPL License -->
  <a href="https://opensource.org/licenses/GPL-3.0/"><img
	src="https://badges.frapsoft.com/os/gpl/gpl.png?v=103"
	border="0"
	alt="GPL Licence"
	title="GPL Licence">
  </a>
  <!-- Open Source Love -->
  <a href="https://opensource.org/licenses/GPL-3.0/"><img
	src="https://badges.frapsoft.com/os/v1/open-source.svg?v=103"
	border="0"
	alt="Open Source Love"
	title="Open Source Love">
  </a>
  <!-- Go Report Card -->
  <a href="https://goreportcard.com/report/github.com/williamchanrico/ali"><img
	src="https://goreportcard.com/badge/github.com/williamchanrico/ali"
	border="0"
	alt="Go Report Card"
	title="Go Report Card">
  </a>
</div>

<div align="center">
  <h3>
    <a href="https://arzhon.id">Website</a>
  </h3>
</div>

## Introduction

Your personal CLI tool to interact with Aliyun or Alibaba Cloud console.

[![asciicast](https://asciinema.org/a/gee4XkKWpvENAuBOaHbSMFIIN.png)](https://asciinema.org/a/gee4XkKWpvENAuBOaHbSMFIIN)

## Setup

```sh
> /home/william/.ali.yaml
---
ALICLOUD_ACCESS_KEY: <YOUR_ACCESS_KEY>
ALICLOUD_SECRET_KEY: <YOUR_SECRET_KEY>
ALICLOUD_REGION_ID: <REGION_ID>
```

It can also be read from Env variable.

Supports `dep ensure -v`.

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
