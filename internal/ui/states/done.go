package states

import (
	"strings"

	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/components"
	"fmt"
)

func RenderDone(m model.Model) string {
	var b strings.Builder
	b.WriteString("Search complete! Press 'F' to view results\n")
	b.WriteString(fmt.Sprintf("%d businesses found.\n", len(m.Businesses)))

	if m.ShowResults {
		b.WriteString(components.RenderResults(m.Results))
	}

	if m.Err != "" {
		b.WriteString("\nError: " + m.Err + "\n")
	}

	return b.String()
}

