package constants

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

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
