package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Constants and layout values
const (
	BarOffset        = 8
	MaxWidth         = 150
	MaxHeight        = 70
	MinWidthForSplit = 100
)

var WinSize tea.WindowSizeMsg

// Color definitions
var (
	SubtleColor         = lipgloss.Color("250") // Background color, always 250
	HighlightColor      = lipgloss.Color("197") // Single highlight color used for both text and background
	LightTextColor      = lipgloss.Color("250") // White text for light themes
	LightTextColorMuted = lipgloss.Color("245") // White text for light themes
	DarkTextColor       = lipgloss.Color("238") // Dark text for dark themes

	// Styles
	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

	DocStyle = lipgloss.NewStyle().Width(BodyWidth()).Height(BodyHeight())

	PadBodyContent = lipgloss.NewStyle().Width(BodyWidth()).Padding(2, 0)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(DarkTextColor).
			Background(SubtleColor).
			Width(BodyWidth()).
			Height(1).
			MaxHeight(1)

	StatusStyle = lipgloss.NewStyle().
			Inherit(StatusBarStyle).
			Foreground(LightTextColor).
			Background(HighlightColor).
			Padding(0, 1).
			MarginRight(1)

	FishCakeStyle = lipgloss.NewStyle().
			Foreground(LightTextColor).
			Background(HighlightColor).
			Padding(0, 1)

	StatusText = lipgloss.NewStyle().Inherit(StatusBarStyle)
)

// Layout and rendering functions
func BodyHeight() int {
	// Ensure height does not exceed MaxHeight
	if WinSize.Height > MaxHeight {
		return MaxHeight - 8
	}
	return WinSize.Height - 8
}

func BodyWidth() int {
	// Ensure width does not exceed MaxWidth
	if WinSize.Width > MaxWidth {
		return MaxWidth
	}
	return WinSize.Width
}
