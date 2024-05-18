package main

import (
	"context"
	"log"
	"time"

	"github.com/aes421/cliStandup/db/dbmodel"
	tea "github.com/charmbracelet/bubbletea"

	_ "embed"

	"github.com/aes421/cliStandup/models"
	"github.com/aes421/cliStandup/state"
)

type GeneratedReport string
type FatalError string
type LoadedUpdates []models.Update

//go:embed db/schema.sql
var ddl string

func LoadListCmd() tea.Msg {
	log.Print("loading list...")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// query for updates
	dbValues, err := dbmodel.New(state.Db).GetActiveUpdates(ctx)
	if err != nil {
		return FatalError(err.Error())
	}
	state.Updates = models.DbToUpdate(dbValues)

	return LoadedUpdates(state.Updates)
}

func (m ListModel) SaveUpdateCmd(description string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Saving update: %v", description)
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		r, err := dbmodel.New(state.Db).CreateUpdate(ctx, description)
		if err != nil {
			return FatalError(err.Error())
		}

		log.Print("update saved.")
		state.Updates = append([]models.Update{models.NewUpdate(r.ID, description)}, state.Updates...)

		// TODO find somewhere to do this
		m.updateList.Select(0)
		return LoadedUpdates(state.Updates)
	}
}

func (m ListModel) DeleteUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		log.Printf("Deleting update: %v", m.updateList.Index())
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		updateModel := m.updateList.SelectedItem().(models.Update)
		err := dbmodel.New(state.Db).ArchiveUpdate(ctx, updateModel.Id)
		if err != nil {
			return FatalError(err.Error())
		}

		index := m.updateList.Index()

		log.Print(state.Updates[:index])
		log.Print(state.Updates[index+1:])
		log.Print(state.Updates)
		state.Updates = append(state.Updates[:index], state.Updates[index+1:]...)
		return LoadedUpdates(state.Updates)
	}
}

func (m ListModel) GenerateReportCmd() tea.Cmd {
	return func() tea.Msg {
		log.Print("Generating report...")
		if state.Config.ExternalCallsEnabled == false {
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
