// defines Lipgloss styles for consistent UI formatting.
package components

import(
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff58c4")).
		Padding(1, 2).
		Align(lipgloss.Center)

	LabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#24ccdc")).
		Align(lipgloss.Center)

	InputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#81edef")).
		PaddingLeft(1).
		Align(lipgloss.Center)

	StatusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5deba4")).
		Italic(true).
		Align(lipgloss.Center)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e8404c")).
		Bold(true).
		Align(lipgloss.Center)

	FooterStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true).
		Align(lipgloss.Center)
)

