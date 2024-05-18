package tui

import (
	"context"
	"log"
	"time"

	"github.com/aes421/cliStandup/db/dbmodel"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aes421/cliStandup/models"
	"github.com/aes421/cliStandup/state"
)

type GeneratedReport string
type FatalError string
type LoadedUpdates int

func LoadFromDb() tea.Msg {
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

	return LoadedUpdates(0)
}

func SaveUpdate(description string) tea.Cmd {
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

		return LoadedUpdates(0)
	}
}

func DeleteUpdate(index int, item models.Update) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Deleting update: %v", index)
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err := dbmodel.New(state.Db).ArchiveUpdate(ctx, item.Id)
		if err != nil {
			return FatalError(err.Error())
		}

		state.Updates = append(state.Updates[:index], state.Updates[index+1:]...)
		return LoadedUpdates(index)
	}
}

func GenerateReport() tea.Cmd {
	return func() tea.Msg {
		log.Print("Generating report...")
		if state.Config.ExternalCallsEnabled == false {
			time.Sleep(5 * time.Second)
			return GeneratedReport("external calls are disabled")
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		content, err := state.LLMConnector.Generate(ctx)
		if err != nil {
			return FatalError(err.Error())
		}

		return GeneratedReport(content)
	}
}
