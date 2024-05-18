package main

import (
	"fmt"
	"strings"

	"github.com/aes421/cliStandup/state"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var titleStyle = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
}()

type outputModel struct {
	content    string
	generating bool
	viewport   viewport.Model
	help       help.Model
	spinner    spinner.Model
}

func NewOutputModel() tea.Model {

	viewport := viewport.New(state.WindowSize.Width, state.WindowSize.Height-1)
	m := outputModel{
		generating: true,
		viewport:   viewport,
		help:       help.New(),
		spinner:    spinner.New(),
	}
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	m.spinner.Spinner = spinner.Dot
	m.SetViewport()
	return m
}

func (m outputModel) Init() tea.Cmd {
	return nil
}

func (m outputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return NewListModel(), nil
		}
	case tea.WindowSizeMsg:
		state.WindowSize = msg
		m.SetViewport()

	case GeneratedReport:
		m.content = string(msg)
		m.SetViewport()
		m.generating = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.generating {
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m outputModel) View() string {
	if m.generating {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.headerView(),
			m.spinner.View(),
			m.help.View(outputModelkeyMap),
		)
	}
	s := fmt.Sprintf("%s\n%s\n%s\n", m.headerView(), m.viewport.View(), m.help.View(outputModelkeyMap))
	return s
}

func (m *outputModel) headerView() string {
	title := titleStyle.Render("Standup")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *outputModel) SetViewport() {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.help.View(outputModelkeyMap))
	m.viewport.Width = state.WindowSize.Width
	m.viewport.Height = state.WindowSize.Height - headerHeight - footerHeight
	m.viewport.YPosition = headerHeight
	m.viewport.SetContent(m.content)
	m.help.Width = state.WindowSize.Width
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
