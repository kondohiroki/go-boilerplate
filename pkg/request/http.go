package request

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type HttpRequest struct {
	HttpClient *http.Client
	Url        string
	Method     string
	Headers    map[string]string
	Body       []byte
	Query      map[string]string
	Params     map[string]string
}

func NewHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
}

func MakeHTTPRequest(req HttpRequest) (*http.Response, error) {
	httpReq, err := http.NewRequest(req.Method, req.Url, http.NoBody)
	if err != nil {
		return nil, err
	}

	for key, value := range req.Headers {
		httpReq.Header.Add(key, value)
	}

	if req.Body != nil {
		httpReq.Body = io.NopCloser(bytes.NewReader(req.Body))
	}

	if req.Query != nil {
		query := httpReq.URL.Query()
		for key, value := range req.Query {
			query.Add(key, value)
		}
		httpReq.URL.RawQuery = query.Encode()
	}

	if req.Params != nil {
		path := httpReq.URL.Path
		for key, value := range req.Params {
			path = strings.Replace(path, ":"+key, value, -1)
		}
		httpReq.URL.Path = path
	}

	if req.HttpClient.Timeout == 0 {
		req.HttpClient.Timeout = 30 * time.Second
	}

	resp, err := req.HttpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
