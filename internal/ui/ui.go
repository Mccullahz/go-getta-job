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
	//"cliscraper/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

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
		// q sends to previous state or quits if at home
		case "q", "Q", "esc":
			if u.CurrentState == model.StateHome {
				return u, tea.Quit
				} else if u.CurrentState == model.StateTitleInput{
					//do nothing, prevent going back to radius input -- below is going back on q, need for testing starred
					u.CurrentState = model.PreviousState(u.CurrentState)
				} else {
					u.CurrentState = model.PreviousState(u.CurrentState)
					
			return u, nil
		}
		// ctrl+c always quits completely
		case "ctrl+c":
			return u, tea.Quit
		case "f", "F":
		    if u.CurrentState == model.StateDone && !u.ShowResults {
			// fetch results from server instead of lfs
		        results, err := u.Service().Results()
		        if err != nil {
		            u.Err = "Failed to load results: " + err.Error()
		        } else {
		            u.Results = results
		            u.ResultsList = components.NewResultsList(results, u.Width, u.Height-2)
		            u.ShowResults = true
		        }
		        return u, nil
		    }
		// using s for the time being to toggle starred state
    		case "s":
		        if u.CurrentState == model.StateDone && u.ShowResults {
				idx := u.ResultsList.Index()
				if idx >= 0 && len(u.ResultsList.Items()) > 0 {
					if it, ok := u.ResultsList.SelectedItem().(components.JobItem); ok {
						it.Starred = !it.Starred
						u.ResultsList.SetItem(idx, it)
						if it.Starred {
							// add
							u.Starred = append(u.Starred, it)
						} else {
							// remove
							for i := range u.Starred {
								if u.Starred[i].URL == it.URL {
									u.Starred = append(u.Starred[:i], u.Starred[i+1:]...)
									break
								}
							}
						}
					}
				}
            return u, nil
        }
    }


		// delegate to state updates
		switch u.CurrentState {
		case model.StateHome:
			u.Model, cmd = states.UpdateHome(u.Model, msg)
		case model.StateZipInput:
			u.Model, cmd = states.UpdateZip(u.Model, msg)
		case model.StateRadiusInput:
			u.Model, cmd = states.UpdateRadius(u.Model, msg)
		case model.StateTitleInput:
			u.Model, cmd = states.UpdateTitle(u.Model, msg)
		case model.StateSearching:
			u.Model, cmd = states.UpdateSearching(u.Model, msg)
		case model.StateStarred:
			var c tea.Cmd
			u.StarredList, c = u.StarredList.Update(msg)
			return u, c
		case model.StateDone:
		    if u.ShowResults {
		        var c tea.Cmd
		        u.ResultsList, c = u.ResultsList.Update(msg)
		        return u, c
		    }
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
	case model.StateHome:
		b.WriteString(states.ViewHome(u.Model))
	case model.StateZipInput:
		b.WriteString(states.ViewZip(u.Model))
	case model.StateRadiusInput:
		b.WriteString(states.ViewRadius(u.Model))
	case model.StateTitleInput:
		b.WriteString(states.ViewTitle(u.Model))
	case model.StateSearching:
		b.WriteString(states.ViewSearching(u.Model))
	case model.StateStarred:
		if len(u.StarredList.Items()) == 0 {
			b.WriteString(components.StatusStyle.Render("No starred jobs yet.\n"))
		} else {
			b.WriteString(u.StarredList.View())
		}
	case model.StateDone:
	    if u.ShowResults {
	        b.WriteString(u.ResultsList.View())
	    } else {
	        b.WriteString(states.ViewDone(u.Model))
	    }
	}

	if u.Err != "" {
		b.WriteString("\n" + components.ErrorStyle.Render("Error: "+u.Err) + "\n")
	}

	// footer content
	tips := "q / ctrl + c : quit     f : show results (if any)     j / k + h / l + arrow keys : scroll results"
	footer := ("\n" + components.FooterStyle.Render(tips) + "\n")

	// padding footer to bottom of screen
	contentHeight := strings.Count(b.String(), "\n") + 1
	paddingLines := u.Height - contentHeight - 2 // -2 for footer
	if paddingLines > 0 {
		b.WriteString(strings.Repeat("\n", paddingLines))
	}
	b.WriteString(footer)

	return b.String()
}

// entry point passed to main.go
func Run(svc model.Service) {
    p := tea.NewProgram(UI{Model: model.InitialModel(svc)}, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}
