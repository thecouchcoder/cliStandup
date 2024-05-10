package main

import "github.com/charmbracelet/bubbles/list"

type UpdateItem struct {
	title       string
	description string
}

func (u UpdateItem) Title() string {
	return u.title
}

func (u UpdateItem) Description() string {
	return u.description
}

func (u UpdateItem) FilterValue() string {
	return u.description
}

type Model struct {
	updates list.Model
}

func NewModel() Model {
	sampleUpdate := UpdateItem{title: "Update", description: "This is a sample update"}

	m := Model{
		updates: list.New(
			[]list.Item{sampleUpdate, sampleUpdate, sampleUpdate},
			list.NewDefaultDelegate(),
			0,
			0),
	}

	m.updates.Title = "Sprint Updates"
	return m
}
