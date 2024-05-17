package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const padding = 2

var titleStyle = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
}()

type outputModel struct {
	width        int
	height       int
	content      string
	generating   bool
	viewport     viewport.Model
	help         help.Model
	progress     progress.Model
	incrementing bool
}

func NewOutputModel(width int, height int) tea.Model {

	viewport := viewport.New(width, height-1)
	m := outputModel{
		width:        width,
		height:       height,
		generating:   true,
		viewport:     viewport,
		help:         help.New(),
		progress:     progress.New(progress.WithDefaultGradient()),
		incrementing: true,
	}
	m.SetViewport()
	m.progress.Width = m.width - padding*2 - 4
	m.progress.ShowPercentage = false
	return m
}

type tickMsg time.Time

func (m outputModel) Init() tea.Cmd {
	return nil
}

func (m outputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return models["list"], nil
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.progress.Width = msg.Width - padding*2 - 4
		m.SetViewport()

	case GeneratedReport:
		m.content = string(msg)
		m.SetViewport()
		m.generating = false
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			m.incrementing = false
		}
		if m.progress.Percent() == 0.0 {
			m.incrementing = true
		}

		var cmd tea.Cmd
		switch {
		case m.incrementing:
			cmd = m.progress.IncrPercent(1)
		case !m.incrementing:
			cmd = m.progress.DecrPercent(1)
		}
		return m, tea.Batch(cmd, tickCmd())

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m outputModel) View() string {
	if m.generating {
		pad := strings.Repeat(" ", padding)
		return "\n" +
			pad + m.progress.View() + "\n\n" +
			pad + m.help.View(outputModelkeyMap)
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
	m.viewport.Width = m.width
	m.viewport.Height = m.height - headerHeight - footerHeight
	m.viewport.YPosition = headerHeight
	m.viewport.SetContent(m.content)
	m.help.Width = m.width
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
