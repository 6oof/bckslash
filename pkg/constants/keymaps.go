package constants

import "github.com/charmbracelet/bubbles/key"

//HOME KEYMAP

var HomeHelpString string = "↑/↓: navigate • enter: select •  esc/q: back • /: filter • crtl+c: quit"

var HomeAltHelpString string = "↑/↓: navigate • enter: select •  esc/q: back • /: filter"

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

//MODAL/VIEWPORT KEYMAP

var ModalHelpString string = "↑/↓: navigate  • esc/q: back"

type modalKeymap struct {
	Back key.Binding
}

var ModalKeymap = modalKeymap{
	Back: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "back"),
	),
}

//FORM KEYMAP

var FormHelpString string = "esc/q: abort"

type formKeymap struct {
	Back key.Binding
}

var FormKeymap = formKeymap{
	Back: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "back"),
	),
}
