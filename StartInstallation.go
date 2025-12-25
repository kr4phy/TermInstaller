package main

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

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
