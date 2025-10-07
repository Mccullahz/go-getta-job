// simliar to the progress bar, following the bubbletea example here: https://github.com/charmbracelet/bubbletea/blob/main/examples/spinner/main.go and repurposing it for the searching state. might use in a separate state to avoid clutter during searching
package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Spin struct {
	Spinner  spinner.Model
	Quitting bool
	Err      error
}

func InitialSpinner() Spin {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Spin{Spinner: s}
}

func (m *Spin) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m *Spin) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return cmd
}

func (m *Spin) View() string {
	return m.Spinner.View()
}

