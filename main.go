package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var models = make(map[string]tea.Model)

func InitModels() {
	models["list"] = NewModel()
	models["add"] = NewAddModel()
}
func main() {
	os.Remove("debug.log")
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	InitModels()
	p := tea.NewProgram(models["list"], tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
