package e3m

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"time"
)


const (
	e3_user_root = "root"
	client_dial_timeout = 5 * time.Second
)

type E3Config struct {
	Organization 	string
	Endpoints 		[]string
	RootPassword 	string
}
var (
	_config   E3Config
	_rootcli  *clientv3.Client
	_rootauth clientv3.Auth
)

func Setup(config E3Config)  {
	_config = config
	var err error
	if _rootcli, err = clientv3.New(clientv3.Config{
		Endpoints:   _config.Endpoints,
		DialTimeout: client_dial_timeout,
		Username:    e3_user_root,
		Password:    _config.RootPassword,
	}); err != nil {
		panic(err)
	}
	_rootauth = clientv3.NewAuth(_rootcli)
	if _, err := _rootauth.AuthEnable(context.Background()); err != nil {
		panic(err)
	}
}

func Client() *clientv3.Client {
	return _rootcli
}

func Auth() clientv3.Auth {
	return _rootauth
}

func Organization() string {
	return _config.Organization
}

func Endpoints() []string {
	return _config.Endpoints
}