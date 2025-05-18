package cms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	apiKey     string
	httpClient HTTPDoer
	baseURL    string
}

type Content struct {
	ID string `json:"id"`
}

type PublishRequest struct {
	Title   string `json:"title"`
	Tags    string `json:"tags"`
	QiitaID string `json:"qiitaId"`
	Content string `json:"content"`
}

type CheckExistsResponse struct {
	TotalCount int       `json:"totalCount"`
	Contents   []Content `json:"contents"`
}

func NewClient(serviceID, apiKey, endpoint string, httpClient HTTPDoer) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: httpClient,
		baseURL:    fmt.Sprintf("https://%s.microcms.io/api/v1/%s", serviceID, endpoint),
	}
}

func (c *Client) Create(ctx context.Context, title, tags, qiitaID, content string) error {
	req := PublishRequest{
		Title:   title,
		Tags:    tags,
		QiitaID: qiitaID,
		Content: content,
	}

	return c.sendRequest(ctx, http.MethodPost, c.baseURL, req, nil)
}

func (c *Client) Update(ctx context.Context, id, title, tags, qiitaId, content string) error {
	apiUrl := fmt.Sprintf("%s/%s", c.baseURL, id)
	req := PublishRequest{
		Title:   title,
		Tags:    tags,
		QiitaID: qiitaId,
		Content: content,
	}

	return c.sendRequest(ctx, http.MethodPatch, apiUrl, req, nil)
}

func (c *Client) CheckExists(ctx context.Context, qiitaID string) (bool, string, error) {
	rawFilter := fmt.Sprintf("qiitaId[equals]%s", qiitaID)
	encodedFilter := url.QueryEscape(rawFilter)
	apiUrl := fmt.Sprintf("%s?filters=%s", c.baseURL, encodedFilter)

	var response CheckExistsResponse
	if err := c.sendRequest(ctx, http.MethodGet, apiUrl, nil, &response); err != nil {
		return false, "", err
	}

	if response.TotalCount > 0 && len(response.Contents) > 0 {
		return true, response.Contents[0].ID, nil
	}

	return false, "", nil
}

func (c *Client) sendRequest(ctx context.Context, method, url string, requestBody, responseBody interface{}) error {
	var body io.Reader
	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-MICROCMS-API-KEY", c.apiKey)
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if responseBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}
	return nil
}
