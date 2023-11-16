package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestGetbody(t *testing.T) {
	var (
		want        = []byte(gofakeit.LoremIpsumWord())
		inputReader = strings.NewReader(string(want))
		data        = io.NopCloser(inputReader)
	)

	got, err := getBody(data)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
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

	client := New()

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
