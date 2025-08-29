package states

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/components"
)
//need to feed in SearchForJobPages function from states/searching.go not from model
func UpdateRadius(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if isValidRadius(m.Radius) {
				m.CurrentState = model.StateSearching
				m.Err = ""
				return m, StartSearchCmd(m.Zip, m.Radius)
			} else {
				m.Err = "Radius must be a number"
			}
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.Radius) > 0 {
				m.Radius = m.Radius[:len(m.Radius)-1]
			}
		default:
			// numeric input only
			if msg.String() >= "0" && msg.String() <= "9" {
				m.Radius += msg.String()
			} else if msg.String() == "." && !strings.Contains(m.Radius, ".") {
				m.Radius += msg.String()
			}
		}
	}
	return m, nil
}

func ViewRadius(m model.Model) string {
	return components.LabelStyle.Render("Enter Search Radius (miles): ") +
		components.InputStyle.Render(m.Radius) + "\n"
}

func isValidRadius(r string) bool {
	for _, ch := range r {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return r != ""
}

