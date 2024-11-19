package main

import (
	"flag"
	"log"
	"magisterium/utils"
	"os/exec"
	"time"
)

func main() {
	var size int
	flag.IntVar(&size, "s", 0, "size of RPC to send")
	flag.Parse()

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
