package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
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

func InitClient(baseURL string) Request {
	return Request{
		Headers:    make(map[string]string),
		HTTPClient: DefaultClient,
		BaseURL:    baseURL,
	}
}

func (r *Request) Execute(method string, path string, body interface{}) (string, error) {
	url := r.getBase(path)
	payload := getBody(body)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", errors.E(err, errors.Params{"url": url, "method": method})
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	res, err := r.HTTPClient.Do(req)
	if err != nil {
		return "", errors.E(err, errors.Params{"url": url, "method": method})
	}
	defer res.Body.Close()
	read, err := ioutil.ReadAll(res.Body)
	return string(read), nil
}

func (r *Request) getBase(path string) string {
	if path == "" {
		return fmt.Sprintf("%s", r.BaseURL)
	}
	return fmt.Sprintf("%s/%s", r.BaseURL, path)
}

func getBody(body interface{}) (buf io.ReadWriter) {
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return
		}
	}
	return
}
