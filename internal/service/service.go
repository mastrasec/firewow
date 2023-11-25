package service

import (
	"encoding/json"
	"fmt"

	"github.com/mastrasec/firewow/internal/adapter/client"
	"github.com/mastrasec/firewow/internal/service/canary"
)

type Contract interface {
	OpenAIV1ChatCompletionsPipe(
		header map[string][]string,
		body []byte,
	) (resHeaders map[string][]string, resBody []byte, err error)
}

type Service struct {
	httpCliAdapter client.Contract
	canaryService  canary.Contract
}

func New(httpCli client.Contract, canaryService canary.Contract) *Service {
	return &Service{
		httpCliAdapter: httpCli,
		canaryService:  canaryService,
	}
}

func (svc *Service) OpenAIV1ChatCompletionsPipe(
	header map[string][]string,
	body []byte,
) (resHeaders map[string][]string, resBody []byte, err error) {
	url := "https://api.openai.com/v1/chat/completions"

	payload := &oaiV1ChatCompletionsReq{}

	if err := json.Unmarshal(body, payload); err != nil {
		return nil, nil, err
	}

	lastReqIndex := len(payload.Messages) - 1
	msg, canaryToken := svc.canaryService.AddToken(payload.Messages[lastReqIndex].Content)
	payload.Messages[lastReqIndex].Content = msg

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	resHeaders, resBody, err = svc.httpCliAdapter.Post(header, reqBody, url)
	if err != nil {
		return nil, nil, err
	}

	resPayload := &oaiV1ChatCompletionsRes{}
	if err := json.Unmarshal(resBody, resPayload); err != nil {
		fmt.Println("kkkkkkkkkkkkkkkkkkkkkkkkk")

		return nil, nil, err
	}

	resMsg := resPayload.Choices[0].Message.Content

	if leaked := svc.canaryService.HasLeakage(resMsg, canaryToken); leaked {
		cleanedResMsg := svc.canaryService.HandleLeakage(resMsg, canaryToken)
		resPayload.Choices[0].Message.Content = cleanedResMsg
	}

	finalResBody, err := json.Marshal(resPayload)
	if err != nil {
		return nil, nil, err
	}

	return resHeaders, finalResBody, nil
}
