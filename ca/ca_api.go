package ca

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/gokit/validate"
	"github.com/sean-tech/webkit/gohttp"
	"net/http"
)

var password = GenerateKey([]byte("ajkhzjc!bzb323654zklsjfd4zxc6"))

type RsaClientVersionsGetParamter struct {
	Product string	`json:"product" validate:"required,gte=1"`
}

type RsaClientFilesGetParamter struct {
	Product string	`json:"product" validate:"required,gte=1"`
	Version string	`json:"version" validate:"required,gte=1"`
}

type RSAClientFileDownloadParamter struct {
	Product 	string	`json:"product" validate:"required,gte=1"`
	FileName 	string	`json:"fileName" validate:"required,gte=1"`
}

type ICaApi interface {
	RsaClientVersionsGet(ctx *gin.Context)
	RsaClientVersionFilesGet(ctx *gin.Context)
	RsaClientFileDownload(ctx *gin.Context)
	RsaNewCertFile(ctx *gin.Context)
}

type apiImpl struct {
}
var Api ICaApi = &apiImpl{}

func (this *apiImpl) RsaClientVersionsGet(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter RsaClientVersionsGetParamter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if files, err := GetRSAClientVersions(parameter.Product); err != nil {
		g.ResponseError(err)
	} else {
		g.ResponseData(files)
	}
}

func (this *apiImpl) RsaClientVersionFilesGet(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter RsaClientFilesGetParamter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if files, err := GetRSAClientVersionFiles(parameter.Product, parameter.Version); err != nil {
		g.ResponseError(err)
	} else {
		g.ResponseData(files)
	}
}

func (this *apiImpl) RsaClientFileDownload(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter RSAClientFileDownloadParamter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if data, err := GetRsaKeyPairFile(parameter.Product, parameter.FileName); err != nil {
		g.ResponseError(err)
	} else {
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Header("Content-Disposition", "attachment; filename=" + parameter.FileName)
		ctx.Header("Content-Type", "application/text/plain")
		ctx.Header("Accept-Length", fmt.Sprintf("%d", len(data)))
		ctx.Writer.Write(data)
	}
}

func (this *apiImpl) RsaNewCertFile(ctx *gin.Context) {
	g := gohttp.Gin{Ctx:ctx}
	var parameter RsaClientFilesGetParamter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := validate.ValidateParameter(parameter); err != nil {
		g.ResponseError(err)
		return
	}
	if err := NewRsaKeyPair(parameter.Product, parameter.Version, password); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(nil)
}