package main

import "github.com/charmbracelet/bubbles/key"

type listModelKeys struct {
	Add      key.Binding
	Delete   key.Binding
	Generate key.Binding
}

var listModelKeyMap = listModelKeys{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Generate: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "generate"),
	),
}

type addModelKeys struct {
	EscWriteMode key.Binding
	EscViewMode  key.Binding
	Write        key.Binding
	Save         key.Binding
}

var addModelKeyMap = addModelKeys{
	EscWriteMode: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "view")),
	EscViewMode: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back")),
	Write: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "write mode"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save"),
	),
}
