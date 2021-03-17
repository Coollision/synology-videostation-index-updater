package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"synology-videostation-reindexer/synology/config"
	"synology-videostation-reindexer/synology/internal/data"
)

type Api interface {
	//Request
	//{url} needs to start with % sign, because it will be Sprinted with the baseURL
	//{req} the struct containing the request data, preferably form inside the data package.
	Request(urlDest string, req interface{}, resp interface{}, options ...Options) error
}

type api struct {
	client *http.Client
	config  *config.Config
	encoder *form.Encoder
}

func NewSynoAPI(config *config.Config) *api {
	encoder := form.NewEncoder()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SkipSSLVerification},
	}
	client := &http.Client{Transport: tr}
	return &api{config: config, encoder: encoder, client: client}

}

//goland:noinspection ALL
func (api *api) Request(urlDest string, req interface{}, resp interface{}, optionss ...Options) error {
	options := Options{}
	if len(optionss) == 1 {
		options = optionss[0]
	}

	dataEncoded, err := api.encoder.Encode(req)
	if err != nil {
		return err
	}
	options.addParams(&dataEncoded)

	url := fmt.Sprintf(urlDest, api.config.Url)
	logrus.WithField("syno","api").Trace(dataEncoded.Encode())
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(dataEncoded.Encode()))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//post, err := api.client.PostForm(url, dataEncoded)
	post, err := api.client.Do(request)
	if err != nil {
		return err
	}
	response := &data.Resp{}

	////for debugging all responses from syno
	//dump, err := httputil.DumpResponse(post, true)
	//if err != nil {
	//	panic(err.Error())
	//}
	//logrus.WithField("syno","api").Traceln(string(dump))


	err = json.NewDecoder(post.Body).Decode(response)
	if err != nil {
		return err
	}

	if !response.Success {
		logrus.WithField("syno", "api").Errorf("failed to get a successful response but got reason: %s", response.Reason)
		return fmt.Errorf("failed to get a successful response but got reason: %s", response.Reason)
	}

	err = mapstructure.Decode(response.Data, resp)
	logrus.WithField("syno", "api").Tracef("%#v", response)

	return err
}
