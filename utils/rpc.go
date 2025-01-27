package utils

import (
	"context"
	"log"
	sendmessage "magisterium/sendmess"
	"time"

	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type rpc struct {
	prio
	isLowered bool
	elapsed   time.Duration
	size      int
}

type prio struct {
	prio            string
	latency         time.Duration
	target_pctl     int
	incr_window     time.Duration
	p_admit         float64
	t_last_increase time.Time
}

var prios = []prio{{"hi", 20 * time.Millisecond, 99, 0, 1, time.Now()}, {"lo", 15 * time.Millisecond, 85, 0, 1, time.Now()}}

func (r rpc) send() (bool, time.Duration, int) {
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
		log.Fatalf("failed to connect to gRPC server at %v: %v", sock, err)
		return false, 0, 0
	}
	defer conn.Close()
	c := sendmessage.NewSendMessageServiceClient(conn)

	header := metadata.New(map[string]string{
		"prio": r.prio.prio,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctxWithMD := metadata.NewOutgoingContext(ctx, header)

	start := time.Now()
	resp, err := c.SendMessage(ctxWithMD, &sendmessage.SendMessageRequest{})
	log.Printf("sending RPC with priority: %v to %v \n", r.prio.prio, sock)
	if err != nil {
		log.Printf("error calling function SendMessage: %v", err)
		return false, 0, 0
	}
	r.elapsed = time.Since(start)

	return resp.GetAliveResp(), r.elapsed, int(resp.GetSize())
}

func SendRPC(use_64kb_payload bool) {
	var rpc rpc
	prio_to_assign := prios[rand.Intn(len(prios))]
	if use_64kb_payload {
		rpc.size = 64
	} else {
		rpc.size = 32
	}

	for {
		if rand.Float64() <= prio_to_assign.p_admit {
			rpc.prio.prio = prio_to_assign.prio
		} else {
			rpc.prio.prio = "be"
			rpc.isLowered = true
		}

		completed, elapsed, size := rpc.send()
		rpc.elapsed = elapsed
		rpc.size = size

		if completed {
			log.Printf("completed an RPC of size %v with prio %v", rpc.size, rpc.prio.prio)
			rpc.admit()
		} else {
			log.Printf("falied to complete an RPC of size %v with prio %v, lowering priority", rpc.size, rpc.prio.prio)
			rpc.isLowered = true
		}
		// if rpc.prio.prio == "hi" {
		// 	time.Sleep(5 * time.Second)
		// } else {
		// 	time.Sleep(time.Millisecond)
		// }
	}
}

func SendRPCNoAequitas(use_64kb_payload bool) {
	var rpc rpc
	prio_to_assign := prios[rand.Intn(len(prios))].prio
	rpc.prio.prio = prio_to_assign
	if use_64kb_payload {
		rpc.size = 64
	} else {
		rpc.size = 32
	}

	for {
		completed, elapsed, size := rpc.send()
		rpc.elapsed = elapsed
		rpc.size = size

		if completed {
			log.Printf("completed an RPC with prio %v", rpc.prio.prio)
		} else {
			log.Printf("falied to complete an RPC with prio %v", rpc.prio.prio)
		}
	}
}
