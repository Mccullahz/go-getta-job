// this file handles the starred state, where users can view their starred job pages
package states

import (
	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/components"
	//"cliscraper/internal/utils"
)

func ViewStarred(m model.Model) string {
	if len(m.Starred) == 0 {
		return components.StatusStyle.Render("No starred jobs yet.\n")
	}

	// temporary list.Model
	l := components.NewStarredList(m.Starred, m.Width, m.Height-2)
	return l.View()
}



