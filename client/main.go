package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"magisterium/utils"
	"os"
	"os/exec"
	"time"
)

var noAequitas bool
var use64 bool
var stderr bytes.Buffer

func main() {

	flag.BoolVar(&noAequitas, "n", false, "do not use the aequitas algorithm")
	flag.BoolVar(&use64, "u", false, "use 64kb payload for messages")

	logFile, err := os.Create("client.log")
	if err != nil {
		log.Fatalf("failed to create logfile, error: %v", err)
	}
	log.SetOutput(logFile)

	flag.Parse()

	log.Printf("shaping traffic...")
	tc := exec.Command("./tc-on-host.sh")
	if err = tc.Run(); err != nil {
		log.Fatalf("failed to apply traffic control, error: %v", err)
	}
	log.Printf("traffic control added")
	utils.AequitasInit(15, 98)

	//remember to run the container with --mount flag. It will enable you to collect the generated logs.
	tcpdump := exec.Command("./run-tcpdump.sh")
	tcpdump.Stderr = &stderr

	go func() {
		err = tcpdump.Run()

		if err != nil {
			fmt.Printf("error: %v: %v", err, stderr.String())
			log.Fatalf("failed to start capturing traffic data, error: %v", err)
		}
	}()

	log.Printf("sending RPCs...")
	//weighted random selection of priorities required
	for {

		if noAequitas {
			go utils.SendRPCNoAequitas(use64)
		} else if !noAequitas {
			go utils.SendRPC(use64)
		}

		time.Sleep(time.Millisecond)
	}
}
