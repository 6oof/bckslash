package layout

import (
	"lg/views/constants"

	"github.com/charmbracelet/lipgloss"
)

func RenderBar(location string) string {
	// Status bar
	w := lipgloss.Width

	statusKey := constants.StatusStyle.Render(`\`)
	fishCake := constants.FishCakeStyle.Render(location)
	statusVal := constants.StatusText.
		Width(constants.BodyWidth() - w(statusKey) - w(fishCake)).
		Render("Bckslash")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		fishCake,
	)

	return constants.StatusBarStyle.Width(constants.BodyWidth()).Render(bar)
}

func RenderHelpBar(wi int, helpstring string) string {
	// Status bar

	statusVal := constants.StatusText.Padding(0, 1).Render(helpstring)

	return constants.StatusBarStyle.Width(constants.BodyWidth()).Render(statusVal)
}

func Layout(location, helpstring, children string) string {

	return lipgloss.JoinVertical(lipgloss.Left, RenderBar(location), constants.DocStyle.Render(children), RenderHelpBar(constants.BodyWidth(), helpstring))
}
