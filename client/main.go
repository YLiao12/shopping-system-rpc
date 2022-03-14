package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	demo "shopping-system-rpc/data"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "localhost:50057", "服务端地址，格式： host:port")
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

func findProducts(client demo.DemoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	emp := &demo.Empty{}
	getStreamClient, err := client.GetProducts(ctx, emp)
	if err != nil {
		log.Fatalf("error: %v: ", err)
	}
	for {
		resp, err := getStreamClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetStream错误=%v", client, err)
		}
		fmt.Printf("本次返回结果:%v\n", resp.Name)
	}
}

func makeOrder(client demo.DemoClient, orderInfo *demo.OrderInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.MakeOrder(ctx, orderInfo)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("返回内容: %v", resp.Result)
}

func main() {
	flag.Parse()
	//var opts []grpc.DialOption
	//
	//opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := demo.NewDemoClient(conn)

	//fmt.Printf("#############User RPC test########\n")
	//findUsers(client, &demo.UserId{UserId: 1})
	//fmt.Printf("\n\n")

	//fmt.Printf("#############Product RPC test########\n")
	//findProducts(client)

	fmt.Printf("#############order RPC test########\n")
	makeOrder(client, &demo.OrderInfo{ProductId: 1, UserId: 1})

}
