package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/config"
	"github.com/speed1405/arch-installer-go/ui"
)

// PartFormModel is the partition configuration screen.
type PartFormModel struct {
	bootSize textinput.Model
	rootSize textinput.Model
	homeSize textinput.Model
	focused  int
}

// NewPartForm creates a new PartFormModel with default values.
func NewPartForm() PartFormModel {
	boot := textinput.New()
	boot.Placeholder = "512M"
	boot.CharLimit = 16
	boot.Focus()

	root := textinput.New()
	root.Placeholder = "50G"
	root.CharLimit = 16

	home := textinput.New()
	home.Placeholder = "rest"
	home.CharLimit = 16

	return PartFormModel{
		bootSize: boot,
		rootSize: root,
		homeSize: home,
	}
}

// Update handles keyboard input for the partition form.
func (m PartFormModel) Update(msg tea.Msg) (PartFormModel, bool, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.focused = (m.focused + 1) % 3
			m.refocus()
		case "shift+tab", "up":
			m.focused = (m.focused + 2) % 3
			m.refocus()
		case "enter":
			if m.focused == 2 {
				return m, true, nil
			}
			m.focused = (m.focused + 1) % 3
			m.refocus()
		}
	}

	var cmd tea.Cmd
	switch m.focused {
	case 0:
		m.bootSize, cmd = m.bootSize.Update(msg)
	case 1:
		m.rootSize, cmd = m.rootSize.Update(msg)
	case 2:
		m.homeSize, cmd = m.homeSize.Update(msg)
	}
	return m, false, cmd
}

func (m *PartFormModel) refocus() {
	m.bootSize.Blur()
	m.rootSize.Blur()
	m.homeSize.Blur()
	switch m.focused {
	case 0:
		m.bootSize.Focus()
	case 1:
		m.rootSize.Focus()
	case 2:
		m.homeSize.Focus()
	}
}

// Result returns the configured partitions.
func (m PartFormModel) Result() []config.Partition {
	bootSz := m.bootSize.Value()
	if bootSz == "" {
		bootSz = "512M"
	}
	rootSz := m.rootSize.Value()
	if rootSz == "" {
		rootSz = "50G"
	}
	return []config.Partition{
		{MountPoint: "/boot", Filesystem: "fat32", Size: bootSz},
		{MountPoint: "/", Filesystem: "ext4", Size: rootSz},
		{MountPoint: "/home", Filesystem: "ext4", Size: "rest"},
	}
}

// View renders the partition form.
func (m PartFormModel) View() string {
	title := ui.StyleTitle.Render("Partition Configuration")

	fields := "\n" +
		field("/boot (EFI)", m.bootSize.View(), m.focused == 0) + "\n" +
		field("/     (root)", m.rootSize.View(), m.focused == 1) + "\n" +
		field("/home", m.homeSize.View(), m.focused == 2) + "\n\n" +
		ui.StyleMuted.Render("Tab/↓ next field • Shift+Tab/↑ prev • Enter on last field to confirm")

	return ui.StyleBorder.Render(title + fields)
}

func field(label, input string, focused bool) string {
	lbl := ui.StyleMuted.Render(label + ": ")
	if focused {
		lbl = ui.StyleHighlight.Render(label + ": ")
	}
	return lbl + input
}
