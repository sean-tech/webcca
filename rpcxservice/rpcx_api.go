package rpcxservice

import (
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/gokit/validate"
	"github.com/sean-tech/webkit/gohttp"
)

type apiImpl struct {}

func (this *apiImpl) ServicesGet(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter ServiceGetParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if services, err := reg.fetchServices(parameter.Product); err != nil {
		g.ResponseError(err)
	} else {
		g.ResponseData(services)
	}
}


func (this *apiImpl) ServiceActive(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter ServiceActiveParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var err error
	switch parameter.Actived {
	case 0:
		err = reg.activateService(parameter.Product, parameter.Name, parameter.Address)
	case 1:
		err = reg.deactivateService(parameter.Product, parameter.Name, parameter.Address)
	}
	if err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData("")
}

func (this *apiImpl) ServiceMetadataUpdate(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter ServiceMetaDataUpdateParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := reg.updateMetadata(parameter.Product, parameter.Name, parameter.Address, parameter.MetaData); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData("")
}