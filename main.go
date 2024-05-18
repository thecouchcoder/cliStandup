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
	"github.com/aes421/cliStandup/state"
	tea "github.com/charmbracelet/bubbletea"
)

var deprecatedmodels = make(map[string]tea.Model)

func Init() (map[string]tea.Model, error) {
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
	var err error
	state.Db, err = sql.Open("sqlite", "cliStandup.db")
	if err != nil {
		return nil, err
	}

	if _, err := state.Db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	log.Print("reading config...")
	configFile, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)

	if err := jsonParser.Decode(&state.Config); err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Print("initializing llm...")
	state.LLMConnector = llm.NewChatGPT(
		state.Config.ChatGPT.Endpoint,
		state.Config.ChatGPT.APIKey,
		state.Config.ChatGPT.Model,
		state.Config.ChatGPT.Temp,
		state.Config.ChatGPT.MaxTokens)

	log.Print("initializing models...")
	initModels := make(map[string]tea.Model)
	initModels["list"] = NewListModel()

	return initModels, nil
}
func main() {
	os.Remove("debug.log")
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	models, err := Init()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer state.Db.Close()
	p := tea.NewProgram(models["list"], tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
