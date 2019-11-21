package etcdGrpcClient

import (
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func TestMain1(t *testing.T) {

	etcdConfig := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 3 * time.Second,
		//Username:  "",
		//Password:  "",
	}

	etcdClient, err := clientv3.New(etcdConfig)

	t.Log(etcdClient, err)
}
