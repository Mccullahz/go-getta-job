// defines Lipgloss styles for consistent UI formatting.
package ui

import(
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff58c4")).
		Padding(1, 2)

	LabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#24ccdc"))

	InputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#81edef")).
		PaddingLeft(1)

	StatusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA")).
		Italic(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e8404c")).
		Bold(true)
)

