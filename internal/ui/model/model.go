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
    StateStarred
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
    StarredList list.Model

    InnerCursor int
    TopCursor int

    Width  int
    Height int
}

func PreviousState(s state) state {
	switch s {
	case StateZipInput:
		return StateHome
	case StateRadiusInput:
		return StateZipInput
	case StateTitleInput:
		return StateRadiusInput
	case StateSearching:
		return StateTitleInput
	case StateStarred:
		return StateHome
	case StateDone:
		return StateHome
	default:
		return StateHome
	}
}


func InitialModel() Model {
    return Model{
	CurrentState: StateHome,
	Starred: []components.JobItem{},
	TopCursor: 0,
	InnerCursor: 0,
	}
}

