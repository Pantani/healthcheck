package client

import (
	"io"
	"reflect"
	"testing"
)

func TestRequest_getBase(t *testing.T) {
	tests := []struct {
		name    string
		BaseURL string
		path    string
		want    string
	}{
		{"get api url with path", "http://blockatlas.com", "api", "http://blockatlas.com/api"},
		{"get api url", "http://blockatlas.com", "", "http://blockatlas.com"},
		{"get api without url", "", "api", "/api"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				BaseURL: tt.BaseURL,
			}
			if got := r.getBase(tt.path); got != tt.want {
				t.Errorf("getBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBody(t *testing.T) {
	tests := []struct {
		name    string
		body    interface{}
		wantBuf io.ReadWriter
	}{
		{"test nil", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBuf := getBody(tt.body); !reflect.DeepEqual(gotBuf, tt.wantBuf) {
				t.Errorf("getBody() = %v, want %v", gotBuf, tt.wantBuf)
			}
		})
	}
}
