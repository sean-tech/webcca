package cc

import (
	"cca/e3m"
	"github.com/sean-tech/gokit/foundation"
	"github.com/sean-tech/gokit/pubsub"
	"github.com/sean-tech/webkit/config"
	"github.com/sean-tech/webkit/database"
	"github.com/sean-tech/webkit/gohttp"
	"github.com/sean-tech/webkit/gorpc"
	"github.com/sean-tech/webkit/logging"
	"strings"
	"sync"
	"time"
)

const (
	module_admin = "admin"
	pubsub_topic_worker_add = "/cca/cc/worker/add"
	pubsub_topic_worker_del = "/cca/cc/worker/del"
)

var (
	whitelistips []string
	wlock sync.RWMutex
	workeradd_pub = pubsub.NewPublisher(pubsub_topic_worker_add, 10*time.Second)
	workerdel_pub = pubsub.NewPublisher(pubsub_topic_worker_del, 10*time.Second)
	workeradd_sub = pubsub.SubscribeTopic(pubsub_topic_worker_add, 100)
	workerdel_sub = pubsub.SubscribeTopic(pubsub_topic_worker_del, 100)
)

func Serving(port int)  {
	go config.ConfigCernterServing(new(center), port, ipFitter)
}

type center struct {}

func (this *center) AppConfigLoad(worker *config.Worker, appcfg *config.AppConfig) error {
	appconfig, err := GetAppConfig(worker.Product)
	if err != nil {
		return err
	}
	// get set workerid
	w, err := GetWorker(worker.Product, worker.Module, worker.Ip)
	if err != nil {
		return err
	}
	appconfig.Http.WorkerId = w.WorkerId
	appconfig.Mysql.WorkerId = w.WorkerId
	// get set http rsa
	if appconfig.Http.RsaOpen {
		if rsaConfigMap, err := GetRsaConfigMap(worker.Product); err != nil {
			logging.Error(err)
		} else {
			appcfg.Http.RsaMap = rsaConfigMap
		}
	}
	// get set rpc tls
	if appconfig.Rpc.TlsOpen {
		if tlscfg, err := NewServiceTLSCert(worker.Product, worker.Module); err != nil {
			logging.Error(err)
		} else {
			appconfig.Rpc.Tls = tlscfg
		}
	}
	// get set rpc registry
	rpcuser, err := e3m.RpcUser(worker.Product)
	if err != nil {
		logging.Error(err)
	} else {
		appconfig.Rpc.Registry = &gorpc.EtcdRegistry{
			EtcdEndPoints:   e3m.Endpoints(),
			EtcdRpcBasePath: rpcuser.Basepath,
			EtcdRpcUserName: rpcuser.Username,
			EtcdRpcPassword: rpcuser.Password,
		}
	}
	// get set config etcd user
	var etcdAuthUser *e3m.AuthUser
	sepModule := strings.Split(worker.Module, ".")
	if len(sepModule) == 2 && sepModule[1] == module_admin {
		etcdAuthUser, err = e3m.ConfigModuleUser(worker.Product, worker.Module, e3m.AuthPermReadWrite)
	} else if worker.Module == module_admin {
		etcdAuthUser, err = e3m.ConfigProductUser(worker.Product, e3m.AuthPermReadWrite)
	} else {
		etcdAuthUser, err = e3m.ConfigModuleUser(worker.Product, worker.Module, e3m.AuthPermReadOnly)
	}
	if err != nil {
		return err
	}
	appconfig.CE = &config.ConfigEtcd{
		EtcdEndPoints:      e3m.Endpoints(),
		EtcdConfigBasePath: etcdAuthUser.Basepath,
		EtcdConfigUserName: etcdAuthUser.Username,
		EtcdConfigPassword: etcdAuthUser.Password,
	}
	*appcfg = *appconfig
	return appcfg.Validate()
}

func ipFitter(clientIp string) bool {
	if whitelistips == nil {
		wlock.Lock()
		whitelistips = whiteListLoad()
		wlock.Unlock()
	}
	if whitelistips == nil || len(whitelistips) == 0 {
		return true
	}
	for _, ip := range whitelistips {
		if clientIp == ip {
			return true
		}
	}
	return false
}

func whiteListLoad() []string {
	var whiteList []string
	products, err := GetProducts()
	if err != nil {
		panic(err)
	}
	for _, product := range products {
		modules, err := GetAllModules(product)
		if err != nil {
			logging.Error(err)
			continue
		}
		for _, module := range modules {
			workers, err := GetAllWorkers(product, module)
			if err != nil {
				logging.Error(err)
				continue
			}
			for _, worker := range workers {
				whiteList = append(whiteList, worker.Ip)
			}
		}
	}
	whiteListSubscribing()
	return whiteList
}

func whiteListSubscribing()  {
	go func() {
		for message := range workeradd_sub.Message {
			if workerIp, ok := message.(string); ok {
				wlock.Lock()
				whitelistips = append(whitelistips, workerIp)
				wlock.Unlock()
			}
		}
	}()
	go func() {
		for message := range workerdel_sub.Message {
			if workerIp, ok := message.(string); ok {
				for idx, ipstr := range whitelistips {
					if workerIp == ipstr {
						wlock.Lock()
						whitelistips = append(whitelistips[:idx], whitelistips[idx+1:]...)
						wlock.Unlock()
					}
				}
			}
		}
	}()
}

var appconfig_template = &config.AppConfig{
	Log: nil,
	Http:    &gohttp.HttpConfig{
		RunMode:      "release",
		WorkerId:     0,
		HttpPort:     9022,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		RsaOpen: false,
		RsaMap:     nil,
	},
	Rpc:     &gorpc.RpcConfig{
		RunMode:              "release",
		RpcPort:              9021,
		RpcPerSecondConnIdle: 500,
		ReadTimeout:          60 * time.Second,
		WriteTimeout:         60 * time.Second,
		TokenSecret:          foundation.RandString(12),
		TokenIssuer:          "/sean-tech/webkit/auth",
		TlsOpen:              false,
		Tls:                  nil,
		WhiteListOpen:        false,
		WhiteListIps:         []string{"127.0.0.1"},
		Registry: nil,
	},
	Mysql:   &database.MysqlConfig{
		WorkerId:    0,
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
	CE: nil,
}