package rpcxservice

import (
	"cca/e3m"
	"encoding/base64"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/docker/libkv"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

type EtcdRegistry struct {
	kv kvstore.Store
}

func (r *EtcdRegistry) initRegistry() {
	etcd.Register()

	kv, err := libkv.NewStore(kvstore.ETCD, serverConfig.RegistryURLs, nil)
	if err != nil {
		log.Printf("cannot create etcd registry: %v", err)
		return
	}
	r.kv = kv

	return
}

func (r *EtcdRegistry) fetchServices(product string) ([]*RpcService, error) {
	rpcuser, err := e3m.RpcUser(product)
	if err != nil {
		return nil, err
	}
	var ServiceBaseURL = rpcuser.Basepath
	if !strings.HasSuffix(ServiceBaseURL, "/") {
		ServiceBaseURL += "/"
	}
	var services []*RpcService
	kvs, err := r.kv.List(rpcuser.Basepath)
	if err != nil {
		log.Printf("failed to list services %s: %v", ServiceBaseURL, err)
		return services, nil
	}

	for _, value := range kvs {

		nodes, err := r.kv.List(value.Key)
		if err != nil {
			log.Printf("failed to list %s: %v", value.Key, err)
			continue
		}

		for _, n := range nodes {
			key := string(n.Key[:])
			i := strings.LastIndex(key, "/")
			serviceName := strings.TrimPrefix(key[0:i], ServiceBaseURL)
			var serviceAddr string
			fields := strings.Split(key, "/")
			if fields != nil && len(fields) > 1 {
				serviceAddr = fields[len(fields)-1]
			}
			v, err := url.ParseQuery(string(n.Value[:]))
			if err != nil {
				log.Println("etcd value parse failed. error: ", err.Error())
				continue
			}
			state := "n/a"
			group := ""
			if err == nil {
				state = v.Get("state")
				if state == "" {
					state = "active"
				}
				group = v.Get("group")
			}
			id := base64.StdEncoding.EncodeToString([]byte(serviceName + "@" + serviceAddr))
			service := &RpcService{ID: id, Name: serviceName, Address: serviceAddr, Metadata: string(n.Value[:]), State: state, Group: group}
			services = append(services, service)
		}

	}

	return services, nil
}

func (r *EtcdRegistry) deactivateService(product, name, address string) error {
	rpcuser, err := e3m.RpcUser(product)
	if err != nil {
		return err
	}
	var ServiceBaseURL = rpcuser.Basepath
	if !strings.HasSuffix(ServiceBaseURL, "/") {
		ServiceBaseURL += "/"
	}
	key := path.Join(ServiceBaseURL, name, address)

	kv, err := r.kv.Get(key)

	if err != nil {
		return err
	}

	v, err := url.ParseQuery(string(kv.Value[:]))
	if err != nil {
		log.Println("etcd value parse failed. err ", err.Error())
		return err
	}
	v.Set("state", "inactive")
	err = r.kv.Put(kv.Key, []byte(v.Encode()), &kvstore.WriteOptions{IsDir: false})
	if err != nil {
		log.Println("etcd set failed, err : ", err.Error())
	}

	return err
}

func (r *EtcdRegistry) activateService(product, name, address string) error {
	rpcuser, err := e3m.RpcUser(product)
	if err != nil {
		return err
	}
	var ServiceBaseURL = rpcuser.Basepath
	key := path.Join(ServiceBaseURL, name, address)
	kv, err := r.kv.Get(key)

	v, err := url.ParseQuery(string(kv.Value[:]))
	if err != nil {
		log.Println("etcd value parse failed. err ", err.Error())
		return err
	}
	v.Set("state", "active")
	err = r.kv.Put(kv.Key, []byte(v.Encode()), &kvstore.WriteOptions{IsDir: false})
	if err != nil {
		log.Println("etcdv3 put failed. err: ", err.Error())
	}

	return err
}

func (r *EtcdRegistry) updateMetadata(product, name, address string, metadata string) error {
	rpcuser, err := e3m.RpcUser(product)
	if err != nil {
		return err
	}
	var ServiceBaseURL = rpcuser.Basepath
	key := path.Join(ServiceBaseURL, name, address)
	return r.kv.Put(key, []byte(metadata), &kvstore.WriteOptions{IsDir: false})
}
