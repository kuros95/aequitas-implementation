package main

import (
	"context"
	"log"
	"magisterium/stayalive"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
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

	log.Printf("gRPC server response: %v, with elapsed time: %v", r.GetAliveResp(), elapsed)
}
