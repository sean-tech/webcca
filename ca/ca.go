package ca

import (
	"fmt"
	"github.com/sean-tech/gokit/fileutils"
	"github.com/square/certstrap/pkix"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

const (
	keybits = 2048
	days = time.Duration(1024)

	filename_rootca_pem = "rootca.pem"
	filename_rootca_key = "rootca.key"
	filename_signca_csr = "signca.csr"
	filename_signca_pem = "signca.pem"
	filename_signca_key = "signca.key"
)

type CASetting struct {
	CompanyName    	string
	CASavePath 		string
}

type CAConfig struct {
	ProductName    	string
	Ips            	[]net.IP	 // net.ParseIP("127.0.0.1")
	Domains        	[]string 	// "*." + "ex.liansheng"
	KeyPassword    	[]byte
	SignCACert 		*pkix.Certificate
	SignCAKey  		*pkix.Key
}

var (
	_setting CASetting
)

func Setup(setting CASetting)  {
	// setting
	_setting = setting
	if _setting.CompanyName == "" {
		panic("CompanyName for ca could not be nil.")
	}
	_setting.CASavePath = strings.ReplaceAll(_setting.CASavePath, " ", "")
	if strings.HasSuffix(_setting.CASavePath, "/") == false {
		_setting.CASavePath += "/"
	}
	if fileutils.CheckExist(_setting.CASavePath) == false {
		panic("CASavePath for ca not exist.")
	}
}

func LoadSignCA(product string, keypassword []byte, ips []net.IP, domains []string) ([]byte, *pkix.Certificate, *pkix.Key, error) {
	var signcapem []byte; var signcacert *pkix.Certificate; var signcakey *pkix.Key; var err error
	if fileutils.CheckExist(signca_pem_filepath(product) + filename_signca_pem) &&
		fileutils.CheckExist(signca_key_filepath(product) + filename_signca_key) {
		signcapem, signcacert, signcakey, err = loadSignCAFromFile(product, keypassword)
		return signcapem, signcacert, signcakey, err
	}
	var rootcacert *pkix.Certificate; var rootcakey *pkix.Key
	if fileutils.CheckExist(rootca_pem_filepath(product) + filename_rootca_pem) &&
		fileutils.CheckExist(rootca_key_filepath(product) + filename_rootca_key) {
		if rootcapem, err := ioutil.ReadFile(rootca_pem_filepath(product) + filename_rootca_pem); err != nil {
			return nil, nil, nil, err
		} else if rootcacert, err = pkix.NewCertificateFromPEM(rootcapem); err != nil {
			return nil, nil, nil, err
		}
		if rootcakeybts, err := ioutil.ReadFile(rootca_key_filepath(product) + filename_rootca_key); err != nil {
			return nil, nil, nil, err
		} else if rootcakey, err = pkix.NewKeyFromEncryptedPrivateKeyPEM(rootcakeybts, keypassword); err != nil {
			return nil, nil, nil, err
		}
	} else if rootcacert, rootcakey, err = newRootCA(product, keypassword); err != nil {
		return nil, nil, nil, err
	}
	if err := newSignCA(rootcacert, rootcakey, product, ips, domains, keypassword); err != nil {
		return nil, nil, nil, err
	}
	if signcapem, signcacert, signcakey, err = loadSignCAFromFile(product, keypassword); err != nil {
		return nil, nil, nil, err
	} else {
		return signcapem, signcacert, signcakey, nil
	}
}

func loadSignCAFromFile(product string, keypassword []byte) ([]byte, *pkix.Certificate, *pkix.Key, error) {
	var signcapem []byte; var signcacert *pkix.Certificate; var signcakey *pkix.Key; var err error
	if signcapem, err = ioutil.ReadFile(signca_pem_filepath(product) + filename_signca_pem); err != nil {
		return nil, nil, nil, err
	} else if signcacert, err = pkix.NewCertificateFromPEM(signcapem); err != nil {
		return nil, nil, nil, err
	}
	if signcakeybts, err := ioutil.ReadFile(signca_key_filepath(product) + filename_signca_key); err != nil {
		return nil, nil, nil, err
	} else if signcakey, err = pkix.NewKeyFromEncryptedPrivateKeyPEM(signcakeybts, keypassword); err != nil {
		return nil, nil, nil, err
	}
	return signcapem, signcacert, signcakey, nil
}

func newRootCA(product string, password []byte) (cert *pkix.Certificate, key *pkix.Key, err error) {
	// init root ca and export
	key, err = pkix.CreateRSAKey(keybits)
	if err != nil {
		return nil, nil, err
	}
	expiry := time.Now().Add(days * 24 *time.Hour)
	cert, err = pkix.CreateCertificateAuthority(key, "", expiry, "", "", "", "", _setting.CompanyName)
	if err != nil {
		return nil, nil, err
	}
	if pemData, err := cert.Export(); err != nil {
		return nil, nil, err
	} else if err := writeToFile(rootca_pem_filepath(product), filename_rootca_pem, pemData); err != nil {
		return nil, nil, err
	}
	if keyData, err := key.ExportEncryptedPrivate(password); err != nil {
		return nil, nil, nil
	} else if err := writeToFile(rootca_key_filepath(product), filename_rootca_key, keyData); err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

func newSignCA(authCert *pkix.Certificate, authKey *pkix.Key, product string, ips []net.IP, domains []string, password []byte) error {
	key, err := pkix.CreateRSAKey(keybits)
	if err != nil {
		return err
	}
	commonName := GetSignCACommonName(product)
	csr, err := pkix.CreateCertificateSigningRequest(key, "", ips, domains, nil, "", "", "", "", commonName)
	if err != nil {
		return err
	}
	if csrData, err := csr.Export(); err != nil {
		return err
	} else if err := writeToFile(signca_csr_filepath(product), filename_signca_csr, csrData); err != nil {
		return err
	}
	if keyData, err := key.ExportEncryptedPrivate(password); err != nil {
		return err
	} else if err := writeToFile(signca_key_filepath(product), filename_signca_key, keyData); err != nil {
		return err
	}

	cert, err := pkix.CreateIntermediateCertificateAuthority(authCert, authKey, csr, time.Now().Add(days * 24 *time.Hour))
	if err != nil {
		return err
	}
	if pemData, err := cert.Export(); err != nil {
		return err
	} else if err := writeToFile(signca_pem_filepath(product), filename_signca_pem, pemData); err != nil {
		return err
	}
	return nil
}

func GetSignCACertData(product string) ([]byte, error) {
	var signcacert *pkix.Certificate
	if signcapem, err := ioutil.ReadFile(signca_pem_filepath(product) + filename_signca_pem); err != nil {
		return nil, err
	} else if signcacert, err = pkix.NewCertificateFromPEM(signcapem); err != nil {
		return nil, err
	}
	if certData, err := signcacert.Export(); err != nil {
		return nil, err
	} else {
		return certData, nil
	}
}

func GetSignCACommonName(product string) string {
	return product + "." + _setting.CompanyName
}

func NewServerCert(product, service string, caconfig CAConfig) (certData []byte, keyData []byte, err error) {
	var signCACert *pkix.Certificate = caconfig.SignCACert
	var signCAKey *pkix.Key = caconfig.SignCAKey
	var ips []net.IP = caconfig.Ips
	var domains []string = caconfig.Domains
	// gen key & cert
	serverkey, err := pkix.CreateRSAKey(keybits)
	if err != nil {
		return nil, nil, err
	}
	commonName := service + "." + product + "." + _setting.CompanyName
	serverCsr, err := pkix.CreateCertificateSigningRequest(serverkey, "", ips, domains, nil, "", "", "", "", commonName)
	if err != nil {
		return nil, nil, err
	}
	if csrData, err := serverCsr.Export(); err != nil {
		return nil, nil, err
	} else {
		writeToFile(serviceca_csr_filepath(product), service + ".csr", csrData)
	}
	if keyData, err = serverkey.ExportPrivate(); err != nil {
		return nil, nil, err
	}
	cert, err := pkix.CreateCertificateHost(signCACert, signCAKey, serverCsr, time.Now().Add(days * 24 *time.Hour))
	if err != nil {
		return nil, nil, err
	}
	if certData, err = cert.Export(); err != nil {
		return nil, nil, err
	}
	return certData, keyData, nil
}

func rootca_pem_filepath(product string) string {
	return _setting.CASavePath + "tls/" + product + "/"
}
func rootca_key_filepath(product string) string {
	return _setting.CASavePath + "tls/" + product + "/"
}
func signca_pem_filepath(product string) string {
	return _setting.CASavePath + "tls/" + product + "/"
}
func signca_key_filepath(product string) string {
	return _setting.CASavePath + "tls/" + product + "/"
}
func signca_csr_filepath(product string) string {
	return _setting.CASavePath + "tls/" + product + "/"
}
func serviceca_csr_filepath(product string) string {
	return _setting.CASavePath + "tls/" + product + "/"
}



func GenerateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 32)
	copy(genKey, key)
	for i := 32; i < len(key); {
		for j := 0; j < 32 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func writeToFile(filePath, fileName string, data []byte) error {
	file, err := openFile(filePath, fileName)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

func openFile(filePath, fileName string) (*os.File, error) {
	src := filePath
	perm := fileutils.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err := fileutils.MKDirIfNotExist(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := fileutils.Open(src + fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}
	return f, nil
}