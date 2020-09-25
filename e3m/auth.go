package e3m

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sean-tech/gokit/encrypt"
	"github.com/sean-tech/gokit/foundation"
	"go.etcd.io/etcd/clientv3"
)

const (
	e3_auth_user_password_length = 16
)

type AuthPermission int
const (
	AuthPermReadOnly  AuthPermission = 0
	AuthPermReadWrite AuthPermission = 2
)

type AuthUser struct {
	Username string
	Password string
	Basepath string
}

func ConfigProductUser(product string, perm AuthPermission) (*AuthUser, error) {
	var rolename = config_product_rolename(product, perm)
	var username = config_product_username(product, perm)
	var basepath = config_product_path(product)
	return GetUser(product, rolename, username, basepath, perm)
}

func ConfigModuleUser(product, module string, perm AuthPermission) (*AuthUser, error) {
	var rolename = config_module_rolename(product, module, perm)
	var username = config_module_username(product, module, perm)
	var basepath = config_module_path(product, module)
	return GetUser(product, rolename, username, basepath, perm)
}

func RpcUser(product string) (*AuthUser, error) {
	var rolename = rpc_rolename(product)
	var username = rpc_username(product)
	var basepath = rpc_path(product)
	return GetUser(product, rolename, username, basepath, AuthPermReadWrite)
}

func GetUser(product, rolename, username, basepath string, perm AuthPermission) (*AuthUser, error) {
	// search user
	var users []string; var err error
	if users, err = GetUsers(); err != nil {
		return nil, err
	}
	for _, name := range users {
		if username != name {
			continue
		}
		// found username, get password
		if password, err := usergetpassword(product, username); err != nil {
			return nil, err
		} else {
			return &AuthUser{
				Username: username,
				Password: password,
				Basepath: basepath,
			}, nil
		}
	}
	// not found user in etcd auth users
	// search role
	var roles []string
	if roles, err = GetRoles(); err != nil {
		return nil, err
	}
	for _, name := range roles {
		if rolename != name {
			continue
		}
		// found role, create user
		var password = foundation.RandString(e3_auth_user_password_length)
		_, err := _rootauth.UserAdd(context.Background(), username, password)
		if err != nil {
			return nil, err
		}
		// grant role to user, if failed, delete user created
		if _, err := _rootauth.UserGrantRole(context.Background(), username, rolename); err != nil {
			_rootauth.UserDelete(context.Background(), username)
			return nil, err
		}
		// save password, if failed, delete user created and granted
		if err := usersavepassword(product, username, password); err != nil {
			_rootauth.UserDelete(context.Background(), username)
			return nil, err
		}
		return &AuthUser{
			Username: username,
			Password: password,
			Basepath: basepath,
		}, nil
	}
	// not found role in etcd auth roles
	// create role
	if _, err := _rootauth.RoleAdd(context.Background(), rolename); err != nil {
		return nil, err
	}
	if _, err := _rootauth.RoleGrantPermission(context.Background(), rolename, basepath, basepath+"0", clientv3.PermissionType(perm)); err != nil {
		return nil, err
	}
	// create user
	var password = foundation.RandString(e3_auth_user_password_length)
	if _, err := _rootauth.UserAdd(context.Background(), username, password); err != nil {
		return nil, err
	}
	// grant role to user
	if _, err := _rootauth.UserGrantRole(context.Background(), username, rolename); err != nil {
		return nil, err
	}
	// save password, if failed, delete user created and granted
	if err := usersavepassword(product, username, password); err != nil {
		_rootauth.UserDelete(context.Background(), username)
		return nil, err
	}
	return &AuthUser{
		Username: username,
		Password: password,
		Basepath: basepath,
	}, nil
}

func GetUsers() ([]string, error) {
	resp, err := _rootauth.UserList(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}

func GetRoles() ([]string, error) {
	resp, err := _rootauth.RoleList(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Roles, nil
}

func usersavepassword(product, username, password string) error {
	var passwordsecret string
	if encryptData, err := encrypt.GetAes().EncryptCBC([]byte(password), pwdkey(product, username)); err != nil {
		return err
	} else {
		passwordsecret = base64.StdEncoding.EncodeToString(encryptData)
	}
	var path = authuserpath(product, username)
	if _, err := _rootcli.Put(context.Background(), path, passwordsecret, clientv3.WithPrevKV()); err != nil {
		return err
	} else {
		return nil
	}
}

func usergetpassword(product, username string) (string, error) {
	var path = authuserpath(product, username)
	if resp, err := _rootcli.Get(context.Background(), path); err != nil {
		return "", err
	} else if len(resp.Kvs) != 1 {
		return "", errors.New("password not found for username " + username)
	} else if encryptData, err := base64.StdEncoding.DecodeString(string(resp.Kvs[0].Value)); err != nil {
		return "", err
	} else if decryptData, err := encrypt.GetAes().DecryptCBC(encryptData, pwdkey(product, username)); err != nil {
		return "", err
	} else if  decryptData == nil {
		return "", errors.New("password decrypt failed for username " + username)
	} else {
		return string(decryptData), nil
	}
}

func authuserpath(product, username string) string {
	return fmt.Sprintf("/%s/%s/authuser/%s", _config.Organization, product, username)
}

func pwdkey(product, username string) []byte {
	var originkey = fmt.Sprintf("%s_%s_authuser_%s_pwdkey", _config.Organization, product, username)
	md5Value := encrypt.GetMd5().Encode([]byte(originkey))
	return generateKey([]byte(md5Value))
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 32)
	copy(genKey, key)
	for i := 32; i < len(key); {
		for j := 0; j < 32 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func config_product_path(product string) string {
	return fmt.Sprintf("/%s/%s/config/", _config.Organization, product)
}
func config_product_username(product string, perm AuthPermission) string {
	switch perm {
	case AuthPermReadOnly:
		return fmt.Sprintf("%s_config_rouser", product)
	case AuthPermReadWrite:
		return fmt.Sprintf("%s_config_rwuser", product)
	}
	return ""
}
func config_product_rolename(product string, perm AuthPermission) string {
	switch perm {
	case AuthPermReadOnly:
		return fmt.Sprintf("%s_config_rorole", product)
	case AuthPermReadWrite:
		return fmt.Sprintf("%s_config_rwrole", product)
	}
	return ""
}
func config_module_path(product, module string) string {
	return fmt.Sprintf("/%s/%s/config/%s", _config.Organization, product, module)
}
func config_module_username(product, module string, perm AuthPermission) string {
	switch perm {
	case AuthPermReadOnly:
		return fmt.Sprintf("%s_%s_config_rouser", product, module)
	case AuthPermReadWrite:
		return fmt.Sprintf("%s_%s_config_rwuser", product, module)
	}
	return ""
}
func config_module_rolename(product, module string, perm AuthPermission) string {
	switch perm {
	case AuthPermReadOnly:
		return fmt.Sprintf("%s_%s_config_rorole", product, module)
	case AuthPermReadWrite:
		return fmt.Sprintf("%s_%s_config_rwrole", product, module)
	}
	return ""
}
func rpc_path(product string) string {
	return fmt.Sprintf("/%s/%s/rpc", _config.Organization, product)
}
func rpc_username(product string) string {
	return fmt.Sprintf("%s_rpcx_rpc_user", product)
}
func rpc_rolename(product string) string {
	return fmt.Sprintf("%s_rpcx_rpc_role", product)
}
