package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StringArray []string


// Marshal >>>
func (sr StringArray) Marshal() ([]byte, error) {
	data, err := json.Marshal(sr)
	if err != nil {
		return nil, fmt.Errorf("failed to Marshal ReportConfig: %v", err)
	}

	return data, nil
}

// UnmarshalReportConfig >>>
func UnmarshalStringArrayResponse(data []byte) (*StringArray, error) {
	var test []string
	var config = StringArray(test)
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal ReportConfig: %v", err)
	}

	return &config, nil
}


// Render >>>
func (sr StringArray) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
