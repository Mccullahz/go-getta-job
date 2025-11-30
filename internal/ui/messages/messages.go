// common message types used to communicate between different packages

package messages

import (
    "cliscraper/internal/backend/geo"
    "cliscraper/internal/utils"
    "cliscraper/internal/ui/model"
)

type DoneMsg struct {
    Businesses []geo.Business
    Results    []utils.JobPageResult
    Err        error
}

type AuthDoneMsg struct {
    User *model.UserInfo
    Err  error
}
