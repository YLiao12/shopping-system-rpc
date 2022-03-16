# shopping-system-rpc

![sho](https://user-images.githubusercontent.com/62742611/158542354-97409069-2131-4345-8dc3-57b915aaa422.png)


## How to build this rpc application

This application could be build and deployed on localhost, and different servers (including a user server, a product server and an order server) will run on different ports. We have also deployed this application on our own server, you can test the application on (vue ip address). 

Here is the localhost deploy steps:
 

### 1. run servers on different terminals

`cd shopping-system-rpc`

`go run .\user_server\user_server.go`

`go run .\product_server\product_server.go`

`go run .\order_server\order_server.go`


### 2. run the rpc client


`go run client/main.go`


### 3.  deploy Vue module

```
cd web/shopping

npm install

npm run serve
```

Access the system via your browser 

