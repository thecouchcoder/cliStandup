package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"time"

	_ "modernc.org/sqlite"

	"github.com/aes421/cliStandup/db/dbmodel"
	"github.com/aes421/cliStandup/llm"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	updateList    list.Model
	width, height int
	loaded        bool
	db            *sql.DB
	llm           llm.LLM
}

func NewListModel(db *sql.DB, llm llm.LLM) ListModel {

	m := ListModel{
		updateList: list.New(
			[]list.Item{},
			list.NewDefaultDelegate(),
			0,
			0),
		loaded: false,
		db:     db,
		llm:    llm,
	}

	m.updateList.Title = "Sprint Updates"
	m.updateList.AdditionalShortHelpKeys = getListModelKeys()
	m.updateList.AdditionalFullHelpKeys = getListModelKeys()
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
			models["list"] = m
			return NewAddModel(m.width, m.height), nil
		case listModelkeyMap.Generate.Keys()[0]:
			models["list"] = m
			model := NewOutputModel(m.width, m.height)
			return model, tea.Batch(model.(outputModel).spinner.Tick, m.GenerateReportCmd())
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

//go:embed db/schema.sql
var ddl string

func (m ListModel) LoadListCmd() tea.Msg {
	log.Print("loading list...")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// query for updates
	dbValues, err := dbmodel.New(m.db).GetActiveUpdates(ctx)
	if err != nil {
		return FatalError(err.Error())
	}
	updates = dbToUpdate(dbValues)

	return LoadedUpdates(updates)
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
		updates = append([]Update{NewUpdate(r.ID, description)}, updates...)

		// TODO find somewhere to do this
		m.updateList.Select(0)
		return LoadedUpdates(updates)
	}
}

func (m ListModel) DeleteUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		log.Printf("Deleting update: %v", m.updateList.Index())
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		mdl := m.updateList.SelectedItem().(Update)
		err := dbmodel.New(m.db).ArchiveUpdate(ctx, mdl.id)
		if err != nil {
			return FatalError(err.Error())
		}

		index := m.updateList.Index()

		log.Print(updates[:index])
		log.Print(updates[index+1:])
		log.Print(updates)
		updates = append(updates[:index], updates[index+1:]...)
		return LoadedUpdates(updates)
	}
}

func (m ListModel) GenerateReportCmd() tea.Cmd {
	return func() tea.Msg {
		log.Print("Generating report...")
		if GetConfig().ExternalCallsEnabled == false {
			time.Sleep(5 * time.Second)
			return GeneratedReport("external calls are disabled")
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		content, err := m.llm.Generate(ctx)
		if err != nil {
			return FatalError(err.Error())
		}

		return GeneratedReport(content)
	}
}

type GeneratedReport string
