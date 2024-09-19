package constants

import "github.com/charmbracelet/lipgloss"

func BodyHalfWidth() int {
	// If the screen is too narrow, return the full width
	if BodyWidth() < MinWidthForSplit {
		return BodyWidth()
	}

	// Otherwise, return half of the available width
	return BodyWidth() / 2
}

func HalfAndHalfComposition(left, right string) string {
	// Check if the screen width is smaller than the defined minimum width for splitting
	if BodyWidth() < MinWidthForSplit {
		// If the screen is too narrow, only display the left half
		return lipgloss.NewStyle().Width(BodyWidth()).Height(BodyHeight()).Render(left)
	}

	// Otherwise, display both the left and right halves
	leftList := lipgloss.NewStyle().Width(BodyWidth() / 2).Height(BodyHeight()).Render(left)
	rightList := lipgloss.NewStyle().Width(BodyWidth() / 2).Height(BodyHeight()).Render(right)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftList, rightList)
}

func RenderBar(location string) string {
	w := lipgloss.Width

	statusKey := StatusStyle.Render(`\`)
	fishCake := FishCakeStyle.Render(location)
	statusVal := StatusText.Width(BodyWidth() - w(statusKey) - w(fishCake)).Render("Bckslash")

	bar := lipgloss.JoinHorizontal(lipgloss.Top, statusKey, statusVal, fishCake)

	return StatusBarStyle.Width(BodyWidth()).Render(bar)
}

func RenderHelpBar(wi int, helpstring string) string {
	statusVal := StatusText.Padding(0, 1).Render(helpstring)
	return StatusBarStyle.Width(BodyWidth()).Render(statusVal)
}

func Layout(location, helpstring, children string) string {
	content := lipgloss.JoinVertical(lipgloss.Top,
		RenderBar(location),
		DocStyle.Render(children),
		RenderHelpBar(BodyWidth(), helpstring),
	)

	// Apply border style to the entire layout
	return lipgloss.Place(WinSize.Width, WinSize.Height, lipgloss.Center, lipgloss.Center, content)
}

func Card(content, background string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content, lipgloss.WithWhitespaceChars(background), lipgloss.WithWhitespaceForeground(lipgloss.Color("240")))

}
