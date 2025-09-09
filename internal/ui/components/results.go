// results lists components for displaying search results and starred jobs lists

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

// general styling for list items
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

	marker := ""
	if item.Starred {
		marker = starStyle.Render("⭐ ")
	}

	title := marker + item.BusinessName
	desc := item.URL

	if index == m.Index() {
		// highlighting the selected item
		fmt.Fprintf(w, "%s\n%s", selectedStyle.Render(title), desc)
	} else {
		fmt.Fprintf(w, "%s\n%s", dimStyle.Render(title), desc)
	}
}

// builder for both results && starred lists
func newJobList(items []JobItem, title string, width, height int, showHelp, filter bool) list.Model {
	raw := make([]list.Item, len(items))
	for i := range items {
		raw[i] = items[i]
	}

	l := list.New(raw, jobDelegate{}, width, height)
	l.Title = title
	l.SetShowHelp(showHelp)
	l.SetFilteringEnabled(filter)
	return l
}

// results list builder
func NewResultsList(results []utils.JobPageResult, width, height int) list.Model {
	if len(results) == 0 {
		return newJobList(
			[]JobItem{{BusinessName: "No job pages found.", URL: ""}},
			"Job Search Results", width, height, false, false,
		)
	}

	items := make([]JobItem, 0, len(results))
	for _, r := range results {
		items = append(items, JobItem{BusinessName: r.BusinessName, URL: r.URL})
	}

	return newJobList(items, "Job Search Results", width, height, true, true)
}

// starred list builder 
func NewStarredList(items []JobItem, width, height int) list.Model {
	if len(items) == 0 {
		return newJobList(
			[]JobItem{{BusinessName: "No starred jobs yet.", URL: ""}},
			"⭐ Starred Jobs", width, height, false, false,
		)
	}
	return newJobList(items, "⭐ Starred Jobs", width, height, true, false)
}

