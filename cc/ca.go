package cc

import (
	"bytes"
	"cca/ca"
	"cca/e3m"
	"context"
	"fmt"
	"github.com/sean-tech/gokit/foundation"
	"github.com/sean-tech/webkit/gohttp"
	"github.com/sean-tech/webkit/gorpc"
	"net"
)

func NewServiceTLSCert(product, service string) (*gorpc.TlsConfig, error) {
	// load ca config saved
	saved_caconfig, err := GetCaConfig(product)
	if err != nil {
		return nil, err
	}
	// ips convert
	var ips []net.IP
	for _, ip := range saved_caconfig.Ips {
		ips = append(ips, net.ParseIP(ip))
	}
	// load signca from local file path
	var keypassword = ca.GenerateKey([]byte(saved_caconfig.KeyPassword))
	_, signcacert, signcakey, err := ca.LoadSignCA(product, keypassword, ips, saved_caconfig.Domains)
	if err != nil {
		return nil, err
	}
	var caconfig = &ca.CAConfig{
		ProductName: saved_caconfig.ProductName,
		Ips:         ips,
		Domains:     saved_caconfig.Domains,
		KeyPassword: keypassword,
		SignCACert:	 signcacert,
		SignCAKey:   signcakey,
	}
	// gen from ca to tls
	signCACertData, err := ca.GetSignCACertData(product)
	if err != nil {
		return nil, err
	}
	certData, keyData, err := ca.NewServerCert(product, service, *caconfig)
	if err != nil {
		return nil, err
	}
	return &gorpc.TlsConfig{
		CACert:       string(signCACertData),
		CACommonName: ca.GetSignCACommonName(product),
		ServerCert:   string(certData),
		ServerKey:    string(keyData),
	}, nil
}

func GetRsaConfigMap(product string) (map[string]*gohttp.RsaConfig, error) {
	keypairMap, err := ca.GetRsaKeyPairMap(product)
	if err != nil {
		return nil, err
	}
	if len(keypairMap) == 0 {
		return nil, nil
	}
	var rsaConfigMap = make(map[string]*gohttp.RsaConfig)
	for version, keypair := range keypairMap {
		rsaConfigMap[version] = &gohttp.RsaConfig{
			ServerPubKey: string(keypair.ServerPubKey),
			ServerPriKey: string(keypair.ServerPriKey),
			ClientPubKey: string(keypair.ClientPubKey),
		}
	}
	return rsaConfigMap, nil
}




type CAConfig struct {
	ProductName    string	`json:"product_name"`
	Ips            []string	`json:"ips"`
	Domains        []string `json:"domains"`
	KeyPassword    string	`json:"key_password"`
	RsaKeyPassword string	`json:"rsa_key_password"`
}

func newCaConfig(product string) error {
	var caconfig, err = GetCaConfig(product)
	if err != nil {
		return err
	}
	if caconfig != nil {
		return nil
	}
	caconfig = &CAConfig{
		ProductName: product,
		Ips:         nil,
		Domains:     nil,
		KeyPassword: foundation.RandString(16),
		RsaKeyPassword: foundation.RandString(16),
	}
	if buf, err := encode(caconfig); err != nil {
		return err
	} else if _, err := e3m.Client().Put(context.Background(), productcapath(product), string(buf.Bytes())); err != nil {
		return err
	} else {
		return nil
	}
}

func SetCaConfig(product string, ips, domains []string) error {
	var caconfig, err = GetCaConfig(product)
	if err != nil {
		return err
	}
	if caconfig == nil {
		caconfig = &CAConfig{
			ProductName: product,
			Ips:         ips,
			Domains:     domains,
			KeyPassword: foundation.RandString(16),
			RsaKeyPassword: foundation.RandString(16),
		}
	} else {
		caconfig.Ips = ips
		caconfig.Domains = domains
	}
	if buf, err := encode(caconfig); err != nil {
		return err
	} else if _, err := e3m.Client().Put(context.Background(), productcapath(product), string(buf.Bytes())); err != nil {
		return err
	} else {
		return nil
	}
}

func GetCaConfig(product string) (*CAConfig, error) {
	var caconfig = new(CAConfig)
	resp, err := e3m.Client().Get(context.Background(), productcapath(product))
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	if  err := decode(bytes.NewBuffer(resp.Kvs[0].Value), caconfig); err != nil {
		return nil, err
	} else {
		return caconfig, nil
	}
}

func productcapath(product string) string {
	return fmt.Sprintf("/%s/webkit.ca/%s", e3m.Organization(), product)
}
