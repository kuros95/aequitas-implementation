package main

import (
	"context"
	"flag"
	"log"
	"magisterium/stayalive"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type gRPCServer struct {
	stayalive.UnimplementedStayAliveServiceServer
}

func (m *gRPCServer) StayAlive(ctx context.Context, request *stayalive.StayAliveRequest) (*stayalive.StayAliveResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	prio := md.Get("prio")[0]
	log.Printf("Got prio: %v. Sending response...", prio)
	return &stayalive.StayAliveResponse{AliveResp: bool(true)}, nil
}

func main() {
	var port int
	flag.IntVar(&port, "p", 1, "port on which to run rpc server")
	flag.Parse()

	if port == 1 {
		log.Fatalf("exiting, no port provided. Please provide a port on which to run server (-p flag).")
	}

	portNumber := strconv.Itoa(port)
	addr := "localhost:" + portNumber
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	s := grpc.NewServer()
	gRPCServer := &gRPCServer{}
	stayalive.RegisterStayAliveServiceServer(s, gRPCServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
