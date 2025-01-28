package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	sendmessage "magisterium/sendmess"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var port int

type gRPCServer struct {
	sendmessage.UnimplementedSendMessageServiceServer
}

func (m *gRPCServer) SendMessage(ctx context.Context, request *sendmessage.SendMessageRequest) (*sendmessage.SendMessageResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	prio := md.Get("prio")[0]
	size := request.GetSize()

	bufferSize := size * 1024
	file, err := os.Open(request.Payload)
	buff := make([]byte, bufferSize)
	if err != nil {
		log.Printf("failed to open response payload: %v", err)
		return &sendmessage.SendMessageResponse{
			AliveResp: bool(false),
			Size:      size,
			Payload:   "",
			MessChunk: make([]byte, 0),
		}, err
	}
	defer file.Close()

	var messChunk []byte

	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("failed to read response payload: %v", err)
			}
			break
		}
		messChunk = buff[:bytesRead]
	}

	log.Printf("Got prio: %v. Sending response...", prio)
	return &sendmessage.SendMessageResponse{
		AliveResp: bool(true),
		Size:      size,
		Payload:   file.Name(),
		MessChunk: messChunk,
	}, nil
}

func getIP(iface string) string {
	var ipv4 net.IP

	nic, err := net.InterfaceByName(iface)
	if err != nil {
		log.Fatalf("could not get network interface info, error: %v", err)
	}

	addrs, err := nic.Addrs()
	if err != nil {
		log.Fatalf("could not get addresses from interface, error: %v", err)
	}

	for _, addr := range addrs {
		if ipv4 = addr.(*net.IPNet).IP.To4(); ipv4 != nil {
			break
		}
	}

	if ipv4 == nil {
		return ""
	}

	return ipv4.String()
}

func main() {

	flag.IntVar(&port, "p", 1, "port on which to run rpc server")
	flag.Parse()

	if port == 1 {
		log.Fatalf("exiting, no port provided. Please provide a port on which to run server (-p flag).")
	}

	ip := getIP("eth0")

	portNumber := strconv.Itoa(port)
	addr := ip + ":" + portNumber
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	s := grpc.NewServer()
	gRPCServer := &gRPCServer{}
	sendmessage.RegisterSendMessageServiceServer(s, gRPCServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
