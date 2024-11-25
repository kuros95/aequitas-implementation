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

var port int
var size int

type gRPCServer struct {
	stayalive.UnimplementedStayAliveServiceServer
}

func (m *gRPCServer) StayAlive(ctx context.Context, request *stayalive.StayAliveRequest) (*stayalive.StayAliveResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	prio := md.Get("prio")[0]
	log.Printf("Got prio: %v. Sending response...", prio)
	return &stayalive.StayAliveResponse{
		AliveResp: bool(true),
		Size:      int32(size),
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
	flag.IntVar(&size, "s", 32, "max rpc message size in kilobytes")
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

	realSize := size * 1024
	s := grpc.NewServer(grpc.MaxRecvMsgSize(realSize))
	gRPCServer := &gRPCServer{}
	stayalive.RegisterStayAliveServiceServer(s, gRPCServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
