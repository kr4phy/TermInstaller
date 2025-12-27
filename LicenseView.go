package main

import (
	_ "embed"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

//go:embed resources/LICENSE.md
var content string

func setLicenseViewport(m model) (viewport.Model, error) {
	wh, wv := windowStyle.GetFrameSize()
	mh, mv := mainStyle.GetFrameSize()
	helpHeight := lipgloss.Height(helpView())
	vp := viewport.New(m.vpWidth-wh-mh, m.vpHeight-wv-mv-helpHeight)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)
	vp.Width -= vp.Style.GetHorizontalFrameSize()
	vp.Height -= vp.Style.GetVerticalFrameSize()

	// We need to adjust the width of the glamour render from our main width
	// to account for a few things:
	//
	//  * The LicenseViewport border width
	//  * The LicenseViewport padding
	//  * The LicenseViewport margins
	//  * The gutter glamour applies to the left side of the content
	//
	const glamourGutter = 4
	glamourRenderWidth := m.vpWidth - wh - mh - vp.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		return viewport.Model{
			Width:             0,
			Height:            0,
			KeyMap:            viewport.KeyMap{},
			MouseWheelEnabled: false,
			MouseWheelDelta:   0,
			YOffset:           0,
			YPosition:         0,
			Style:             lipgloss.Style{},
		}, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return viewport.Model{
			Width:             0,
			Height:            0,
			KeyMap:            viewport.KeyMap{},
			MouseWheelEnabled: false,
			MouseWheelDelta:   0,
			YOffset:           0,
			YPosition:         0,
			Style:             lipgloss.Style{},
		}, err
	}

	vp.SetContent(str)

	return vp, nil
}

//
//func (m model) LicenseViewInit() tea.Cmd {
//	return nil
//}

func updateLicenseView(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.CurrentView = stateLicenseAccept
			return m, nil
		default:
			var cmd tea.Cmd
			m.LicenseViewport, cmd = m.LicenseViewport.Update(msg)
			return m, cmd
		}
	default:
		return m, nil
	}
}

func LicenseView(m model) string {
	return m.LicenseViewport.View() + helpView()
}

func helpView() string {
	return subtleStyle.Render("\n  ↑/↓: Navigate • q: Quit\n")
}
