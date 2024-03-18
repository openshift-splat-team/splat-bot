package util

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const (
	PROMPT_RESPONSE_TIMEOUT = time.Second * 120
)

// GenerateResponse generates a response from an ollama API endpoint
func GenerateResponse(ctx context.Context, prompt string) (string, error) {
	endpoint := os.Getenv("OLLAMA_ENDPOINT")
	if len(endpoint) == 0 {
		return "", errors.New("OLLAMA_ENDPOINT must be exported")
	}

	llm, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		log.Fatal(err)
	}

	timedCtx, cancel := context.WithTimeout(ctx, PROMPT_RESPONSE_TIMEOUT)
	defer cancel()
	completion, err := llms.GenerateFromSinglePrompt(timedCtx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	return completion, nil
}
