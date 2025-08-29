// common message types used to communicate between different packages

package messages

import (
    "cliscraper/internal/geo"
    "cliscraper/internal/utils"
)

type DoneMsg struct {
    Businesses []geo.Business
    Results    []utils.JobPageResult
    Err        error
}
