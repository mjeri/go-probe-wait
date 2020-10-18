# go-probe-wait

![Go](https://github.com/mjeri/go-probe-wait/workflows/Go/badge.svg)

## Description

go-probe-wait is a utility which probes an HTTP endpoint until it returns a 2xx response code, or a timeout is reached and optionally starts a program.
It's inspired by the famous [wait-for-it.sh](https://github.com/vishnubob/wait-for-it) script, but works on HTTP endpoints.
The meaning of the options are slightly different, however, hopefully a bit more explanatory.

Note: if you specify to run a program on success the implementation uses `syscall.Exec` to replace the go-probe-wait process with the program you specified. 

## Usage

```
> ./go-wait-probe --help
Usage: go-wait-probe [OPTION]... [CMD]...

Examples:

  go-wait-probe --endpoint http://localhost:8080/ready
  go-wait-probe --endpoint http://localhost:8080/ready echo 'ready to run anything :)'
  go-wait-probe --endpoint http://localhost:8080/ready --programTimeout 2s --runCommandOnTimeout echo 'ready to run anything :)'
  go-wait-probe -e http://localhost:8080/ready -i 2s -t 10s -c echo 'ready to run anything :)'

Options:

  -e, --endpoint string         REQUIRED - The endpoint to probe.
  -h, --help                    OPTIONAL - Show this online help.
  -i, --probeInterval string    OPTIONAL - The interval at which the probe is executed. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s (default "1s")
  -t, --programTimeout string   OPTIONAL - Timeout after the program is considered unsuccessful and it exits with 1. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s (default "15s")
  -c, --runCommandOnTimeout     OPTIONAL - Run the specified command also on a programTimeout.
```

Precompiled binaries for unix based systems are available in the releases of the GitHub project.

If you want to mess it with on your own, clone the project and a simple `go build` does the job.

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
