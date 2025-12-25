package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)

	// Gradient colors we'll use for the progress bar
	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
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
	initialModel := model{0, 0, false, 0, 10, 0, 0, false, false, nil, false}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

type model struct {
	StartInstallationChoice int
	LicenseAcceptChoice     int
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
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		switch k {
		case "ctrl+c", "q", "esc":
			m.Quiting = true
			return m, tea.Quit
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
	default:
		s = AfterInstallationView(m)
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

func updateStartInstallationChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.StartInstallationChoice++
			if m.StartInstallationChoice > 1 {
				m.StartInstallationChoice = 1
			}
		case "k", "up":
			m.StartInstallationChoice--
			if m.StartInstallationChoice < 0 {
				m.StartInstallationChoice = 0
			}
		case "enter":
			if m.StartInstallationChoice == 1 {
				m.Quiting = true
				return m, tea.Quit
			}
			m.CurrentView++
			return m, frame()
		}
	case tickMsg:
		if m.Ticks == 0 {
			m.Quiting = true
			return m, tea.Quit
		}
		m.Ticks--
		return m, tick()
	}

	return m, nil
}

func updateLicenseAcceptChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.LicenseAcceptChoice++
			if m.LicenseAcceptChoice > 1 {
				m.LicenseAcceptChoice = 1
			}
		case "k", "up":
			m.LicenseAcceptChoice--
			if m.LicenseAcceptChoice < 0 {
				m.LicenseAcceptChoice = 0
			}
		case "enter":
			if m.LicenseAcceptChoice == 1 {
				m.ExitSetup = true
			}
			m.CurrentView++
			return m, frame()
		}
	}
	return m, nil
}

func updateInstallation(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if !m.InstallStarted {
		err := Install()
		m.InstallStarted = true
		if err != nil {
			m.InstallationError = err
		}
	}
	if m.InstallationError != nil {
		m.CurrentView++
	}
	switch msg.(type) {
	case frameMsg:
		if !m.Installed {
			// Placeholder for installation progress: To be changed.
			m.Frames++
			m.Progress = float64(m.Frames) / float64(100)
			if m.Progress >= 1 {
				m.Progress = 1
				m.Installed = true
				m.CurrentView++
				return m, nil
			}
			return m, frame()
		}
	}
	return m, nil
}

func updateAfterInstallation(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func StartInstallationView(m model) string {
	c := m.StartInstallationChoice
	tpl := "Welcome to TermInstaller!\n\n"
	tpl += "%s\n\n"
	tpl += "Setup wizard quits in %s seconds\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s",
		checkbox("Next", c == 0),
		checkbox("Quit setup wizard", c == 1),
	)
	return fmt.Sprintf(tpl, choices, ticksStyle.Render(strconv.Itoa(m.Ticks)))
}

func LicenseAcceptView(m model) string {
	c := m.LicenseAcceptChoice
	tpl := "You need to accept license to install.\n\n"
	tpl += "%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s",
		checkbox("Yes, I accept the license", c == 0),
		checkbox("No, I don't accept the license", c == 1),
	)
	return fmt.Sprintf(tpl, choices)
}

func InstallationView(m model) string {
	var msg string
	msg = "Please wait while setup is installing software to your system..."
	label := "Installing..."
	if m.Installed {
		label = fmt.Sprintf("Done.")
	}

	return msg + "\n\n" + label + "\n" + progressbar(m.Progress) + "%"
}

func AfterInstallationView(m model) string {
	var msg string
	switch m.InstallationError {
	case nil:
		msg = "Installation complete!\n\n"
	default:
		msg = "Error during installation: " + keywordStyle.Render(m.InstallationError.Error()) + "\n\n"
		msg += "Installation is canceled!\n\n"
	}
	return msg + "Press " + keywordStyle.Render("enter") + " to exit wizard."
}

func Install() error {
	if false {
		return errors.New("Installation failed!!!")
	} else {
		return nil
	}
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
