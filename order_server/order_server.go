package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	demo "shopping-system-rpc/data"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var orderDb *sql.DB

type demoServer struct {
	demo.UnimplementedDemoServer
	savedResults []*demo.Response //用于服务端流
}

// 定义一个初始化数据库的函数
func initOrderDB() (err error) {
	// DSN:Data Source Name
	dsn := "leo:21516114@tcp(139.196.187.198:3306)/rpc_order?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	orderDb, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = orderDb.Ping()
	if err != nil {
		return err
	}
	return nil
}

var (
	portOrder         = flag.Int("port", 50057, "服务端口")
	productServerAddr = flag.String("product_server_addr", "localhost:50056", "服务端地址，格式： host:port")
	userServerAddr    = flag.String("user_server_addr", "localhost:50055", "服务端地址，格式： host:port")
)

func findProduct(client demo.DemoClient, productId *demo.ProductId) *demo.Product {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetProduct(ctx, productId)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("返回内容: %v", resp.Name)
	return resp
}

func findUser(client demo.DemoClient, userId *demo.UserId) *demo.UserInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetUserByUserId(ctx, userId)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("返回内容: %v", resp.Name)
	return resp
}

func decreaseBalance(client demo.DemoClient, userId int32, price int32, balance int32) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	decreaseBalanceMsg := &demo.DecreaseBalance{
		UserId:  userId,
		Price:   price,
		Balance: balance,
	}
	resp, err := client.DecreaseUserBalance(ctx, decreaseBalanceMsg)
	if err != nil {
		return
	}
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("用户: %v 的余额已减少", resp.Result)
}

func decreaseStock(client demo.DemoClient, productId *demo.ProductId) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.DecreaseProductStock(ctx, productId)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("商品库存已减少 \n", resp.Result)
}

func (s *demoServer) MakeOrder(ctx context.Context, order *demo.OrderInfo) (*demo.Response, error) {
	// 根据productId查询产品具体信息
	flag.Parse()
	conn, err := grpc.Dial(*productServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	productClient := demo.NewDemoClient(conn)

	productId := &demo.ProductId{ProductId: order.ProductId}
	product := findProduct(productClient, productId)
	price := product.Price
	stock := product.Stock

	// 检查用户余额
	userConn, userErr := grpc.Dial(*userServerAddr, grpc.WithInsecure())
	if userErr != nil {
		log.Fatalf("fail to dial: %v", userErr)
	}
	defer userConn.Close()
	userClient := demo.NewDemoClient(userConn)

	userId := &demo.UserId{UserId: order.UserId}
	user := findUser(userClient, userId)
	balance := user.Balance

	var response *demo.Response

	fmt.Printf("商品的库存为：%d.\n ", stock)

	// 更新order表
	if stock > 0 && balance > price {
		sqlStr := "insert into orders(user_id, product_id) values (?,?)"
		ret, err := orderDb.Exec(sqlStr, order.UserId, order.ProductId)
		if err != nil {
			fmt.Printf("insert failed, err:%v\n", err)
		}
		theID, err := ret.LastInsertId() // 新插入数据的id
		if err != nil {
			fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		}
		fmt.Printf("insert success, the id is %d.\n", theID)

		// 减少用户余额
		decreaseBalance(userClient, order.UserId, price, balance)

		// 减少商品库存
		decreaseStock(productClient, productId)
		response = &demo.Response{
			Result: 0,
		}
	} else if stock == 0 {
		response = &demo.Response{
			Result: 2,
		}
	} else {
		response = &demo.Response{
			Result: 1,
		}
	}

	return response, nil
}

func main() {
	flag.Parse()

	errDb := initOrderDB() // 调用输出化数据库的函数
	if errDb != nil {
		fmt.Printf("init db failed,err:%v\n", errDb)
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portOrder))
	if err != nil {
		log.Fatalln(err)
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)
	demo.RegisterDemoServer(s, &demoServer{})
	reflection.Register(s)
	log.Printf("Server listening at :%v\n", *portOrder)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
