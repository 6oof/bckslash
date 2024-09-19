package constants

import "github.com/charmbracelet/bubbles/key"

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
