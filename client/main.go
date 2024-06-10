package main

import (
	"context"
	"log"
	aeuqitas "magisterium/aequitas"
	"magisterium/stayalive"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connect(prio int) time.Duration {
	conn, err := grpc.Dial("localhost:2222", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:2222: %v", err)
	}
	defer conn.Close()
	c := stayalive.NewStayAliveServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	start := time.Now()
	r, err := c.StayAlive(ctx, &stayalive.StayAliveRequest{})
	if err != nil {
		log.Fatalf("error calling function SayHello: %v", err)
	}
	elapsed := time.Since(start)
	log.Printf("gRPC server response: %v, with elapsed time: %v, and priority: %v", r.GetAliveResp(), elapsed, prio)

	return elapsed
}

func main() {
	prio := 3
	for {
		elapsed := connect(prio)
		reduce := aeuqitas.TimeCheck(time.Millisecond, elapsed)
		if reduce && prio > 0 {
			prio = prio - 1
		}
		time.Sleep(time.Second)
	}

}
