# go-probe-wait

![Go](https://github.com/mjeri/go-probe-wait/workflows/Go/badge.svg)

## WIP

The project is still work in progress and needs a bit more work to be actually useful.

## Description

go-probe-wait is a utility which probes an HTTP endpoint until it returns a 2xx response code or a timeout is reached.
It's inspired by the famous [wait-for-it.sh](https://github.com/vishnubob/wait-for-it) script, but works on HTTP endpoints.
The meaning of the options are slightly different, however, hopefully a bit more explanatory.

## Installation

Precompiled binaries for linux, macos and windows are available as releases.

If you want to mess it with on your own, the `build-releases.sh` script is your friend to build the program on your own.
Only standard go going on here :)

## Usage

```
> ./go-wait-probe --help

-endpoint string
      The endpoint to probe.
-probeInterval string
      The interval at which the probe is executed. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s (default "1s")
-programTimeout string
      Timeout after the program is considered unsuccessful and the tool exits with 1. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s (default "15s")
```

### Examples

```
> ./go-wait-probe --endpoint https://my-service.com/ready --probeInterval 300ms --programTimeout 5s
```

## Advantages over bash scripts

While it may seem too much to use a binary for such an easy task, it also can help, as the go binary doesn't have any requirements on the actual environment it's run in.
In other words no `curl` or `wget` or similar programs are required to issue HTTP requests and thus the solution is actually more light-weight.
But you can do as you prefer of course :)
