package tui

import (
	"log"

	_ "modernc.org/sqlite"

	"github.com/aes421/cliStandup/models"
	"github.com/aes421/cliStandup/state"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listModel struct {
	updateList list.Model
	loaded     bool
}

func NewListModel(loaded bool) listModel {
	m := listModel{
		updateList: list.New(
			[]list.Item{},
			list.NewDefaultDelegate(),
			0,
			0),
		loaded: loaded,
	}

	m.updateList.Title = "Sprint Updates"
	m.updateList.AdditionalShortHelpKeys = func() []key.Binding { return getListModelKeys() }
	m.updateList.AdditionalFullHelpKeys = func() []key.Binding { return getListModelKeys() }
	m.updateList.SetItems(UpdatesToListItems(state.Updates))
	m.updateList.SetSize(state.WindowSize.Width, state.WindowSize.Height)
	return m
}

func (m listModel) Init() tea.Cmd {
	return LoadFromDb
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(FatalError); ok {
		log.Fatal(msg)
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		state.WindowSize = msg
		m.updateList.SetSize(state.WindowSize.Width, state.WindowSize.Height)
	case tea.KeyMsg:
		if key.Matches(msg, Keymap.Esc) || msg.String() == "ctrl+c" {
			return m, tea.Quit
		} else if key.Matches(msg, Keymap.Delete) {
			return m, DeleteUpdate(m.updateList.Index(), m.updateList.SelectedItem().(models.Update))
		} else if key.Matches(msg, Keymap.Add) {
			return NewAddModel(), nil
		} else if key.Matches(msg, Keymap.Generate) {
			model := NewOutputModel()
			return model, tea.Batch(model.(outputModel).spinner.Tick, GenerateReport())
		}
	case LoadedUpdates:
		m.updateList.SetItems(UpdatesToListItems(state.Updates))
		m.updateList.Select(int(msg))
		m.loaded = true
		return m, nil
	}

	var cmd tea.Cmd
	m.updateList, cmd = m.updateList.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	if !m.loaded {
		return "Loading..."
	}
	return m.updateList.View()
}

func UpdatesToListItems(updates []models.Update) []list.Item {
	items := make([]list.Item, len(updates))
	for i, u := range updates {
		items[i] = u
	}
	return items
}
