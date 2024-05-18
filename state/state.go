package state

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/aes421/cliStandup/models"
	tea "github.com/charmbracelet/bubbletea"
)

var Updates []models.Update
var Db *sql.DB
var LLMConnector models.LLM
var Config models.Config
var WindowSize tea.WindowSizeMsg

func InitState(ddl string, llmConnector models.LLM) error {
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
	Db, err = sql.Open("sqlite", "cliStandup.db")
	if err != nil {
		return err
	}

	if _, err := Db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	log.Print("reading config...")
	configFile, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)

	if err := jsonParser.Decode(&Config); err != nil {
		log.Fatal(err)
		return err
	}

	log.Print("initializing llm...")
	LLMConnector = llmConnector

	log.Print("initializing models...")

	return nil
}
