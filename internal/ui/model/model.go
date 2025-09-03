// global state management, all states are defined, all sub-models will operate on this model

package model

import (
    "cliscraper/internal/geo"
    "cliscraper/internal/utils"
    "cliscraper/internal/ui/components"
    "github.com/charmbracelet/bubbles/list"
)

type state int

const (
    StateHome state = iota
    StateZipInput
    StateRadiusInput
    StateTitleInput
    StateSearching
    StateDone
)

type Model struct {
    CurrentState state
    Zip          string
    Radius       string
    Title        string
    Err          string
    Businesses   []geo.Business

    Results      []utils.JobPageResult
    ShowResults  bool
    ResultsList list.Model
    Starred    []components.JobItem

    InnerCursor int
    TopCursor int

    Width  int
    Height int
}

func InitialModel() Model {
    return Model{
	CurrentState: StateHome,
	Starred: []components.JobItem{},
	TopCursor: 0,
	InnerCursor: 0,
	}
}

