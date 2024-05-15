package llm

import (
	"context"
	"log"
	"os"
)

type LLM interface {
	Generate(context.Context) (string, error)
}

type chatgpt struct{}

func NewChatGPT() LLM {
	return chatgpt{}
}

func (chatgpt) Generate(ctx context.Context) (string, error) {
	content, err := os.ReadFile("assets/prompt.txt")
	if err != nil {
		return "", err
	}

	log.Print(string(content))

	return "", nil
}
