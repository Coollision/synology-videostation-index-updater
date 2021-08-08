package testing

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"synology-videostation-reindexer/synology/internal/api"
	"synology-videostation-reindexer/synology/internal/data"
)

var log = logrus.WithField("", "")

type testingAPI struct {
	api    api.Api
}

func NewTestRequests(api api.Api) *testingAPI {
	return &testingAPI{ api: api}
}

func (v *testingAPI) test() ([]byte, error) {
	url := "%s/webapi/entry.cgi"
	req := data.Req{
		Api:     "SYNO.Core.Desktop.Timeout",
		Method:  "get",
		Version: 1,
	}
	resp := &map[string]interface{}{}
	err := v.api.Request(url, req, resp)
	if err != nil {
		return nil , fmt.Errorf("failed to list Library Types: %w", err)
	}

	str, _ := json.Marshal(resp)

	return str , nil
}
