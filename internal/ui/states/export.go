package states

import (
	"strings"

	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/components"
	//"cliscraper/internal/utils"
)

func ViewExportSettings(m model.Model) string {
	var b strings.Builder
	// use the StatusStyle for consistent styling
	b.WriteString(components.TitleStyle.Render("Export Settings\n"))
	b.WriteString(components.LabelStyle.Render("Export settings coming soon! Press 'q' to go back! \n"))

	return b.String()
}
