package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type AddModel struct {
	textArea textarea.Model
}

// TODO better colors and fix sizing
func NewAddModel() AddModel {
	m := AddModel{
		textArea: textarea.New(),
	}
	m.textArea.Placeholder = "Enter your update here"
	m.textArea.Focus()

	return m
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("AddModel: %v", msg)
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

		case "ctrl+s":
			return models["list"], nil
		case "enter":
			if !m.textArea.Focused() {
				return models["list"], nil
			}
		case "w":
			if !m.textArea.Focused() {
				m.textArea.Focus()
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

func (m AddModel) View() string { return m.textArea.View() }
