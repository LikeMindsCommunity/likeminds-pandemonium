package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"likeminds-pandemonium/common"
	"net/http"
	"time"
)

type APIClientCaravan struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAPIClient initializes a new API client instance
func NewAPIClientCaravan() *APIClientCaravan {
	return &APIClientCaravan{
		BaseURL: common.GoDotEnvVariable(common.DotEnvVarBaseUrlCaravan),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Post makes a POST request to the external API
func (c *APIClientCaravan) Post(endpoint string, body any) error {
	reqURL := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Execute the request
	_, err = c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
