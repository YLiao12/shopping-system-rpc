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

	_ "github.com/go-sql-driver/mysql"
)

var productDb *sql.DB

type demoServer struct {
	demo.UnimplementedDemoServer
	savedResults []*demo.Response //用于服务端流
}

// 定义一个初始化数据库的函数
func initProductDB() (err error) {
	// DSN:Data Source Name
	dsn := "leo:21516114@tcp(139.196.187.198:3306)/rpc_product?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	productDb, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = productDb.Ping()
	if err != nil {
		return err
	}
	return nil
}

type Product struct {
	// 对应id表字段
	id int32 `db:"id"`
	// 对应name表字段
	name string `db:"name"`
	// 库存
	stock int32 `db:"stock"`
	// 单价
	price int32 `db:"price"`
}

func (s *demoServer) GetProduct(ctx context.Context, in *demo.ProductId) (*demo.Product, error) {
	fmt.Println("查询成功12")
	sqlStr := "select id, name, stock, price from products where id = ?;"
	var p Product
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := productDb.QueryRow(sqlStr, in.ProductId).Scan(&p.id, &p.name, &p.stock, &p.price)

	cRes := &demo.Product{
		Id:    p.id,
		Name:  p.name,
		Stock: p.stock,
		Price: p.price,
	}
	if err != nil {
		fmt.Println("err=", err)
	}
	fmt.Println("查询成功", p)
	return cRes, err
}

//查询所有商品
func (s *demoServer) GetProducts(emp *demo.Empty, pipe demo.Demo_GetProductsServer) error {
	sqlStr := "select id, name, stock, price from products;"
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	fmt.Println("查询成功")
	rows, err := productDb.Query(sqlStr)
	defer rows.Close()

	cRes := &[]demo.Product{}

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.id, &p.name, &p.stock, &p.price)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
		}

		errSend := pipe.Send(&demo.Product{
			Id:    p.id,
			Name:  p.name,
			Stock: p.stock,
			Price: p.price,
		})

		if errSend != nil {
			return errSend
		}
	}

	if err != nil {
		fmt.Println("err=", err)
	}
	fmt.Println("查询成功", cRes)

	return nil
}

func (s *demoServer) DecreaseProductStock(ctx context.Context, in *demo.ProductId) (*demo.Response, error) {
	sqlStr := "update products set stock = stock - 1 where id = ?"
	ret, err := productDb.Exec(sqlStr, in.ProductId)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
	}
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
	}
	fmt.Printf("product table update success, affected rows:%d\n", ret)

	response := &demo.Response{
		Result: 1,
	}

	return response, err
}

var (
	portProduct = flag.Int("port", 50056, "服务端口")
)

func main() {
	flag.Parse()

	errDb := initProductDB() // 调用输出化数据库的函数
	if errDb != nil {
		fmt.Printf("init db failed,err:%v\n", errDb)
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portProduct))
	if err != nil {
		log.Fatalln(err)
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)
	demo.RegisterDemoServer(s, &demoServer{})
	reflection.Register(s)
	log.Printf("Server listening at :%v\n", *portProduct)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
