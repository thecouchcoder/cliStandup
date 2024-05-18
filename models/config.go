package models

type Config struct {
	ChatGPT struct {
		Endpoint  string  `json:"endpoint"`
		APIKey    string  `json:"api_key"`
		Model     string  `json:"llmModel"`
		Temp      float32 `json:"temperature"`
		MaxTokens int     `json:"max_tokens"`
	}
	ExternalCallsEnabled bool `json:"externalCallsEnabled"`
}
