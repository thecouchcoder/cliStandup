package main

import (
	"log"

	_ "modernc.org/sqlite"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	updateList    list.Model
	width, height int
	loaded        bool
}

func NewListModel() ListModel {

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
			m.width, m.height = msg.Width, msg.Height
			m.updateList.SetSize(m.width, m.height)

		// TODO this is duplicated
		case LoadedUpdates:
			log.Printf("received %d items\n", len(msg))
			items := make([]list.Item, len(msg))
			for i, u := range msg {
				items[i] = u
			}
			m.updateList.SetItems(items)
			m.loaded = true
			return m, nil
		}

		return m, nil
	}

	m.updateList.SetSize(m.width, m.height)
	switch msg := msg.(type) {
	case LoadedUpdates:
		log.Printf("received %d items\n", len(msg))
		items := make([]list.Item, len(msg))
		for i, u := range msg {
			items[i] = u
		}
		m.updateList.SetItems(items)
		m.loaded = true
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case listModelkeyMap.Delete.Keys()[0]:
			return m, m.DeleteUpdateCmd()
		case listModelkeyMap.Add.Keys()[0]:
			deprecatedmodels["list"] = m
			return NewAddModel(m.width, m.height), nil
		case listModelkeyMap.Generate.Keys()[0]:
			deprecatedmodels["list"] = m
			model := NewOutputModel(m.width, m.height)
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
