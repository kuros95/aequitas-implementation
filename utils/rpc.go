package utils

import (
	"context"
	"log"
	"magisterium/stayalive"
	"time"

	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type rpc struct {
	prio
	isLowered bool
	goal      time.Duration
	elapsed   time.Duration
}

type prio struct {
	prio        string
	latency     time.Duration
	target_pctl int
	incr_window int
}

var prios = []prio{{"hi", 20 * time.Millisecond, 99, 0}, {"lo", 15 * time.Millisecond, 97, 0}}

func (r rpc) send() (bool, time.Duration) {
	if r.isLowered {
		r.prio.prio = "be"
	}
	var sock string
	if r.prio.prio == "hi" {
		sock = "172.17.0.2:2220"
	}
	if r.prio.prio == "lo" {
		sock = "172.17.0.3:2222"
	}
	if r.prio.prio == "be" {
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
		"prio": r.prio.prio,
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

func SendRPC() {
	var rpc rpc
	rpc.prio.prio = prios[rand.Intn(len(prios))].prio
	for {
		completed, elapsed := rpc.send()
		rpc.elapsed = elapsed
		if completed {
			log.Printf("completed an RPC with prio %v", rpc.prio)
		}
		if !rpc.isLowered {
			rpc.isLowered = rpc.admit()
		}
		if rpc.prio.prio == "hi" {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(time.Second)
		}
	}
}
