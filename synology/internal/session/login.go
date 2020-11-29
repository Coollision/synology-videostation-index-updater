package session

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/big"
	"synology-videostation-reindexer/synology/internal/api"
	"synology-videostation-reindexer/synology/internal/data"
)

func (s *synoSession) getEncryptionInfo() (data.RespEncryption, error) {
	const encryptionInfoURL = `%s/webapi/encryption.cgi`
	req := data.ReqEncryption{
		Api:     "SYNO.API.Encryption",
		Method:  "getinfo",
		Version: 1,
		Format:  "module",
	}
	info := &data.RespEncryption{}
	err := s.api.Request(encryptionInfoURL, req, info)
	if err != nil {
		panic(err.Error())
	}
	return *info, nil
}

func createPassAndEncryptedPass(encryptionInfo data.RespEncryption) (pass string, encryptedPass []byte, err error) {
	n, ok := new(big.Int).SetString(encryptionInfo.PublicKey, 16)
	if !ok {
		logrus.WithField("session", "convertPublicKey").Debug("failed to convert hex to bigInt")
		return "", nil, fmt.Errorf("failed to convert hex to bigInt")
	}

	pubkey := &rsa.PublicKey{
		N: n,
		E: 65537, // value plucked from https://github.com/openstack/cinder/blob/61d506880eafcfcfee9047ac182a45d555e32a22/cinder/volume/drivers/synology/synology_common.py#L243
	}

	bytes := make([]byte, 50) //generate a random key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		return "", nil, err
	}
	pass = hex.EncodeToString(bytes)

	encryptedPass, err = rsa.EncryptPKCS1v15(rand.Reader, pubkey, []byte(pass))
	if err != nil {
		logrus.WithField("session", "createEncryptedPass").Debug(err)
		return "", nil, fmt.Errorf("failed to create encrypted pass: %w", err)
	}

	return pass, encryptedPass, nil
}

func (s *synoSession) doLoginRequest(encInfo data.RespEncryption, encPass, encParam []byte) (*data.RespAuth, error) {
	const authUrl = `%s/webapi/auth.cgi`

	req := data.ReqAuth{
		Api:     "SYNO.API.Auth",
		Method:  "login",
		Version: 6,
		Session: s.name,
	}
	options := api.NewOptions().
		AddParam(encInfo.CipherKey, `{"rsa": "`+base64.StdEncoding.EncodeToString(encPass)+`","aes": "`+string(encParam)+`"}`)

	resp := &data.RespAuth{}
	err := s.api.Request(authUrl, req, resp, *options)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return resp, nil
}
