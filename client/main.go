package main

import (
	"flag"
	"log"
	"magisterium/utils"
	"os"
	"os/exec"
	"time"
)

var noAequitas bool

func main() {

	flag.BoolVar(&noAequitas, "n", false, "do not use the aequitas algorithm")
	logFile, err := os.Create("client.log")
	if err != nil {
		log.Fatalf("failed to create logfile, error: %v", err)
	}
	log.SetOutput(logFile)

	flag.Parse()

	log.Printf("shaping traffic...")
	tc := exec.Command("./tc-on-host.sh")
	if err := tc.Run(); err != nil {
		log.Fatalf("failed to apply traffic control, error: %v", err)
	}
	log.Printf("traffic control added")
	utils.AequitasInit(15, 98)

	log.Printf("sending RPCs...")
	//weighted random selection of priorities required
	for {

		if noAequitas {
			go utils.SendRPCNoAequitas()
		} else if !noAequitas {
			go utils.SendRPC()
		}

		time.Sleep(time.Millisecond)
	}
}
