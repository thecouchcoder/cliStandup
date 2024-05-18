package main

import (
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/aes421/cliStandup/llm"
	"github.com/aes421/cliStandup/state"
	"github.com/aes421/cliStandup/tui"
	tea "github.com/charmbracelet/bubbletea"
)

//go:embed db/schema.sql
var ddl string

func main() {
	os.Remove("debug.log")
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	chatgpt := llm.NewChatGPT(
		state.Config.ChatGPT.Endpoint,
		state.Config.ChatGPT.APIKey,
		state.Config.ChatGPT.Model,
		state.Config.ChatGPT.Temp,
		state.Config.ChatGPT.MaxTokens)

	err = state.InitState(ddl, chatgpt)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer state.Db.Close()
	p := tea.NewProgram(tui.NewListModel(false), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
