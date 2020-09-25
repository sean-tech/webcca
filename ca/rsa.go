package ca

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/sean-tech/gokit/fileutils"
	"io/ioutil"
	"os"
	"strings"
)

type RSA_PERM_TYPE int
const (
	RSA_PERM_TYPE_SERVER RSA_PERM_TYPE = 1
	RSA_PERM_TYPE_CLIENT RSA_PERM_TYPE = 2
	server_pubkey_filename = "server_pubkey"
	server_prikey_filename = "server_prikey"
	client_pubkey_filename = "client_pubkey"
	client_prikey_filename = "client_prikey"
	pemkey_version_sep = "-"
)

type RSAKeyPair struct {
	ServerPubKey []byte
	ServerPriKey []byte
	ClientPubKey []byte
	ClientPriKey []byte
}

/**
 * get rsa keypair map
 * map key is version
 * map value is RSAKeyPair pointer value
 */
func GetRsaKeyPairMap(product string) (map[string]*RSAKeyPair, error) {
	fileinfos, err := GetAllRSAKeyFileInfos(product)
	if err != nil {
		return nil, err
	}
	var keypairmap = make(map[string]*RSAKeyPair)
	for _, f := range fileinfos {
		keypairname_version := strings.Split(f.Name(), pemkey_version_sep)
		if len(keypairname_version) != 2 {
			continue
		}
		keypairname := keypairname_version[0]
		version := keypairname_version[1]
		if _, ok := keypairmap[version]; !ok {
			keypairmap[version] = new(RSAKeyPair)
		}
		keypair := keypairmap[version]
		keydata, err := ioutil.ReadFile(rsa_pem_filepath(product) + f.Name())
		if err != nil {
			continue
		}
		switch keypairname {
		case server_pubkey_filename:
			keypair.ServerPubKey = keydata
		case server_prikey_filename:
			keypair.ServerPriKey = keydata
		case client_pubkey_filename:
			keypair.ClientPubKey = keydata
		case client_prikey_filename:
			keypair.ServerPriKey = keydata
		}
	}
	return keypairmap, nil
}

/**
 * get all versions of rsakey for client
 */
func GetRSAClientVersions(product string)  ([]string, error) {
	fileInofs, err :=  GetAllRSAKeyFileInfos(product)
	if err != nil {
		return nil, err
	}
	var versionMaps = make(map[string]string)
	for _, f := range fileInofs {
		if f.IsDir() {
			continue
		}
		if strings.Contains(f.Name(), server_prikey_filename) {
			continue
		}
		keypairname_version := strings.Split(f.Name(), pemkey_version_sep)
		if len(keypairname_version) != 2 {
			continue
		}
		version := keypairname_version[1]
		versionMaps[version] = version
	}
	var versions = []string{"all"}
	for k, _ := range versionMaps {
		versions = append(versions, k)
	}
	return versions, nil
}

/**
 * get all client rsakey files of some version
 */
func GetRSAClientVersionFiles(product, filterVersion string) ([]string, error) {
	fileInofs, err :=  GetAllRSAKeyFileInfos(product)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, f := range fileInofs {
		if f.IsDir() {
			continue
		}
		if strings.Contains(f.Name(), server_prikey_filename) {
			continue
		}
		keypairname_version := strings.Split(f.Name(), pemkey_version_sep)
		if len(keypairname_version) != 2 {
			continue
		}
		version := keypairname_version[1]
		if filterVersion == "all" {
			files = append(files, f.Name())
			continue
		}
		if version == filterVersion {
			files = append(files, f.Name())
		}
	}
	return files, nil
}

func GetRsaKeyPairFile(product, filename string) ([]byte, error) {
	return ioutil.ReadFile(rsa_pem_filepath(product) + filename)
}

/**
 * create rsa keypair for server and client
 * will not create new when some version keypair file exist
 */
func NewRsaKeyPair(product, version string, keypassword []byte) error {
	fileinfos, err := GetAllRSAKeyFileInfos(product)
	if err != nil {
		return err
	}
	for _, f := range fileinfos {
		keypairname_version := strings.Split(f.Name(), pemkey_version_sep)
		if len(keypairname_version) != 2 {
			continue
		}
		version := keypairname_version[1]
		if version == version {
			return nil
		}
	}
	if err := newRsaKeyPair(product, version, RSA_PERM_TYPE_SERVER, keypassword); err != nil {
		return err
	}
	if err := newRsaKeyPair(product, version, RSA_PERM_TYPE_CLIENT, keypassword); err != nil {
		return err
	}
	return nil
}

/**
 * get all fileinfo of rsakey from local file saved path
 */
func GetAllRSAKeyFileInfos(product string) ([]os.FileInfo, error) {
	if fileutils.CheckExist(rsa_pem_filepath(product)) == false {
		return nil, nil
	}
	files, err := ioutil.ReadDir(rsa_pem_filepath(product))
	if err != nil {
		return nil, err
	}
	return files, nil
}

/**
 * create rsa keypair
 */
func newRsaKeyPair(product, version string, perm_type RSA_PERM_TYPE, keypassword []byte) error {
	var pubkey_name string; var prikey_name string;
	switch perm_type {
	case RSA_PERM_TYPE_SERVER:
		pubkey_name = server_pubkey_filename
		prikey_name = server_prikey_filename
	case RSA_PERM_TYPE_CLIENT:
		pubkey_name = client_pubkey_filename
		prikey_name = client_prikey_filename
	}

	// privatekey gen
	privateKey, err := rsa.GenerateKey(rand.Reader, keybits)
	if err != nil {
		return err
	}
	// MarshalPKCS1PrivateKey converts a private key to ASN.1 DER encoded form.
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	//block := &pem.Block{
	//	Type:  "RSA PRIVATE KEY",
	//	Bytes: derStream,
	//}
	block, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", derStream, keypassword, x509.PEMCipher3DES)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := pem.Encode(buf, block); err != nil {
		return err
	}
	writeToFile(rsa_pem_filepath(product), rsa_pem_filename(prikey_name, version), buf.Bytes())

	// publickey gen
	publicKey := &privateKey.PublicKey
	// MarshalPKIXPublicKey serialises a public key to DER-encoded PKIX format.
	//derPub, err := x509.MarshalPKIXPublicKey(publicKey)
	// MarshalPKCS1PublicKey converts an RSA public key to PKCS#1, ASN.1 DER form.
	derPub := x509.MarshalPKCS1PublicKey(publicKey)
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPub,
	}
	buf.Reset()
	if err := pem.Encode(buf, block); err != nil {
		return err
	}
	writeToFile(rsa_pem_filepath(product), rsa_pem_filename(pubkey_name, version), buf.Bytes())
	return nil
}

/**
 * rsa keypair saved path
 */
func rsa_pem_filepath(product string) string {
	return _setting.CASavePath + "rsa/" + product + "/"
}
/**
 * rsa keypair file name
 */
func rsa_pem_filename(filename, version string) string {
	return filename + pemkey_version_sep + version + ".pem"
}