package states

import (
	"strings"

	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/components"
	//"cliscraper/internal/utils"
)

func ViewDone(m model.Model) string {
	var b strings.Builder
	// use the StatusStyle for consistent styling
	b.WriteString(components.StatusStyle.Render("Search Complete!\n\n"))
	b.WriteString(components.LabelStyle.Render("Press 'f' to view results from the latest search.\n"))
	// currently we are just rendering the formatted results directly, will be changing this to a list with further interaction options soon
	if m.ShowResults {
		b.WriteString(components.RenderResults(m.Results))
	}

	if m.Err != "" {
		b.WriteString("\nError: " + m.Err + "\n")
	}

	return b.String()
}

