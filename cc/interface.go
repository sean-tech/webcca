package cc

import (
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/gokit/requisition"
	"github.com/sean-tech/webkit/config"
	"sync"
)

/**参数：项目创建**/
type ProductCreateParameter struct {
	Product string 		`json:"product" validate:"required,gte=1"`
}

/**参数：所有项目获取**/
type ProductsGetParameter struct {
	UserName string 	`json:"userName" validate:"required,gte=1"`
}

/**参数：所有项目获取**/
type AppConfigPutParameter struct {
	Product string 				`json:"product" validate:"required,gte=1"`
	Config *config.AppConfig	`json:"config" validate:"required"`
}

type Worker struct {
	Module  	string 	`json:"module"`
	Ip 			string 	`json:"ip"`
	WorkerId 	int64 	`json:"workerId"`
}

type WorkerAddParameter struct {
	Product string 		`json:"product" validate:"required,gte=1"`
	Module  string 		`json:"module" validate:"required,gte=1"`
	Ip 		string		`json:"ip" validate:"required,ip"`
	WorkerId int64		`json:"workerId" validate:"min=0"`
}

type WorkersGetParameter struct {
	Product string 		`json:"product" validate:"required,gte=1"`
	Module  string 		`json:"module" validate:"required,gte=1"`
}

type WorkerDeleteParameter struct {
	Product string 		`json:"product" validate:"required,gte=1"`
	Module  string 		`json:"module" validate:"required,gte=1"`
	Ip 		string		`json:"ip" validate:"required,ip"`
}

/**参数：服务实例配置获取**/
type ModulesGetParameter struct {
	Product string 		`json:"product" validate:"required,gte=1"`
}



type ICCApi interface {
	ProductCreate(ctx *gin.Context)
	ProductsGet(ctx *gin.Context)
	AppConfigModify(ctx *gin.Context)
	AppConfigGet(ctx *gin.Context)
	WorkerAdd(ctx *gin.Context)
	WorkersGet(ctx *gin.Context)
	WorkerDelete(ctx *gin.Context)
	ModulesGet(ctx *gin.Context)
}

var (
	_apiOnce sync.Once
	_api     ICCApi
)

func Api() ICCApi {
	_apiOnce.Do(func() {
		_api = new(apiImpl)
	})
	return _api
}



const (
	_                                int = 0
	error_code_product_exist             = 12001
	error_code_product_not_exist         = 12002
	error_code_server_exist              = 12003
	error_code_server_not_exist          = 12004
	error_code_server_workerid_exist     = 12005
	error_code_update_failed             = 12007
	error_code_delete_failed             = 12008
)

func init() {
	requisition.SetMsgMap(requisition.LanguageZh, map[int]string{
		error_code_product_exist:             "项目已存在",
		error_code_product_not_exist:         "项目不存在",
		error_code_server_exist             : "实例已存在",
		error_code_server_not_exist         : "实例不存在",
		error_code_server_workerid_exist    : "workerId已存在",
		error_code_update_failed			: "修改失败",
		error_code_delete_failed			: "删除失败",
	})
	requisition.SetMsgMap(requisition.LanguageEn, map[int]string{
		error_code_product_exist:             "product exist already",
		error_code_product_not_exist:         "product not exist",
		error_code_server_exist             : "server instance exist already",
		error_code_server_not_exist         : "server instance not exist",
		error_code_server_workerid_exist    : "workerId exist already",
		error_code_update_failed			: "update failed",
		error_code_delete_failed			: "delete failed",
	})
}