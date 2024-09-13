package constants

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Constants and layout values
const (
	BarOffset        = 8
	MaxWidth         = 100
	MaxHeight        = 40
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

var MainHelpString string = "↑/↓: navigate  • esc: back • /: filter • q: quit"

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
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content, lipgloss.WithWhitespaceChars(background), lipgloss.WithWhitespaceForeground(lipgloss.Color("235")))

}

// bubbles and huh overrides

func ListStyle() list.Styles {
	listStyle := list.DefaultStyles()
	listStyle.Title = listStyle.Title.Background(HighlightColor).MarginTop(1)
	return listStyle
}

func HuhBsTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color(HighlightColor))
	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color(HighlightColor)).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(lipgloss.Color(HighlightColor))
	t.Focused.Directory = t.Focused.Directory.Foreground(lipgloss.Color(HighlightColor))
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.Color(LightTextColorMuted))
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(lipgloss.Color("9"))
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(lipgloss.Color("9"))
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(lipgloss.Color(HighlightColor))
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(lipgloss.Color("3"))
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(lipgloss.Color("3"))
	t.Focused.Option = t.Focused.Option.Foreground(lipgloss.Color(HighlightColor))
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(lipgloss.Color("3"))
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lipgloss.Color(HighlightColor))
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(lipgloss.Color(HighlightColor))
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(lipgloss.Color("7"))
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(lipgloss.Color("7")).Background(lipgloss.Color(HighlightColor))
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(lipgloss.Color("7")).Background(lipgloss.Color("8"))

	t.Focused.TextInput.Cursor.Foreground(lipgloss.Color(HighlightColor))
	t.Focused.TextInput.Placeholder.Foreground(lipgloss.Color("8"))
	t.Focused.TextInput.Prompt.Foreground(lipgloss.Color("3"))

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.NoteTitle = t.Blurred.NoteTitle.Foreground(lipgloss.Color(LightTextColor))
	t.Blurred.Title = t.Blurred.NoteTitle.Foreground(lipgloss.Color(LightTextColor))

	t.Blurred.TextInput.Prompt = t.Blurred.TextInput.Prompt.Foreground(lipgloss.Color("8"))
	t.Blurred.TextInput.Text = t.Blurred.TextInput.Text.Foreground(lipgloss.Color(LightTextColor))

	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}

func CustomDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(HighlightColor).BorderLeftForeground(HighlightColor)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(HighlightColor).BorderLeftForeground(HighlightColor)
	return d
}
