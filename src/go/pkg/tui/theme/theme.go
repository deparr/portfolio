package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	renderer *lipgloss.Renderer

	background lipgloss.TerminalColor
	body       lipgloss.TerminalColor
	border     lipgloss.TerminalColor
	highlight  lipgloss.TerminalColor
	accent     lipgloss.TerminalColor
	error      lipgloss.TerminalColor

	base lipgloss.Style
}

func BaseTheme(renderer *lipgloss.Renderer) Theme {
	base := Theme{
		renderer: renderer,
	}

	base.background = lipgloss.AdaptiveColor{Dark: "#151515", Light: "#dcdcd6"}
	base.body = lipgloss.AdaptiveColor{Dark: "#c5c8c6", Light: "#4d4d4c"}
	base.border = lipgloss.AdaptiveColor{Dark: "#696969", Light: "#696969"}
	base.accent = lipgloss.AdaptiveColor{Dark: "#e78c45", Light: "#d05200"}
	base.base = lipgloss.NewStyle().Foreground(base.body)

	return base
}

func (t Theme) Background() lipgloss.TerminalColor {
	return t.background
}

func (t Theme) Body() lipgloss.TerminalColor {
	return t.body
}

func (t Theme) Base() lipgloss.Style {
	return t.base
}

func (t Theme) Border() lipgloss.TerminalColor {
	return t.border
}

func (t Theme) Highlight() lipgloss.TerminalColor {
	return t.highlight
}

func (t Theme) Accent() lipgloss.TerminalColor {
	return t.accent
}

func (t Theme) Error() lipgloss.TerminalColor {
	return t.error
}
