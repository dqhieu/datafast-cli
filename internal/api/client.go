package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dqhieu/datafast-cli/internal/config"
)

const BaseURL = "https://datafa.st/api/v1"

type Client struct {
	BaseURL    string
	APIKey     string
	WebsiteID  string
	HTTPClient *http.Client
}

type APIResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
	Error  *APIError       `json:"error,omitempty"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

func New(websiteIDOverride string) (*Client, error) {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, err
	}

	wID := websiteIDOverride
	if wID == "" {
		wID, err = config.GetWebsiteID()
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		BaseURL:   BaseURL,
		APIKey:    apiKey,
		WebsiteID: wID,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *Client) do(method, path string, body io.Reader) (json.RawMessage, error) {
	url := c.BaseURL + path
	if strings.Contains(path, "?") {
		url += "&website_id=" + c.WebsiteID
	} else {
		url += "?website_id=" + c.WebsiteID
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiResp APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err == nil && apiResp.Error != nil {
			return nil, apiResp.Error
		}
		return nil, &APIError{Code: resp.StatusCode, Message: string(respBody)}
	}

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return respBody, nil
	}

	if apiResp.Status == "error" && apiResp.Error != nil {
		return nil, apiResp.Error
	}

	if apiResp.Data != nil {
		return apiResp.Data, nil
	}
	return respBody, nil
}

func (c *Client) Get(path string) (json.RawMessage, error) {
	return c.do(http.MethodGet, path, nil)
}

func (c *Client) Post(path string, body io.Reader) (json.RawMessage, error) {
	return c.do(http.MethodPost, path, body)
}

func (c *Client) Delete(path string) (json.RawMessage, error) {
	return c.do(http.MethodDelete, path, nil)
}
