package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	demo "shopping-system-rpc/data"
	"strconv"
	"time"
)

var (
	userServerAddr    = flag.String("user_server_addr", "localhost:50055", "服务端地址，格式： host:port")
	productServerAddr = flag.String("product_server_addr", "localhost:50056", "服务端地址，格式： host:port")
	orderServerAddr   = flag.String("order_server_addr", "localhost:50057", "服务端地址，格式： host:port")
)

func findUsers(client demo.DemoClient, userId *demo.UserId) *demo.UserInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetUserByUserId(ctx, userId)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("返回内容: %v", resp.Name)
	return resp
}

func findProduct(client demo.DemoClient, productId *demo.ProductId) []*demo.Product {
	var products []*demo.Product
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetProduct(ctx, productId)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("返回内容: %v", resp.Name)
	products = append(products, resp)
	return products
}

func findProducts(client demo.DemoClient) []*demo.Product {
	var products []*demo.Product

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
		products = append(products, resp)
	}
	fmt.Println(products[1])
	return products
}

func makeOrder(client demo.DemoClient, orderInfo *demo.OrderInfo) *demo.Response {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.MakeOrder(ctx, orderInfo)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("返回内容: %v \n", resp.Result)
	return resp
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Connection", "close")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func main() {
	flag.Parse()
	//var opts []grpc.DialOption
	//
	//opts = append(opts, grpc.WithBlock())

	userconn, err := grpc.Dial(*userServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer userconn.Close()
	userclient := demo.NewDemoClient(userconn)

	productconn, err := grpc.Dial(*productServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer productconn.Close()
	productclient := demo.NewDemoClient(productconn)

	orderconn, err := grpc.Dial(*orderServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer orderconn.Close()
	orderclient := demo.NewDemoClient(orderconn)

	// fmt.Printf("#############User RPC test########\n")
	// findUsers(client, &demo.UserId{UserId: 1})
	// fmt.Printf("\n\n")

	fmt.Printf("#############Product RPC test########\n")
	findProducts(productclient)

	fmt.Printf("#############order RPC test########\n")
	makeOrder(orderclient, &demo.OrderInfo{ProductId: 2, UserId: 1})

	r := gin.Default()
	r.Use(Cors())

	r.GET("/MakeOrder", func(c *gin.Context) {
		fmt.Printf("#############order RPC test########\n")

		p_id, _ := strconv.ParseInt(c.Query("productId"), 10, 32)
		u_id, _ := strconv.ParseInt(c.Query("userId"), 10, 32)
		productId := int32(p_id)
		userId := int32(u_id)
		res := makeOrder(orderclient, &demo.OrderInfo{ProductId: productId, UserId: userId})
		c.JSON(200, res)
	})
	r.GET("/AllProducts", func(c *gin.Context) {
		fmt.Printf("#############Product RPC test########\n")
		productslist := findProducts(productclient)
		c.JSON(200, productslist)
	})

	r.GET("/OneProduct", func(c *gin.Context) {
		fmt.Printf("#############FindOneProduct RPC test########\n")
		p_id, _ := strconv.ParseInt(c.Query("productId"), 10, 32)
		product_Id := int32(p_id)
		oneproduct := findProduct(productclient, &demo.ProductId{ProductId: product_Id})
		c.JSON(200, oneproduct)
	})

	r.GET("/GetUser", func(c *gin.Context) {
		fmt.Printf("#############GetUser RPC test########\n")
		user := findUsers(userclient, &demo.UserId{UserId: 1})
		c.JSON(200, user)
	})

	r.Run(":5555")

}
