// entry point for the tui, glues everything in the component + states packages together

package ui

import (
	"fmt"
	"os"
	//"strings"

	"cliscraper/internal/ui/components"
	"cliscraper/internal/ui/states"
	"cliscraper/internal/ui/model"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	p := tea.NewProgram(model.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Example View function that uses states
func RenderUI(m model.Model) string {
	switch m.CurrentState {
	case model.StateZipInput:
		return renderZipInput(m)
	case model.StateRadiusInput:
		return renderRadiusInput(m)
	case model.StateTitleInput:
		return renderTitleInput(m)
	case model.StateSearching:
		return renderSearching(m)
	case model.StateDone:
		return states.RenderDone(m)
	default:
		return "Unknown state"
	}
}

