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
var add_inc float64
var mul_dec float64
var min_adm float64
var lat_tgt int
var tgt_pctl int
var stderr bytes.Buffer

func main() {

	flag.BoolVar(&noAequitas, "n", false, "do not use the aequitas algorithm")
	flag.BoolVar(&use64, "u", false, "use 64kb payload for messages")
	flag.Float64Var(&add_inc, "a", 0.01, "set additive increase factor for aequitas algorithm, 0.01 by default")
	flag.Float64Var(&mul_dec, "m", 0.01, "set multiplicative decrease factor for aequitas algorithm, 0.01 by default")
	flag.Float64Var(&min_adm, "d", 0.01, "set minimum admission probability for aequitas algorithm, 0.01 by default - DO NOT SET TO ZERO")
	flag.IntVar(&lat_tgt, "l", 15, "latency target (in ms) for aequitas algorithm, 15 by default")
	flag.IntVar(&tgt_pctl, "p", 98, "target percentile for aequitas algorithm, 98 by default")

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
	utils.AequitasInit(lat_tgt, tgt_pctl)

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
	for {

		if noAequitas {
			go utils.SendRPCNoAequitas(use64)
		} else if !noAequitas {
			go utils.SendRPC(use64, add_inc, mul_dec, min_adm)
		}

		time.Sleep(time.Millisecond)
	}
}
