package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	BaseURL    string
	Headers    map[string]string
	HTTPClient *http.Client
}

var DefaultClient = &http.Client{
	Timeout: time.Second * 15,
}

func InitClient(baseURL string, timeout time.Duration) Request {
	httpClient := DefaultClient
	if timeout > 0 {
		httpClient = &http.Client{Timeout: timeout}
	}
	return Request{
		Headers:    make(map[string]string),
		HTTPClient: httpClient,
		BaseURL:    baseURL,
	}
}

func (r *Request) Execute(ctx context.Context, method string, path string, body interface{}) (string, error) {
	requestURL, err := r.getBase(path)
	if err != nil {
		return "", err
	}
	payload, err := getBody(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, method, requestURL, payload)
	if err != nil {
		return "", fmt.Errorf("create request %s %s: %w", method, requestURL, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	res, err := r.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute request %s %s: %w", method, requestURL, err)
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read response %s %s: %w", method, requestURL, err)
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("request %s %s returned %s: %s", method, requestURL, res.Status, strings.TrimSpace(string(read)))
	}
	return string(read), nil
}

func (r *Request) getBase(requestPath string) (string, error) {
	base, err := url.Parse(r.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL %q: %w", r.BaseURL, err)
	}
	if base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("invalid base URL %q", r.BaseURL)
	}
	if requestPath == "" {
		return base.String(), nil
	}
	relative, err := url.Parse(requestPath)
	if err != nil {
		return "", fmt.Errorf("invalid request path %q: %w", requestPath, err)
	}
	if relative.IsAbs() || relative.Host != "" {
		return "", fmt.Errorf("request path must be relative: %q", requestPath)
	}
	base.Path = joinURLPath(base.Path, relative.Path)
	base.RawQuery = relative.RawQuery
	return base.String(), nil
}

func getBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, fmt.Errorf("encode request body: %w", err)
	}
	return buf, nil
}

func joinURLPath(basePath, requestPath string) string {
	if basePath == "" {
		return "/" + strings.TrimLeft(requestPath, "/")
	}
	return strings.TrimRight(basePath, "/") + "/" + strings.TrimLeft(requestPath, "/")
}
