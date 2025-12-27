package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// General stuff for styling the view
var (
	windowStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)

	// Gradient colors we'll use for the progress bar
	ramp                   = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
	stateStartInstallation = 0
	stateLicenseAccept     = 1
	stateInstallation      = 2
	stateAfterInstallation = 3
	stateLicenseView       = 4
)

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

func main() {
	initialModel := model{0, 0, 0, 0, viewport.New(0, 0), false, stateStartInstallation, 10, 0, 0, false, false, nil, false}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

type model struct {
	vpWidth                 int
	vpHeight                int
	StartInstallationChoice int
	LicenseAcceptChoice     int
	LicenseViewport         viewport.Model
	ExitSetup               bool
	CurrentView             int
	Ticks                   int
	Frames                  int
	Progress                float64
	InstallStarted          bool
	Installed               bool
	InstallationError       error
	Quiting                 bool
}

func (m model) Init() tea.Cmd { return tick() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.vpWidth = msg.Width
		m.vpHeight = msg.Height
		licenseViewport, err := setLicenseViewport(m)
		if err != nil {
			log.Fatal(err)
		}
		m.LicenseViewport = licenseViewport
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			if m.CurrentView != stateLicenseView {
				m.Quiting = true
				return m, tea.Quit
			}
		}
	}

	if m.ExitSetup {
		m.Quiting = true
		return m, tea.Quit
	}

	switch m.CurrentView {
	case 0:
		return updateStartInstallationChoices(msg, m)
	case 1:
		return updateLicenseAcceptChoices(msg, m)
	case 2:
		return updateInstallation(msg, m)
	case 3:
		return updateAfterInstallation(msg, m)
	case 4:
		return updateLicenseView(msg, m)
	default:
		return updateAfterInstallation(msg, m)
	}
}

func (m model) View() string {
	var s string
	if m.Quiting {
		return "\n  Exiting setup wizard\n\n"
	}

	switch m.CurrentView {
	case 0:
		s = StartInstallationView(m)
	case 1:
		s = LicenseAcceptView(m)
	case 2:
		s = InstallationView(m)
	case 3:
		s = AfterInstallationView(m)
	case 4:
		s = LicenseView(m)
	default:
		s = AfterInstallationView(m)
	}
	wh, wv := windowStyle.GetFrameSize()
	mh, mv := mainStyle.GetFrameSize()
	s = lipgloss.Place(
		m.vpWidth-wh-mh,
		m.vpHeight-wv-mv,
		lipgloss.Center,
		lipgloss.Center,
		s,
	)
	return windowStyle.Render(mainStyle.Render(s))
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += ramp[i].Render(progressFullChar)
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Utils

// Generate a blend of colors.
func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
