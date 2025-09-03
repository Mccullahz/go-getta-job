// this file handles the job title input state, and title viewing.
package states

import (
	tea "github.com/charmbracelet/bubbletea"
	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/components"
)

// handle incoming messages while in the title input state
func UpdateTitle(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Title != "" {
				m.CurrentState = model.StateSearching
				m.Err = ""
				return m, StartSearchCmd(m.Zip, m.Radius, m.Title)
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

// render the title input view for the ui
func ViewTitle(m model.Model) string {
	return components.LabelStyle.Render("Enter Job Title/Keyword: ") +
		components.InputStyle.Render(m.Title) + "\n"
}
