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
	BaseUrl    string
	Headers    map[string]string
	HttpClient *http.Client
}

var DefaultClient = &http.Client{
	Timeout: time.Second * 15,
}

func InitClient(baseUrl string) Request {
	return Request{
		Headers:    make(map[string]string),
		HttpClient: DefaultClient,
		BaseUrl:    baseUrl,
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

	res, err := r.HttpClient.Do(req)
	if err != nil {
		return "", errors.E(err, errors.Params{"url": url, "method": method})
	}
	defer res.Body.Close()
	read, err := ioutil.ReadAll(res.Body)
	return string(read), nil
}

func (r *Request) getBase(path string) string {
	if path == "" {
		return fmt.Sprintf("%s", r.BaseUrl)
	}
	return fmt.Sprintf("%s/%s", r.BaseUrl, path)
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
