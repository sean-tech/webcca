package ca

import (
	"fmt"
	"github.com/sean-tech/gokit/foundation"
	"testing"
)

func TestRsa(t *testing.T) {
	Setup(CASetting{
		CompanyName: "sean-tech",
		CASavePath:  "/Users/sean/Desktop/",
	})
	product := "testca"
	fileinfos, err := GetAllRSAKeyFileInfos(product)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fileinfos)
	err = NewRsaKeyPair(product, "1.0.0", []byte(foundation.RandString(16)))
	if err != nil {
		t.Error(err)
	}
	fileinfos, err = GetAllRSAKeyFileInfos(product)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fileinfos)

}

//type CAConfig2 struct {
//	CompanyName    string
//	ProductName    string
//	Ips            []net.IP	`json:"ips"` // net.ParseIP("127.0.0.1")
//	Domains        []string `json:"domains"` // "*." + "ex.liansheng"
//	ExportFilePath string	`json:"export_file_path"`
//	KeyPassword    []byte	`json:"key_password"`
//}
//var (
//	_config CAConfig2
//	_caCert *pkix.Certificate
//	_caKey  *pkix.Key
//	_caPem  []byte
//)
//
//func Setup2(config CAConfig2) {
//	_config = config
//	_config.KeyPassword = GenerateKey(_config.KeyPassword)
//	if fileutils.CheckExist(config.ExportFilePath +filename_signca_pem) == true {
//		if fileutils.CheckExist(config.ExportFilePath +filename_signca_key) == true {
//			loadCA2()
//			return
//		}
//	}
//	// init root ca and export
//	rootCert, rootKey := newRootCA2(config.CompanyName, config.ExportFilePath, _config.KeyPassword)
//	initCA(rootCert, rootKey, config.ProductName + "." + config.CompanyName, config.Ips, config.Domains, config.ExportFilePath, _config.KeyPassword)
//	loadCA2()
//}
//
//func loadCA2() {
//	if _caPem, err := ioutil.ReadFile(_config.ExportFilePath + filename_signca_pem); err != nil {
//		panic(err)
//	} else if _caCert, err = pkix.NewCertificateFromPEM(_caPem); err != nil {
//		panic(err)
//	}
//	if caKey, err := ioutil.ReadFile(_config.ExportFilePath + filename_signca_key); err != nil {
//		panic(err)
//	} else if _caKey, err = pkix.NewKeyFromEncryptedPrivateKeyPEM(caKey, _config.KeyPassword); err != nil {
//		panic(err)
//	}
//}
//
//func newRootCA2(commonName, exportFilePath string, password []byte) (cert *pkix.Certificate, key *pkix.Key) {
//	// init root ca and export
//	key, err := pkix.CreateRSAKey(keybits)
//	if err != nil {
//		panic(err)
//	}
//	cert, err = pkix.CreateCertificateAuthority(key, "", time.Now().Add(days * 24 *time.Hour), "", "", "", "", commonName)
//	if err != nil {
//		panic(err)
//	}
//	if pemData, err := cert.Export(); err != nil {
//		panic(err)
//	} else {
//		writeToFile(exportFilePath, filename_rootca_pem, pemData)
//	}
//	if keyData, err := key.ExportEncryptedPrivate(password); err != nil {
//		panic(err)
//	} else {
//		writeToFile(exportFilePath, filename_rootca_key, keyData)
//	}
//	return cert, key
//}
//
//func initCA(authCert *pkix.Certificate, authKey *pkix.Key, commonName string, ips []net.IP, domains []string, exportFilePath string, password []byte) {
//	key, err := pkix.CreateRSAKey(keybits)
//	if err != nil {
//		panic(err)
//	}
//	csr, err := pkix.CreateCertificateSigningRequest(key, "", ips, domains, nil, "", "", "", "", commonName)
//	if err != nil {
//		panic(err)
//	}
//	if csrData, err := csr.Export(); err != nil {
//		panic(err)
//	} else {
//		writeToFile(exportFilePath, filename_signca_csr, csrData)
//	}
//	if keyData, err := key.ExportEncryptedPrivate(password); err != nil {
//		panic(err)
//	} else {
//		writeToFile(exportFilePath, filename_signca_key, keyData)
//	}
//
//	cert, err := pkix.CreateIntermediateCertificateAuthority(authCert, authKey, csr, time.Now().Add(days * 24 *time.Hour))
//	if err != nil {
//		panic(err)
//	}
//	if pemData, err := cert.Export(); err != nil {
//		panic(err)
//	} else {
//		writeToFile(exportFilePath, filename_signca_pem, pemData)
//	}
//}
//
//func NewServerCert2(serviceName string) (certData []byte, keyData []byte, err error) {
//	serverkey, err := pkix.CreateRSAKey(keybits)
//	if err != nil {
//		return nil, nil, err
//	}
//	commonName := serviceName + "." + _config.ProductName + "." + _config.CompanyName
//	serverCsr, err := pkix.CreateCertificateSigningRequest(serverkey, "", _config.Ips, _config.Domains, nil, "", "", "", "", commonName)
//	if err != nil {
//		return nil, nil, err
//	}
//	if csrData, err := serverCsr.Export(); err != nil {
//		return nil, nil, err
//	} else {
//		writeToFile(_config.ExportFilePath, serviceName + ".csr", csrData)
//	}
//	if keyData, err = serverkey.ExportPrivate(); err != nil {
//		return nil, nil, err
//	}
//	cert, err := pkix.CreateCertificateHost(_caCert, _caKey, serverCsr, time.Now().Add(days * 24 *time.Hour))
//	if err != nil {
//		return nil, nil, err
//	}
//	if certData, err = cert.Export(); err != nil {
//		return nil, nil, err
//	}
//	return certData, keyData, nil
//}
//
//func CAPem() []byte {
//	return _caPem
//}