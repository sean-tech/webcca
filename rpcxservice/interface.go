package rpcxservice

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
)

type RpcService struct {
	ID       string	`json:"id"`
	Name     string	`json:"name"`
	Address  string	`json:"address"`
	Metadata string	`json:"metadata"`
	State    string	`json:"state"`
	Group    string	`json:"group"`
}

type ServiceGetParameter struct {
	Product string 		`json:"product" validate:"required,gte=1"`
}

type ServiceActiveParameter struct {
	Product string 		`json:"product" validate:"required,gte=1"`
	Name 	string		`json:"name" validate:"required,gte=1"`
	Address string		`json:"address" validate:"required,gt=1"`
	Actived int			`json:"actived" validate:"oneof=0 1"`
}

type ServiceMetaDataUpdateParameter struct {
	Product string 		`json:"product" validate:"required,gte=1"`
	Name string			`json:"name" validate:"required,gte=1"`
	Address string		`json:"address" validate:"required,gt=1"`
	MetaData string		`json:"metaData" validate:"gte=0"`
}

type IRpcxApi interface {
	ServicesGet(ctx *gin.Context)
	ServiceActive(ctx *gin.Context)
	ServiceMetadataUpdate(ctx *gin.Context)
}

var (
	_apiOnce     sync.Once
	_apiInstance     IRpcxApi
)

func Api() IRpcxApi {
	_apiOnce.Do(func() {
		_apiInstance = new(apiImpl)
	})
	return _apiInstance
}

type Registry interface {
	initRegistry()
	fetchServices(product string)([]*RpcService, error)
	deactivateService(product, name, address string) error
	activateService(product, name, address string) error
	updateMetadata(product, name, address string, metadata string) error
}

// Config parameters
var serverConfig = Configuration{
	RegistryType:   "etcd",
	RegistryURLs:    []string{"localhost:2379"},
}
var reg Registry

func Setup(endpoints []string)  {
	serverConfig.RegistryURLs = endpoints
	switch serverConfig.RegistryType {
	case "etcd":
		reg = &EtcdRegistry{}
	default:
		fmt.Printf("unsupported registry: %s\n", serverConfig.RegistryType)
		os.Exit(2)
	}
	reg.initRegistry()
}

// Configuration is configuration strcut refects the config.json
type Configuration struct {
	RegistryType   string 	`json:"registry_type"`
	RegistryURLs   []string `json:"registry_urls"`
}
