package main

import "github.com/aes421/cliStandup/db/dbmodel"

type UpdateItems []UpdateItem

func NewUpdateItems(updates []dbmodel.Update) UpdateItems {
	items := make([]UpdateItem, len(updates))
	for i, u := range updates {
		items[i] = UpdateItem{description: u.Description}
	}
	return UpdateItems(items)
}

type UpdateItem struct {
	description string
}

func (u UpdateItem) Title() string {
	return ""
}

func (u UpdateItem) Description() string {
	return u.description
}

func (u UpdateItem) FilterValue() string {
	return u.description
}
