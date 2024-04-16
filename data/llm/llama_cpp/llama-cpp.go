package llamacpp

type ChatCompletion struct {
    Choices  []Choice `json:"choices"`
    Created  int64    `json:"created"`
    ID       string   `json:"id"`
    Model    string   `json:"model"`
    Object   string   `json:"object"`
    Usage    Usage    `json:"usage"`
}

type Choice struct {
    FinishReason string  `json:"finish_reason"`
    Index        int     `json:"index"`
    Message      Message `json:"message"`
}

type Message struct {
    Content string `json:"content"`
    Role    string `json:"role"`
}

type Usage struct {
    CompletionTokens int `json:"completion_tokens"`
    PromptTokens     int `json:"prompt_tokens"`
    TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionRequest struct {
    Model string `json:"model"`
    Messages []Message `json:"messages"`
}