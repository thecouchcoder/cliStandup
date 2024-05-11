package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/aes421/cliStandup/db/dbmodel"

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
	db            *sql.DB
}

func NewModel(db *sql.DB) ListModel {

	m := ListModel{
		updates: list.New(
			[]list.Item{},
			list.NewDefaultDelegate(),
			0,
			0),
		loaded: false,
		db:     db,
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
	return m.LoadListCmd
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.loaded {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width, m.height = msg.Width, msg.Height
			m.updates.SetSize(m.width, m.height)
		case fatalError:
			log.Fatal(msg)
			return m, tea.Quit
		case InitiallyLoadedUpdates:
			log.Printf("leceived %d items\n", len(msg))
			items := make([]list.Item, len(msg))
			for i, u := range msg {
				items[i] = u
			}
			m.updates.SetItems(items)
			m.loaded = true
		}

		return m, nil
	}

	m.updates.SetSize(m.width, m.height)
	switch msg := msg.(type) {
	case NewUpdate:
		return m, m.SaveUpdateCmd(msg)
	case ListModel:
		m = msg
		return m, nil
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
	if !m.loaded {
		return "Loading..."
	}
	return m.updates.View()
}

//go:embed db/schema.sql
var ddl string

func (m ListModel) LoadListCmd() tea.Msg {
	log.Print("loading list...")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// query for updates
	updates, err := dbmodel.New(m.db).GetUpdates(ctx)
	if err != nil {
		return fatalError(err.Error())
	}

	return NewUpdateItems(updates)
}

type fatalError string

func (m ListModel) SaveUpdateCmd(msg NewUpdate) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Saving update: %v", msg.Description())
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err := dbmodel.New(m.db).CreateUpdate(ctx, msg.Description())
		if err != nil {
			return fatalError(err.Error())
		}

		log.Print("update saved.")
		m.updates.InsertItem(0, msg)
		return m
	}
}
