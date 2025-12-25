package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

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

func InstallationView(m model) string {
	var msg string
	msg = "Please wait while setup is installing software to your system..."
	label := "Installing..."
	if m.Installed {
		label = fmt.Sprintf("Done.")
	}

	return msg + "\n\n" + label + "\n" + progressbar(m.Progress) + "%"
}
