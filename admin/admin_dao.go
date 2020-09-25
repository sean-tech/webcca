package admin

import (
	"bytes"
	"cca/e3m"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/sean-tech/gokit/foundation"
	"github.com/sean-tech/gokit/requisition"
	"go.etcd.io/etcd/clientv3"
)

type daoImpl struct {
	worker *foundation.Worker
}

func (this *daoImpl) roleAdd(role AdminRole) error {
	if savedRole, err := this.roleGet(role.RoleName); err != nil {
		return err
	} else if savedRole != nil {
		return requisition.NewError(nil, error_code_role_exist)
	}
	if buf, err := encode(role); err != nil {
		return err
	} else if _, err := e3m.Client().Put(context.Background(), admin_role_base_path() + role.RoleName, string(buf.Bytes())); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *daoImpl) roleUpdate(role AdminRole) error {
	if savedRole, err := this.roleGet(role.RoleName); err != nil {
		return err
	} else if savedRole == nil {
		return requisition.NewError(nil, error_code_role_not_exist)
	}
	if buf, err := encode(role); err != nil {
		return err
	} else if _, err := e3m.Client().Put(context.Background(), admin_role_base_path() + role.RoleName, string(buf.Bytes())); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *daoImpl) roleDelete(rolename string) error {
	users, err := this.adminList()
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Role == rolename {
			this.adminDelete(user.UserName)
		}
	}
	if _, err := e3m.Client().Delete(context.Background(), admin_role_base_path() + rolename); err != nil {
		return err
	}
	return nil
}

func (this *daoImpl) roleGet(rolename string) (*AdminRole, error) {
	var role = new(AdminRole)
	if resp, err := e3m.Client().Get(context.Background(), admin_role_base_path() + rolename); err != nil {
		return nil, err
	} else if len(resp.Kvs) == 0 {
		return nil, nil
	} else if err := decode(bytes.NewBuffer(resp.Kvs[0].Value), role); err != nil {
		return nil, err
	} else {
		return role, nil
	}
}

func (this *daoImpl) roleGetAll() ([]*AdminRole, error) {
	resp, err := e3m.Client().Get(context.Background(), admin_role_base_path(), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var roles []*AdminRole
	for _, kv := range resp.Kvs {
		var role = new(AdminRole)
		if err := decode(bytes.NewBuffer(kv.Value), &role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (this *daoImpl) adminExistCheck(userName string) (bool, error) {
	if user, err := this.adminGet(userName); err != nil {
		return false, err
	} else if user == nil {
		return false, nil
	}
	return true, nil
}

func (this *daoImpl) adminAdd(rolename, username, password string) error {
	role, err := this.roleGet(rolename)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("not found role " + rolename + " for user " + username)
	}
	if resp, err := e3m.Client().Get(context.Background(), admin_user_base_path() + username); err != nil {
		return err
	} else if len(resp.Kvs) > 0 {
		return requisition.NewError(nil, error_code_user_exist)
	}

	var user = &AdminUser{
		UserId:   this.worker.GetId(),
		UserName: username,
		Password: password,
		Role:     rolename,
		Enabled:  true,
	}
	if buf, err := encode(user); err != nil {
		return err
	} else if _, err := e3m.Client().Put(context.Background(), admin_user_base_path() + username, string(buf.Bytes())); err != nil {
		return err
	}
	return nil
}

func (this *daoImpl) adminModifyPassword(username, oldpassword, newpassword string) error {
	user, err := this.adminGet(username)
	if err != nil {
		return err
	}
	if user == nil {
		return requisition.NewError(nil, error_code_user_not_exist)
	}
	if user.Password != oldpassword {
		return requisition.NewError(nil, error_code_notfilter_username_password)
	}
	user.Password = newpassword
	if buf, err := encode(user); err != nil {
		return err
	} else if _, err := e3m.Client().Put(context.Background(), admin_user_base_path() + username, string(buf.Bytes())); err != nil {
		return err
	}
	return nil
}

func (this *daoImpl) adminDelete(username string) error {
	if _, err := e3m.Client().Delete(context.Background(), admin_user_base_path() + username); err != nil {
		return err
	}
	return nil
}

func (this *daoImpl) adminGetByUserNameAndPassword(username, password string) (*AdminUser, error) {
	var user = new(AdminUser)
	if resp, err := e3m.Client().Get(context.Background(), admin_user_base_path() + username); err != nil {
		return nil, err
	} else if len(resp.Kvs) == 0 {
		return nil, requisition.NewError(nil, error_code_user_not_exist)
	} else if err := decode(bytes.NewBuffer(resp.Kvs[0].Value), user); err != nil {
		return nil, err
	} else if password != user.Password {
		return nil, requisition.NewError(nil, error_code_notfilter_username_password)
	} else if user.Enabled == false {
		return nil, requisition.NewError(nil, error_code_user_disenabled)
	} else {
		return user, nil
	}
}

func (this *daoImpl) adminGet(username string) (*AdminUser, error) {
	var user = new(AdminUser)
	if resp, err := e3m.Client().Get(context.Background(), admin_user_base_path() + username); err != nil {
		return nil, err
	} else if len(resp.Kvs) == 0 {
		return nil, nil
	} else if err := decode(bytes.NewBuffer(resp.Kvs[0].Value), user); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (this *daoImpl) adminList() ([]*AdminUser, error) {
	resp, err := e3m.Client().Get(context.Background(), admin_user_base_path(), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var users []*AdminUser
	for _, kv := range resp.Kvs {
		var user = new(AdminUser)
		if err := decode(bytes.NewBuffer(kv.Value), user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (this *daoImpl) adminEnabled(username string, enabled bool) error {
	var user = new(AdminUser)
	if resp, err := e3m.Client().Get(context.Background(), admin_user_base_path() + username); err != nil {
		return err
	} else if len(resp.Kvs) == 0 {
		return requisition.NewError(nil, error_code_user_not_exist)
	} else if err := decode(bytes.NewBuffer(resp.Kvs[0].Value), user); err != nil {
		return err
	} else {
		user.Enabled = enabled
	}
	// resave
	if buf, err := encode(user); err != nil {
		return err
	} else if _, err := e3m.Client().Put(context.Background(), admin_user_base_path() + username, string(buf.Bytes())); err != nil {
		return err
	} else {
		return nil
	}
}



func admin_user_base_path() string {
	return fmt.Sprintf("/%s/cca/admin/user/", e3m.Organization())
}
func admin_role_base_path() string {
	return fmt.Sprintf("/%s/cca/admin/role/", e3m.Organization())
}

func encode(data interface{}) (*bytes.Buffer, error) {
	//Buffer类型实现了io.Writer接口
	var buf bytes.Buffer
	//得到编码器
	enc := gob.NewEncoder(&buf)
	//调用编码器的Encode方法来编码数据data
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	//编码后的结果放在buf中
	return &buf, nil
}

func decode(buf *bytes.Buffer, data interface{}) error {
	//获取一个解码器，参数需要实现io.Reader接口
	dec := gob.NewDecoder(buf)
	//调用解码器的Decode方法将数据解码，用Q类型的q来接收
	return dec.Decode(data)
}