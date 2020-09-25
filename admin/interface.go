package admin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/gokit/foundation"
	"github.com/sean-tech/gokit/requisition"
	"github.com/sean-tech/webkit/auth"
	"github.com/sean-tech/webkit/gohttp"
	"sync"
)

const (
	workerid = 1
)

type AdminUser struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	Password string `json:"-"`
	Role     string `json:"role"`
	Enabled  bool   `json:"enabled"`
}
type AdminRole struct {
	RoleName  		string   	`json:"roleName"`
	Description  	string		`json:"description"`
	AllowApis 		[]string 	`json:"allowApis"`
	AllowProducts 	[]string	`json:"allowProducts"`
	FrontRoutes		[]string	`json:"frontRoutes"`
}

type AdminInfo struct {
	*AdminUser
	*auth.AuthResult
}

/**参数：角色创建**/
type RoleAddParameter struct {
	RoleName  		string   	`json:"roleName" validate:"required,gte=1"`
	Description  	string   	`json:"description" validate:"required,gte=1"`
	AllowApis 		[]string 	`json:"allowApis" validate:"required,gte=1,dive,gte=1"`
	AllowProducts 	[]string 	`json:"allowProducts" validate:"required,gte=1,dive,gte=1"`
	FrontRoutes		[]string	`json:"frontRoutes"`
}

/**参数：角色操作**/
type RoleNameParameter struct {
	RoleName string 	`json:"roleName" validate:"required,gte=1"`
}

/**参数：用户创建**/
type AdminAddParameter struct {
	UserName string 	`json:"userName" validate:"required,gte=2"`
	RoleName string 	`json:"roleName" validate:"required,gte=1"`
	Password string		`json:"password" validate:"required,md5"`
}
/**参数：用户登录**/
type AdminLoginParameter struct {
	UserName string 	`json:"userName" validate:"required,gte=2"`
	Password string		`json:"password" validate:"required,md5"`
	UUID	 string		`json:"uuid" validate:"required,gte=6"`
	Client	 string		`json:"client" validate:"required,gte=2"`
}

/**参数：获取用户信息**/
type AdminGetParameter struct {
	UserId   int64	`json:"userId" validate:"required,gte=1"`
	UserName string `json:"userName" validate:"required,gte=2"`
}

/**参数：用户名存在校验**/
type AdminCheckParameter struct {
	UserName string 	`json:"userName" validate:"required,gte=2"`
}

type AdminDeleteParameter struct {
	UserName string `json:"userName" validate:"required,gte=2"`
}

type AdminEnabledParameter struct {
	UserName string `json:"userName" validate:"required,gte=2"`
	Enabled bool	`json:"enabled"`
}

/**参数：用户修改密码**/
type AdminPassworddUpdateParameter struct {
	UserId   int64	`json:"userId" validate:"required,gte=1"`
	UserName string 	`json:"userName" validate:"required,gte=2"`
	OldPassword string	`json:"oldPassword" validate:"required,md5"`
	NewPassword string	`json:"newPassword" validate:"required,md5"`
}



type IAdminApi interface {
	RoleAdd(ctx *gin.Context)
	RoleUpdate(ctx *gin.Context)
	RoleGet(ctx *gin.Context)
	RoleGetAll(ctx *gin.Context)
	RoleDelete(ctx *gin.Context)

	AdminCheck(ctx *gin.Context)
	AdminAdd(ctx *gin.Context)
	AdminLogin(ctx *gin.Context)
	AdminGet(ctx *gin.Context)
	AdminGetList(ctx *gin.Context)
	AdminDelete(ctx *gin.Context)
	AdminEnable(ctx *gin.Context)
	AdminPasswordModify(ctx *gin.Context)
}

