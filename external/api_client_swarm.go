package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/constant"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/common"
	"io"
	"net/http"
	"time"
)

type APIClientSwarm struct {
	BaseURL    string
	HTTPClient *http.Client
}

type APIResponseSwarm struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"-"`
}

// NewAPIClient initializes a new API client instance
func NewAPIClientSwarm() *APIClientSwarm {
	return &APIClientSwarm{
		BaseURL: common.GoDotEnvVariable(common.DotEnvVarBaseUrlSwarm),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Post makes a POST request to the external API
func (c *APIClientSwarm) Post(endpoint string, headers constant.ApiHeaders, body any) (*APIResponseSwarm, error) {
	reqURL := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse swarm request data, err=%s", err)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create swarm request, err=%s", err)
	}

	req.Header.Add(constant.HeadersApiKey, headers.HeadersApiKey)
	req.Header.Add(constant.HeadersMemberID, headers.HeadersMemberID)
	req.Header.Add(constant.HeadersPlatformType, headers.HeadersPlatformType)

	// Execute the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to receive swarm response, err=%s", err)
	}
	defer resp.Body.Close()

	return parseResponse(resp)
}

func parseResponse(resp *http.Response) (*APIResponseSwarm, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read swarm response data, err=%s", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed at swarm, err=%s", body)
	}

	var apiResp APIResponseSwarm
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse swarm response, err=%s", err)
	}

	if err := json.Unmarshal(body, &apiResp.Data); err != nil {
		return nil, fmt.Errorf("failed to parse swarm response data, err=%s", err)
	}

	return &apiResp, nil
}
