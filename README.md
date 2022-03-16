# shopping-system-rpc

![image](https://user-images.githubusercontent.com/62742611/158137046-8b00b82f-bd5f-4c84-8061-d73360541438.png)


All server can run on localhost with different ports

user server port: 50055

product server port: 50056

order server port: 50057

The MySQL is on a Aliyun Server



## 1. run servers on different terminals

`go run .\user_server\user_server.go`

`go run .\product_server\product_server.go`

`go run .\order_server\order_server.go`


## 2. the client is use to test order


`go run client/main.go`


## 3.  deploy Vue module

```
cd web/shopping

npm install

npm run serve
```

Access the system via your browser 

