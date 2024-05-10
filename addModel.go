package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	EscWriteMode key.Binding
	EscViewMode  key.Binding
	Write        key.Binding
	Save         key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		EscWriteMode: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "view")),
		EscViewMode:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Write: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "write mode"),
		),
		Save: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "save"),
		),
	}
}

type AddModel struct {
	textArea textarea.Model
	help     help.Model
	KeyMap   KeyMap
}

func (m AddModel) ShortHelp() []key.Binding {
	if m.textArea.Focused() {
		return []key.Binding{m.KeyMap.EscViewMode}
	}
	return []key.Binding{m.KeyMap.EscViewMode, m.KeyMap.Write, m.KeyMap.Save}
}

// Noop to satisfy the interface
func (k AddModel) FullHelp() [][]key.Binding { return nil }

func NewAddModel(width, height int) AddModel {
	m := AddModel{
		textArea: textarea.New(),
		help:     help.New(),
		KeyMap:   DefaultKeyMap(),
	}
	m.textArea.Placeholder = "Enter your update here"
	m.textArea.Focus()
	m.textArea.SetWidth(width)
	m.textArea.SetHeight(height - 1)
	m.help.Width = width
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

func (m AddModel) View() string {
	helpView := m.help.View(m)
	textAreaView := m.textArea.View()

	// TODO we could dynamically determine height incase of window resize
	return lipgloss.JoinVertical(lipgloss.Left, textAreaView, helpView)
}

func (m AddModel) SaveUpdateCmd() tea.Msg {
	return UpdateItem{description: m.textArea.Value()}
}
