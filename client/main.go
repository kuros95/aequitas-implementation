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

func rpcWithPrio(prio string, lowered bool) bool {
	if lowered {
		prio = "be"
	}
	conn, err := grpc.NewClient("localhost:2222", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:2222: %v", err)
		return false
	}
	defer conn.Close()
	c := stayalive.NewStayAliveServiceClient(conn)

	header := metadata.New(map[string]string{
		"prio": prio,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctxWithMD := metadata.NewOutgoingContext(ctx, header)

	resp, err := c.StayAlive(ctxWithMD, &stayalive.StayAliveRequest{})
	fmt.Printf("sending RPC with priority: %v \n", prio)
	if err != nil {
		log.Fatalf("error calling function StayAlive: %v", err)
	}

	return resp.GetAliveResp()
}

func main() {
	prios := [2]string{"hi", "lo"}
	var done int
	var fail int
	lowered := false
	//send rpcs in parallel with different prios, then measure each individually and individually lower prio, retain lowered prio
	//needs a category of rpc to send, so that lowered prio may be preserved (interface?), and they may be differently served
	//weighted random selection of priorities required
	for {
		prio := prios[rand.Intn(len(prios))]
		start := time.Now()
		completed := rpcWithPrio(prio, lowered)
		elapsed := time.Since(start)
		if completed {
			done++
		} else {
			fail++
		}

		reduce := aequitas.TimeCheck(time.Millisecond, elapsed)
		if reduce && prio == "hi" || prio == "lo" {
			lowered = true
		}
		time.Sleep(time.Second)
	}

}
