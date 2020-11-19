package synology

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Luzifer/go-openssl/v4"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
)

type SynoStuff struct {
	config *Config
}

type synoSession struct {
}

func NewSynoStuff(config *Config) *SynoStuff {
	return &SynoStuff{config: config}
}

func (syno *SynoStuff) Update() error {
	_, _ = syno.login()
	return nil
}

func (syno *SynoStuff) login() (*synoSession, error) {
	encryptionInfo := syno.getEncryptionInfo()

	n, ok := new(big.Int).SetString(encryptionInfo.PublicKey, 16)
	if !ok {
		panic("failed to convert hex to bigInt")
	}

	pubkey := &rsa.PublicKey{
		N: n,
		E: 65537, // value plucked from https://github.com/openstack/cinder/blob/61d506880eafcfcfee9047ac182a45d555e32a22/cinder/volume/drivers/synology/synology_common.py#L243
	}

	bytes := make([]byte, 50) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}
	pass := hex.EncodeToString(bytes)

	encryptedPass, err := rsa.EncryptPKCS1v15(rand.Reader, pubkey, []byte(pass))
	if err != nil {
		log.Error(err)
	}

	params := url.Values{}
	params.Add("account", syno.config.UserName)
	params.Add("passwd", syno.config.UserPassword)
	params.Add("session", "dsm")
	params.Add("format", "sid")
	params.Add(encryptionInfo.Ciphertoken, fmt.Sprintf("%d", encryptionInfo.ServerTime))
	encodedParams := []byte(params.Encode())

	o := openssl.New()

	enc, err := o.EncryptBytes(pass, encodedParams, openssl.BytesToKeyMD5)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
	}

	const loginUrl = `%s/webapi/auth.cgi`

	data := url.Values{}
	data.Add("api", "SYNO.API.Auth")
	data.Add("method", "login")
	data.Add("version", "6")
	data.Add(encryptionInfo.Cipherkey,
		`{"rsa": "`+base64.StdEncoding.EncodeToString(encryptedPass)+`","aes": "`+string(enc)+`"}`)


	resp, err := http.PostForm(fmt.Sprintf(loginUrl, syno.config.Url), data)
	if err != nil {
		panic(err.Error())
	}

	respData, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respData))

	return nil, nil
}

func (syno *SynoStuff) getEncryptionInfo() *encryptionInfo {
	const encryptionInfoURL = `%s/webapi/encryption.cgi`

	data := url.Values{}
	data.Add("api", "SYNO.API.Encryption")
	data.Add("method", "getinfo")
	data.Add("version", "1")
	data.Add("format", "module")

	get, err := http.PostForm(fmt.Sprintf(encryptionInfoURL, syno.config.Url), data)
	if err != nil {
		panic(err.Error())
	}
	info := &encryptionInfo{}
	err = json.NewDecoder(get.Body).Decode(info)
	if err != nil {
		panic(err.Error())
	}
	return info
}
