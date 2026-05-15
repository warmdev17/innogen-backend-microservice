package piston

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// pistonRequest is the JSON body sent to the Piston execute endpoint.
type pistonRequest struct {
	Language string       `json:"language"`
	Version  string       `json:"version"`
	Files    []pistonFile `json:"files"`
	Stdin    string       `json:"stdin"`
}

// pistonFile represents a single source file in a Piston request.
type pistonFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// pistonResponse is the top-level JSON response from the Piston execute endpoint.
type pistonResponse struct {
	Language string       `json:"language"`
	Version  string       `json:"version"`
	Run      *pistonStage `json:"run"`
	Compile  *pistonStage `json:"compile"`
}

// pistonStage represents a single stage (compile or run) in a Piston response.
type pistonStage struct {
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	Code     int     `json:"code"`
	Signal   *string `json:"signal"`
	Output   string  `json:"output"`
	CPUTime  float64 `json:"cpu_time"`
	WallTime float64 `json:"wall_time"`
}

// Client is an HTTP client for the Piston code execution API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Piston API client with the given base URL.
// The base URL trailing slash is trimmed.
func NewClient(baseURL string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Execute sends a code execution request to the Piston API.
func (c *Client) Execute(ctx context.Context, language, version, fileName, code, stdin string) (*pistonResponse, error) {
	req := pistonRequest{
		Language: language,
		Version:  version,
		Files: []pistonFile{
			{Name: fileName, Content: code},
		},
		Stdin: stdin,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("piston: failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v2/execute", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("piston: failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("piston: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("piston: failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("piston: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var result pistonResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("piston: failed to decode response: %w", err)
	}

	return &result, nil
}
