package main

import (
	"context"
	"log"
	"time"

	pb "gRPC/proto"
	"google.golang.org/grpc"
)

func main() {
	// 连接到服务端
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	// 调用服务
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("Failed to call SayHello: %v", err)
	}

	log.Printf("Response: %s", resp.Message)
}
