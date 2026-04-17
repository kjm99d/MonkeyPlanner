package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// mcpHTTPClient forwards MCP tool calls to the running HTTP server.
// This ensures SSE events are published through the server's broker.
type mcpHTTPClient struct {
	base string
	hc   *http.Client
}

func (c *mcpHTTPClient) get(path string, q url.Values) (json.RawMessage, error) {
	u := c.base + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	resp, err := c.hc.Get(u)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	return readHTTPBody(resp)
}

func (c *mcpHTTPClient) post(path string, body any) (json.RawMessage, error) {
	return c.doJSON("POST", path, body)
}

func (c *mcpHTTPClient) patch(path string, body any) (json.RawMessage, error) {
	return c.doJSON("PATCH", path, body)
}

func (c *mcpHTTPClient) doJSON(method, path string, body any) (json.RawMessage, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, c.base+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer resp.Body.Close()
	return readHTTPBody(resp)
}

func readHTTPBody(resp *http.Response) (json.RawMessage, error) {
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return json.RawMessage(b), nil
}
