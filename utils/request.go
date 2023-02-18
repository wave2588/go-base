package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cast"
)

type HttpMethod string

const (
	HttpMethodPost HttpMethod = "POST"
	HttpMethodGet             = "GET"
)

var DefaultTimeout = 30 * time.Second

func HttpRequest(method HttpMethod, endpoint string, header, values map[string]interface{}, timeout time.Duration) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if timeout != 0 {
		client.Timeout = timeout
	}

	if len(values) != 0 {
		endpoint += "?"
		for k, v := range values {
			endpoint = fmt.Sprintf("%s%s=%s&", endpoint, k, cast.ToString(v))
		}
	}
	requestBody, _ := json.Marshal(values)
	body := bytes.NewBuffer(requestBody)
	r, err := http.NewRequest(string(method), endpoint, body)
	if err != nil {
		return nil, err
	}

	for key, value := range header {
		r.Header.Add(key, cast.ToString(value))
	}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
