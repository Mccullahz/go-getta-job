// global state management, all states are defined, all sub-models will operate on this model

package model

import (
    "cliscraper/internal/geo"
    "cliscraper/internal/utils"
)

type state int

const (
    StateZipInput state = iota
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
}

func InitialModel() Model {
    return Model{CurrentState: StateZipInput}
}

