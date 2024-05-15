package llm

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aes421/cliStandup/db/dbmodel"
)

type LLM interface {
	Generate(context.Context) (string, error)
}

type chatGptCompletionRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type chatGptCompletionResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type chatgpt struct {
	db          *sql.DB
	endpoint    string
	api_key     string
	llmModel    string
	temperature float32
	maxTokens   int
}

func NewChatGPT(db *sql.DB, endpoint string, api_key string, llmModel string, temperature float32, maxTokens int) LLM {
	return chatgpt{
		db:          db,
		endpoint:    endpoint,
		api_key:     api_key,
		llmModel:    llmModel,
		temperature: temperature,
		maxTokens:   maxTokens,
	}
}

func (c chatgpt) Generate(ctx context.Context) (string, error) {
	prompt, err := os.ReadFile("assets/prompt.txt")
	if err != nil {
		return "", err
	}

	// read the db
	updates, err := dbmodel.New(c.db).GetActiveUpdates(ctx)
	if err != nil {
		return "", err
	}

	promptString := string(prompt)
	for u := range updates {
		promptString += " " + updates[u].Description
	}

	// construct the request
	reqObj := chatGptCompletionRequest{
		Model:       c.llmModel,
		Prompt:      promptString,
		Temperature: c.temperature,
		MaxTokens:   c.maxTokens,
	}

	log.Printf("prompt: %s", promptString)

	jsonBody, err := json.Marshal(reqObj)
	if err != nil {
		return "", err
	}

	log.Printf("jsonBody: %s", jsonBody)
	reader := bytes.NewReader(jsonBody)

	// call openai
	client := &http.Client{}
	req, err := http.NewRequest("POST", c.endpoint+"/openai/v1/completions", reader)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.api_key)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// extract the response
	var response chatGptCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", nil
	}

	// write to txt file (better to do this somewhere else later)
	err = os.WriteFile("assets/generated.txt", []byte(response.Choices[0].Text), 0644)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Text, nil
}
