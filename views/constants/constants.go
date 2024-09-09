package constants

import (
	"github.com/charmbracelet/bubbles/key"
)

var MainHelpString string = "↑/↓: navigate  • esc: back • tab: switch focus • /: filter (where avaliable) • q: quit"

type keymap struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Tab   key.Binding
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
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
