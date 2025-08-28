package components

import (
	"fmt"
	"strings"

	"cliscraper/internal/utils"
)

func RenderResults(results []utils.JobPageResult) string {
	if len(results) == 0 {
		return "No job pages found."
	}

	var b strings.Builder
	for _, r := range results {
		b.WriteString(fmt.Sprintf("%s - %s\n", r.BusinessName, r.URL))
	}
	return b.String()
}
