package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRequest_getBase(t *testing.T) {
	tests := []struct {
		name    string
		BaseURL string
		path    string
		want    string
		wantErr bool
	}{
		{"get api url with path", "http://blockatlas.com", "api", "http://blockatlas.com/api", false},
		{"get api url", "http://blockatlas.com", "", "http://blockatlas.com", false},
		{"keep base path", "http://blockatlas.com/base", "api?x=1", "http://blockatlas.com/base/api?x=1", false},
		{"reject missing host", "", "api", "", true},
		{"reject absolute request path", "http://blockatlas.com", "https://example.com/api", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				BaseURL: tt.BaseURL,
			}
			got, err := r.getBase(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getBase() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBody(t *testing.T) {
	tests := []struct {
		name string
		body interface{}
		want string
	}{
		{"nil body", nil, ""},
		{"json body", map[string]interface{}{"ok": true}, "{\"ok\":true}\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuf, err := getBody(tt.body)
			if err != nil {
				t.Fatalf("getBody() error = %v", err)
			}
			if gotBuf == nil {
				if tt.want != "" {
					t.Fatalf("getBody() = nil, want %q", tt.want)
				}
				return
			}
			got, err := io.ReadAll(gotBuf)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("getBody() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestRequest_Execute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api" {
			t.Errorf("path = %s, want /api", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %s, want application/json", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		if string(body) != "{\"probe\":true}\n" {
			t.Errorf("body = %q", string(body))
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	request := InitClient(server.URL, time.Second)
	got, err := request.Execute(context.Background(), http.MethodPost, "api", map[string]bool{"probe": true})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got != `{"ok":true}` {
		t.Errorf("Execute() = %q", got)
	}
}

func TestRequest_ExecuteStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "down", http.StatusBadGateway)
	}))
	defer server.Close()

	request := InitClient(server.URL, time.Second)
	_, err := request.Execute(context.Background(), http.MethodGet, "", nil)
	if err == nil {
		t.Fatal("Execute() error = nil, want status error")
	}
	if !strings.Contains(err.Error(), "502 Bad Gateway") {
		t.Fatalf("Execute() error = %v, want status in error", err)
	}
}
