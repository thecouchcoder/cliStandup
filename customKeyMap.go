package main

import "github.com/charmbracelet/bubbles/key"

type listModelKeys struct {
	Add      key.Binding
	Delete   key.Binding
	Generate key.Binding
}

type addModelKeys struct {
	EscWriteMode key.Binding
	EscViewMode  key.Binding
	Write        key.Binding
	Save         key.Binding
}

type outputModelKeys struct {
	Esc key.Binding
}

var listModelkeyMap = listModelKeys{
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

func getListModelKeys() func() []key.Binding {
	return func() []key.Binding {
		return []key.Binding{
			listModelkeyMap.Add,
			listModelkeyMap.Delete,
			listModelkeyMap.Generate,
		}
	}
}

var addModelkeyMap = addModelKeys{
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

var outputModelkeyMap = outputModelKeys{
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}

func (m outputModelKeys) ShortHelp() []key.Binding {
	return []key.Binding{outputModelkeyMap.Esc}
}

func (m outputModelKeys) FullHelp() [][]key.Binding {
	return nil
}
