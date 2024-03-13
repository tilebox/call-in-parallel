package main

import (
	"context"
	"github.com/integrii/flaggy"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var version = "0.1.1"

// call-in-parallel runs a given command multiple times in parallel.
// Example Usage:
// > call-in-parallel -n 3 -- echo hello
func main() {
	n, delay, command := parseArgs()
	wg := new(sync.WaitGroup)
	wg.Add(n)

	// Let's create a context that will be canceled when an interrupt signal is received
	// That way we can stop starting new commands when the user wants to exit during startup
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			err := runCmd(command)
			if err != nil {
				log.Println("ERROR:", err)
			}
		}()

		// avoid time.Sleep() here, since it doesn't interrupt in case of received signals such as SIGINT
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop() // interrupt signal received, stop starting commands
			return
		case <-timer.C: // the delay time is over, start the next command
		}
	}
	// Wait for all running commands to finish
	wg.Wait()
}

// runCmd runs a command and waits for it to finish
// It returns an error if the command fails to start or run
// Stdout and Stderr of the command are forwarded to the parent process
func runCmd(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// parseArgs parses command line arguments and returns the number of instances to start,
// the delay between starting each instance, and the command to run
// If the command is missing, it prints an error message and exits
func parseArgs() (int, time.Duration, []string) {
	flaggy.SetName("call-in-parallel")
	flaggy.SetDescription("Run a command multiple times in parallel")
	flaggy.SetVersion(version)

	n := 1
	delay := 10 * time.Millisecond
	flaggy.Int(&n, "n", "", "number of instances to start")
	flaggy.Duration(&delay, "d", "delay", "delay between starting each instance")
	flaggy.Parse()
	if len(flaggy.TrailingArguments) == 0 {
		log.Fatalln("Missing command! Supply command after trailing --. Example: call-in-parallel -n 3 -- echo hello")
	}

	return n, delay, flaggy.TrailingArguments
}
