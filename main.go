package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
)

func main() {
	var help, succeedAnyway bool
	var endpoint, probeIntervalString, programTimeoutString string

	flag.BoolVarP(&help, "help", "h", false, "OPTIONAL - Show this online help.")

	flag.StringVarP(&endpoint, "endpoint", "e", "", "REQUIRED - The endpoint to probe.")
	flag.StringVarP(&probeIntervalString, "probeInterval", "i", "1s", "OPTIONAL - The interval at which the probe is executed. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s")
	flag.StringVarP(&programTimeoutString, "programTimeout", "t", "15s", "OPTIONAL - Timeout after the program is considered unsuccessful and it exits with 1. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s")
	flag.BoolVarP(&succeedAnyway, "succeedAnyway", "s", false, "OPTIONAL - Even when the timeout occurs, consider the run as success and exit with 0 or run the specified command.")

	flag.Parse()

	if help {
		printHelp("")
		os.Exit(1)
	}

	argsForCommandToExecute := flag.Args()

	if endpoint == "" {
		printHelp("Error: missing --endpoint configuration")
		os.Exit(1)
	}

	probeInterval, err := time.ParseDuration(probeIntervalString)
	if err != nil {
		printHelp(fmt.Sprintf("Error: invalid format specificed for probeInterval: %s", err.Error()))
		os.Exit(1)
	}

	programTimeout, err := time.ParseDuration(programTimeoutString)
	if err != nil {
		printHelp(fmt.Sprintf("Error: invalid format specificed for programTimeout: %s", err.Error()))
		os.Exit(1)
	}

	success := make(chan bool)

	go runProbeLoop(endpoint, probeInterval, success)

	select {
	case <-success:
		if len(argsForCommandToExecute) > 0 {
			execSuccessCommand(argsForCommandToExecute)
		} else {
			os.Exit(0)
		}
	case <-time.After(programTimeout):
		if len(argsForCommandToExecute) > 0 && succeedAnyway {
			execSuccessCommand(argsForCommandToExecute)
		} else if succeedAnyway {
			os.Exit(0)
		} else {
			fmt.Printf("Error: programTimeout of %s reached. Exiting with 1\n", programTimeout.String())
			os.Exit(1)
		}
	}
}

func printHelp(errorMessage string) {
	if len(errorMessage) > 0 {
		fmt.Print(errorMessage, "\n\n")
	}

	fmt.Print("Usage: go-wait-probe [OPTION]... [CMD]...\n\n")
	fmt.Print("Examples:\n\n")
	fmt.Print("  go-wait-probe --endpoint http://localhost:8080/ready\n")
	fmt.Print("  go-wait-probe --endpoint http://localhost:8080/ready echo 'ready to run anything :)'\n")
	fmt.Print("  go-wait-probe --endpoint http://localhost:8080/ready --programTimeout 2s --runCommandOnTimeout echo 'ready to run anything :)'\n")
	fmt.Print("  go-wait-probe -e http://localhost:8080/ready -i 2s -t 10s -c echo 'ready to run anything :)'\n\n")

	fmt.Print("Options:\n\n")
	flag.PrintDefaults()
}

func runProbeLoop(endpoint string, probeInterval time.Duration, success chan<- bool) {
	ticker := time.NewTicker(probeInterval)
	defer ticker.Stop()

	// Start the first probe immediately.
	// The ticker ticks the first time after the respective time.Duration, but we want to start probing immediately.
	go runProbe(endpoint, success)

	for range ticker.C {
		go runProbe(endpoint, success)
	}
}

func runProbe(endpoint string, success chan<- bool) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return
	}

	defer func() {
		// ignore the error
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		success <- true
	}
}

func execSuccessCommand(args []string) {
	cmdName := args[0]
	resolvedCmdName, err := exec.LookPath(cmdName)
	if err != nil {
		fmt.Printf("Error: resolving command name failed")
		os.Exit(1)
	}

	if err := syscall.Exec(resolvedCmdName, args, os.Environ()); err != nil {
		fmt.Printf("Error: exec of %s failed: %s", resolvedCmdName, err.Error())
		os.Exit(1)
	}
}
