package synology

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/Luzifer/go-openssl/v4"
	"github.com/prometheus/common/log"
	"math/big"
	"net/url"
	"synology-videostation-reindexer/synology/api"
	config2 "synology-videostation-reindexer/synology/config"
	"synology-videostation-reindexer/synology/data"
)

type SynoStuff struct {
	config *config2.Config
	api    api.Api
}

type synoSession struct {
}

func NewSynoStuff(config *config2.Config) *SynoStuff {
	api:= api.NewSynoAPI(config)
	return &SynoStuff{config: config, api:api}
}

func (syno *SynoStuff) Update() error {
	_, err := syno.login()
	return err
}

func (syno *SynoStuff) login() (*synoSession, error) {
	encryptionInfo, err := syno.getEncryptionInfo()
	if err != nil{
		return nil, err
	}

	n, ok := new(big.Int).SetString(encryptionInfo.PublicKey, 16)
	if !ok {
		panic("failed to convert hex to bigInt")
	}

	pubkey := &rsa.PublicKey{
		N: n,
		E: 65537, // value plucked from https://github.com/openstack/cinder/blob/61d506880eafcfcfee9047ac182a45d555e32a22/cinder/volume/drivers/synology/synology_common.py#L243
	}

	bytes := make([]byte, 50) //generate a random key for AES-256
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
	params.Add(encryptionInfo.CipherToken, fmt.Sprintf("%d", encryptionInfo.ServerTime))
	encodedParams := []byte(params.Encode())

	o := openssl.New()
	enc, err := o.EncryptBytes(pass, encodedParams, openssl.BytesToKeyMD5)
	if err != nil {
		panic(fmt.Sprintf("An error occurred: %s\n", err))
	}

	const loginUrl = `%s/webapi/auth.cgi`

	req := data.ReqLogin{
		Api:     "SYNO.API.Auth",
		Method:  "login",
		Version: 6,
	}
	options:= api.NewOptions().AddParam(
		encryptionInfo.CipherKey,
		`{"rsa": "`+base64.StdEncoding.EncodeToString(encryptedPass)+`","aes": "`+string(enc)+`"}`,
		)

	resp:= &data.RespLogin{}
	err = syno.api.Request(loginUrl, req, resp, options)
	if err!= nil{
		return nil,err
	}
	fmt.Printf("%+v\n",*resp)

	return nil, nil
}

func (syno *SynoStuff) getEncryptionInfo() (data.RespEncryption, error) {
	const encryptionInfoURL = `%s/webapi/encryption.cgi`
	req := data.ReqEncryption{
		Api:     "SYNO.API.Encryption",
		Method:  "getinfo",
		Version: 1,
		Format:  "module",
	}
	info := &data.RespEncryption{}
	err := syno.api.Request(encryptionInfoURL, req, info)
	if err != nil {
		panic(err.Error())
	}
	return *info , nil
}
