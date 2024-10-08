package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deparr/portfolio/go/pkg/tui/theme"
)

type page int
type size int

const (
	splashPage page = iota
	homePage
	contactPage
	projectsPage
	experiencePage
	debugPage
)

const (
	undersized size = iota
	small
	medium
	large
)

type model struct {
	page            page
	theme           theme.Theme
	renderer        *lipgloss.Renderer
	state           state
	viewport        viewport.Model
	termWidth       int
	termHeight      int
	size            size
	contentWidth    int
	contentHeight   int
	containerWidth  int
	containerHeight int
	switched        bool
	ready           bool
	hasScroll       bool
}

type state struct {
	splash  splashState
	footer  footerState
	spinner spinnerState
}

func NewModel(renderer *lipgloss.Renderer) tea.Model {
	return model{
		page:     splashPage,
		theme:    theme.BaseTheme(renderer),
		renderer: renderer,
		state: state{
			splash: splashState{delay: false},
			footer: footerState{
				binds: []footerBinding{
					{key: "j/k", action: "scroll"},
					{key: "q", action: "quit"},
				},
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return m.splashInit()
}

func (m model) switchPage(newPage page) model {
	m.page = newPage
	m.switched = true
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

		switch {
		case m.termWidth < 80 || m.termHeight < 30:
			m.size = undersized
			m.containerWidth = m.termWidth
			m.containerHeight = m.termHeight
		default:
			m.size = large
			m.containerWidth = 80
			m.containerHeight = min(msg.Height, 30)
		}

		m.contentWidth = m.containerWidth - 4
		m = m.updateViewport()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.page {
	case splashPage:
		m, cmd = m.splashUpdate(msg)
	case contactPage:
		m, cmd = m.contactUpdate(msg)
	case projectsPage:
		m, cmd = m.projectsUpdate(msg)
	}

	m, headerCmd := m.headerUpdate(msg)
	cmds := []tea.Cmd{headerCmd}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.viewport.SetContent(m.getContent())
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport.SetContent(m.getContent())
	m.viewport, cmd = m.viewport.Update(msg)
	if m.switched {
		m = m.updateViewport()
		m.switched = false
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) updateViewport() model {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight + 2

	width := m.containerWidth - 4
	m.contentHeight = m.containerHeight - verticalMarginHeight

	if !m.ready {
		m.viewport = viewport.New(width, m.contentHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = false
		m.viewport.KeyMap = viewport.DefaultKeyMap()
		m.ready = true
	} else {
		m.viewport.Width = width
		m.viewport.Height = m.contentHeight
		m.viewport.GotoTop()
	}

	m.hasScroll = m.viewport.VisibleLineCount() < m.viewport.TotalLineCount()

	return m
}

func (m model) View() string {
	if m.size == undersized {
		return m.resizeView()
	}

	switch m.page {
	case splashPage:
		return m.splashView()
	default:
		header := m.headerView()
		footer := m.footerView()

		var view string
		if m.hasScroll {
			view = lipgloss.JoinVertical(
				lipgloss.Right,
				m.viewport.View(),
				m.locView(),
			)
		} else {
			view = m.getContent()
		}

		height := m.containerHeight
		height -= lipgloss.Height(header)
		height -= lipgloss.Height(footer)

		boxedView := lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			m.theme.Base().
				Width(m.containerWidth).
				Height(height).
				Padding(0, 1).
				Render(view),
			footer,
		)

		return m.renderer.Place(
			m.termWidth,
			m.termHeight,
			lipgloss.Center,
			lipgloss.Center,
			m.theme.Base().
				MaxWidth(m.termWidth).
				MaxHeight(m.termHeight).
				Render(boxedView),
		)
	}
}

func (m model) locView() string {
	lines := m.viewport.TotalLineCount()
	if m.viewport.VisibleLineCount() >= lines {
		return "ALL"
	}
	y := m.viewport.YOffset
	percent := int(m.viewport.ScrollPercent() * 100)
	var view string
	switch percent {
	case 0:
		view = "TOP"
	case 100:
		view = "BOT"
	default:
		view = fmt.Sprintf("%d%% %d/%d", percent, y, m.viewport.TotalLineCount())
	}
	return m.theme.TextAccent().Bold(true).Render(view)
}

func (m model) getContent() string {
	content := "none :("
	switch m.page {
	case homePage:
		content = m.homeView()
	case contactPage:
		content = m.contactView()
	case projectsPage:
		content = m.projectsView()
	case experiencePage:
		content = m.experienceView()
	case debugPage:
		content = m.debugView()
	}

	return content
}
