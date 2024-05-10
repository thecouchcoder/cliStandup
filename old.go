package main

// import (
// 	"fmt"
// 	"os"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// type model struct {
// 	choices   []string
// 	cursor    int
// 	state     state
// 	updates   []string
// 	userInput string
// }

// const ADD_UPDATE = "add update"
// const LIST_UPDATE = "list update"
// const REMOVE_UPDATE = "remove update"
// const CLEAR = "clear updates"
// const GENERATE = "generate standup"

// type state int

// const (
// 	main_state state = iota
// 	add_state
// 	list_state
// 	remove_state
// 	clear_state
// 	generate_state
// )

// func initialModel() model {
// 	return model{
// 		choices: []string{ADD_UPDATE, LIST_UPDATE, REMOVE_UPDATE, CLEAR, GENERATE},
// 		state:   main_state,
// 		updates: []string{"bob", "said", "it"}, //make([]string, 0),
// 	}
// }

// func (m model) Init() tea.Cmd {
// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

// 	switch m.state {
// 	case main_state:
// 		switch msg := msg.(type) {

// 		// Is it a key press?
// 		case tea.KeyMsg:

// 			// Cool, what was the actual key pressed?
// 			switch msg.String() {

// 			// These keys should exit the program.
// 			case "ctrl+c", "esc":
// 				return m, tea.Quit

// 			// The "up" and "k" keys move the cursor up
// 			case "up", "k":
// 				if m.cursor > 0 {
// 					m.cursor--
// 				}

// 			// The "down" and "j" keys move the cursor down
// 			case "down", "j":
// 				if m.cursor < len(m.choices)-1 {
// 					m.cursor++
// 				}

// 			// The "enter" key and the spacebar (a literal space) toggle
// 			// the state _state for the item that the cursor is pointing at.
// 			case "enter", " ":
// 				m.state = state(m.cursor + 1)
// 				m.cursor = 0
// 			}
// 		}
// 	case add_state:
// 		switch msg := msg.(type) {

// 		// Is it a key press?
// 		case tea.KeyMsg:
// 			switch val := msg.String(); val {
// 			case "esc":
// 				m.cursor = 0
// 				m.state = main_state
// 			case "enter":
// 				m.updates = append(m.updates, m.userInput)
// 				m.userInput = ""
// 				m.cursor = 0
// 				m.state = main_state
// 			case "backspace":
// 				if len(m.userInput) > 0 {
// 					m.userInput = m.userInput[:len(m.userInput)-1]
// 				}
// 			default:
// 				m.userInput += val
// 			}

// 		}
// 	case remove_state:
// 		switch msg := msg.(type) {

// 		// Is it a key press?
// 		case tea.KeyMsg:
// 			switch msg.String() {
// 			case "esc":
// 				m.cursor = 0
// 				m.state = main_state
// 			case "up", "k":
// 				if m.cursor > 0 {
// 					m.cursor--
// 				}
// 			case "down", "j":
// 				if m.cursor < len(m.updates)-1 {
// 					m.cursor++
// 				}
// 			case "enter", " ":
// 				m.updates = append(m.updates[:m.cursor], m.updates[m.cursor+1:]...)
// 			}

// 		}
// 	case clear_state:
// 		switch msg := msg.(type) {

// 		// Is it a key press?
// 		case tea.KeyMsg:
// 			switch msg.String() {
// 			case "esc":
// 				m.cursor = 0
// 				m.state = main_state

// 			// The "enter" key and the spacebar (a literal space) toggle
// 			// the state _state for the item that the cursor is pointing at.
// 			case "enter", " ":
// 				m.updates = make([]string, 0)
// 				m.cursor = 0
// 				m.state = main_state
// 			}

// 		}

// 	default:
// 		switch msg := msg.(type) {

// 		// Is it a key press?
// 		case tea.KeyMsg:
// 			switch msg.String() {
// 			case "esc":
// 				m.cursor = 0
// 				m.state = main_state
// 			}

// 		}

// 	}

// 	// Return the updated model to the Bubble Tea runtime for processing.
// 	// Note that we're not returning a command.
// 	return m, nil
// }

// func (m model) View() string {
// 	s := ""
// 	if m.state == main_state {
// 		s += m.mainMenuView()
// 		s += "\nPress esc to quit.\n"
// 	} else {
// 		s += m.choices[m.state-1] + "\n"

// 		switch m.state {
// 		case add_state:
// 			s += m.addUpdateView()
// 		case remove_state:
// 			s += m.removeUpdateView()
// 		case list_state:
// 			s += m.listUpdateView()
// 		case clear_state:
// 			s += m.clearUpdatesView()
// 		case generate_state:
// 			s += m.notImplementedView()
// 		}

// 		s += "\n\nPress esc for main menu.\n"
// 	}

// 	// Send the UI for rendering
// 	return s
// }

// func (m model) notImplementedView() string {
// 	return "\nNot Implemented\n"
// }
// func (m model) addUpdateView() string {
// 	s := "Type your update. \n"
// 	s += m.userInput
// 	return s
// }

// func (m model) removeUpdateView() string {
// 	s := ""

// 	// Iterate over our choices
// 	for i, u := range m.updates {

// 		// Is the cursor pointing at this choice?
// 		cursor := " " // no cursor
// 		if m.cursor == i {
// 			cursor = ">" // cursor!
// 		}

// 		// Render the row
// 		s += fmt.Sprintf("%d: %s %s\n", i, cursor, u)
// 	}
// 	return s
// }

// func (m model) clearUpdatesView() string {
// 	s := "Are you sure? (enter)"
// 	return s
// }

// func (m model) listUpdateView() string {
// 	if len(m.updates) == 0 {
// 		return "no updates saved yet"
// 	}
// 	s := ""
// 	for i, val := range m.updates {
// 		s += fmt.Sprintf("%d: %s\n", i, val)
// 	}

// 	return s
// }

// func (m model) mainMenuView() string {
// 	s := "Choose an option\n"

// 	// Iterate over our choices
// 	for i, choice := range m.choices {

// 		// Is the cursor pointing at this choice?
// 		cursor := " " // no cursor
// 		if m.cursor == i {
// 			cursor = ">" // cursor!
// 		}

// 		// Render the row
// 		s += fmt.Sprintf("%s %s\n", cursor, choice)
// 	}
// 	return s
// }
