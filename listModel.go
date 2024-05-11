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
			listModelKeyMap.Add,
			listModelKeyMap.Delete,
			listModelKeyMap.Generate,
		}
	}
	m.updates.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listModelKeyMap.Add,
			listModelKeyMap.Delete,
			listModelKeyMap.Generate,
		}
	}
	return m
}

func (m ListModel) Init() tea.Cmd {
	return m.LoadListCmd
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
			m.updates.SetSize(m.width, m.height)
		case InitiallyLoadedUpdates:
			log.Printf("received %d items\n", len(msg))
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
	case UpdatedModel:
		m = ListModel(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case listModelKeyMap.Delete.Keys()[0]:
			return m, m.DeleteUpdateCmd()
		case listModelKeyMap.Add.Keys()[0]:
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
	updates, err := dbmodel.New(m.db).GetActiveUpdates(ctx)
	if err != nil {
		return FatalError(err.Error())
	}

	return NewUpdateItems(updates)
}

type FatalError string

func (m ListModel) SaveUpdateCmd(description string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Saving update: %v", description)
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		r, err := dbmodel.New(m.db).CreateUpdate(ctx, description)
		if err != nil {
			return FatalError(err.Error())
		}

		log.Print("update saved.")
		m.updates.InsertItem(0, NewUpdate(r.ID, description))
		m.updates.Select(0)
		return UpdatedModel(m)
	}
}

func (m ListModel) DeleteUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		log.Printf("Deleting update: %v", m.updates.Index())
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		mdl := m.updates.SelectedItem().(Update)
		err := dbmodel.New(m.db).ArchiveUpdate(ctx, mdl.id)
		if err != nil {
			return FatalError(err.Error())
		}

		m.updates.RemoveItem(m.updates.Index())
		return UpdatedModel(m)
	}
}

type UpdatedModel ListModel
