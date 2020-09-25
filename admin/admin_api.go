package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/webkit/gohttp"
)

type apiImpl struct{

}

/**
 * Api-添加角色
 */
func (this *apiImpl) RoleAdd(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter RoleAddParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var result = new(bool)
	if err := Service().RoleAdd(ctx, &parameter, result); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData("")
}

func (this *apiImpl) RoleUpdate(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter RoleAddParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var result = new(bool)
	if err := Service().RoleUpdate(ctx, &parameter, result); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData("")
}

/**
 * Api-获取角色信息
 */
func (this *apiImpl) RoleGet(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter RoleNameParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var role = new(AdminRole)
	if err := Service().RoleGet(ctx, &parameter, role); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(role)
}

/**
 * Api-获取所有角色
 */
func (this *apiImpl) RoleGetAll(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter RoleNameParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var roles = make([]*AdminRole, 0)
	if err := Service().RoleGetAll(ctx, &parameter, &roles); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(roles)
}

/**
 * Api-删除角色
 */
func (this *apiImpl) RoleDelete(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter RoleNameParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var result = new(bool)
	if err := Service().RoleDelete(ctx, &parameter, result); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(result)
}

/**
 * Api-用户名存在校验
 */
func (this *apiImpl) AdminCheck(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter AdminCheckParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var exist = new(bool)
	if err := Service().AdminCheck(ctx, &parameter, exist); err != nil {
		g.ResponseError(err)
		return
	}
	var resp = map[string]interface{}{"exist":exist}
	g.ResponseData(resp)
}

/**
 * Api-用户注册
 */
func (this *apiImpl) AdminAdd(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter AdminAddParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var result = new(bool)
	if err := Service().AdminAdd(ctx, &parameter, result); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(map[string]interface{}{"result":result})
}

/**
 * Api-用户登录
 */
func (this *apiImpl) AdminLogin(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter AdminLoginParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var adminInfo = new(AdminInfo)
	if err := Service().AdminLogin(ctx, &parameter, adminInfo); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(adminInfo)
}

/**
 * Api-获取用户信息
 */
func (this *apiImpl) AdminGet(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter AdminGetParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var admin = new(AdminUser)
	if err := Service().AdminGet(ctx, &parameter, admin); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(admin)
}

/**
 * Api-获取用户信息
 */
func (this *apiImpl) AdminGetList(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter AdminCheckParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var admins = &[]*AdminUser{}
	if err := Service().AdminGetAll(ctx, &parameter, admins); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(admins)
}

/**
 * Api-删除用户
 */
func (this *apiImpl) AdminDelete(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter AdminDeleteParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var result = new(bool)
	if err := Service().AdminDelete(ctx, &parameter, result); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(map[string]interface{}{"result":result})
}

/**
 * Api-禁用用户
 */
func (this *apiImpl) AdminEnable(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter AdminEnabledParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var result = new(bool)
	if err := Service().AdminEnabled(ctx, &parameter, result); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(map[string]interface{}{"result":result})
}

/**
 * Api-用户修改密码
 */
func (this *apiImpl) AdminPasswordModify(ctx *gin.Context) {
	g := gohttp.Gin{
		Ctx: ctx,
	}
	var parameter AdminPassworddUpdateParameter
	if err := g.BindParameter(&parameter); err != nil {
		g.ResponseError(err)
		return
	}
	var result = new(bool)
	if err := Service().AdminModifyPassword(ctx, &parameter, result); err != nil {
		g.ResponseError(err)
		return
	}
	g.ResponseData(map[string]interface{}{"result":result})
}
