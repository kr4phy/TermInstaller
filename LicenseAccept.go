package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

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