type IAdminService interface {
	RoleAdd(ctx context.Context, parameter *RoleAddParameter, result *bool) error
	RoleUpdate(ctx context.Context, parameter *RoleAddParameter, result *bool) error
	RoleGet(ctx context.Context, parameter *RoleNameParameter, role *AdminRole) error
	RoleGetAll(ctx context.Context, parameter *RoleNameParameter, roles *[]*AdminRole) error
	RoleDelete(ctx context.Context, parameter *RoleNameParameter, result *bool) error

	AdminCheck(ctx context.Context, parameter *AdminCheckParameter, exist *bool) error
	AdminAdd(ctx context.Context, parameter *AdminAddParameter, result *bool) error
	AdminLogin(ctx context.Context, parameter *AdminLoginParameter, admininfo *AdminInfo) error
	AdminGet(ctx context.Context, parameter *AdminGetParameter, admin *AdminUser) error
	AdminGetAll(ctx context.Context, parameter *AdminCheckParameter, admins *[]*AdminUser) error
	AdminDelete(ctx context.Context, parameter *AdminDeleteParameter, result *bool) error
	AdminEnabled(ctx context.Context, parameter *AdminEnabledParameter, result *bool) error
	AdminModifyPassword(ctx context.Context, parameter *AdminPassworddUpdateParameter, result *bool) error
}

type iAdminDao interface {
	roleAdd(role AdminRole) error
	roleUpdate(role AdminRole) error
	roleDelete(rolename string) error
	roleGet(rolename string) (*AdminRole, error)
	roleGetAll() ([]*AdminRole, error)

	adminExistCheck(userName string) (bool, error)
	adminAdd(rolename, userName, password string) error
	adminGetByUserNameAndPassword(userName, password string) (*AdminUser, error)
	adminGet(userName string) (*AdminUser, error)
	adminModifyPassword(username, oldpassword, newpassword string) error
	adminDelete(userName string) error
	adminEnabled(userName string, enabled bool) error
	adminList() ([]*AdminUser, error)
}

var (
	_api         IAdminApi
	_apiOnce     sync.Once
	_service     IAdminService
	_serviceOnce sync.Once
	_dao         iAdminDao
	_daoOnce     sync.Once
)

func Api() IAdminApi {
	_apiOnce.Do(func() {
		_api = new(apiImpl)
	})
	return _api
}

func Service() IAdminService {
	_serviceOnce.Do(func() {
		_service = new(serviceImpl)
	})
	return _service
}

func dao() iAdminDao {
	_daoOnce.Do(func() {
		worker, _ :=  foundation.NewWorker(workerid)
		_dao = &daoImpl{worker:worker}
	})
	return _dao
}



const (
	_                                     int 	= 0
	error_code_permission_denied				= gohttp.STATUS_CODE_PERMISSION_DENIED
	error_code_role_exist                       = 11101
	error_code_role_not_exist					= 11102
	error_code_user_exist                       = 11001
	error_code_user_not_exist                   = 11002
	error_code_notfilter_username_password      = 11003
	error_code_user_disenabled					= 11009
	error_code_update_failed			        = 11004
	error_code_delete_failed			        = 11005
	error_code_enable_failed			        = 11006
	error_code_disenabled_failed			    = 11007
)

func init() {
	requisition.SetMsgMap(requisition.LanguageZh, map[int]string{
		error_code_role_exist                    : "角色已存在",
		error_code_role_not_exist				 : "角色不存在",
		error_code_user_exist                    : "用户已存在",
		error_code_user_not_exist                : "用户不存在",
		error_code_notfilter_username_password   : "用户名或密码不正确",
		error_code_user_disenabled				 : "用户已被禁用",
		error_code_update_failed			     : "修改失败",
		error_code_delete_failed			     : "删除失败",
		error_code_enable_failed			     : "启用失败",
		error_code_disenabled_failed			 : "禁用失败",
	})
	requisition.SetMsgMap(requisition.LanguageEn, map[int]string{
		error_code_role_exist                    : "role exist already",
		error_code_role_not_exist				 : "role is not exist",
		error_code_user_exist                    : "user exist already",
		error_code_user_not_exist                : "user not exist",
		error_code_notfilter_username_password   : "username or password not right",
		error_code_user_disenabled				 : "user is disenabled",
		error_code_update_failed			     : "update failed",
		error_code_delete_failed			     : "delete failed",
		error_code_enable_failed			     : "user enable failed",
		error_code_disenabled_failed			 : "user disenable failed",
	})
}