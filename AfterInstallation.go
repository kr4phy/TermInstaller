package main

import tea "github.com/charmbracelet/bubbletea"

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
