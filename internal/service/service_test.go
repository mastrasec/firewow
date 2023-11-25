package service

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/mastrasec/firewow/internal/service/canary"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/req_openai_v1_chat_completions_pipe.json
var reqOpenAIV1ChatCompletionsPipe string

//go:embed testdata/res_openai_v1_chat_completions_pipe.json
var resOpenAIV1ChatCompletionsPipe string

func TestOpenAIV1ChatCompletionsPipe(t *testing.T) {
	client := &mockClient{}
	can := canary.New()

	svc := New(client, can)

	header := map[string][]string{
		"Content-Type":    {"gzip"},
		"Accept-Language": {"ru", "ja"},
	}
	body := []byte(reqOpenAIV1ChatCompletionsPipe)

	resHeader, resBody, err := svc.OpenAIV1ChatCompletionsPipe(header, body)
	assert.NoError(t, err)

	// Workaround to remove unwanted characters and facilitate the comparison.
	tmpPayload := &oaiV1ChatCompletionsRes{}
	err = json.Unmarshal([]byte(resOpenAIV1ChatCompletionsPipe), tmpPayload)
	assert.NoError(t, err)
	wantBody, err := json.Marshal(tmpPayload)
	assert.NoError(t, err)

	// It ensures all sent header values return,
	// so they got there in the first place.
	for key, values := range header {
		assert.Equal(t, values, resHeader[key])
	}

	assert.Equal(t, string(wantBody), string(resBody))
}

type mockClient struct {
}

func (mc *mockClient) Post(
	header map[string][]string,
	body []byte,
	url string,
) (resHeader map[string][]string, resBody []byte, err error) {
	return header, []byte(resOpenAIV1ChatCompletionsPipe), nil
}
