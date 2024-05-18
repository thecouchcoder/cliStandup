package main

import (
	"github.com/aes421/cliStandup/state"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AddModel struct {
	textArea textarea.Model
	help     help.Model
	keyMap   addModelKeys
}

func (m AddModel) ShortHelp() []key.Binding {
	if m.textArea.Focused() {
		return []key.Binding{m.keyMap.EscWriteMode}
	}
	return []key.Binding{m.keyMap.EscViewMode, m.keyMap.Write, m.keyMap.Save}
}

// Noop to satisfy the interface
func (k AddModel) FullHelp() [][]key.Binding { return nil }

func NewAddModel() AddModel {
	m := AddModel{
		textArea: textarea.New(),
		help:     help.New(),
		keyMap:   addModelkeyMap,
	}
	m.textArea.Placeholder = "Enter your update here"
	m.textArea.Focus()
	m.textArea.SetWidth(state.WindowSize.Width)
	m.textArea.SetHeight(state.WindowSize.Height - 1)
	m.help.Width = state.WindowSize.Width
	m.textArea.CharLimit = 0 // unlimited
	m.help.ShowAll = false

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
				return NewListModel(), nil
			}
		// Can fix this functionality later
		// case "ctrl+s":
		// 	log.Printf("Saving update: %v", m.textArea.Value())
		// 	return models["list"], m.SaveUpdateCmd
		case "enter":
			if !m.textArea.Focused() {
				return NewListModel(), SaveUpdateCmd(m.textArea.Value())
			}
		case "w":
			if !m.textArea.Focused() {
				m.textArea.Focus()
				return m, textarea.Blink
			}
		}
	case tea.WindowSizeMsg:
		m.textArea.SetWidth(state.WindowSize.Width)
		m.textArea.SetHeight(state.WindowSize.Height - 1)
		m.help.Width = state.WindowSize.Width
	}

	var cmd tea.Cmd
	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

func (m AddModel) View() string {
	helpView := m.help.View(m)
	textAreaView := m.textArea.View()

	return lipgloss.JoinVertical(lipgloss.Left, textAreaView, helpView)
}
