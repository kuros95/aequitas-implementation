package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"magisterium/utils"
	"os"
	"os/exec"
	"sync"
	"time"
)

var noAequitas bool
var noTc bool
var use64 bool
var add_inc float64
var mul_dec float64
var min_adm float64
var lat_tgt int
var tgt_pctl int
var maxRpcs int
var stderr bytes.Buffer

func main() {

	flag.BoolVar(&noAequitas, "n", false, "do not use the aequitas algorithm")
	flag.BoolVar(&noTc, "t", false, "do not use traffic control")
	flag.BoolVar(&use64, "u", false, "use 64kb payload for messages")
	flag.Float64Var(&add_inc, "a", 0.01, "set additive increase factor for aequitas algorithm, 0.01 by default")
	flag.Float64Var(&mul_dec, "m", 0.01, "set multiplicative decrease factor for aequitas algorithm, 0.01 by default")
	flag.Float64Var(&min_adm, "d", 0.01, "set minimum admission probability for aequitas algorithm, 0.01 by default - DO NOT SET TO ZERO")
	flag.IntVar(&lat_tgt, "l", 15, "latency target (in ms) for aequitas algorithm, 15 by default")
	flag.IntVar(&tgt_pctl, "p", 98, "target percentile for aequitas algorithm, 98 by default")
	flag.IntVar(&maxRpcs, "r", 10000, "maximum number of RPCs to send, 10000 by default")

	logFile, err := os.Create("client.log")
	if err != nil {
		log.Fatalf("failed to create logfile, error: %v", err)
	}
	log.SetOutput(logFile)

	flag.Parse()

	waitChan := make(chan struct{}, maxRpcs)
	var wg sync.WaitGroup

	if !noTc {
		log.Printf("shaping traffic...")
		tc := exec.Command("./tc-on-host.sh")
		if err = tc.Run(); err != nil {
			log.Fatalf("failed to apply traffic control, error: %v", err)
		}
		log.Printf("traffic control added")
	}

	if !noAequitas {
		utils.AequitasInit(lat_tgt, tgt_pctl)
		log.Printf("Aequitas initiated with latency target %vms, and target percentile of completed RPCs %v", lat_tgt, tgt_pctl)
	}

	tcpdump := exec.Command("./run-tcpdump.sh")
	tcpdump.Stderr = &stderr

	go func() {
		err = tcpdump.Run()

		if err != nil {
			fmt.Printf("error: %v: %v", err, stderr.String())
			log.Fatalf("failed to start capturing traffic data, error: %v", err)
		}
	}()

	if noAequitas {
		log.Printf("sending RPCs...")
	} else {
		log.Printf("sending RPCs with additive increase set to %v, multiplicative decrease set to %v, and minimum admission probability set to %v", add_inc, mul_dec, min_adm)
	}

	for i := range maxRpcs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			waitChan <- struct{}{}
			utils.SendRPC(use64, noAequitas, add_inc, mul_dec, min_adm)
			time.Sleep(time.Millisecond)
			<-waitChan
		}(i)
	}
	wg.Wait()
}
