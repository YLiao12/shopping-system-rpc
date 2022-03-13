package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	demo "shopping-system-rpc/data"
	"time"
)

var (
	tls                = flag.Bool("tls", false, "是否使用tls")
	serverAddr         = flag.String("server_addr", "localhost:50055", "服务端地址，格式： host:port")
	serverHostOverride = flag.String("server_host_override", "a.grpc.test.com", "验证TLS握手返回的主机名的服务器名称。需要和服务端证书中dns段落匹配")
)

func findUsers(client demo.DemoClient, userId *demo.UserId) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetUserByUserId(ctx, userId)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("返回内容: %v", resp.Name)
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	if *tls {
		creds, err := credentials.NewClientTLSFromFile("keys/ca.crt", *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := demo.NewDemoClient(conn)

	fmt.Printf("#############第1次请求，简单模式########\n")
	findUsers(client, &demo.UserId{UserId: 1})
	fmt.Printf("\n\n")

}
