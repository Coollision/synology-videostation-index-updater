package request

import (
	"encoding/json"
	"net/http"
)

type Test map[string]interface{}

func (t Test) Bind(r *http.Request) error {
	return nil
}

func (t Test) String() string {
	marshal, _ := json.Marshal(t)
	return string(marshal)
}
