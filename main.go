package main

import (
	"github.com/integrii/flaggy"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	n := 1
	flaggy.Int(&n, "n", "", "number of instances to start")
	flaggy.Parse()
	if len(flaggy.TrailingArguments) == 0 {
		log.Fatalln("Supply command after trailing --")
	}
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wg.Done()
			err := runCmd(flaggy.TrailingArguments)
			if err != nil {
				log.Println("ERROR:", err)
			}
		}()
	}
	wg.Wait()
}

func runCmd(args []string) error {
	cmd := exec.Command(flaggy.TrailingArguments[0], flaggy.TrailingArguments[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
