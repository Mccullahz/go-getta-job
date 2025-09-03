package components

import (
	"fmt"
	"io"

	"cliscraper/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type JobItem struct {
	BusinessName string
	URL          string
	Starred      bool
}

// implement list.Item
func (j JobItem) Title() string       { return j.BusinessName }
func (j JobItem) Description() string { return j.URL }
func (j JobItem) FilterValue() string { return j.BusinessName }

// results specific styling will likely move to styles.go at some point for consistency
var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	starStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("227")).Bold(true)
)

type jobDelegate struct{}

func (d jobDelegate) Height() int                               { return 2 }
func (d jobDelegate) Spacing() int                              { return 1 }
func (d jobDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d jobDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(JobItem)
	if !ok {
		return
	}

	title := item.BusinessName
	desc := item.URL

	if index == m.Index() {
		// hightlighting the selected item
		fmt.Fprintf(w, "%s\n%s", selectedStyle.Render(title), desc)
	} else {
		fmt.Fprintf(w, "%s\n%s", dimStyle.Render(title), desc)
	}
}

// build a new list from results
func NewResultsList(results []utils.JobPageResult, width, height int) list.Model {
	if len(results) == 0 {
		items := []list.Item{JobItem{BusinessName: "No job pages found.", URL: ""}}
		l := list.New(items, jobDelegate{}, width, height)
		l.SetShowHelp(false)
		return l
	}

	items := make([]list.Item, 0, len(results))
	for _, r := range results {
		items = append(items, JobItem{BusinessName: r.BusinessName, URL: r.URL})
	}

	l := list.New(items, jobDelegate{}, width, height)
	l.Title = "Job Search Results"
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)

	return l
}

