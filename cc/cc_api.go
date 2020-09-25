package cc

import (
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/gokit/foundation"
	"github.com/sean-tech/gokit/validate"
	"github.com/sean-tech/webkit/config"
	"github.com/sean-tech/webkit/gohttp"
)

type apiImpl struct {
}

func (this *apiImpl) ProductCreate(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter ProductCreateParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := NewProduct(parameter.Product); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(nil)
}


func (this *apiImpl) ProductsGet(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter ProductsGetParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var products []string; var err error
	if products, err = GetProducts(); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(products)
}

func (this *apiImpl) AppConfigModify(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter = new(AppConfigPutParameter)
	if err := g.BindParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if parameter.Config.Http != nil && parameter.Config.Http.RunMode == "" {
		parameter.Config.Http.RunMode = foundation.RUN_MODE_RELEASE
	}
	if parameter.Config.Rpc != nil && parameter.Config.Rpc.RunMode == "" {
		parameter.Config.Rpc.RunMode = foundation.RUN_MODE_RELEASE
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := PutAppConfig(parameter.Product, true, parameter.Config); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(nil)
}

func (this *apiImpl) AppConfigGet(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter = new(ProductCreateParameter)
	if err := g.BindParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var cfg *config.AppConfig; var err error
	if cfg, err = GetAppConfig(parameter.Product); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(cfg)
}

func (this *apiImpl) WorkerAdd(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter WorkerAddParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := PutWorker(parameter.Product, parameter.Module, parameter.Ip, parameter.WorkerId); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(nil)
}

func (this *apiImpl) WorkersGet(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter WorkersGetParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if parameter.Module != "all" {
		if workers, err := GetAllWorkers(parameter.Product, parameter.Module); err != nil {
			g.ResponseError(err)
		} else {
			g.ResponseData(workers)
		}
		return
	}
	// all workers
	modules, err := GetAllModules(parameter.Product)
	if err != nil {
		g.ResponseError(err)
		return
	}
	var workers []Worker
	for _, module := range modules {
		moduleWorkers, err := GetAllWorkers(parameter.Product, module)
		if err != nil {
			g.ResponseError(err)
			return
		}
		workers = append(workers, moduleWorkers...)
	}
	g.ResponseData(workers)
}

func (this *apiImpl) WorkerDelete(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter WorkerDeleteParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := DeleteWorker(parameter.Product, parameter.Module, parameter.Ip); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(nil)
}

func (this *apiImpl) ModulesGet(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter ModulesGetParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if modules, err := GetAllModules(parameter.Product); err != nil {
		g.ResponseError(err)
	} else {
		g.ResponseData(modules)
	}
}
