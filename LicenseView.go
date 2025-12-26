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

func setLicenseViewport() (viewport.Model, error) {
	const width = 78

	vp := viewport.New(width, 8)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	// We need to adjust the width of the glamour render from our main width
	// to account for a few things:
	//
	//  * The Viewport border width
	//  * The Viewport padding
	//  * The Viewport margins
	//  * The gutter glamour applies to the left side of the content
	//
	const glamourGutter = 2
	glamourRenderWidth := width - vp.Style.GetHorizontalFrameSize() - glamourGutter

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
			m.Viewport, cmd = m.Viewport.Update(msg)
			return m, cmd
		}
	default:
		return m, nil
	}
}

func LicenseView(m model) string {
	return m.Viewport.View() + helpView()
}

func helpView() string {
	return subtleStyle.Render("\n  ↑/↓: Navigate • q: Quit\n")
}

//
//func main() {
//	model, err := newLicenseViewport()
//	if err != nil {
//		fmt.Println("Could not initialize Bubble Tea model:", err)
//		os.Exit(1)
//	}
//
//	if _, err := tea.NewProgram(model).Run(); err != nil {
//		fmt.Println("Bummer, there's been an error:", err)
//		os.Exit(1)
//	}
//}
