package main

import (
	"log"
	"magisterium/utils"
	"os"
	"os/exec"
	"time"
)

func main() {
	logFile, err := os.Create("client.log")
	if err != nil {
		log.Fatalf("failed to create logfile, error: %v", err)
	}
	log.SetOutput(logFile)
	log.Printf("shaping traffic...")
	tc := exec.Command("./tc-on-host.sh")
	if err := tc.Run(); err != nil {
		log.Fatalf("failed to apply traffic control, error: %v", err)
	}
	log.Printf("traffic control added")

	log.Printf("sending RPCs...")
	//weighted random selection of priorities required
	for {
		go utils.SendRPC()
		time.Sleep(time.Millisecond)
	}
}
