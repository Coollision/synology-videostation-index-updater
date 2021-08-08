package response

import (
	"net/http"
)

type Bytes []byte


// Marshal >>>
func (js Bytes) Marshal() ([]byte, error) {
	return js, nil
}

// Render >>>
func (js Bytes) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
