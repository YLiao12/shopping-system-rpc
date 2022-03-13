package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	demo "shopping-system-rpc/data"

	_ "github.com/go-sql-driver/mysql"
)

type demoServer struct {
	demo.UnimplementedDemoServer
	savedResults []*demo.Response //用于服务端流
}

type User struct {
	// 对应name表字段
	name string `db:"name"`
	// 对应age表字段
	balance int32 `db:"balance"`
}

var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "leo:21516114@tcp(139.196.187.198:3306)/rpc_user?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

//查询用户
func (s *demoServer) GetUserByUserId(ctx context.Context, in *demo.UserId) (*demo.UserInfo, error) {
	sqlStr := "select name, balance from users where id = ?;"
	var u User
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := db.QueryRow(sqlStr, in.UserId).Scan(&u.name, &u.balance)

	cRes := &demo.UserInfo{
		UserId:  in.UserId,
		Name:    u.name,
		Balance: u.balance,
	}
	if err != nil {
		fmt.Println("err=", err)
	}
	fmt.Println("查询成功", u)
	return cRes, err
}

var (
	tls  = flag.Bool("tls", false, "使用启用tls") //默认false
	port = flag.Int("port", 50055, "服务端口")    //默认50055
)

func main() {
	flag.Parse()

	errDb := initDB() // 调用输出化数据库的函数
	if errDb != nil {
		fmt.Printf("init db failed,err:%v\n", errDb)
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}

	var opts []grpc.ServerOption

	if *tls {
		creds, err := credentials.NewServerTLSFromFile("keys/server.crt", "keys/server.key")
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	s := grpc.NewServer(opts...)
	demo.RegisterDemoServer(s, &demoServer{})
	reflection.Register(s)
	log.Printf("Server listeing at :%v\n", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
