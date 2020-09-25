package main

import (
	"cca/admin"
	"cca/ca"
	"cca/cc"
	"cca/rpcxservice"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/gokit/storage"
	"github.com/sean-tech/webkit/auth"
	"github.com/sean-tech/webkit/config"
	"github.com/sean-tech/webkit/gohttp"
	"github.com/sean-tech/webkit/logging"
	"github.com/storyicon/grbac"
	"runtime"
	"sync"
)

func main()  {
	// cmd set
	cmdset("sean.tech", "127.0.0.1:2379", "etcd.user.root.pwd", "/Users/Sean/Desktop/")
	// concurrent
	runtime.GOMAXPROCS(runtime.NumCPU())
	// config
	config.Setup("", "cca", 9965, 9966, "/Users/Sean/Desktop/", appConfig, func(appConfig *config.AppConfig) {
		// logging set
		logging.Setup(*appConfig.Log)
		// e3m start
		e3mStartWithCmd()
		// ca start
		caStartWithCmd()
		// cc rpc start
		cc.Serving(appConfig.Rpc.RpcPort)
		// auth start
		auth.Setup(authConfig, storage.Hash())
		// server start
		gohttp.HttpServerServe(*appConfig.Http, logging.Logger(), RegisterApi)
	})
}

func RegisterApi(engine *gin.Engine)  {
	engine.Use(cors.Default())
	var tokenHandler =  gohttp.InterceptToken(func(ctx context.Context, token string) (userId uint64, userName, role, key string, err error) {
		var parameter = &auth.AccessTokenAuthParameter{AccessToken: token}
		var accessTokenItem = new(auth.TokenItem)
		if err := auth.Service().AccessTokenAuth(ctx, parameter, accessTokenItem); err != nil {
			return 0, "", "", "", err
		}
		return accessTokenItem.UserId, accessTokenItem.UserName,  accessTokenItem.Role, accessTokenItem.Key, nil
	})
	apiv1 := engine.Group("api/v1/")
	{
		auth_v1 := apiv1.Group("auth")
		{
			auth_v1.POST("refresh", auth.Api().AuthRefresh)
		}
		role_v1 := apiv1.Group("role").Use(tokenHandler)
		{
			role_v1.POST("get", admin.Api().RoleGet)
			role_v1.POST("list", admin.Api().RoleGetAll)
			role_v1.POST("add", admin.Api().RoleAdd)
			role_v1.POST("update", admin.Api().RoleUpdate)
			role_v1.POST("delete", admin.Api().RoleDelete)
		}
		admin_v1 := apiv1.Group("admin")
		{
			admin_v1.POST("check", admin.Api().AdminCheck)
			admin_v1.POST("login", admin.Api().AdminLogin)
			admin_v1.Use(tokenHandler)
			{
				admin_v1.POST("get", admin.Api().AdminGet)
				admin_v1.POST("list", admin.Api().AdminGetList)
				admin_v1.POST("add", admin.Api().AdminAdd)
				admin_v1.POST("delete", admin.Api().AdminDelete)
				admin_v1.POST("enabled", admin.Api().AdminEnable)
			}
		}
		ca_v1 := apiv1.Group("ca").Use(tokenHandler)
		{
			ca_v1.POST("rsaClientFileVersions", ca.Api.RsaClientVersionsGet)
			ca_v1.POST("rsaClientFiles", ca.Api.RsaClientVersionFilesGet)
			ca_v1.POST("rsaClientFileDownload", ca.Api.RsaClientFileDownload)
			ca_v1.POST("rsaNewCertFile", ca.Api.RsaNewCertFile)
		}
		cc_v1 := apiv1.Group("cc").Use(tokenHandler)
		{
			cc_v1.POST("productCreate", cc.Api().ProductCreate)
			cc_v1.POST("products", cc.Api().ProductsGet)
			cc_v1.POST("config", cc.Api().AppConfigGet)
			cc_v1.POST("configModify", cc.Api().AppConfigModify)
			cc_v1.POST("workerAdd", cc.Api().WorkerAdd)
			cc_v1.POST("workers", cc.Api().WorkersGet)
			cc_v1.POST("workerDelete", cc.Api().WorkerDelete)
			cc_v1.POST("modules", cc.Api().ModulesGet)
		}
		services_v1 := apiv1.Group("services").Use(tokenHandler)
		{
			services_v1.POST("get", rpcxservice.Api().ServicesGet)
			services_v1.POST("active", rpcxservice.Api().ServiceActive)
			services_v1.POST("metadata", rpcxservice.Api().ServiceMetadataUpdate)
		}
	}
}

