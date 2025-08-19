// contains the TUI logic using Bubble Tea and Lipgloss.
package ui

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"cliscraper/internal/geo"
	"cliscraper/internal/web"
	"cliscraper/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateZipInput state = iota
	stateRadiusInput
	stateSearching
	stateDone
)

type model struct {
	currentState state
	zip          string
	radius       string
	err          string
	businesses   []geo.Business
	results      []utils.JobPageResult
	showResults  bool
}

func initialModel() model {
	return model{
		currentState: stateZipInput,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
// main update loop to handle user input and state transitions
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "ctrl+c":
			return m, tea.Quit
		case "f", "F":
			if m.currentState == stateDone {
				results, err := utils.LoadLatestResults("./output")
				if err != nil {
					m.err = "Failed to load results: " + err.Error()
				} else {
					m.results = results
					m.showResults = true
				}
				return m, nil
			}
		}

		switch m.currentState {
		case stateZipInput:
			if msg.Type == tea.KeyEnter {
				if isValidZip(m.zip) {
					m.currentState = stateRadiusInput
					m.err = ""
				} else {
					m.err = "Invalid ZIP code"
				}
			} else if msg.Type == tea.KeyBackspace || msg.Type == tea.KeyDelete {
				if len(m.zip) > 0 {
					m.zip = m.zip[:len(m.zip)-1]
				}
			} else {
				m.zip += msg.String()
			}

	case stateRadiusInput:
	    if msg.Type == tea.KeyEnter {
	        if isValidRadius(m.radius) {
	            m.currentState = stateSearching
	            m.err = ""
	            r, _ := strconv.Atoi(m.radius) // convert string → int
	            return m, searchForJobPages(m.zip, strconv.Itoa(r))
	        } else {
	            m.err = "Radius must be a number"
	        }
		}
		}
	// DONE state handling, builder for the results string used above
	case doneMsg:
		if msg.Err != nil {
			m.err = msg.Err.Error()
			m.currentState = stateDone
		} else {
			m.currentState = stateDone
			m.businesses = msg.Businesses
			if len(msg.Results) == 0 {
				m.err = "No job pages found"
			} else {
			results, err := utils.LoadLatestResults("./output")
			if err == nil {
				m.results = results
			}
			}

		}
	}
	return m, nil
}


// renders the current state of the UI with lipgloss styles
func (m model) View() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("CLI Job Scraper") + "\n\n")

	switch m.currentState {
	case stateZipInput:
		b.WriteString(LabelStyle.Render("Enter ZIP Code: "))
		b.WriteString(InputStyle.Render(m.zip) + "\n")
	case stateRadiusInput:
		b.WriteString(LabelStyle.Render("Enter Search Radius (miles): "))
		b.WriteString(InputStyle.Render(m.radius) + "\n")
	case stateSearching:
		b.WriteString(StatusStyle.Render("Searching for job pages near " + m.zip + "...\n"))
	case stateDone:
		b.WriteString("Search complete! Press 'F' to view results\n")
		b.WriteString(fmt.Sprintf("%d businesses found.\n", len(m.businesses)))
		if m.showResults {
        	for _, r := range m.results {
	            b.WriteString(fmt.Sprintf("%s - %s\n", r.BusinessName, r.URL))
		}
		}
	}


	if m.err != "" {
		b.WriteString("\n" + ErrorStyle.Render("Error: " + m.err) + "\n")
	}

	return b.String()
}

// similarly to the searchForJobPages function below, might move isValid funcs to the utils package
func isValidZip(zip string) bool {
	return len(zip) == 5
}

func isValidRadius(r string) bool {
	for _, ch := range r {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return r != ""
}

type doneMsg struct {
	Businesses []geo.Business
	Results    []utils.JobPageResult
	Err        error
}

// now we tie geo + scraper + detector together
func searchForJobPages(zip, radius string) tea.Cmd {
	return func() tea.Msg {
		r, err := strconv.Atoi(radius)
		if err != nil {
			return doneMsg{Err: fmt.Errorf("invalid radius: %w", err)}
		}

		lat, lon, err := geo.GetCoordinatesFromZip(zip)
		if err != nil {
			return doneMsg{Err: fmt.Errorf("failed to get coordinates for ZIP %s: %w", zip, err)}
		}

		businesses, err := geo.LocateBusinesses(lat, lon, r)
		if err != nil {
			return doneMsg{Err: err}
		}

		// scrape each business for job pages
		var results []utils.JobPageResult
		for i, b := range businesses {
			if b.URL == "" {
				continue
			}

			jobURL, err := web.ScrapeWebsite(b.URL)
			if err != nil || jobURL == "" {
				continue
			}

			// attach job URL back to business
			businesses[i].URL = jobURL

			results = append(results, utils.JobPageResult{
				BusinessName: b.Name,
				URL:          jobURL,
				Description:  "Auto-discovered from scan",
			})
		}

		if err := utils.WriteResults(results, "./output"); err != nil {
			return doneMsg{Err: fmt.Errorf("failed to write results: %w", err)}
		}

		return doneMsg{
			Businesses: businesses,
			Err:        nil,
		}
	}
}
func Run() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
