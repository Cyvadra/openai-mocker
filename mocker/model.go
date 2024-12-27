package mocker

type ChatRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewSystemPrompt(prompt string) ChatRequestMessage {
	return ChatRequestMessage{
		Role:    "system",
		Content: prompt,
	}
}
func NewUserMessage(msg string) ChatRequestMessage {
	return ChatRequestMessage{
		Role:    "user",
		Content: msg,
	}
}
func NewAssistantMessage(msg string) ChatRequestMessage {
	return ChatRequestMessage{
		Role:    "assistant",
		Content: msg,
	}
}

type ChatRequest struct {
	Stream   bool                 `json:"stream"`
	Model    string               `json:"model"`
	Messages []ChatRequestMessage `json:"messages"`
}

type Completion struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionsStreamResponseChoiceDelta struct {
	Content string `json:"content"`
}

type ChatCompletionsStreamResponseChoice struct {
	Delta        ChatCompletionsStreamResponseChoiceDelta `json:"delta"`
	FinishReason *string                                  `json:"finish_reason"`
}

type ChatCompletionsStreamResponse struct {
	Id      string                                `json:"id"`
	Object  string                                `json:"object"`
	Created int64                                 `json:"created"`
	Model   string                                `json:"model"`
	Choices []ChatCompletionsStreamResponseChoice `json:"choices"`
}
