// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package requests

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/XiaoMi/Gaea/log"
)

// Supported http methods
const (
	Get    string = "GET"
	Post   string = "POST"
	Put    string = "PUT"
	Patch  string = "PATCH"
	Delete string = "DELETE"
)

// default global client，safe for concurrent use by multiple goroutines
var defaultClient *http.Client

func init() {
	var dials uint64
	tr := &http.Transport{}
	tr.Dial = func(network, addr string) (net.Conn, error) {
		c, err := net.DialTimeout(network, addr, time.Second*10)
		if err == nil {
			log.Debug("rpc: dial new connection to %s， dials: %d", addr, atomic.AddUint64(&dials, 1)-1)
		}
		return c, err
	}
	defaultClient = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}
	go func() {
		for {
			time.Sleep(time.Minute)
			tr.CloseIdleConnections()
		}
	}()
}

// Request request info
type Request struct {
	User     string
	Password string
	Method   string
	URL      string
	Header   map[string]string
	Params   map[string]string
	Body     []byte
}

// NewRequest return Request
func NewRequest(url, method string, header map[string]string, params map[string]string, data []byte) *Request {
	return &Request{Method: method, URL: url, Header: header, Params: params, Body: data}
}

// SetBasicAuth set basic auth of request
func (req *Request) SetBasicAuth(user, password string) {
	req.User = user
	req.Password = password
}

// Response response info
type Response struct {
	StatusCode int
	Header     map[string][]string
	Body       []byte
}

// AddParameters adds query parameters to the URL.
func AddParameters(baseURL string, queryParams map[string]string) string {
	baseURL += "?"
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return baseURL + params.Encode()
}

// BuildHTTPRequest build a http request object
func BuildHTTPRequest(request *Request) (*http.Request, error) {
	// handle parameters
	if len(request.Params) > 0 {
		request.URL = AddParameters(request.URL, request.Params)
	}
	// build http request
	httpReq, err := http.NewRequest(request.Method, request.URL, bytes.NewReader(request.Body))
	if err != nil {
		return nil, err
	}
	// set basic auth
	if request.User != "" && request.Password != "" {
		httpReq.SetBasicAuth(request.User, request.Password)
	}

	// build http header
	for k, v := range request.Header {
		httpReq.Header.Set(k, v)
	}

	// default json
	_, ok := request.Header["Content-Type"]
	if len(request.Body) > 0 && !ok {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	return httpReq, nil
}

func buildResponse(res *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	r := &Response{
		StatusCode: res.StatusCode,
		Body:       body,
		Header:     res.Header,
	}
	return r, nil
}

// Send send http request
func Send(request *Request) (*Response, error) {
	var start = time.Now()

	// build http request
	httpReq, err := BuildHTTPRequest(request)
	if err != nil {
		return nil, err
	}

	// send http request
	rsp, err := defaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, rsp.Body)
		// close http response
		rsp.Body.Close()
		log.Debug("call rpc [%s] %s in %v", httpReq.Method, httpReq.URL, time.Since(start))
	}()

	// build response
	r, err := buildResponse(rsp)
	return r, err
}

// SendPut send put http request
func SendPut(url, user, password string) error {
	req := NewRequest(url, Put, nil, nil, nil)
	req.SetBasicAuth(user, password)

	resp, err := Send(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(resp.Body))
	}
	return nil
}

// SendGet send get http request
func SendGet(url, user, password string) (*Response, error) {
	req := NewRequest(url, Get, nil, nil, nil)
	req.SetBasicAuth(user, password)

	resp, err := Send(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return resp, nil
}

// EncodeURL encode url
func EncodeURL(host string, format string, args ...interface{}) string {
	var u url.URL
	u.Scheme = "http"
	u.Host = host
	u.Path = fmt.Sprintf(format, args...)
	return u.String()
}
