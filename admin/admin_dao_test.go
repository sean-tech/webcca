package admin

import (
	"cca/e3m"
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"testing"
)

const (
	adminrole_super = "超级管理员"
	adminrole_manage = "管理员"
	adminrole_normal = "查阅员"

	adminuser_super = "super1"
	adminuser_super_pwd = "super1password"
	adminuser_manage = "manage1"
	adminuser_manage_pwd = "manage1password"
	adminuser_normal = "normal1"
	adminuser_normal_pwd = "normal1password"
)


var (
	e3config = e3m.E3Config{
		Organization: "sean-tech",
		Endpoints:    []string{"127.0.0.1:2379"},
		RootPassword: "etcd.user.root.pwd",
	}
	product = "ccatest"
	module = "testmodule"
	ip = "12.36.45.818"
)

func TestAdminRoleAdd(t *testing.T) {
	e3m.Setup(e3config)
	if err := dao().roleAdd(AdminRole{
		RoleName:      adminrole_super,
		AllowApis:     []string{"/admin/adduser",
			"/admin/addrole", "/admin/userlist", "/admin/rolelist", "/config/put", "/config/get", "/config/getall"},
		AllowProducts: nil,
	}); err != nil {
		t.Error(err)
	}
	if err := dao().roleAdd(AdminRole{
		RoleName:      adminrole_manage,
		AllowApis:     []string{"/admin/adduser",
			"/admin/addrole", "/admin/userlist", "/admin/rolelist", "/config/put", "/config/get", "/config/getall"},
		AllowProducts: nil,
	}); err != nil {
		t.Error(err)
	}
	if err := dao().roleAdd(AdminRole{
		RoleName:      adminrole_normal,
		AllowApis:     []string{"/admin/userlist", "/admin/rolelist", "/config/put",
			"/config/get", "/config/getall"},
		AllowProducts: nil,
	}); err != nil {
		t.Error(err)
	}
	fmt.Println("roles add success")
}

func TestAdminRolesGet(t *testing.T) {
	e3m.Setup(e3config)
	if roles, err := dao().roleGetAll(); err != nil {
		t.Error(err)
	} else {
		for _, role := range roles {
			fmt.Println(role.RoleName, ":", role.AllowApis)
		}
	}
}

func TestAdminRoleGet(t *testing.T) {
	e3m.Setup(e3config)
	if role, err := dao().roleGet(adminrole_normal); err != nil {
		t.Error(err)
	} else if role != nil {
		fmt.Println("role ", role.RoleName, "get success. allowapis:", role.AllowApis)
	} else {
		fmt.Println("role ", role, "get nil.")
	}
}

func TestAdminRoleDeleteAll(t *testing.T) {
	e3m.Setup(e3config)
	var roles = []string{
		adminrole_super,
		adminrole_manage,
		adminrole_normal,
	}
	for _, role := range roles {
		if err := dao().roleDelete(role); err != nil {
			t.Error(err)
		} else {
			fmt.Println("role", role, "delete success")
		}
	}
}

func TestAdminUserAdd(t *testing.T) {
	e3m.Setup(e3config)
	var tests = []AdminUser{
		{123, adminuser_super, adminuser_super_pwd, adminrole_super, true},
		{124, adminuser_manage, adminuser_manage_pwd, adminrole_manage, true},
		{125, adminuser_normal, adminuser_normal_pwd, adminrole_normal, true},
	}
	for _, user := range tests {
		if err := dao().adminAdd(user.Role, user.UserName, user.Password); err != nil {
			t.Error(err)
		} else {
			fmt.Println(user.UserName, "add success")
		}
	}
}

func TestAdminUserDeleteAll(t *testing.T) {
	e3m.Setup(e3config)
	if _, err := e3m.Client().Delete(context.Background(), admin_user_base_path(), clientv3.WithPrefix()); err != nil {
		t.Error(err)
	} else {
		fmt.Println("admin user delete all success")
	}
}

func TestAdminUserGetAll(t *testing.T) {
	e3m.Setup(e3config)
	if users, err := dao().adminList(); err != nil {
		t.Error(err)
	} else {
		for _, user := range users {
			fmt.Println(user.UserName, "--", user.Role, "--", user.Password)
		}
	}
}

func TestAdminUserGetByPwd(t *testing.T) {
	e3m.Setup(e3config)
	if user, err := dao().adminGetByUserNameAndPassword(adminuser_super, adminuser_super_pwd); err != nil {
		t.Error(err)
	} else {
		fmt.Println(user.UserName, "--", user.Role, "--", user.Password)
	}
}

func TestAdminUserGetByWrongPwd(t *testing.T) {
	e3m.Setup(e3config)
	if user, err := dao().adminGetByUserNameAndPassword(adminuser_manage, adminuser_normal_pwd); err != nil {
		t.Error(err)
	} else {
		fmt.Println(user.UserName, "--", user.Role, "--", user.Password)
	}
}
