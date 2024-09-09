package compositions

import (
	"lg/views/constants"

	"github.com/charmbracelet/lipgloss"
)

func HalfAndHalfComposition(left, right string, height int) string {

	// Set the width for the lists and render them
	leftList := lipgloss.NewStyle().Width(constants.BodyWidth() / 2).Height(height).Render(left)
	rightList := lipgloss.NewStyle().Width(constants.BodyWidth() / 2).Height(height).Render(right)

	// Return the layout with two columns, each taking 50% of the barWidth
	return lipgloss.JoinHorizontal(lipgloss.Top, leftList, rightList)
}
