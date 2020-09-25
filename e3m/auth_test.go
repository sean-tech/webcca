package e3m

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"testing"
	"time"
)

var (
	_cli *clientv3.Client
)

var (
	e3config = E3Config{
		Organization: "sean-tech",
		Endpoints:    []string{"127.0.0.1:2379"},
		RootPassword: "etcd.user.root.pwd",
	}
	product = "ccatest"
	module = "testmodule"
	ip = "12.36.45.818"
)

func InitCli()  {
	Setup(e3config)
	var err error
	if _cli, err = clientv3.New(clientv3.Config{
		Endpoints:   _config.Endpoints,
		DialTimeout: client_dial_timeout,
		Username:    e3_user_root,
		Password:    _config.RootPassword,
	}); err != nil {
		panic(err)
	}
}

func TestNoUserGetRoles(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	auth := clientv3.NewAuth(cli)
	resp, err := auth.RoleList(context.Background())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.Roles)
}

func TestGetRoles(t *testing.T) {
	InitCli()
	auth := clientv3.NewAuth(_cli)
	resp, err := auth.RoleList(context.Background())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.Roles)
}

func TestGetUsers(t *testing.T) {
	InitCli()
	auth := clientv3.NewAuth(_cli)
	resp, err := auth.UserList(context.Background())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.Users)
}

func TestNewRole(t *testing.T) {
	InitCli()
	auth := clientv3.NewAuth(_cli)
	var role_name = "exampleproduct_user_config_read"
	_, err := auth.RoleAdd(context.Background(), role_name)
	if err != nil {
		t.Error(err)
	}
	var key = "/sean-tech/webkit/config/user/"
	var key_end = "/sean-tech/webkit/config/user0"
	if _, err := auth.RoleGrantPermission(context.Background(), role_name, key, key_end, clientv3.PermissionType(clientv3.PermRead)); err != nil {
		t.Error(err)
	}
	fmt.Println("role ", role_name, "created success.")
}

func TestNewUserWithNewRole(t *testing.T) {
	InitCli()
	auth := clientv3.NewAuth(_cli)
	var username = "exampleproduct_user_config_read_user"
	var password = "exampleproduct_user_config_read_user_password"
	_, err := auth.UserAdd(context.Background(), username, password)
	if err != nil {
		t.Error(err)
	}
	var role_name = "exampleproduct_user_config_read"
	if _, err := auth.UserGrantRole(context.Background(), username, role_name); err != nil {
		t.Error(err)
	}
	fmt.Println("user of role ", role_name, "created success")
}

func TestUserPutGetInPermission(t *testing.T) {
	var username = "exampleproduct_user_config_read_user"
	var password = "exampleproduct_user_config_read_user_password"
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Error(err)
	}
	InitCli()
	var key_prefix = "/sean-tech/webkit/config/user/"
	if resp, err := _cli.Put(context.Background(), key_prefix, "this is config for module user", clientv3.WithPrevKV()); err != nil {
		t.Error(err)
	} else if resp, err = _cli.Put(context.Background(), key_prefix + "1983427", "worker id is 1", clientv3.WithPrevKV()); err != nil {
		t.Error(err)
	} else {
		fmt.Println("user ", username, "put in permission success:")
		fmt.Println(resp.PrevKv.String())
	}

	if resp, err := cli.Get(context.Background(), key_prefix, clientv3.WithPrefix()); err != nil {
		t.Error(err)
	} else {
		fmt.Println("all value in key prefix ", key_prefix, "is :")
		fmt.Println(resp.Kvs)
	}
}

func TestUserGetInPermission(t *testing.T) {
	var username = "exampleproduct_user_config_read_user"
	var password = "exampleproduct_user_config_read_user_password"
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Error(err)
	}
	var key_prefix = "/sean-tech/webkit/config/user/"
	if resp, err := cli.Get(context.Background(), key_prefix, clientv3.WithPrefix()); err != nil {
		t.Error(err)
	} else {
		fmt.Println("all value in key prefix ", key_prefix, "is :")
		fmt.Println(resp.Kvs)
	}
}

func TestUserGetOutPermission(t *testing.T) {
	var username = "exampleproduct_user_config_read_user"
	var password = "exampleproduct_user_config_read_user_password"
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Error(err)
	}
	var key_prefix = "/sean-tech/webkit/config/notuser/"
	if resp, err := cli.Get(context.Background(), key_prefix, clientv3.WithPrefix()); err != nil {
		t.Error(err)
	} else {
		fmt.Println("all value in key prefix ", key_prefix, "is :")
		fmt.Println(resp.Kvs)
	}
}

func TestSaveGetPassword(t *testing.T) {
	Setup(e3config)
	if err := usersavepassword(product, "testusersave", "testusersavepwd"); err != nil {
		t.Error(err)
	} else {
		fmt.Println("user password save success")
	}
	if password, err := usergetpassword(product, "testusersave"); err != nil {
		t.Error(err)
	} else {
		fmt.Println("user password get success: ", password)
		if password == "testusersavepwd" {
			fmt.Println("pwd is equal")
		} else {
			fmt.Println("pwd not equal")
		}
	}
}

func TestDeleteSavedUserPassword(t *testing.T) {
	Setup(e3config)
	if _, err :=_rootcli.Delete(context.Background(), authuserpath(product, "ccatest_testmodule_config_rwuser")); err != nil {
		t.Error(err)
	} else {
		fmt.Println("saved user delete success")
	}
}

func TestDeletedAllSavedUserPassword(t *testing.T) {
	Setup(e3config)
	var path = fmt.Sprintf("/%s/%s/authuser/", _config.Organization, product)
	if _, err :=_rootcli.Delete(context.Background(), path, clientv3.WithPrefix()); err != nil {
		t.Error(err)
	} else {
		fmt.Println("saved user delete success")
	}
}

func TestAuthConfigRWUser(t *testing.T) {
	Setup(e3config)
	var user *AuthUser; var err error
	if user, err = ConfigModuleUser(product, module, AuthPermReadWrite); err != nil {
		t.Error(err)
	}
	fmt.Println("config rw user get success: ", user.Username, user.Password, user.Basepath)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e3config.Endpoints,
		DialTimeout: client_dial_timeout,
		Username:    user.Username,
		Password:    user.Password,
	})
	if err != nil {
		t.Error(err)
	}
	_, err = cli.Put(context.Background(), user.Basepath, "test config for product " +product+ " module " +module)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("rw user put config success")
	}
	if resp, err := cli.Get(context.Background(), user.Basepath); err != nil {
		t.Error(err)
	} else {
		fmt.Println("ro user get config success:")
		fmt.Println(resp.Kvs)
	}
}

func TestAuthConfigROUser(t *testing.T) {
	Setup(e3config)
	var user *AuthUser; var err error
	if user, err = ConfigModuleUser(product, module, AuthPermReadOnly); err != nil {
		t.Error(err)
	}
	fmt.Println("config ro user get success: ", user.Username, user.Password, user.Basepath)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e3config.Endpoints,
		DialTimeout: client_dial_timeout,
		Username:    user.Username,
		Password:    user.Password,
	})
	if err != nil {
		t.Error(err)
	}
	if resp, err := cli.Get(context.Background(), user.Basepath); err != nil {
		t.Error(err)
	} else {
		fmt.Println("ro user get config success:")
		fmt.Println(resp.Kvs)
	}
}