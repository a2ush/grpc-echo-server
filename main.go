package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/a2ush/grpc-echo-server/rpc"
)

type EchoServerHandler struct{}

type GreetingServerHandler struct {
	name string
}

func NewServerHandler() *GreetingServerHandler {
	return &GreetingServerHandler{
		name: "GreetingBot",
	}
}

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen; %v", err)
	}

	// Add logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(zapLogger),
			),
		),
	)

	rpc.RegisterEchoServer(
		server,
		&EchoServerHandler{},
	)

	rpc.RegisterGreatingServer(
		server,
		NewServerHandler(),
	)

	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		server.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()
}

func (s *GreetingServerHandler) Reply(
	ctx context.Context,
	req *rpc.GreatingClientRequest,
) (*rpc.ServerGreatingResponse, error) {

	wording := ""

	if req.ClientGreeting == rpc.Format_Unknown || req.ClientGreeting > rpc.Format_Casual {
		wording = "hummm :( WHO ARE YOU?"
	}

	switch req.ClientGreeting {
	case rpc.Format_Formal:
		log.Printf("Formal")
		wording = "How are you?"
	case rpc.Format_Normal:
		log.Printf("Normal")
		wording = "Hello"
	case rpc.Format_Casual:
		log.Printf("Casual")
		wording = "Hi :)"
	}

	return &rpc.ServerGreatingResponse{
		Name: s.name,
		Format: &rpc.Format{
			Echo: wording,
		},
		CreateTime: &timestamp.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   int32(time.Now().Nanosecond()),
		},
	}, nil
}

func (s *EchoServerHandler) Reply(
	ctx context.Context,
	req *rpc.ClientRequest,
) (*rpc.ServerResponse, error) {

	return &rpc.ServerResponse{
		Name:    "EchoBot",
		Message: req.Message + " by EchoBot",
		CreateTime: &timestamp.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   int32(time.Now().Nanosecond()),
		},
	}, nil
}
