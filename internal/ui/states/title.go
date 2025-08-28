package states

import (
	tea "github.com/charmbracelet/bubbletea"
	"cliscraper/internal/ui"
	"cliscraper/internal/ui/components"
)

func UpdateTitle(m ui.Model, msg tea.Msg) (ui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Title != "" {
				m.CurrentState = ui.StateSearching
				m.Err = ""
				return m, ui.SearchForJobPages(m.Zip, m.Radius)
			} else {
				m.Err = "Job title cannot be empty"
			}
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.Title) > 0 {
				m.Title = m.Title[:len(m.Title)-1]
			}
		default:
			m.Title += msg.String()
		}
	}
	return m, nil
}

func ViewTitle(m ui.Model) string {
	return components.LabelStyle.Render("Enter Job Title: ") +
		components.InputStyle.Render(m.Title) + "\n"
}

