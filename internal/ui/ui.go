// contains the TUI logic using Bubble Tea and Lipgloss.
package ui

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"cliscraper/internal/geo"
	"cliscraper/internal/web"
	"cliscraper/internal/output"
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
	results      []output.JobPageResult
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
		// F key to show results SUBOPTIMAL WAY TO DO THIS, STRING BUILDER BROKEN
		case "f", "F":
			if m.currentState == stateDone {
				results, err := utils.LoadLatestResults("./output") //borked
				if err != nil {
					m.err = "Failed to load results: " + err.Error()
				} else {
				m.results = results
				m.showResults = true

				if m.showResults && len(m.businesses) > 0 {
					for _, b := range m.results {
						out := fmt.Sprintf("%s - %s\n", b.BusinessName, b.URL)
						fmt.Printf("\n" + out)
			/*SORRY ME, CLOSING THESE {} ARE TERRIBLE*/
					}
				}

				}
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
					return m, searchForJobPages(m.zip, m.radius)
				} else {
					m.err = "Radius must be a number"
				}
			} else if msg.Type == tea.KeyBackspace || msg.Type == tea.KeyDelete {
				if len(m.radius) > 0 {
					m.radius = m.radius[:len(m.radius)-1]
				}
			} else {
				m.radius += msg.String()
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
			var results []output.JobPageResult
			for _, b := range msg.Businesses {
				results = append(results, output.JobPageResult{
					BusinessName: b.Name,
					URL:          b.URL,
					Description:  "Auto-discovered from scan",
				})
			}
			if err := output.WriteResults(results, "./output"); err != nil {
				m.err = "Failed to write results: " + err.Error()
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

type doneMsg struct{
	Businesses []geo.Business
	Err        error
}

// it might be cleaner to move this function into the utils package and just call it with zip and radius args from the model
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
        return doneMsg{
            Businesses: businesses,
            Err:        err,
        }
	
	var results []output.JobPageResult
	for _, b := range businesses {
			jobURL, err := web.ScrapeWebsite(b.URL)
			if err != nil || jobURL == "" {
				continue
			}
			results = append(results, output.JobPageResult{
				BusinessName: b.Name,
				URL:          jobURL,
				Description:  "Auto-discovered from scan",
			})
			}
		if err := output.WriteResults(results, "./output"); err != nil {
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
