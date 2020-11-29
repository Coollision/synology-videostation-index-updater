package api

import (
	"net/http"
	"net/url"
)

//Options ...
type Options struct {
	//authToken         string // auth field can just be in params...
	additionalParams  map[string]string
	additionalHeaders map[string]string
}

func NewOptions() *Options { return &Options{} }

func (o *Options) AddParam(key, value string) *Options {
	if o.additionalParams == nil {
		o.additionalParams = map[string]string{}
	}
	o.additionalParams[key] = value
	return o
}

func (o *Options) AddHeader(key, value string) *Options {
	if o.additionalHeaders == nil {
		o.additionalHeaders = map[string]string{}
	}
	o.additionalHeaders[key] = value
	return o
}

//func (o *Options) SetAuth(authToken string) *Options {
//	o.authToken = authToken
//	return o
//}

func (o *Options) addParams(values *url.Values) {
	if o.additionalParams != nil {
		for ak, ap := range o.additionalParams {
			values.Add(ak, ap)
		}
	}
}

func (o *Options) addHeaders(req *http.Request) {
	//if o.authToken != "" {
	//	req.Header.Add("x-syno-token", o.authToken)
	//}
	if o.additionalHeaders != nil {
		for ahk, ah := range o.additionalHeaders {
			req.Header.Add(ahk, ah)
		}
	}
}
