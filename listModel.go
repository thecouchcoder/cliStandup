package main

import (
	"log"

	_ "modernc.org/sqlite"

	"github.com/aes421/cliStandup/models"
	"github.com/aes421/cliStandup/state"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	updateList list.Model
	loaded     bool
}

func InitListModel() ListModel {
	m := ListModel{
		updateList: list.New(
			[]list.Item{},
			list.NewDefaultDelegate(),
			0,
			0),
		loaded: false,
	}

	m.updateList.Title = "Sprint Updates"
	m.updateList.AdditionalShortHelpKeys = getListModelKeys()
	m.updateList.AdditionalFullHelpKeys = getListModelKeys()
	m.updateList.SetItems(UpdatesToListItems(state.Updates))
	return m
}

func NewListModel() ListModel {
	m := ListModel{
		updateList: list.New(
			[]list.Item{},
			list.NewDefaultDelegate(),
			0,
			0),
		loaded: true,
	}

	m.updateList.Title = "Sprint Updates"
	m.updateList.AdditionalShortHelpKeys = getListModelKeys()
	m.updateList.AdditionalFullHelpKeys = getListModelKeys()
	m.updateList.SetItems(UpdatesToListItems(state.Updates))
	return m
}

func (m ListModel) Init() tea.Cmd {
	return LoadListCmd
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(FatalError); ok {
		log.Fatal(msg)
		return m, tea.Quit
	}

	if !m.loaded {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			state.WindowSize = msg
			m.updateList.SetSize(state.WindowSize.Width, state.WindowSize.Height)

		// TODO this is duplicated
		case LoadedUpdates:
			log.Printf("received %d items\n", len(msg))
			m.updateList.SetItems(UpdatesToListItems(state.Updates))
			m.loaded = true
			return m, nil
		}

		return m, nil
	}

	m.updateList.SetSize(state.WindowSize.Width, state.WindowSize.Height)
	switch msg := msg.(type) {
	case LoadedUpdates:
		log.Printf("received %d items\n", len(msg))
		m.updateList.SetItems(UpdatesToListItems(state.Updates))
		m.loaded = true
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case listModelkeyMap.Delete.Keys()[0]:
			return m, m.DeleteUpdateCmd()
		case listModelkeyMap.Add.Keys()[0]:
			return NewAddModel(), nil
		case listModelkeyMap.Generate.Keys()[0]:
			model := NewOutputModel()
			return model, tea.Batch(model.(outputModel).spinner.Tick, GenerateReportCmd())
		}
	}

	var cmd tea.Cmd
	m.updateList, cmd = m.updateList.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
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
