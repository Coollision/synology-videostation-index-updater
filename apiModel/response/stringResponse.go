package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type String string

// Marshal >>>
func (sr String) Marshal() ([]byte, error) {
	data, err := json.Marshal(sr)
	if err != nil {
		return nil, fmt.Errorf("failed to Marshal ReportConfig: %v", err)
	}

	return data, nil
}

// Render >>>
func (sr String) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
