package main

import (
	"context"
	"fmt"
	pb "gRPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

// MyResolverBuilder 实现 resolver.Builder 接口
type MyResolverBuilder struct {
	resolver.Builder
}

// Build 方法用于构建自定义解析器
func (b *MyResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &MyResolver{
		cc:     cc,
		target: target,
	}
	store := make(map[string][]string)
	store["demo"] = []string{"127.0.0.1:50051"}
	r.addrsStore = store
	r.start()
	return r, nil
}

// Scheme 方法返回此解析器的标识符，类似于 HTTP 中的 http:// 或 https://
func (b *MyResolverBuilder) Scheme() string {
	return "mygrpc" // 自定义 Scheme
}

// MyResolver 实现 resolver.Resolver 接口
type MyResolver struct {
	resolver.Resolver
	cc         resolver.ClientConn
	target     resolver.Target
	addrsStore map[string][]string
}

func (r *MyResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (r *MyResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	// 模拟从配置文件或其他方式获取服务地址
	// 假设解析到 127.0.0.1:50051
	fmt.Println(r.target)
	address := "127.0.0.1:50051"
	addresses := []resolver.Address{
		{Addr: address},
	}
	r.cc.UpdateState(resolver.State{Addresses: addresses})

}

func (r *MyResolver) Close() {
	// 清理资源
}

// 注册自定义解析器
func init() {
	// 注册我们的解析器
	resolver.Register(&MyResolverBuilder{})
}

func main() {
	// 连接到服务端
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	// 调用服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("Failed to call SayHello: %v", err)
	}

	log.Printf("Response: %s", resp.Message)
}
