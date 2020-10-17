# go-probe-wait

go-probe-wait is a utility that probes and blocks an HTTP endpoint until it returns a 200 response code.
It's inspired by the famous wait-for-it.sh script, but works on HTTP endpoints.
The meaning of the options are slightly different, however, more hopefully a bit more explanatory.

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