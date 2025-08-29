package states

import (
	"strings"
	"fmt"

	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/components"
	//"cliscraper/internal/utils"
)

func ViewDone(m model.Model) string {
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

