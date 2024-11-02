package main

import (
	"context"
	"fmt"
	"log"
	"magisterium/aequitas"
	"magisterium/stayalive"
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

func (r rpc) send() bool {
	if r.isLowered {
		r.prio = "be"
	}
	var sock string
	if r.prio == "hi" {
		sock = "localhost:2220"
	}
	if r.prio == "lo" {
		sock = "localhost:2222"
	}
	if r.prio == "be" {
		sock = "localhost:2224"
	}
	conn, err := grpc.NewClient(sock, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:2222: %v", err)
		return false
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
	fmt.Printf("sending RPC with priority: %v \n", r.prio)
	if err != nil {
		log.Fatalf("error calling function StayAlive: %v", err)
	}
	r.elapsed = time.Since(start)
	fmt.Printf("this rpc took %v to complete \n", r.elapsed)

	return resp.GetAliveResp()
}

func (r rpc) admit() bool {
	if r.prio == "hi" {
		r.goal = 50 * time.Millisecond
	} else {
		r.goal = 25 * time.Millisecond
	}
	return aequitas.TimeCheck(r.goal, r.elapsed)
}

func sendRPC(done int, sent int) {
	var rpc rpc
	rpc.prio = prios[rand.Intn(len(prios))]
	for {
		sent++
		completed := rpc.send()
		if completed {
			done++
		}
		if !rpc.isLowered {
			rpc.isLowered = rpc.admit()
			fmt.Println("priority lowered")
		}
		if rpc.prio == "hi" {
			time.Sleep(time.Second)
		}
	}
}

func main() {
	var done int
	var sent int
	//weighted random selection of priorities required
	for {
		go sendRPC(done, sent)
	}
}
