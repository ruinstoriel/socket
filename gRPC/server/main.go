package main

import (
	"context"
	"log"
	"net"

	pb "gRPC/proto"
	"google.golang.org/grpc"
)

// 定义服务实现
type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %s", req.Name)
	return &pb.HelloReply{Message: "Hello " + req.Name}, nil
}

func main() {
	// 创建监听器
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 创建 gRPC 服务
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &greeterServer{})

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
