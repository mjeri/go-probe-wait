# go-wait-for-it

![Go](https://github.com/mjeri/go-wait-for-it/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/mjeri/go-wait-for-it)](https://goreportcard.com/report/mjeri/go-wait-for-it)

## Description

go-wait-for-it is a utility which probes an HTTP endpoint until it returns a 2xx response code, or a timeout is reached and optionally starts a program.
It's inspired by the famous [wait-for-it.sh](https://github.com/vishnubob/wait-for-it) script, but works on HTTP endpoints.
The meaning of the options are slightly different, however, hopefully a bit more explanatory.

Note: if you specify to run a program on success the implementation uses `syscall.Exec` to replace the go-wait-for-it process with the program you specified. 

## Installation

Binaries built from the latest main:
 
- [Linux](https://github.com/mjeri/go-wait-for-it/blob/main/bin/amd64/linux/go-wait-for-it)
- [Darwin](https://github.com/mjeri/go-wait-for-it/blob/main/bin/amd64/darwin/go-wait-for-it)

## Usage

```
Usage: go-wait-for-it [OPTION]... [CMD]...

Examples:

  go-wait-for-it --endpoint http://localhost:8080/ready
  go-wait-for-it --endpoint http://localhost:8080/ready echo 'ready to run anything :)'
  go-wait-for-it --endpoint http://localhost:8080/ready --programTimeout 2s --runCommandOnTimeout echo 'ready to run anything :)'
  go-wait-for-it -e http://localhost:8080/ready -i 2s -t 10s -s echo 'ready to run anything :)'

Options:

  -e, --endpoint string         REQUIRED - The endpoint to probe.
  -h, --help                    OPTIONAL - Show this online help.
  -i, --probeInterval string    OPTIONAL - The interval at which the probe is executed. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s (default "1s")
  -t, --programTimeout string   OPTIONAL - Timeout after the program is considered unsuccessful and it exits with 1. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s (default "15s")
  -s, --succeedAnyway           OPTIONAL - Even when the timeout occurs, consider the run as success and exit with 0 or run the specified command.
```

If you want to mess it with on your own, clone the project and a simple `go build` does the job.

## Use cases

### docker-compose

Given you have a service which needs to wait with its startup until another service becomes ready, `go-wait-probe` is perfect for that.
Consider the example where `service_two` requires `service_one` to be ready before it should startup.
Note: as `go-wait-probe` uses the `syscall.Exec` the `yarn start` still runs as PID 1 in the second docker-compose service.

```yml
version: 3
services:
  service-one:
    image: 'node:lts'
    container_name: 'service-one'
    ports:
      - '3000:3000'
    working_dir: /app
    volumes:
      - '.:/app'
    command: [ "yarn", "start" ]
  service-two:
    image: 'node:lts'
    container_name: 'service-one'
    ports:
      - '3010:3010'
    working_dir: /app
    volumes:
      - '.:/app'
    command: [ "./go-wait-probe", "--endpoint", "http://service-one:3000/ready", "yarn", "start" ]
```

## Further improvement ideas

- add unit tests
- support for TCP and UDP probes
- support for probe timeouts
- provide docker images
- make HTTP status codes that are considered successes configurable


## Advantages over bash scripts

While it may seem too much to use a binary for such an easy task, it also can help, as the go binary has fewer requirements on the actual environment it's run in.
For example no `curl` or `wget` or similar programs are required to issue HTTP requests and thus the solution is actually more light-weight.
But you can do as you prefer of course :)
