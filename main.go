package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	var help, runCommandOnTimeout bool
	var endpoint, probeIntervalString, programTimeoutString string

	flag.BoolVar(&help, "help", false, "OPTIONAL - Show this online help.")

	flag.StringVar(&endpoint, "endpoint", "", "REQUIRED - The endpoint to probe.")
	flag.StringVar(&endpoint, "e", "", "(shorthand for -endpoint)")

	flag.StringVar(&probeIntervalString, "probeInterval", "1s", "OPTIONAL - The interval at which the probe is executed. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s")
	flag.StringVar(&probeIntervalString, "i", "1s", "(shorthand for -probeInterval)")

	flag.StringVar(&programTimeoutString, "programTimeout", "15s", "OPTIONAL - Timeout after the program is considered unsuccessful and the tool exits with 1. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s")
	flag.StringVar(&programTimeoutString, "t", "15s", "(shorthand for -programTimeout)")

	flag.BoolVar(&runCommandOnTimeout, "runCommandOnTimeout", false, "OPTIONAL - Run the specified command also on a programTimeout.")
	flag.BoolVar(&runCommandOnTimeout, "c", false, "(shorthand for -runCommandOnTimeout)")

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
		if len(argsForCommandToExecute) > 0 && runCommandOnTimeout {
			execSuccessCommand(argsForCommandToExecute)
		} else {
			fmt.Printf("Error: programTimeout of %s reached. Exiting with 1\n", programTimeout.String())
			os.Exit(1)
		}
	}
}

func printHelp(errorMessage string) {
	if len(errorMessage) > 0 {
		fmt.Print("\n", errorMessage, "\n")
	}

	fmt.Print("\nUsage:   go-wait-probe [OPTION]... [CMD]...\n\n")
	fmt.Print("Examples:\n\n")
	fmt.Print("  go-wait-probe --endpoint http://localhost:8080/ready\n")
	fmt.Print("  go-wait-probe --endpoint http://localhost:8080/ready echo 'ready to run anything :)'\n")
	fmt.Print("  go-wait-probe --endpoint http://localhost:8080/ready -- echo 'ready to run anything :)'\n")
	fmt.Print("  go-wait-probe --endpoint http://localhost:8080/ready --programTimeout 2s --runCommandOnTimeout echo 'ready to run anything :)'\n")

	fmt.Print("Flags:\n\n")
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
