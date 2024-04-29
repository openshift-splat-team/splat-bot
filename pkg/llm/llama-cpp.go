package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	llamacpp "github.com/openshift-splat-team/splat-bot/data/llm/llama_cpp"
)


const DEFAULT_SYSTEM_PROMPT = "You are a helpful assistant. Your top priority is accuracy using the information you're given."
const DEFAULT_PORT = 8080
//const DEFAULT_HOST = "192.168.0.145"
const DEFAULT_HOST = "localhost"

// 
var urlTemplate = "http://%s:%d/v1/chat/completions"

func Completion(content string) (*llamacpp.ChatCompletion, error) {
	request := llamacpp.ChatCompletionRequest{
		Model: "tinyllama",
		Messages: []llamacpp.Message{
			{
				Role:    "user",
				Content: content,
			},
			{
				Role:    "system",
				Content: "This is a conversation between User and Llama, a friendly chatbot. Llama is helpful, kind, honest, good at writing, and never fails to answer any requests immediately and with precision. ",
			},
		},
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf(urlTemplate, DEFAULT_HOST, DEFAULT_PORT), bytes.NewBuffer(requestBytes))

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer no-key")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	completion := &llamacpp.ChatCompletion{}
	err = json.NewDecoder(res.Body).Decode(completion)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	defer res.Body.Close()
	return completion, nil
}