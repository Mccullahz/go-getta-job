// simliar to the progress bar, following the bubbletea example here: https://github.com/charmbracelet/bubbletea/blob/main/examples/spinner/main.go and repurposing it for the searching state. might use in a separate state to avoid clutter during searching
package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/spinner"
)

type spin struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func initialSpinner() spin {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return spin{spinner: s}
}

func (m spin) Init() tea.Cmd {
	return m.spinner.Tick
}
