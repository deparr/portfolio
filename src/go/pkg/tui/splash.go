package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type splashState struct {
	delay bool
}

type DelayCompleteMsg struct{}

func (m model) splashInit() tea.Cmd {
	cmd := tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return DelayCompleteMsg{}
	})
	spin := m.spinStartCmd()
	return tea.Batch(cmd, spin)
}

func (m model) splashUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg.(type) {
	case DelayCompleteMsg:
		m.state.splash.delay = true
		m.spinning = false
	}

	if m.state.splash.delay {
		return m.projectsSwitch()
	}

	if m.spinning {
		return m.spinnerAdvance()
	}

	return m, nil
}

func (m model) splashView() string {
	symbol := m.theme.TextHighlight().Italic(true).Bold(true)
	name := m.theme.TextAccent().Italic(true).Bold(true)
	view := lipgloss.JoinHorizontal(
		lipgloss.Center,
		symbol.Render("@"),
		name.Render("dparrott"),
	)

	return m.renderer.Place(
		m.viewportWidth,
		m.viewportHeight,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}