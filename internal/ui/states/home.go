package states

import (
	"strings"
	"fmt"

	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var headers = []string{"Search", "Starred Jobs", "Settings"}

var options = map[string][]string{
	"Search":    {"Start New Search", "View Last Results"},
	"Starred Jobs": {"View All", "Export"},
	"Settings":   {"Account Settings", "Output - Export Preferences"},
}

// Styles
var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(2)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("63")).Bold(true).Padding(0, 1)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
)

func UpdateHome(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			if m.TopCursor > 0 {
				m.TopCursor--
				m.InnerCursor = 0
			}
		case "l", "right":
			if m.TopCursor < len(headers)-1 {
				m.TopCursor++
				m.InnerCursor = 0
			}
		case "j", "down":
			curHeader := headers[m.TopCursor]
			if m.InnerCursor < len(options[curHeader])-1 {
				m.InnerCursor++
			}
		case "k", "up":
			if m.InnerCursor > 0 {
				m.InnerCursor--
			}
		case "enter":
			curHeader := headers[m.TopCursor]
			curOption := options[curHeader][m.InnerCursor]
			// switch states depending on option
			if curHeader == "Search" && curOption == "Start New Search" {
				m.CurrentState = model.StateZipInput
			}
			if curHeader == "Search" && curOption == "View Last Results" {
				m.CurrentState = model.StateDone
			}
			if curHeader == "Starred Jobs" && curOption == "View All" {
				m.StarredList = components.NewStarredList(m.Starred, m.Width, m.Height-2)
				m.CurrentState = model.StateStarred
			}
			// other options to be handled later

		}
	}
	return m, nil
}

func ViewHome(m model.Model) string {
	var headerParts []string
	for i, h := range headers {
		if i == m.TopCursor {
			headerParts = append(headerParts,
				lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render(h))
		} else {
			headerParts = append(headerParts, h)
		}
	}
	topBar := strings.Join(headerParts, "   ")

	var opts []string
	curHeader := headers[m.TopCursor]
	for i, opt := range options[curHeader] {
		cursor := "  "
		if i == m.InnerCursor {
			cursor = "> "
		}
		opts = append(opts, fmt.Sprintf("%s %s", cursor, opt))
	}
	content := strings.Join(opts, "\n")

	// place everything centered horizontally, this actually does nothing ;(
	ui := lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Top,
		topBar+"\n\n"+content,
	)

	return ui
}

