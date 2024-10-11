package constants

import "github.com/charmbracelet/bubbles/key"

var HomeHelpString string = "↑/↓: navigate  • esc/q: back • /: filter • crtl+c: quit"

// Keymap struct and key bindings
type homeKeymap struct {
	Enter  key.Binding
	Back   key.Binding
	Delete key.Binding
	Quit   key.Binding
}

// HomeKeymap reusable key mappings shared across models
var HomeKeymap = homeKeymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "back"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
