package states

import (
	//"strings"

	tea "github.com/charmbracelet/bubbletea"
	"cliscraper/internal/ui/components"
	"cliscraper/internal/ui/model"
)

func UpdateZip(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if isValidZip(m.Zip) {
				m.CurrentState = model.StateRadiusInput
				m.Err = ""
			} else {
				m.Err = "Invalid ZIP code"
			}
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.Zip) > 0 {
				m.Zip = m.Zip[:len(m.Zip)-1]
			}
		default:
			m.Zip += msg.String()
		}
	}
	return m, nil
}

func ViewZip(m model.Model) string {
	return components.LabelStyle.Render("Enter ZIP Code: ") +
		components.InputStyle.Render(m.Zip) + "\n"
}

func isValidZip(zip string) bool {
	return len(zip) == 5
}
