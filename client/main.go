package main

import (
	"context"
	"log"
	"magisterium/aequitas"
	"magisterium/stayalive"
	"os/exec"
	"time"

	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var prios = []string{"hi", "lo"}

type rpc struct {
	prio      string
	isLowered bool
	goal      time.Duration
	elapsed   time.Duration
}

func (r rpc) send() (bool, time.Duration) {
	if r.isLowered {
		r.prio = "be"
	}
	var sock string
	if r.prio == "hi" {
		sock = "172.17.0.2:2220"
	}
	if r.prio == "lo" {
		sock = "172.17.0.3:2222"
	}
	if r.prio == "be" {
		sock = "172.17.0.4:2224"
	}

	conn, err := grpc.NewClient(sock, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:2222: %v", err)
		return false, 0
	}
	defer conn.Close()
	c := stayalive.NewStayAliveServiceClient(conn)

	header := metadata.New(map[string]string{
		"prio": r.prio,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctxWithMD := metadata.NewOutgoingContext(ctx, header)

	start := time.Now()
	resp, err := c.StayAlive(ctxWithMD, &stayalive.StayAliveRequest{})
	log.Printf("sending RPC with priority: %v to %v \n", r.prio, sock)
	if err != nil {
		log.Printf("error calling function StayAlive: %v", err)
		return false, 0
	}
	r.elapsed = time.Since(start)

	return resp.GetAliveResp(), r.elapsed
}

func (r rpc) admit() bool {
	if r.prio == "hi" {
		r.goal = 20 * time.Millisecond
	} else {
		r.goal = 15 * time.Millisecond
	}
	return aequitas.LowerPrio(r.goal, r.elapsed)
}

func sendRPC() {
	var rpc rpc
	rpc.prio = prios[rand.Intn(len(prios))]
	for {
		_, rpc.elapsed = rpc.send()

		if !rpc.isLowered {
			rpc.isLowered = rpc.admit()
		}
		if rpc.prio == "hi" {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func main() {

	log.Printf("shaping traffic...")
	tc := exec.Command("./tc-on-host.sh")
	if err := tc.Run(); err != nil {
		log.Fatalf("failed to apply traffic control, error: %v", err)
	}
	log.Printf("traffic control added")

	log.Printf("sending RPCs...")
	//weighted random selection of priorities required
	for {
		go sendRPC()
		time.Sleep(time.Millisecond)
	}
}
