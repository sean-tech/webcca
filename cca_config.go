package main

import (
	"cca/ca"
	"cca/e3m"
	"cca/rpcxservice"
	"flag"
	"github.com/sean-tech/gokit/foundation"
	"github.com/sean-tech/webkit/auth"
	"github.com/sean-tech/webkit/config"
	"github.com/sean-tech/webkit/gohttp"
	"github.com/sean-tech/webkit/gorpc"
	"strings"
	"time"
)

const (
	organization_usage = "please use -organization to pointing at company name."
	etcd_endpoints_usage = "please use -etcd.endpoints to pointing at etcd endpoints, use ',' to sep."
	etcd_rootpwd_usage = "please use -etcd.rootpwd to pointing at etcd password of root user."
	casavepath_usage = "please use -casavepath to pointing at ca file saved path."
)

var (
	organization *string
	etcd_endpoints *string
	etcd_rootpwd *string
	casavepath *string
)

func cmdset(company, etcdEndpoints, etcdRootPwd, caSavePath string)  {
	// organization
	organization = flag.String("organization", company, organization_usage)
	// etcd.endpoints
	etcd_endpoints = flag.String("etcd.endpoints", etcdEndpoints, etcd_endpoints_usage)
	// etcd.rootpwd
	etcd_rootpwd = flag.String("etcd.rootpwd", etcdRootPwd, etcd_rootpwd_usage)
	// casavepath
	casavepath = flag.String("casavepath", caSavePath, casavepath_usage)
}

func e3mStartWithCmd() {
	*organization = strings.ReplaceAll(*organization, " ", "")
	*etcd_endpoints = strings.ReplaceAll(*etcd_endpoints, " ", "")
	*etcd_rootpwd = strings.ReplaceAll(*etcd_rootpwd, " ", "")
	if organization == nil || *organization == "" {
		panic(organization_usage)
	}
	if etcd_endpoints == nil || *etcd_endpoints == "" {
		panic(etcd_endpoints_usage)
	}
	if etcd_rootpwd == nil || *etcd_rootpwd == "" {
		panic(etcd_rootpwd_usage)
	}
	endpoints := strings.Split(*etcd_endpoints, ",")
	e3m.Setup(e3m.E3Config{
		Organization: *organization,
		Endpoints:    endpoints,
		RootPassword: *etcd_rootpwd,
	})
	rpcxservice.Setup(e3m.Endpoints())
}

func caStartWithCmd() {
	*organization = strings.ReplaceAll(*organization, " ", "")
	*casavepath = strings.ReplaceAll(*casavepath, " ", "")
	if organization == nil || *organization == "" {
		panic(organization_usage)
	}
	if casavepath == nil || *casavepath == "" {
		panic(casavepath_usage)
	}
	ca.Setup(ca.CASetting{
		CompanyName: *organization,
		CASavePath:  *casavepath,
	})
}

var authConfig = auth.AuthConfig{
	WorkerId:                0,
	TokenSecret:             "jkznqkjn!12k@20200824",
	TokenIssuer:             "sean.tech/webkit/cca",
	RefreshTokenExpiresTime: time.Hour * 3,
	AccessTokenExpiresTime:  time.Minute * 15,
	AuthCode: 				 "kzncknqnzmc9843yu!z#xv123",
}

var appConfig = &config.AppConfig{
	Http:    &gohttp.HttpConfig{
		RunMode:      "release",
		WorkerId:     0,
		HttpPort:     9965,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		RsaOpen: false,
		RsaMap:     nil,
	},
	Rpc:     &gorpc.RpcConfig{
		RunMode:              "release",
		RpcPort:              9966,
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
}