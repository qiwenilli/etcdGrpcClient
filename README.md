# etcdGrpcClient

###注册服务、服务死掉，将自动delete etcd key
```
improt (
    "go.etcd.io/etcd/clientv3"
    "github.com/qiwenilli/etcdGrpcClient"
    "golang.org/x/net/context"
)

etcdConfig := clientv3.Config{
    Endpoints:   "127.0.0.1:2379",
    DialTimeout: 3 * time.Second,
    //Username:  "",
    //Password:  "",
}
etcdClient, _ := clientv3.New(etcdConfig)

//
serverName:="/etcd3_naming/user.service.grpc"
grpcproxy.Register(etcdClient, serverName, "127.0.0.1:8090", 3)

//在etcd中key的格式是：/etcd3_naming/user.service.grpc/127.0.0.1:8090

```

###发现并使用服务,实现负载均衡
```
discovery := etcdGrpcClient.DiscoveryService{
	Prefix:          "/etcd3_naming",
	ServiceNameList: []string{"user.service.grpc"},
	EtcdClient:      etcdClient,
}

discovery.InitGrpcService()
//
conn,_ := etcdGrpcClient.NewGrpcClientService("user.service.grpc")

//
c := pb.NewGrpcClient(conn)
msg := &pb.FunClientRequest{
    Name: "23",
}

result, err := c.FunClient(context.Background(), msg)

fmt.Println(result, err)

```
