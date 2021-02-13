package session

import (
	"fmt"
	"github.com/Luzifer/go-openssl/v4"
	"github.com/sirupsen/logrus"
	"net/url"
	"synology-videostation-reindexer/synology/config"
	"synology-videostation-reindexer/synology/internal/api"
	"synology-videostation-reindexer/synology/internal/data"
	"time"
)

//synoSession this automates the login part of requests made to the nas box.
//And will also prevent extreme long login-times
type synoSession struct {
	isLoggedIn   bool
	sid          string
	keepLoggedIn chan struct{}
	endSession   chan struct{}
	api          api.Api
	config       *config.Config
	name         string
}

func NewSynoSession(config *config.Config, name string) (ses *synoSession) {
	ses = &synoSession{
		config:       config,
		api:          api.NewSynoAPI(config),
		isLoggedIn:   false,
		keepLoggedIn: make(chan struct{}, 0),
		endSession:   make(chan struct{}, 0),
		name:         name,
	}
	go ses.logoutTimer()
	return
}

func (s *synoSession) Request(urlDest string, req interface{}, resp interface{}, optionss ...api.Options) error {
	if !s.isLoggedIn {
		err := s.login()
		if err != nil {
			return fmt.Errorf("failed to do request, login failed: %w", err)
		}
	}
	options := api.Options{}
	if len(optionss) == 1 {
		options = optionss[0]
	}

	options.AddParam("_sid", s.sid)

	s.keepLoggedIn <- struct{}{}

	return s.api.Request(urlDest, req, resp, options)
}

func (s *synoSession) EndSession() {
	if err := s.logout(); err != nil {
		panic(err)
	}
	s.endSession <- struct{}{}
	s.name = ""
	close(s.keepLoggedIn)
	close(s.endSession)
	s.api = nil
	s.config = nil
}

func (s *synoSession) login() error {
	encryptionInfo, err := s.getEncryptionInfo()
	if err != nil {
		return err
	}
	pass, encPass, err := createPassAndEncryptedPass(encryptionInfo)
	if err != nil {
		return err
	}

	// create login params
	params := url.Values{}
	params.Add("account", s.config.UserName)
	params.Add("passwd", s.config.UserPassword)
	params.Add("session", "dsm")
	params.Add("format", "sid")
	params.Add(encryptionInfo.CipherToken, fmt.Sprintf("%d", encryptionInfo.ServerTime))
	encodedParams := []byte(params.Encode())

	// Encode the login parameters
	o := openssl.New()
	encParam, err := o.EncryptBytes(pass, encodedParams, openssl.BytesToKeyMD5)
	if err != nil {
		panic(fmt.Sprintf("An error occurred: %s\n", err))
	}

	request, err := s.doLoginRequest(encryptionInfo, encPass, encParam)
	if err != nil {
		return err
	}
	s.sid = request.Sid
	s.isLoggedIn = true
	s.keepLoggedIn <- struct{}{}
	logrus.Infof("successfully logged in to dsm with session name: %s", s.name)
	return nil
}

func (s *synoSession) logout() error {
	if !s.isLoggedIn {
		return nil
	}
	const authUrl = `%s/webapi/auth.cgi`

	req := data.ReqAuth{
		Api:     "SYNO.API.Auth",
		Method:  "logout",
		Version: 6,
		Session: s.name,
	}

	resp := &data.RespAuth{}
	err := s.api.Request(authUrl, req, resp)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	logrus.Infof("successfully logged out of to dsm for session name: %s", s.name)
	s.isLoggedIn = false
	s.sid = ""
	return nil
}

func (s *synoSession) logoutTimer() {
	logout := time.Minute * time.Duration(s.config.AutoLogOutSession)
	for {
		select {
		case <-time.After(logout):
			err := s.logout()
			if err != nil {
				panic(err)
			}
		case <-s.keepLoggedIn:
		case <-s.endSession:
			return
		}
	}
}
