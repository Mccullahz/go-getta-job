package states

import (
	"cliscraper/internal/ui/model"

	tea "github.com/charmbracelet/bubbletea"
)

func UpdateSearching(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case model.DoneMsg:
		if msg.Err != nil {
			m.Err = msg.Err.Error()
			m.CurrentState = model.StateDone
		} else {
			m.CurrentState = model.StateDone
			m.Businesses = msg.Businesses
			m.Results = msg.Results
		}
	}
	return m, nil
}
