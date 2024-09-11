package constants

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Constants and layout values
const (
	BarOffset = 8
	MaxWidth  = 100
	MaxHeight = 40
)

var WinSize tea.WindowSizeMsg

// Color definitions
var (
	SubtleColor    = lipgloss.Color("250") // Background color, always 250
	HighlightColor = lipgloss.Color("197") // Single highlight color used for both text and background
	LightTextColor = lipgloss.Color("15")  // White text for light themes
	DarkTextColor  = lipgloss.Color("238") // Dark text for dark themes

	// Styles
	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

	DocStyle = lipgloss.NewStyle().Padding(1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(DarkTextColor).
			Background(SubtleColor).
			Height(1)

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

// Layout and rendering functions
func BodyHeight() int {
	// Ensure height does not exceed MaxHeight
	if WinSize.Height > MaxHeight {
		return MaxHeight
	}
	return WinSize.Height
}

func BodyWidth() int {
	// Ensure width does not exceed MaxWidth
	if WinSize.Width > MaxWidth {
		return MaxWidth
	}
	return WinSize.Width
}

var MainHelpString string = "↑/↓: navigate  • esc: back • tab: switch focus • /: filter • d: delete • q: quit"

// Keymap struct and key bindings
type keymap struct {
	Enter  key.Binding
	Back   key.Binding
	Delete key.Binding
	Quit   key.Binding
	Tab    key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}

// Composition and layout rendering functions
func HalfAndHalfComposition(left, right string, height int) string {
	leftList := lipgloss.NewStyle().Width(BodyWidth() / 2).Height(height).Render(left)
	rightList := lipgloss.NewStyle().Width(BodyWidth() / 2).Height(height).Render(right)
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
	return content
}

func ListStyle() list.Styles {
	listStyle := list.DefaultStyles()

	listStyle.Title = listStyle.Title.Background(HighlightColor)

	return listStyle
}
