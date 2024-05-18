package models

import "github.com/aes421/cliStandup/db/dbmodel"

func DbToUpdate(updates []dbmodel.Update) []Update {
	items := make([]Update, len(updates))
	for i, u := range updates {
		items[i] = NewUpdate(u.ID, u.Description)
	}
	return items
}

type Update struct {
	Id          int64
	description string
}

func NewUpdate(id int64, description string) Update {
	return Update{Id: id, description: description}
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
