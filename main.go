package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateItem struct {
	title       string
	description string
}

func (u UpdateItem) Title() string {
	return u.title
}

func (u UpdateItem) Description() string {
	return u.description
}

func (u UpdateItem) FilterValue() string {
	return u.description
}

type Model struct {
	updates list.Model
}

func NewModel() Model {
	sampleUpdate := UpdateItem{title: "Update", description: "This is a sample update"}

	m := Model{
		updates: list.New(
			[]list.Item{sampleUpdate, sampleUpdate, sampleUpdate},
			list.NewDefaultDelegate(),
			0,
			0),
	}

	m.updates.Title = "Sprint Updates"
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updates.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.updates, cmd = m.updates.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.updates.View()
}

func main() {
	os.Remove("debug.log")
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
