package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	var endpoint, probeIntervalString, programTimeoutString string

	flag.StringVar(&endpoint, "endpoint", "", "The endpoint to probe.")
	flag.StringVar(&probeIntervalString, "probeInterval", "1s", "The interval at which the probe is executed. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s")
	flag.StringVar(&programTimeoutString, "programTimeout", "15s", "Timeout after the program is considered unsuccessful and the tool exits with 1. The format needs to be parsable by time.ParseDuration. Examples: 300ms, 3s")

	flag.Parse()

	if endpoint == "" {
		printHelp("Error: missing --endpoint configuration.")
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
		os.Exit(0)
	case <-time.After(programTimeout):
		fmt.Printf("Error: programTimeout of %s reached. Exiting with 1.\n", programTimeout.String())
		os.Exit(1)
	}
}

func printHelp(errorMessage string) {
	fmt.Print(errorMessage, "\n\n")
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
