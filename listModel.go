package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	updates list.Model
}

func NewModel() ListModel {
	sampleUpdate := UpdateItem{description: "This is a sample update"}

	m := ListModel{
		updates: list.New(
			[]list.Item{sampleUpdate, sampleUpdate, sampleUpdate},
			list.NewDefaultDelegate(),
			0,
			0),
	}

	m.updates.Title = "Sprint Updates"
	return m
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case UpdateItem:
		log.Printf("Adding update: %v", msg)
		m.updates.InsertItem(0, msg)
	case tea.WindowSizeMsg:
		m.updates.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "d":
			m.updates.RemoveItem(m.updates.Index())
		case "a":
			models["add"] = NewAddModel()
			return models["add"].Update(nil)
		}
	}

	var cmd tea.Cmd
	m.updates, cmd = m.updates.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	return m.updates.View()
}
