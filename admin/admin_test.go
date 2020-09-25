package admin

import (
	"cca/e3m"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/gokit/storage"
	"github.com/sean-tech/webkit/auth"
	"github.com/sean-tech/webkit/config"
	"github.com/sean-tech/webkit/gohttp"
	"github.com/sean-tech/webkit/logging"
	"github.com/storyicon/grbac"
	"runtime"
	"testing"
	"time"
)

var testconfig = &config.AppConfig{
	Log: &logging.LogConfig{
		RunMode:     "debug",
		LogSavePath: "/Users/sean/Desktop/",
		LogPrefix:   "auth",
	},
	Http:    &gohttp.HttpConfig{
		RunMode:      "debug",
		WorkerId:     3,
		HttpPort:     9965,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		RsaOpen: false,
		RsaMap:     nil,
		CorsAllow: true,
		CorsAllowOrigins: nil,
	},
}

func TestAdminRun(t *testing.T) {
	e3m.Setup(e3config)

	// concurrent
	runtime.GOMAXPROCS(runtime.NumCPU())
	// log start
	logging.Setup(*testconfig.Log)
	// auth start
	auth.Setup(auth.AuthConfig{
		WorkerId:                workerid,
		TokenSecret:             "szksangknxucv!20@km",
		TokenIssuer:             "sean-tech/cca",
		RefreshTokenExpiresTime: 3 * 24 * time.Hour,
		AccessTokenExpiresTime:  30 * time.Minute,
		AuthCode:                AUTH_CODE,
	}, storage.Hash())
	// server start
	gohttp.HttpServerServe(*testconfig.Http, logging.Logger(), RegisterApi)
}

func RegisterApi(engine *gin.Engine)  {
	apiv1 := engine.Group("api/v1/")
	{
		role_v1 := apiv1.Group("role")
		{
			role_v1.POST("get", Api().RoleGet)
			role_v1.POST("list", Api().RoleGetAll)
			role_v1.POST("add", Api().RoleAdd)
			role_v1.POST("delete", Api().RoleDelete)
		}
		user_v1 := apiv1.Group("user")
		{
			user_v1.POST("login", Api().AdminLogin)
			user_v1.POST("check", Api().AdminCheck)
			user_v1.POST("get", Api().AdminGet)
			user_v1.POST("list", Api().AdminGetList)
			user_v1.POST("add", Api().AdminAdd)
			user_v1.POST("delete", Api().AdminDelete)
			user_v1.POST("enabled", Api().AdminEnable)
			user_v1.POST("resful/:userId", userget)
		}
		other_v1 := apiv1.Group("other").Use(gohttp.Authorization(LoadAuthorizationRules))
		{
			other_v1.POST("resful/:userId", userget)
		}
	}
}

func userget(ctx *gin.Context) {
	g := gohttp.Gin{Ctx: ctx}
	type UserGetParameter struct {
		Real	string	`json:"real" validate:"required, gte=1"`
	}
	var parameter = new(UserGetParameter)
	if err := g.BindParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var resp = make(map[string]string)
	resp["userId"] = ctx.Params.ByName("userId")
	resp["realValue"] = parameter.Real
	g.ResponseData(resp)
}

func LoadAuthorizationRules() (rules grbac.Rules, err error) {

	// 在这里实现你的逻辑
	// ...
	// 你可以从数据库或文件加载授权规则
	// 但是你需要以 grbac.Rules 的格式返回你的身份验证规则
	// 提示：你还可以将此函数绑定到golang结构体
	rules = []*grbac.Rule{
		&grbac.Rule{ID:1001, Resource:&grbac.Resource{Host:"localhost:9965", Path:"/api/v1/user/**", Method: "POST"}, Permission:&grbac.Permission{AuthorizedRoles:[]string{"superadmin"}, AllowAnyone:false}},
		&grbac.Rule{ID:1002, Resource:&grbac.Resource{Host:"localhost:9965", Path:"**/123456", Method: "POST"}, Permission:&grbac.Permission{AuthorizedRoles:[]string{"superadmin"}, AllowAnyone:false}},
	}
	err = nil
	return
}