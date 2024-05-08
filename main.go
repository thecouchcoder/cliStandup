package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices   []string
	cursor    int
	selected  choice
	updates   []string
	userInput string
}

const ADD_UPDATE = "add update"
const LIST_UPDATE = "list update"
const REMOVE_UPDATE = "remove update"
const CLEAR = "clear updates"
const GENERATE = "generate standup"

type choice int

const (
	mainchoice choice = iota
	addchoice
	listchoice
	removechoice
	clearchoice
	generatechoice
)

func initialModel() model {
	return model{
		choices:  []string{ADD_UPDATE, LIST_UPDATE, REMOVE_UPDATE, CLEAR, GENERATE},
		selected: mainchoice,
		updates:  []string{"bob", "said", "it"}, //make([]string, 0),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch m.selected {
	case mainchoice:
		switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyMsg:

			// Cool, what was the actual key pressed?
			switch msg.String() {

			// These keys should exit the program.
			case "ctrl+c", "esc":
				return m, tea.Quit

			// The "up" and "k" keys move the cursor up
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			// The "down" and "j" keys move the cursor down
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				m.selected = choice(m.cursor + 1)
				m.cursor = 0
			}
		}
	case addchoice:
		switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyMsg:
			switch val := msg.String(); val {
			case "esc":
				m.cursor = 0
				m.selected = mainchoice
			case "enter":
				m.updates = append(m.updates, m.userInput)
				m.userInput = ""
				m.cursor = 0
				m.selected = mainchoice
			case "backspace":
				if len(m.userInput) > 0 {
					m.userInput = m.userInput[:len(m.userInput)-1]
				}
			default:
				m.userInput += val
			}

		}
	case removechoice:
		switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.cursor = 0
				m.selected = mainchoice
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.updates)-1 {
					m.cursor++
				}
			case "enter", " ":
				m.updates = append(m.updates[:m.cursor], m.updates[m.cursor+1:]...)
			}

		}
	case clearchoice:
		switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.cursor = 0
				m.selected = mainchoice

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				m.updates = make([]string, 0)
				m.cursor = 0
				m.selected = mainchoice
			}

		}

	default:
		switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.cursor = 0
				m.selected = mainchoice
			}

		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := ""
	if m.selected == mainchoice {
		s += m.mainMenuView()
		s += "\nPress esc to quit.\n"
	} else {
		s += m.choices[m.selected-1] + "\n"

		switch m.selected {
		case addchoice:
			s += m.addUpdateView()
		case removechoice:
			s += m.removeUpdateView()
		case listchoice:
			s += m.listUpdateView()
		case clearchoice:
			s += m.clearUpdatesView()
		case generatechoice:
			s += m.notImplementedView()
		}

		s += "\n\nPress esc for main menu.\n"
	}

	// Send the UI for rendering
	return s
}

func (m model) notImplementedView() string {
	return "\nNot Implemented\n"
}
func (m model) addUpdateView() string {
	s := "Type your update. \n"
	s += m.userInput
	return s
}

func (m model) removeUpdateView() string {
	s := ""

	// Iterate over our choices
	for i, u := range m.updates {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%d: %s %s\n", i, cursor, u)
	}
	return s
}

func (m model) clearUpdatesView() string {
	s := "Are you sure? (enter)"
	return s
}

func (m model) listUpdateView() string {
	if len(m.updates) == 0 {
		return "no updates saved yet"
	}
	s := ""
	for i, val := range m.updates {
		s += fmt.Sprintf("%d: %s\n", i, val)
	}

	return s
}

func (m model) mainMenuView() string {
	s := "Choose an option\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
