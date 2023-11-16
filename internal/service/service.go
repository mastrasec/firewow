package service

import (
	"github.com/mastrasec/firewow/internal/adapter/client"
)

type Contract interface {
	OpenAIV1ChatCompletionsPipe(
		header map[string][]string,
		body []byte,
	) (resHeaders map[string][]string, resBody []byte, err error)
}

type Service struct {
	httpCli client.Contract
}

func New(httpCli client.Contract) *Service {
	return &Service{
		httpCli: httpCli,
	}
}

func (svc *Service) OpenAIV1ChatCompletionsPipe(
	header map[string][]string,
	body []byte,
) (resHeaders map[string][]string, resBody []byte, err error) {
	// unmarshall json
	// do some stuff
	// marshall json
	// call client
	// process response

	url := "https://api.openai.com/v1/chat/completions"

	resHeaders, resBody, err = svc.httpCli.Post(header, body, url)
	if err != nil {
		return nil, nil, err
	}

	return resHeaders, resBody, nil
}
