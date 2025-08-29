package states

import (
	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/messages"
	"cliscraper/internal/ui/components"
	"cliscraper/internal/geo"
	"cliscraper/internal/web"
	"cliscraper/internal/utils"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// when the search complete, done message is sent
type DoneMsg = messages.DoneMsg

// handle incoming messages while in the searching state
func UpdateSearching(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DoneMsg:
		if msg.Err != nil {
			m.Err = msg.Err.Error()
		} else {
			m.Businesses = msg.Businesses
			m.Results = msg.Results
		}
		m.CurrentState = model.StateDone
		return m, nil
	default:
		// searching state, all other messages are ignored
		return m, nil
	}
}
// render the searching view for the ui
func ViewSearching(m model.Model) string {
	return components.StatusStyle.Render(fmt.Sprintf(
		"Searching for job pages near ZIP %s with radius %s miles...\n",
		m.Zip, m.Radius,
	))
}

// return a tea.Cmd that will run the search asynchronously
func StartSearchCmd(zip, radius string) tea.Cmd {
	return func() tea.Msg {
		r, err := strconv.Atoi(radius)
		if err != nil {
			return DoneMsg{Err: fmt.Errorf("invalid radius: %w", err)}
		}

		lat, lon, err := geo.GetCoordinatesFromZip(zip)
		if err != nil {
			return DoneMsg{Err: fmt.Errorf("failed to get coordinates for ZIP %s: %w", zip, err)}
		}

		businesses, err := geo.LocateBusinesses(lat, lon, r)
		if err != nil {
			return DoneMsg{Err: err}
		}

		var results []utils.JobPageResult
		for i, b := range businesses {
			if b.URL == "" {
				continue
			}
			jobURL, err := web.ScrapeWebsite(b.URL)
			if err != nil || jobURL == "" {
				continue
			}
			businesses[i].URL = jobURL
			results = append(results, utils.JobPageResult{
				BusinessName: b.Name,
				URL:          jobURL,
				Description:  "Auto-discovered from scan",
			})
		}

		if err := utils.WriteResults(results, "./output"); err != nil {
			return DoneMsg{Err: fmt.Errorf("failed to write results: %w", err)}
		}

		return DoneMsg{
			Businesses: businesses,
			Results:    results,
			Err:        nil,
		}
	}
}

