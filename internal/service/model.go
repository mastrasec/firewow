package service

type oaiV1ChatCompletionsReq struct {
	Messages    []message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature"`
	// Functions   []function `json:"functions"`
	// functionCall
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// type function struct {
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	// Parameters  []any  `json:"parameters"`
// }

type oaiV1ChatCompletionsRes struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []choice `json:"choices"`
	Usage             usage    `json:"usage"`
}

type choice struct {
	Index        int     `json:"index"`
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
