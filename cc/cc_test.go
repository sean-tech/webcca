package cc

import (
	"cca/e3m"
	"fmt"
	"github.com/sean-tech/gokit/fileutils"
	"github.com/sean-tech/gokit/pubsub"
	"github.com/sean-tech/webkit/config"
	"github.com/sean-tech/webkit/database"
	"github.com/sean-tech/webkit/gohttp"
	"github.com/sean-tech/webkit/gorpc"
	"github.com/sean-tech/webkit/logging"
	"testing"
	"time"
)

var (
	e3config = e3m.E3Config{
		Organization: "sean-tech",
		Endpoints:    []string{"127.0.0.1:2379"},
		RootPassword: "etcd.user.root.pwd",
	}
	product = "ccatest"
	module = "testmodule"
	ip = "12.36.45.818"
)

func TestClearProducts(t *testing.T) {
	e3m.Setup(e3config)
	if err := deleteproducts(); err != nil {
		t.Error(err)
	} else {
		fmt.Println("products clear success")
	}
}

func TestPutConfig(t *testing.T) {
	e3m.Setup(e3config)
	user, err := e3m.RpcUser(product)
	if err != nil {
		t.Error(err)
	}
	err = PutAppConfig(product, false, &config.AppConfig{
		Log:     &logging.LogConfig{
			RunMode:     "test",
			LogSavePath: "/Users/sean/Desktop/",
			LogPrefix:   "e3mtest",
		},
		Http:    &gohttp.HttpConfig{
			RunMode:          "test",
			WorkerId:         0,
			HttpPort:         1024,
			ReadTimeout:      30 * time.Second,
			WriteTimeout:     30 * time.Second,
			CorsAllow:        false,
			CorsAllowOrigins: nil,
			RsaOpen: false,
			RsaMap:     nil,
		},
		Rpc:     &gorpc.RpcConfig{
			RunMode:              "test",
			RpcPort:              1023,
			RpcPerSecondConnIdle: 500,
			ReadTimeout:          30 * time.Second,
			WriteTimeout:         30 * time.Second,
			TokenSecret:          "knzxjcnjzbcj",
			TokenIssuer:          "sean-tech/cca/e3m",
			TlsOpen:              false,
			Tls:                  nil,
			WhiteListOpen:        false,
			WhiteListIps:         nil,
			Registry: &gorpc.EtcdRegistry{
				EtcdEndPoints:        e3m.Endpoints(),
				EtcdRpcBasePath:      user.Basepath,
				EtcdRpcUserName:      user.Username,
				EtcdRpcPassword:      user.Password,
			},
		},
		Mysql:   &database.MysqlConfig{
			WorkerId:    3,
			Type:        "mysql",
			User:        "root",
			Password:    "admin2018",
			Hosts: 		 map[int]string{0:"127.0.0.1:3306"},
			Name:        "etcd_center",
			MaxIdle:     30,
			MaxOpen:     30,
			MaxLifetime: 200 * time.Second,
		},
		Redis:   &database.RedisConfig{
			Host:        "127.0.0.1:6379",
			Password:    "",
			MaxIdle:     30,
			MaxActive:   30,
			IdleTimeout: 200 * time.Second,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("config put success")
}

func TestGetConfig(t *testing.T) {
	e3m.Setup(e3config)
	if cfg, err := GetAppConfig(product); err != nil {
		t.Error(err)
	} else {
		fmt.Println("config get success:")
		fmt.Println(cfg)
	}
}

func TestGetAllModules(t *testing.T) {
	e3m.Setup(e3config)
	if modules, err := GetAllModules(product); err != nil {
		t.Error(err)
	} else {
		fmt.Println("modules get success:")
		fmt.Println(modules)
	}
}

func TestPutWorkerId(t *testing.T) {
	e3m.Setup(e3config)
	if err := PutWorker(product, module, ip, 1); err != nil {
		t.Error(err)
	} else {
		fmt.Println("workerid put success")
	}
}

func TestGetAllWorkers(t *testing.T) {
	e3m.Setup(e3config)
	if workers, err := GetAllWorkers(product, module); err != nil {
		t.Error(err)
	} else {
		fmt.Println("workers get success:")
		fmt.Println(workers)
	}
}

func TestGetWorkerId(t *testing.T) {
	e3m.Setup(e3config)
	if workerid, err := GetWorker(product, module, ip); err != nil {
		t.Error(err)
	} else {
		fmt.Println("workerid get success: ", workerid)
	}
}

func TestDeleteWorkerId(t *testing.T) {
	e3m.Setup(e3config)
	if err := DeleteWorker(product, module, ip); err != nil {
		t.Error(err)
	} else {
		fmt.Println("workerid delete success")
	}
}

func TestNewproduct(t *testing.T) {
	e3m.Setup(e3config)
	if err := NewProduct("webkittest"); err != nil {
		t.Error(err)
	} else {
		fmt.Println("product add success")
	}
	if products, err := GetProducts(); err != nil {
		t.Error(err)
	} else {
		fmt.Println(products)
	}
}

func TestFileExist(t *testing.T) {
	src := "/Users/sean/Desktop"
	fmt.Println(fileutils.CheckExist(src))
}

func TestSliceRemoveObj(t *testing.T) {
	whiteList := []string{"a", "b", "c", "d", "e", "f"}
	var ip = "d"
	for idx, ipstr := range whiteList {
		if ip == ipstr {
			whiteList = append(whiteList[:idx], whiteList[idx+1:]...)
		}
	}
	fmt.Println(whiteList)
}

func TestPubsub(t *testing.T) {
	topic1 := "token"
	topic2 := "some"

	s1 := pubsub.SubscribeTopic(topic1, 1000)
	s2 := pubsub.SubscribeTopic(topic2, 1000)

	p1 := pubsub.NewPublisher(topic1, 10*time.Second)
	defer p1.Close()
	p2 := pubsub.NewPublisher(topic2, 10*time.Second)
	defer p2.Close()



	p1.Publish("hello, world")
	p1.Publish("hello, golang")
	p2.Publish("hello, p2")
	s2.Exit()
	p2.Publish("hello, p22")
	p1.Close()
	p2.Close()
	go func() {
		for Message := range s1.Message {
			fmt.Println("token ", Message)
		}
	}()
	go func() {
		for Message := range s2.Message {
			fmt.Println("some ", Message)
		}
	}()

	time.Sleep(3 * time.Second)
}