package main

import "github.com/aes421/cliStandup/db/dbmodel"

type InitiallyLoadedUpdates []Update

func NewUpdateItems(updates []dbmodel.Update) InitiallyLoadedUpdates {
	items := make([]Update, len(updates))
	for i, u := range updates {
		items[i] = NewUpdate(u.ID, u.Description)
	}
	return InitiallyLoadedUpdates(items)
}

type Update struct {
	id          int64
	description string
}

func NewUpdate(id int64, description string) Update {
	return Update{id: id, description: description}
}

func (u Update) Title() string {
	return ""
}

func (u Update) Description() string {
	return u.description
}

func (u Update) FilterValue() string {
	return u.description
}
