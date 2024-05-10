package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type AddModel struct {
	textArea textarea.Model
}

func NewAddModel(width, height int) AddModel {
	m := AddModel{
		textArea: textarea.New(),
	}
	m.textArea.Placeholder = "Enter your update here"
	m.textArea.Focus()
	m.textArea.SetWidth(width)
	m.textArea.SetHeight(height)

	return m
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.textArea.Focused() {
				m.textArea.Blur()
				return m, nil
			} else {
				return models["list"], nil
			}
		// Can fix this functionality later
		// case "ctrl+s":
		// 	log.Printf("Saving update: %v", m.textArea.Value())
		// 	return models["list"], m.SaveUpdateCmd
		case "enter":
			if !m.textArea.Focused() {
				return models["list"], m.SaveUpdateCmd
			}
		case "w":
			if !m.textArea.Focused() {
				m.textArea.Focus()
				return m, textarea.Blink
			}
		}
	}

	var cmd tea.Cmd
	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

func (m AddModel) View() string { return m.textArea.View() }

func (m AddModel) SaveUpdateCmd() tea.Msg {
	return UpdateItem{description: m.textArea.Value()}
}
