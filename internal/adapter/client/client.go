package client

import (
	"bytes"
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

func New() *HTTPClient {
	return &HTTPClient{
		// TODO: Object pool.
		client: &http.Client{},
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

	resBody, err = getBody(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return resHeaders, resBody, nil
}

func getBody(data io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	return body, nil
}
