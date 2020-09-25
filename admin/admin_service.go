package admin

import (
	"context"
	"github.com/sean-tech/gokit/encrypt"
	"github.com/sean-tech/gokit/requisition"
	"github.com/sean-tech/gokit/validate"
	"github.com/sean-tech/webkit/auth"
)
const (
	AUTH_CODE  		= "kzncknqnzmc9843yu!z#xv123"
	ROLE_SUPER     	= "admin_role_sueper"
	USER_SUPER     	= "superadmin"
	user_super_pwd 	= "superadmin"
)

type serviceImpl struct {
}

/**
 * Service-用户存在校验
 */
func (this *serviceImpl) RoleAdd(ctx context.Context, parameter *RoleAddParameter, result *bool) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	return dao().roleAdd(AdminRole{
		RoleName:      parameter.RoleName,
		Description:   parameter.Description,
		AllowApis:     parameter.AllowApis,
		AllowProducts: parameter.AllowProducts,
		FrontRoutes:   parameter.FrontRoutes,
	})
}

func (this *serviceImpl) RoleUpdate(ctx context.Context, parameter *RoleAddParameter, result *bool) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	return dao().roleUpdate(AdminRole{
		RoleName:      parameter.RoleName,
		Description:   parameter.Description,
		AllowApis:     parameter.AllowApis,
		AllowProducts: parameter.AllowProducts,
		FrontRoutes:   parameter.FrontRoutes,
	})
}

/**
 * Service-用户存在校验
 */
func (this *serviceImpl) RoleGet(ctx context.Context, parameter *RoleNameParameter, role *AdminRole) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	m_role, err := dao().roleGet(parameter.RoleName)
	*role = *m_role
	return err
}

/**
 * Service-用户存在校验
 */
func (this *serviceImpl) RoleGetAll(ctx context.Context, parameter *RoleNameParameter, roles *[]*AdminRole) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	mRoles, err := dao().roleGetAll()
	var superAdminIndex int
	for idx, role := range mRoles {
		if role.RoleName == ROLE_SUPER {
			superAdminIndex = idx
			break
		}
	}
	mRoles = append(mRoles[:superAdminIndex], mRoles[superAdminIndex+1:]...)
	*roles = mRoles
	return err
}

/**
 * Service-
 */
func RoleGetAll() ([]*AdminRole, error) {
	if roles, err := dao().roleGetAll(); err != nil {
		return nil, err
	} else {
		return roles, nil
	}
}

/**
 * Service-用户存在校验
 */
func (this *serviceImpl) RoleDelete(ctx context.Context, parameter *RoleNameParameter, result *bool) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	return dao().roleDelete(parameter.RoleName)
}

/**
 * Service-用户存在校验
 */
func (this *serviceImpl) AdminCheck(ctx context.Context, parameter *AdminCheckParameter, exist *bool) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	user_exist, err := dao().adminExistCheck(parameter.UserName)
	if err != nil {
		return err
	}
	*exist = user_exist
	return nil
}

/**
 * Service-创建用户
 */
func (this *serviceImpl) AdminAdd(ctx context.Context, parameter *AdminAddParameter, result *bool) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	if err := dao().adminAdd(parameter.RoleName, parameter.UserName, parameter.Password); err != nil {
		return err
	}
	*result = true
	return nil
}

/**
 * Service-用户登录验证
 */
func (this *serviceImpl) AdminLogin(ctx context.Context, parameter *AdminLoginParameter, admininfo *AdminInfo) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	if parameter.UserName == USER_SUPER {
		if role, err := dao().roleGet(ROLE_SUPER); err != nil {
			return err
		} else if role == nil {
			if err := dao().roleAdd(AdminRole{
				RoleName:      ROLE_SUPER,
				AllowApis:     nil,
				AllowProducts: nil,
			}); err != nil {
				return err
			}
		}
		if exist, err := dao().adminExistCheck(USER_SUPER); err != nil {
			return err
		} else if exist == false {
			if err := dao().adminAdd(ROLE_SUPER, USER_SUPER, encrypt.GetMd5().Encode([]byte(user_super_pwd))); err != nil {
				return err
			}
		}
	}
	if adminuser, err := dao().adminGetByUserNameAndPassword(parameter.UserName, parameter.Password); err != nil {
		return err
	} else {
		admininfo.AdminUser = adminuser
	}
	var authParameter = &auth.NewAuthParameter{
		AuthCode: AUTH_CODE,
		UUID:     parameter.UUID,
		UserId:   uint64(admininfo.UserId),
		UserName: admininfo.UserName,
		Role: 	  admininfo.Role,
		Client:   parameter.Client,
	}
	var authResult = new(auth.AuthResult)
	if err := auth.Service().NewAuth(ctx, authParameter, authResult); err != nil {
		return err
	}
	admininfo.AuthResult = authResult
	return nil
}

/**
 * 用户信息获取
 */
func (this *serviceImpl) AdminGet(ctx context.Context, parameter *AdminGetParameter, admin *AdminUser) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	checkResult := requisition.CheckTokenUser(ctx, uint64(parameter.UserId), parameter.UserName)
	if checkResult == false {
		return requisition.NewError(nil, error_code_permission_denied)
	}
	model_admin, err := dao().adminGet(parameter.UserName)
	if err != nil {
		return err
	}
	*admin = *model_admin
	return nil
}

/**
 * 用户信息获取
 */
func (this *serviceImpl) AdminGetAll(ctx context.Context, parameter *AdminCheckParameter, admins *[]*AdminUser) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	var user = new(AdminUser)
	if err := this.AdminGet(ctx, &AdminGetParameter{
		UserId:   int64(requisition.GetRequisition(ctx).UserId),
		UserName: requisition.GetRequisition(ctx).UserName,
	}, user); err != nil {
		return err
	}

	var role = new(AdminRole)
	if err := this.RoleGet(ctx, &RoleNameParameter{user.Role}, role); err != nil {
		return err
	}
	if model_admins, err := dao().adminList(); err != nil {
		return err
	} else {
		*admins = model_admins
		return nil
	}
	return nil
}

/**
 * Service-删除用户
 */
func (this *serviceImpl) AdminDelete(ctx context.Context, parameter *AdminDeleteParameter, result *bool) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	if err := dao().adminDelete(parameter.UserName); err != nil {
		return err
	}
	*result = true
	return nil
}

/**
 * Service-用户禁用
 */
func (this *serviceImpl) AdminEnabled(ctx context.Context, parameter *AdminEnabledParameter, result *bool) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	if err := dao().adminEnabled(parameter.UserName, parameter.Enabled); err != nil {
		return err
	}
	*result = true
	return nil
}

/**
 * 修改密码
 */
func (this *serviceImpl) AdminModifyPassword(ctx context.Context, parameter *AdminPassworddUpdateParameter, result *bool) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	checkResult := requisition.CheckTokenUser(ctx, uint64(parameter.UserId), parameter.UserName)
	if checkResult == false {
		return requisition.NewError(nil, error_code_permission_denied)
	}
	if err := dao().adminModifyPassword(parameter.UserName, parameter.OldPassword, parameter.NewPassword); err != nil {
		return err
	}
	*result = true
	return nil
}