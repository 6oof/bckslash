package views

import (
	"lg/views/commands"
	"lg/views/compositions"
	"lg/views/constants"
	"lg/views/layout"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type navigate int

const (
	serverStats navigate = iota
	serverInfo
	serverSettings
	no
)

type item struct {
	title, desc string
	navigation  navigate
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	listLeft  list.Model
	listRight list.Model
	loading   bool
	quitting  bool
	focus     int // 0 for left list, 1 for right list
}

func InitHomeModel() model {
	itemsLeft := []list.Item{
		item{title: "Server info", desc: "Server monitoring dashboard", navigation: serverInfo},
		item{title: "Server stats", desc: "Server monitoring dashboard", navigation: serverStats},
		item{title: "Server settings", desc: "Set default editor and edit firewall settings", navigation: serverSettings},
		item{title: "Domains", desc: "Edit caddy configuration", navigation: no},
	}

	itemsRight := []list.Item{
		item{title: "Logs", desc: "View system logs"},
		item{title: "Services", desc: "Manage server services"},
		item{title: "Storage", desc: "View storage information"},
	}

	m := model{
		listLeft:  list.New(itemsLeft, constants.UnfocusedListDelegate(), 0, 0),
		listRight: list.New(itemsRight, constants.UnfocusedListDelegate(), 0, 0),
		focus:     0, // Start with focus on the left list
		quitting:  false,
		loading:   true,
	}

	m.listLeft.Title = "Server"
	m.listRight.Title = "Projects"
	m.listRight.SetShowHelp(false)
	m.listLeft.SetShowHelp(false)

	return m
}

func GoHome() (tea.Model, tea.Cmd) {
	homeModel := InitHomeModel()
	return homeModel.Update(constants.WinSize)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Tab):
			// Switch focus between lists
			m.focus = (m.focus + 1) % 2
			return m, nil
		case key.Matches(msg, constants.Keymap.Enter):
			switch m.focus {
			case 0:
				switch m.listLeft.SelectedItem().(item).navigation {
				case serverStats:
					return InitServerStatsModel().Update(commands.ExecStartMsg{})

				case serverInfo:
					return InitServerInfoModel().Update(commands.ExecStartMsg{})

				case serverSettings:
					return InitServerSettingsModel().Update(commands.ExecStartMsg{})
				}

			case 1:
			}
		case key.Matches(msg, constants.Keymap.Back):
			return m, nil

		}
	case tea.WindowSizeMsg:
		constants.WinSize = msg
		m.listLeft.SetSize(constants.BodyWidth()/2, constants.BodyHeight())
		m.listRight.SetSize(constants.BodyWidth()/2, constants.BodyHeight())
		m.loading = false
	}

	// Update the focused list
	var cmd tea.Cmd
	if m.focus == 0 {
		m.listLeft, cmd = m.listLeft.Update(msg)
	} else {
		m.listRight, cmd = m.listRight.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {

	if m.quitting {
		return ""
	}

	if m.focus == 0 {
		m.listLeft.SetDelegate(constants.FocusedListDelegate())
	} else {
		m.listRight.SetDelegate(constants.FocusedListDelegate())
	}

	// Set the width for the lists and render them
	if m.loading {
		return layout.Layout("home", "", "loading...")
	} else {
		return layout.Layout("home", constants.MainHelpString, compositions.HalfAndHalfComposition(m.listLeft.View(), m.listRight.View(), constants.BodyHeight()))

	}
}
