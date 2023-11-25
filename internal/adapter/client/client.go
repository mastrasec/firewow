package client

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

type Contract interface {
	Post(
		headers map[string][]string,
		body []byte,
		url string,
	) (resHeaders map[string][]string, resBody []byte, err error)
}

type HTTPClient struct {
	client *http.Client
}

func New(client *http.Client) *HTTPClient {
	return &HTTPClient{
		client: client,
	}
}

func (hc *HTTPClient) Post(
	headers map[string][]string,
	body []byte,
	url string,
) (resHeaders map[string][]string, resBody []byte, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	res, err := hc.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	resHeaders = map[string][]string(res.Header)

	resBody, err = getBody(res)
	if err != nil {
		return nil, nil, err
	}

	return resHeaders, resBody, nil
}

func getBody(res *http.Response) ([]byte, error) {
	if res.Header.Get("Content-Encoding") == "gzip" {
		return gzipBody(res)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return body, nil
}

func gzipBody(res *http.Response) ([]byte, error) {
	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return body, nil
}
