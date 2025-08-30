// entry point for the tui, glues everything in the component + states packages together
// this file should be minimal now, just to unfuck it
package ui

import (
	"fmt"
	"os"
	"strings"

	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/states"
	"cliscraper/internal/ui/components"
	"cliscraper/internal/ui/messages"
	"cliscraper/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

// UI wraps model.Model so we can define methods on it
type UI struct {
	model.Model
	Width       int
	Height      int
}

// satisfies bubbletea.Model interface
func (u UI) Init() tea.Cmd {
	return nil
}

// main update loop, delegates to state updates to modularize ux flow
func (u UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		u.Width = msg.Width
		u.Height = msg.Height
		return u, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "ctrl+c":
			return u, tea.Quit
		case "f", "F":
			if u.CurrentState == model.StateDone && !u.ShowResults {
				results, err := utils.LoadLatestResults("./output")
				if err != nil {
					u.Err = "Failed to load results: " + err.Error()
				} else {
					u.Results = results
					u.ShowResults = true
				}
				return u, nil
			}
		}

		// delegate to state updates
		switch u.CurrentState {
		case model.StateZipInput:
			u.Model, cmd = states.UpdateZip(u.Model, msg)
		case model.StateRadiusInput:
			u.Model, cmd = states.UpdateRadius(u.Model, msg)
		case model.StateTitleInput:
			u.Model, cmd = states.UpdateTitle(u.Model, msg)
		case model.StateSearching:
			u.Model, cmd = states.UpdateSearching(u.Model, msg)
		}

	case messages.DoneMsg:
		u.Model, cmd = states.UpdateSearching(u.Model, msg)
	}

	return u, cmd
}

// main view function, delegates to state views to modularize the ui
func (u UI) View() string {
	var b strings.Builder

	// main content
	switch u.CurrentState {
	case model.StateZipInput:
		b.WriteString(states.ViewZip(u.Model))
	case model.StateRadiusInput:
		b.WriteString(states.ViewRadius(u.Model))
	case model.StateTitleInput:
		b.WriteString(states.ViewTitle(u.Model))
	case model.StateSearching:
		b.WriteString(states.ViewSearching(u.Model))
	case model.StateDone:
		b.WriteString(states.ViewDone(u.Model))
	}

	if u.Err != "" {
		b.WriteString("\n" + components.ErrorStyle.Render("Error: "+u.Err) + "\n")
	}

	// footer content
	tips := "q / ctrl + c : quit     f : show results (if any)     j / k + arrow keys : scroll results"
	footer := ("\n" + components.FooterStyle.Render(tips) + "\n")

	// padding footer to bottom of screen -- currently padding too far down and cannot see the main content
	contentHeight := strings.Count(b.String(), "\n") + 1
	paddingLines := u.Height - contentHeight - 2 // -2 for footer
	if paddingLines > 0 {
		b.WriteString(strings.Repeat("\n", paddingLines))
	}
	b.WriteString(footer)

	return b.String()
}

// entry point passed to main.go
func Run() {
	p := tea.NewProgram(UI{Model: model.InitialModel()}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

