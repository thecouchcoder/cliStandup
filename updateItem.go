package main

import "github.com/aes421/cliStandup/db/dbmodel"

type InitiallyLoadedUpdates []NewUpdate

func NewUpdateItems(updates []dbmodel.Update) InitiallyLoadedUpdates {
	items := make([]NewUpdate, len(updates))
	for i, u := range updates {
		items[i] = NewUpdate{description: u.Description}
	}
	return InitiallyLoadedUpdates(items)
}

type NewUpdate struct {
	description string
}

func (u NewUpdate) Title() string {
	return ""
}

func (u NewUpdate) Description() string {
	return u.description
}

func (u NewUpdate) FilterValue() string {
	return u.description
}
