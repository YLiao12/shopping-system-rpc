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

var userDb *sql.DB

// 定义一个初始化数据库的函数
func initUserDB() (err error) {
	// DSN:Data Source Name
	dsn := "leo:21516114@tcp(139.196.187.198:3306)/rpc_user?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	userDb, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = userDb.Ping()
	if err != nil {
		return err
	}
	return nil
}

//新建用户
func (s *demoServer) CreateUser(ctx context.Context, in *demo.UserInfo) (*demo.Response, error) {
	sqlStr := "insert into users(name, balance) values (?,?)"
	ret, err := userDb.Exec(sqlStr, in.Name, in.Balance)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
	}
	theID, insErr := ret.LastInsertId() // 新插入数据的id
	if insErr != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", insErr)
	}

	fmt.Printf("insert success, the id is %d.\n", theID)
	response := &demo.Response{
		Result: 1,
	}
	return response, err
}

//查询用户
func (s *demoServer) GetUserByUserId(ctx context.Context, in *demo.UserId) (*demo.UserInfo, error) {
	sqlStr := "select name, balance from users where id = ?;"
	var u User
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := userDb.QueryRow(sqlStr, in.UserId).Scan(&u.name, &u.balance)

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

func (s *demoServer) DecreaseUserBalance(ctx context.Context, in *demo.DecreaseBalance) (*demo.Response, error) {
	userId := in.UserId
	price := in.Price
	balance := in.Balance
	sqlStr := "update users set balance = ? where id = ?"
	ret, err := userDb.Exec(sqlStr, balance-price, userId)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
	}
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
	}
	fmt.Printf("user table update success, affected rows:%d \n", ret)

	response := &demo.Response{
		Result: 1,
	}

	return response, err
}

var (
	portUser = flag.Int("port", 50055, "服务端口") //默认50055
)

func main() {
	flag.Parse()

	errDb := initUserDB() // 调用输出化数据库的函数
	if errDb != nil {
		fmt.Printf("init db failed,err:%v\n", errDb)
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portUser))
	if err != nil {
		log.Fatalln(err)
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)
	demo.RegisterDemoServer(s, &demoServer{})
	reflection.Register(s)
	log.Printf("Server listening at :%v\n", *portUser)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
