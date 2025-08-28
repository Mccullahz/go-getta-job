// global state management, all states are defined, all sub-models will operate on this model
package ui

import (
	"cliscraper/internal/geo"
	"cliscraper/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// States for the UI
type State int

const (
	StateZipInput State = iota
	StateRadiusInput
	StateTitleInput
	StateSearching
	StateDone
)

// Model is the main UI state
type Model struct {
	CurrentState State
	Zip          string
	Radius       string
	Title        string
	Err          string
	Businesses   []geo.Business
	Results      []utils.JobPageResult
	ShowResults  bool
}

// Message used when search completes
type DoneMsg struct {
	Businesses []geo.Business
	Results    []utils.JobPageResult
	Err        error
}

// InitialModel returns a fresh Model
func InitialModel() Model {
	return Model{
		CurrentState: StateZipInput,
	}
}

// Implements bubbletea.Model interface
func (m Model) Init() tea.Cmd {
	return nil
}