func LoadAuthorizationRules() (rules grbac.Rules, err error) {
	// 在这里实现你的逻辑
	// ...
	// 你可以从数据库或文件加载授权规则
	// 但是你需要以 grbac.Rules 的格式返回你的身份验证规则
	// 提示：你还可以将此函数绑定到golang结构体
	var id_level = 5005
	var id_low = 1000
	rules = []*grbac.Rule{
		&grbac.Rule{ID:5001, Resource:&grbac.Resource{Host:"*", Path:"/api/v1/admin/check", Method: "POST"}, Permission:&grbac.Permission{AuthorizedRoles:[]string{"*"}, AllowAnyone:true}},
		&grbac.Rule{ID:5002, Resource:&grbac.Resource{Host:"*", Path:"/api/v1/admin/login", Method: "POST"}, Permission:&grbac.Permission{AuthorizedRoles:[]string{"*"}, AllowAnyone:true}},
		&grbac.Rule{ID:5003, Resource:&grbac.Resource{Host:"*", Path:"/api/v1/admin/get", Method: "POST"}, Permission:&grbac.Permission{AuthorizedRoles:[]string{"*"}, AllowAnyone:true}},
		&grbac.Rule{ID:5004, Resource:&grbac.Resource{Host:"*", Path:"**", Method: "*"}, Permission:&grbac.Permission{AuthorizedRoles:[]string{admin.ROLE_SUPER}, AllowAnyone:false}},
	}
	roles, err := admin.RoleGetAll()
	if err != nil {
		return nil, err
	}
	var apiRolesMap = make(map[string][]string)
	var productRolesMap = make(map[string][]string)
	for _, role := range roles {
		for _, api := range role.AllowApis {
			if _, ok := apiRolesMap[api]; !ok {
				apiRolesMap[api] = []string{role.RoleName}
			} else {
				apiRolesMap[api] = append(apiRolesMap[api], role.RoleName)
			}
		}
		for _, product := range role.AllowProducts {
			if _, ok := productRolesMap[product]; !ok {
				productRolesMap[product] = []string{role.RoleName}
			} else {
				productRolesMap[product] = append(productRolesMap[product], role.RoleName)
			}
		}
	}
	for api, roleNames := range apiRolesMap {
		rule := &grbac.Rule{ID: id_level, Resource:&grbac.Resource{Host: "*", Path:api, Method: "*"}, Permission:&grbac.Permission{AuthorizedRoles: roleNames, AllowAnyone:false}}
		rules = append(rules, rule)
		id_level += 1
	}
	for product, roleNames := range productRolesMap {
		id_low += 1
		var path = "**/" + product
		rule := &grbac.Rule{ID: id_low, Resource:&grbac.Resource{Host: "*", Path:path, Method: "*"}, Permission:&grbac.Permission{AuthorizedRoles: roleNames, AllowAnyone:false}}
		rules = append(rules, rule)
	}
	return rules, nil
}

var ApiMap = map[string]string{
	"管理员列表" : "api/v1/admin/list",
	"管理员创建" : "api/v1/admin/add",
	"管理员禁用" : "api/v1/admin/enabled",
	"管理员删除" : "api/v1/admin/delete",

	"角色列表" : "api/v1/role/list",
	"角色信息" : "api/v1/role/get",
	"角色添加" : "api/v1/role/add",
	"角色删除" : "api/v1/role/delete",

	"证书权限" : "api/v1/ca/*",
	"项目列表" : "api/v1/cc/products",
	"项目创建" : "api/v1/cc/productCreate",
	"配置信息" : "api/v1/cc/config",
	"配置修改" : "api/v1/cc/configModify",
	"实例列表" : "api/v1/cc/workers",
	"实例添加" : "api/v1/cc/workerAdd",
	"实例删除" : "api/v1/cc/workerDelete",
	"服务列表" : "api/v1/services/get",
	"服务状态编辑" : "api/v1/services/active",
	"服务数据编辑" : "api/v1/services/metadata",
}

var __reversed_api_map map[string]string
var __reverse_lock sync.Mutex
func ReversedApiMap() map[string]string {
	if __reversed_api_map == nil {
		__reverse_lock.Lock()
		__reversed_api_map = make(map[string]string, len(ApiMap))
		for k, v := range ApiMap {
			__reversed_api_map[v] = k
		}
		__reverse_lock.Unlock()
	}
	return __reversed_api_map
}

