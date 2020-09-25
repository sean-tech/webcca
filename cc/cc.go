package cc

import (
	"bytes"
	"cca/e3m"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/sean-tech/gokit/requisition"
	"github.com/sean-tech/webkit/config"
	"go.etcd.io/etcd/clientv3"
	"strconv"
	"strings"
)

func NewProduct(product string) error {
	products, err := GetProducts()
	if err != nil {
		return err
	}
	for _, p := range products {
		if p == product {
			return nil
		}
	}
	if products == nil {
		products = []string{}
	}
	products = append(products, product)
	if buf, err := encode(products); err != nil {
		return err
	} else if _, err := e3m.Client().Put(context.Background(), productspath(), string(buf.Bytes())); err != nil {
		return err
	}
	if err := PutAppConfig(product, false, appconfig_template); err != nil {
		return err
	} else if err := newCaConfig(product); err != nil {
		return err
	} else  {
		return nil
	}
}

func GetProducts() ([]string, error){
	var products []string
	resp, err := e3m.Client().Get(context.Background(), productspath())
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	if  err := decode(bytes.NewBuffer(resp.Kvs[0].Value), &products); err != nil {
		return nil, err
	} else {
		return products, nil
	}
}

func productExist(product string) (bool, error) {
	products, err := GetProducts()
	if err != nil {
		return false, err
	}
	for _, p := range products {
		if p == product {
			return true, nil
		}
	}
	return false, nil
}

func deleteproducts() error {
	if _, err := e3m.Client().Delete(context.Background(), productspath()); err != nil {
		return err
	}
	return nil
}

func PutAppConfig(product string, checkexist bool, cfg *config.AppConfig) error {
	if cfg == nil {
		cfg = &config.AppConfig{}
	}
	if checkexist {
		if exist, err := productExist(product); err != nil {
			return err
		} else if exist == false {
			return requisition.NewError(nil, error_code_product_not_exist)
		}
	}
	var path = appconfigpath(product)
	if buf, err := encode(cfg); err != nil {
		return err
	} else if _, err := e3m.Client().Put(context.Background(), path, string(buf.Bytes()), clientv3.WithPrevKV()); err != nil {
		return err
	} else {
		return nil
	}
}

func GetAppConfig(product string) (cfg *config.AppConfig, err error) {
	if exist, err := productExist(product); err != nil {
		return nil, err
	} else if exist == false {
		return nil, requisition.NewError(nil, error_code_product_not_exist)
	}
	var path = appconfigpath(product)
	cfg = new(config.AppConfig)
	if resp, err := e3m.Client().Get(context.Background(), path); err != nil {
		return nil, err
	} else if len(resp.Kvs) != 1 {
		return nil, errors.New("config get error:kvs count not only one")
	} else if err := decode(bytes.NewBuffer(resp.Kvs[0].Value), cfg); err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}

func PutWorker(product, module, ip string, workerId int64) error {
	workers, err := GetAllWorkers(product, module)
	if err != nil {
		return err
	}
	for _, worker := range workers {
		if worker.WorkerId == workerId {
			return requisition.NewError(nil, error_code_server_workerid_exist)
		}
	}
	var path = serverpath(product, module, ip)
	if _, err := e3m.Client().Put(context.Background(), path, strconv.FormatInt(workerId, 10), clientv3.WithPrevKV()); err != nil {
		return err
	} else {
		workeradd_pub.Publish(ip)
		return nil
	}
}

func GetWorker(product, module, ip string) (worker *Worker, err error) {
	var path = serverpath(product, module, ip)
	var workerId int64 = 0
	if resp, err := e3m.Client().Get(context.Background(), path); err != nil {
		return nil, err
	} else if len(resp.Kvs) != 1 {
		return nil, errors.New("config get error:kvs count not only one")
	} else if workerId, err = strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64); err != nil {
		return nil, err
	} else {
		return &Worker{
			Module: module,
			Ip:       ip,
			WorkerId: workerId,
		}, nil
	}
}

func GetAllWorkers(product, module string) (workers []Worker, err error) {
	var path = modulepath(product, module) + "/"
	var resp *clientv3.GetResponse
	if resp, err = e3m.Client().Get(context.Background(), path, clientv3.WithPrefix()); err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		var workerId int64
		if workerId, err = strconv.ParseInt(string(kv.Value), 10, 64); err != nil {
			continue
		}
		ip := strings.Replace(string(kv.Key), path, "", 1)
		workers = append(workers, Worker{
			Module: module,
			Ip:       ip,
			WorkerId: workerId,
		})
	}
	return workers, nil
}

func GetAllModules(product string) (modules []string, err error) {
	var path = productpath(product)
	var resp *clientv3.GetResponse
	if resp, err = e3m.Client().Get(context.Background(), path, clientv3.WithPrefix()); err != nil {
		return nil, err
	}
	var moduleKeys = make(map[string]string)
	for _, kv := range resp.Kvs {
		module := strings.Replace(string(kv.Key), path, "", 1)
		if splitmodules := strings.Split(module, "/"); len(splitmodules) >= 2 {
			module = splitmodules[0]
		}
		moduleKeys[module] = module
	}
	for k, _ := range moduleKeys {
		modules = append(modules, k)
	}
	return modules, nil
}

func DeleteWorker(product, module, ip string) error {
	var path = serverpath(product, module, ip)
	if resp, err := e3m.Client().Get(context.Background(), path); err != nil {
		return err
	} else if len(resp.Kvs) == 0 {
		return nil
	} else if _, err := e3m.Client().Delete(context.Background(), path); err != nil {
		return err
	} else {
		workerdel_pub.Publish(ip)
		return nil
	}
}



func productspath() string {
	return fmt.Sprintf("/%s/webkit.products", e3m.Organization())
}
func appconfigpath(product string) string {
	return fmt.Sprintf("/%s/webkit.config/%s", e3m.Organization(), product)
}
func productpath(product string) string {
	return fmt.Sprintf("/%s/webkit.pms/%s/", e3m.Organization(), product)
}
func modulepath(product, module string) string {
	return fmt.Sprintf("/%s/webkit.pms/%s/%s", e3m.Organization(), product, module)
}
func serverpath(product, module, ip string) string {
	return fmt.Sprintf("/%s/webkit.pms/%s/%s/%s", e3m.Organization(), product, module, ip)
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