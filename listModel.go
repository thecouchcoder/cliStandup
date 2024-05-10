package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type keymap struct {
	add      key.Binding
	delete   key.Binding
	generate key.Binding
}

var keyMap = keymap{
	add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	generate: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "generate"),
	),
}

type ListModel struct {
	updates       list.Model
	width, height int
	loaded        bool
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
	m.updates.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyMap.add,
			keyMap.delete,
			keyMap.generate,
		}
	}
	m.updates.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyMap.add,
			keyMap.delete,
			keyMap.generate,
		}
	}
	return m
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.loaded {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width, m.height = msg.Width, msg.Height
			m.updates.SetSize(m.width, m.height)
			m.loaded = true
		}

		return m, nil
	}

	m.updates.SetSize(m.width, m.height)
	switch msg := msg.(type) {
	case UpdateItem:
		m.updates.InsertItem(0, msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case keyMap.delete.Help().Key:
			m.updates.RemoveItem(m.updates.Index())
		case keyMap.add.Help().Key:
			models["list"] = m
			models["add"] = NewAddModel(m.width, m.height)
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
