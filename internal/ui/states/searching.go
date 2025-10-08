// this file contains the ui components for searching state along with the logic for the searching state itself
package states

import (
	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/messages"
	"cliscraper/internal/ui/components"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"fmt"
)

type DoneMsg = messages.DoneMsg

// handle incoming messages while in the searching state
func UpdateSearching(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case DoneMsg:
		if msg.Err != nil {
			m.Err = msg.Err.Error()
		} else {
			m.Results = msg.Results
		}
		m.CurrentState = model.StateDone
		return m, nil

	default:
		// All other messages ignored
		return m, nil
	}
}

// render the searching view for the ui
func ViewSearching(m model.Model) string {
	return components.StatusStyle.Render(fmt.Sprintf(
		"%s Searching for %s job pages near %s within radius of %s miles...\n",
		m.Spinner.View(), m.Title, m.Zip, m.Radius,
	))
}

// return a tea.Cmd that will run the search asynchronously
func StartSearchCmd(m model.Model, zip, radius, title string) tea.Cmd {
    return tea.Batch(
        m.Spinner.Init(), // spinner tick
        func() tea.Msg {
            results, err := m.Service().Search(zip, radius, title)
            if err != nil {
                return DoneMsg{Err: fmt.Errorf("search failed: %w", err)}
            }

            return DoneMsg{
                Results: results,
                // if API doesn’t send them, leave Businesses nil.
            }
        },
    )
}
