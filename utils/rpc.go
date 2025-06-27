package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	sendmessage "magisterium/sendmess"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	wr "github.com/mroth/weightedrand"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type rpc struct {
	prio
	isLowered bool
	elapsed   time.Duration
	size      int32
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
var sock string
var TimeStart time.Time

func setDscp(dscp string) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, addr string) (net.Conn, error) {
		// The int value for 0x20 is 32, and for 0x40 is 64.
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", sock)
		if err != nil {
			return nil, err
		}

		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			err = errors.New("connection is not a TCP connection")
			return nil, err
		}

		tos, err := strconv.Atoi(dscp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert DSCP value %s to int: %w", dscp, err)
		}

		rawConn, err := tcpConn.SyscallConn()
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to get syscall.RawConn: %w", err)
		}

		setErr := rawConn.Control(func(fd uintptr) {
			err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TOS, tos)
		})
		if setErr != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to set DSCP option: %w", setErr)
		}
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to set DSCP option: %w", err)
		}

		return conn, nil
	}
}

func (r rpc) send() (bool, time.Duration, int32, string) {

	var dscp string
	//The int value for 0x20 is 32, and for 0x40 is 64.
	//0x20 is for low priority, 0x40 is for high priority.
	if r.isLowered {
		r.prio.prio = "be"
	}
	if r.prio.prio == "hi" {
		dscp = "64"
		sock = "172.17.0.2:2222"
	}
	if r.prio.prio == "lo" {
		dscp = "32"
		sock = "172.17.0.2:2223"
	}
	if r.prio.prio == "be" {
		dscp = "16"
		sock = "172.17.0.2:2224"
	}

	header := metadata.New(map[string]string{
		"prio": r.prio.prio,
	})

	conn, err := grpc.NewClient("passthrough:///"+sock, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(setDscp(dscp)))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at %v: %v", sock, err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), r.prio.latency)
	defer cancel()
	ctxWithMD := metadata.NewOutgoingContext(ctx, header)

	bufferSize := r.size * 1024
	file, err := os.Open(strconv.Itoa(int(r.size)) + "kb-payload")
	buff := make([]byte, bufferSize)
	if err != nil {
		fmt.Printf("failed to open request payload: %v", err)
	}

	var messChunk []byte

	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Printf("failed to read request payload: %v", err)
			}
			break
		}
		messChunk = buff[:bytesRead]
	}
	err = file.Close()
	if err != nil {
		log.Printf("failed to close request payload file: %v", err)
	}

	c := sendmessage.NewSendMessageServiceClient(conn)

	start := time.Now()
	resp, err := c.SendMessage(ctxWithMD, &sendmessage.SendMessageRequest{
		AliveReq:  "Alive?",
		Size:      r.size,
		Payload:   strconv.Itoa(int(r.size)) + "KB-payload",
		MessChunk: messChunk,
	})
	log.Printf("sending RPC with priority: %v to %v \n", r.prio.prio, sock)
	if err != nil {
		log.Printf("error calling function SendMessage: %v", err)
	}
	r.elapsed = time.Since(start)

	return resp.GetAliveResp(), r.elapsed, resp.GetSize(), r.prio.prio
}

func SendRPC(use_64kb_payload, noAequitas bool, add_inc, mul_dec, min_adm float64) {
	var rpc rpc

	chooser, _ := wr.NewChooser(
		wr.Choice{Item: "hi", Weight: 7},
		wr.Choice{Item: "lo", Weight: 3},
	)
	var indexof int

	prio_name := chooser.Pick().(string)

	if prio_name == "hi" {
		indexof = 0
	} else if prio_name == "lo" {
		indexof = 1
	}

	prio_to_assign := prios[indexof]
	if use_64kb_payload {
		rpc.size = 64
	} else {
		rpc.size = 32
	}

	if !noAequitas {
		if rand.Float64() <= prio_to_assign.p_admit {
			rpc.prio = prio_to_assign
		} else {
			rpc.prio.prio = "be"
			rpc.isLowered = true
		}
	} else if noAequitas {
		rpc.prio = prio_to_assign
	}
	rpcTime := time.Now()
	for {
		completed, elapsed, size, prio := rpc.send()

		if completed {
			log.Printf("completed an RPC of size %vkb with prio %v in %v", size, prio, elapsed)
			fmt.Printf("completed an RPC of size %vkb with prio %v in %v\n", size, prio, elapsed)
			if !noAequitas {
				rpc.admit(add_inc, mul_dec, min_adm)
			}
		} else if !completed {
			log.Printf("failed to complete an RPC of size %vkb with prio %v and latency target %v, because %v was too long... lowering priority", rpc.size, rpc.prio.prio, rpc.prio.latency, rpc.elapsed)
			if !noAequitas {
				rpc.isLowered = true
			}
		}
		if time.Since(rpcTime) > 50*time.Millisecond {
			break
		}
	}
}
