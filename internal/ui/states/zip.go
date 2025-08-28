package states

import (
	//"strings"

	tea "github.com/charmbracelet/bubbletea"
	"cliscraper/internal/ui/components"
	"cliscraper/internal/ui"
)

func UpdateZip(m ui.Model, msg tea.Msg) (ui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if isValidZip(m.Zip) {
				m.CurrentState = ui.StateRadiusInput
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

func ViewZip(m ui.Model) string {
	return components.LabelStyle.Render("Enter ZIP Code: ") +
		components.InputStyle.Render(m.Zip) + "\n"
}

func isValidZip(zip string) bool {
	return len(zip) == 5
}
