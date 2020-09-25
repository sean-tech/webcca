package main

import (
	"cca/ca"
	"cca/cc"
	"cca/e3m"
	"context"
	"fmt"
	"github.com/sean-tech/gokit/requisition"
	"github.com/sean-tech/webkit/config"
	"github.com/sean-tech/webkit/gorpc"
	"github.com/sean-tech/webkit/logging"
	"github.com/smallnest/rpcx/server"
	"go.etcd.io/etcd/clientv3"
	"testing"
)

var (
	e3config = e3m.E3Config{
		Organization: "sean-tech",
		Endpoints:    []string{"127.0.0.1:2379"},
		RootPassword: "etcd.user.root.pwd",
	}
)
const (
	testproduct = "testp"
)

func TestShow(t *testing.T) {
	fmt.Println("show")
}

func TestClearAllPaths(t *testing.T) {
	e3m.Setup(e3config)
	if _, err := e3m.Client().Delete(context.Background(), "/sean.tech/", clientv3.WithPrefix()); err != nil {
		t.Error(err)
	} else {
		fmt.Println("path clear success")
	}
}

func TestCCServing(t *testing.T) {
	logging.Setup(logging.LogConfig{
		RunMode:     "debug",
		LogSavePath: "/Users/Sean/Desktop/",
		LogPrefix:   "cca",
	})
	e3m.Setup(e3config)
	ca.Setup(ca.CASetting{
		CompanyName: "sean-tech",
		CASavePath:  "/Users/sean/Desktop/",
	})
	cc.Serving(9966)
	if err := cc.NewProduct(testproduct); err != nil {
		t.Error(err)
	}
	if err := cc.PutWorker(testproduct, "service1", "192.168.1.21", 0); err != nil {
		t.Error(err)
	}
	if err := cc.PutWorker(testproduct, "service2", "192.168.1.21", 0); err != nil {
		t.Error(err)
	}
	cfg, err := cc.GetAppConfig(testproduct)
	if err != nil {
		t.Error(err)
	}
	cfg.Rpc.TlsOpen = true
	if err := cc.PutAppConfig(testproduct, false, cfg); err != nil {
		t.Error(err)
	}
	select {

	}
}

func TestRpcServer1(t *testing.T) {

	config.TestCCAddress(testproduct, "service1", 10102, 10101, "/Users/Sean/Desktop/", "192.168.1.21:9966", func(appConfig *config.AppConfig) {
		logging.Setup(*appConfig.Log)
		gorpc.ServerServe(*appConfig.Rpc, logging.Logger(), func(server *server.Server) {
			server.RegisterName("s1", Service1, "")
		})
	})
	select {

	}
}

func TestRpcServer2(t *testing.T) {
	config.TestCCAddress(testproduct, "service2", 10104, 10103, "/Users/Sean/Desktop/", "192.168.1.21:9966", func(appConfig *config.AppConfig) {
		logging.Setup(*appConfig.Log)
		gorpc.ServerServe(*appConfig.Rpc, logging.Logger(), func(server *server.Server) {
			server.RegisterName("s2", Service2, "")
		})
	})

	ctx := requisition.NewRequestionContext(context.Background())
	requisition.GetRequisition(ctx).RequestId = 10001
	requisition.GetRequisition(ctx).UserName = "test"
	requisition.GetRequisition(ctx).UserId = 1230001
	var parameter = "this is call to s1.hello"
	var reply = ""
	if err := gorpc.Call("s1", "service1", ctx, "Hello",  &parameter, &reply); err != nil {
		fmt.Printf("err:%s", err.Error())
	} else  {
		fmt.Printf("reply:%s", reply)
	}
}


type service1Impl struct {
}
var Service1 = &service1Impl{}

func (this *service1Impl) Hello(ctx context.Context, parameter *string, reply *string) error {
	*reply = *parameter + "---hello"
	return nil
}

type service2Impl struct {
}
var Service2 = &service2Impl{}

func (this *service2Impl) ShowInfo(ctx context.Context, parameter *string, reply *string) error {
	*reply = *parameter + "---showinfo"
	return nil
}

