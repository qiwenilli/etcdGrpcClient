package etcdGrpcClient

import (
	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/naming"

	//
	"google.golang.org/grpc/balancer/roundrobin"
	// "google.golang.org/grpc/balancer/weightedroundrobin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"

	//
	"log"
	"time"
)

//创建一个grpc client 服务链接
func NewGrpcClientService(serviceName string) (*grpc.ClientConn, error) {
	//使用轮询机制进行负载, 默认2秒超时
	return grpc.Dial(rpcScheme+":///"+serviceName, grpc.WithBalancerName(roundrobin.Name), grpc.WithBlock(), grpc.WithInsecure(), grpc.WithTimeout(2*time.Second))
}

var (
	rpcScheme string
)

//发现服务
type DiscoveryService struct {
	Prefix          string
	ServiceNameList []string
	GRPCResolver    *naming.GRPCResolver
	EtcdClient      *clientv3.Client
}

func (this *DiscoveryService) InitGrpcService() {
	//
	rpcScheme = "qiwen"
	//通过前缀获得同一服务注册的地址
	this.GRPCResolver = &naming.GRPCResolver{Client: this.EtcdClient}
	//
	var addrsStore = make(map[string][]string, 2)
	for _, serviceName := range this.ServiceNameList {
		addrsStore[serviceName] = this.getUpdateAddr(serviceName)
	}

	resolver.Register(&customerResolverBuilder{
		CustomerScheme: rpcScheme,
		AddrsStore:     addrsStore,
	})
}

func (this *DiscoveryService) getUpdateAddr(serviceName string) (addrList []string) {

	//通过前缀获得同一服务注册的地址
	watcher, err := this.GRPCResolver.Resolve(this.Prefix + "/" + serviceName)
	defer watcher.Close()
	//
	if err != nil {
		log.Fatalf("failed to resolve %q (%v)", this.Prefix+"/"+serviceName, err)
		return
	}

	//目前只next一次; 可以重试
	updateList, err := watcher.Next()
	if err != nil {
		return
	}

	// var addrList []string
	for _, update := range updateList {
		addrList = append(addrList, update.Addr)
	}

	return
}

//自定义 ResolverBuilder & Resolver
type customerResolverBuilder struct {
	CustomerScheme string              //"qiwen"
	AddrsStore     map[string][]string // "user.service.grpc": {"192.168.1.176:8091", "192.168.1.176:8090"},
}

func (this *customerResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r := &customizeResolver{
		target:     target,
		cc:         cc,
		addrsStore: this.AddrsStore,
	}
	r.start()
	return r, nil
}
func (this *customerResolverBuilder) Scheme() string { return this.CustomerScheme }

type customizeResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (this *customizeResolver) start() {
	addrStrs := this.addrsStore[this.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	this.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*customizeResolver) ResolveNow(o resolver.ResolveNowOption) {}
func (*customizeResolver) Close()                                 {}
