package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type SpinnerStartMsg struct{}
type SpinnerStopMsg struct{}
type SpinnerFrameMsg struct{}

var frames = []string{
	" ⣷",
	" ⣯",
	" ⣟",
	" ⡿",
	" ⢿",
	" ⣻",
	" ⣽",
	" ⣾",
}

var frameAdvanceCmd = func() tea.Msg {
	return tea.Tick(time.Millisecond*125, func(t time.Time) tea.Msg {
		return SpinnerFrameMsg{}
	})
}

func (m model) spinnerView() string {
	return frames[m.spinnerframe%len(frames)]
}

func (m model) spinnerAdvance() (model, tea.Cmd) {
	m.spinnerframe += 1
	return m, nil
}

func (m model) spinStartCmd() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		return SpinnerStartMsg{}
	},
		frameAdvanceCmd,
	)
}

func (m model) spinStopCmd() tea.Cmd {
	return func() tea.Msg {
		return SpinnerStopMsg{}
	}
}
