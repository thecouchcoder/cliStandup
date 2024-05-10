package main

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
