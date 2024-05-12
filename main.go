package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var models = make(map[string]tea.Model)

func Init() (*sql.DB, map[string]tea.Model, error) {
	log.Print("initializing database...")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	path := "cliStandup.db"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Print("Creating database...")
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

	log.Print("initializing models...")
	initModels := make(map[string]tea.Model)
	initModels["list"] = NewModel(db)
	initModels["add"] = NewAddModel(0, 0)

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
