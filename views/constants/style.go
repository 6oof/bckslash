package constants

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Style definitions.

const (
	BarOffset = 8
)

var WinSize tea.WindowSizeMsg

var (
	SubtleColor     = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	HighlightColor  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	SpecialColor    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	ForegroundColor = lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}
	WhiteTextColor  = lipgloss.Color("#FFFDF5")
	BGHighlightOne  = lipgloss.Color("#FF5F87")
	BGHighlightTwo  = lipgloss.Color("#6124DF")

	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

	DocStyle = lipgloss.NewStyle().Margin(1, 2)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ForegroundColor).
			Background(SubtleColor).
			Margin(1, 2, 1).
			Height(1)

	StatusStyle = lipgloss.NewStyle().
			Inherit(StatusBarStyle).
			Foreground(WhiteTextColor).
			Background(BGHighlightOne).
			Padding(0, 1).
			MarginRight(1)

	FishCakeStyle = lipgloss.NewStyle().
			Foreground(WhiteTextColor).
			Background(BGHighlightTwo).
			Padding(0, 1)

	StatusText = lipgloss.NewStyle().Inherit(StatusBarStyle)
)

func FocusedListDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(HighlightColor).BorderLeftForeground(HighlightColor)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(HighlightColor).BorderLeftForeground(HighlightColor)

	return d
}

func UnfocusedListDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.NormalTitle
	d.Styles.SelectedDesc = d.Styles.NormalDesc

	return d
}

func BodyHeight() int {
	return WinSize.Height - 4 - 4
}

func BodyWidth() int {
	w := WinSize.Width - 4
	return w
}
