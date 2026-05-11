package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/speed1405/arch-installer-go/ui"
)

// DiskDisplay holds display information about a block device.
type DiskDisplay struct {
	Name  string
	Path  string
	Size  string
	Model string
}

// DiskListModel is the disk selection screen component.
type DiskListModel struct {
	disks  []DiskDisplay
	cursor int
	width  int
	height int
}

// NewDiskList creates an empty DiskListModel.
func NewDiskList() DiskListModel {
	return DiskListModel{}
}

// SetDisks populates the list with detected disks.
func (m *DiskListModel) SetDisks(disks []DiskDisplay) {
	m.disks = disks
}

// SetSize updates the available terminal dimensions.
func (m *DiskListModel) SetSize(w, h int) {
	m.width, m.height = w, h
}

// Update handles keyboard navigation and selection.
// Returns the selected disk path when confirmed, or "" if still navigating.
func (m DiskListModel) Update(msg tea.Msg) (DiskListModel, string, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.disks)-1 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.disks) > 0 {
				return m, m.disks[m.cursor].Path, nil
			}
		}
	}
	return m, "", nil
}

// View renders the disk list.
func (m DiskListModel) View() string {
	title := ui.StyleTitle.Render("Select Installation Disk")

	if len(m.disks) == 0 {
		return ui.StyleBorder.Render(title + "\n\n" + ui.StyleMuted.Render("Detecting disks…"))
	}

	var items string
	for i, d := range m.disks {
		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			padRight(d.Path, 12),
			padRight(d.Size, 8),
			d.Model,
		)
		if i == m.cursor {
			items += ui.StyleHighlight.Render("► " + line) + "\n"
		} else {
			items += "  " + line + "\n"
		}
	}

	return ui.StyleBorder.Render(
		title + "\n\n" +
			ui.StyleMuted.Render("NAME        SIZE    MODEL") + "\n" +
			items + "\n" +
			ui.StyleMuted.Render("↑/↓ navigate • Enter select"),
	)
}

func padRight(s string, n int) string {
	for len(s) < n {
		s += " "
	}
	return s
}
