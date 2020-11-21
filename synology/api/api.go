package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"synology-videostation-reindexer/synology/config"
	"synology-videostation-reindexer/synology/data"
)

type Api interface {
	//Request
	//{url} needs to start with % sign, because it will be Sprintfed with the baseURL
	//{req} the struct containing the request data, preferably form inside the data package.
	Request(urlDest string, req interface{}, resp interface{}, options ...options) error
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

//options ...
type options struct {
	auth             string
	AdditionalParams map[string]string
}

func NewOptions() options { return options{} }

func (o options) AddParam(key, value string) options{
	if o.AdditionalParams == nil{
		o.AdditionalParams = map[string]string{}
	}
	o.AdditionalParams[key]=value
	return o
}

func (o options) addParams(values url.Values) url.Values {
	if o.AdditionalParams != nil {
		for ak, ap := range o.AdditionalParams {
			values.Add(ak, ap)
		}
	}
	return values
}

func (api *api) Request(urlDest string, req interface{}, resp interface{}, optionss ...options) error {
	var options options
	if len(optionss) == 1 {
		options = optionss[0]
	}
	dataEncoded, err := api.encoder.Encode(req)
	if err != nil {
		return err
	}

	dataEncoded = options.addParams(dataEncoded)

	url := fmt.Sprintf(urlDest, api.config.Url)

	post, err := api.client.PostForm(url, dataEncoded)
	if err != nil {
		return err
	}
	response := &data.Resp{}

	////for debugging all requests to syno
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
		logrus.WithField("syno", "api").Errorf("response is %+v", response)
		return errors.New("failed to get a successful response")
	}

	err = mapstructure.Decode(response.Data, resp)

	return err
}
