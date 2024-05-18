package main

import (
	"context"
	"log"
	"time"

	"github.com/aes421/cliStandup/db/dbmodel"
	tea "github.com/charmbracelet/bubbletea"

	_ "embed"
)

type GeneratedReport string
type FatalError string
type LoadedUpdates []Update

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
