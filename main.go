package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aes421/cliStandup/llm"
	tea "github.com/charmbracelet/bubbletea"
)

var models = make(map[string]tea.Model)
var config Config

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

func GetConfig() Config {
	return config
}

func Init() (*sql.DB, map[string]tea.Model, error) {
	log.Print("initializing database...")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	path := "cliStandup.db"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Print("creating database...")
		os.Create(path)
	}
	// open and create tables
	db, err := sql.Open("sqlite", "cliStandup.db")
	if err != nil {
		return nil, nil, err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, nil, err
	}

	log.Print("reading config...")
	configFile, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)

	if err := jsonParser.Decode(&config); err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	log.Print("initializing llm...")
	chatgpt := llm.NewChatGPT(
		db,
		config.ChatGPT.Endpoint,
		config.ChatGPT.APIKey,
		config.ChatGPT.Model,
		config.ChatGPT.Temp,
		config.ChatGPT.MaxTokens)

	log.Print("initializing models...")
	initModels := make(map[string]tea.Model)
	initModels["list"] = NewListModel(db, chatgpt)

	return db, initModels, nil
}
func main() {
	os.Remove("debug.log")
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	db, models, err := Init()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer db.Close()
	p := tea.NewProgram(models["list"], tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
