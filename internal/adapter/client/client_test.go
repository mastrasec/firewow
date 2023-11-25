package client

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestGetbody(t *testing.T) {
	testCases := []struct {
		name            string
		contentEncoding string
		body            []byte
	}{
		{
			name:            "raw data",
			contentEncoding: "",
			body:            []byte("raw data"),
		},
		{
			name:            "gzipped data",
			contentEncoding: "gzip",
			body:            []byte("gzipped data"),
		},
	}

	for _, tc := range testCases {
		buf := &bytes.Buffer{}

		if tc.contentEncoding == "gzip" {
			writer := gzip.NewWriter(buf)
			_, err := writer.Write(tc.body)
			assert.NoError(t, err)
			writer.Close()
		} else {
			buf = bytes.NewBuffer(tc.body)
		}

		res := &http.Response{
			Body:   io.NopCloser(buf),
			Header: http.Header{"Content-Encoding": []string{tc.contentEncoding}},
		}

		got, err := getBody(res)
		assert.NoError(t, err)
		assert.Equal(t, tc.body, got)
	}
}

func TestPost(t *testing.T) {
	testCases := []struct {
		name       string
		headers    map[string][]string
		resHeaders map[string][]string
		body       []byte
		resBody    []byte
	}{
		{
			name: "happy path",
			headers: map[string][]string{
				"Accept-Language": {"ru", "ja"},
			},
			resHeaders: map[string][]string{
				"Accept-Language": {"en-us", "pt-br"},
			},
			body:    []byte(gofakeit.LoremIpsumWord()),
			resBody: []byte(gofakeit.LoremIpsumWord()),
		},
		{
			name:    "no request headers",
			headers: map[string][]string{},
			resHeaders: map[string][]string{
				"Accept-Language": {"en-us", "pt-br"},
			},
			body:    []byte(gofakeit.LoremIpsumWord()),
			resBody: []byte(gofakeit.LoremIpsumWord()),
		},
		{
			name: "no request body",
			headers: map[string][]string{
				"Accept-Language": {"ru", "ja"},
			},
			resHeaders: map[string][]string{
				"Accept-Language": {"en-us", "pt-br"},
			},
			body:    []byte{},
			resBody: []byte(gofakeit.LoremIpsumWord()),
		},
		{
			name: "no response headers",
			headers: map[string][]string{
				"Accept-Language": {"ru", "ja"},
			},
			resHeaders: map[string][]string{},
			body:       []byte(gofakeit.LoremIpsumWord()),
			resBody:    []byte(gofakeit.LoremIpsumWord()),
		},
		{
			name: "no response body",
			headers: map[string][]string{
				"Accept-Language": {"ru", "ja"},
			},
			resHeaders: map[string][]string{
				"Accept-Language": {"en-us", "pt-br"},
			},
			body:    []byte(gofakeit.LoremIpsumWord()),
			resBody: []byte{},
		},
	}

	client := New(&http.Client{})

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					for key, values := range tc.resHeaders {
						for _, value := range values {
							w.Header().Add(key, value)
						}
					}

					_, err := w.Write(tc.resBody)
					assert.NoError(t, err)
				}))
				defer server.Close()

				resHeaders, resBody, err := client.Post(tc.headers, tc.body, server.URL)

				assert.NoError(t, err)

				// Looping over resHeaders would not work because some headers
				// are automatically added to the response.
				for key, values := range tc.resHeaders {
					assert.Equal(t, values, resHeaders[key])
				}

				assert.Equal(t, tc.resBody, resBody)
			},
		)

	}

}
