package tui

import "github.com/charmbracelet/bubbles/key"

type keys struct {
	Add          key.Binding
	Delete       key.Binding
	Generate     key.Binding
	EscWriteMode key.Binding
	EscViewMode  key.Binding
	Write        key.Binding
	Save         key.Binding
	Esc          key.Binding
	Quit         key.Binding
}

var Keymap = keys{
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
	), EscWriteMode: key.NewBinding(
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
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "back"),
	),
}

func getListModelKeys() []key.Binding {
	return []key.Binding{
		Keymap.Add,
		Keymap.Delete,
		Keymap.Generate,
	}
}

func getAddModelViewModeKeys() []key.Binding {
	return []key.Binding{Keymap.EscViewMode, Keymap.Write, Keymap.Save}
}

func getAddModelWriteModeKeys() []key.Binding {
	return []key.Binding{Keymap.EscWriteMode}
}

func getOutputModelsKeys() []key.Binding {
	return []key.Binding{Keymap.Esc, Keymap.Quit}
}
