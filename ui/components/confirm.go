package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/ui"
)

// ConfirmModel is a simple Yes/No confirmation dialog.
type ConfirmModel struct {
	message string
	yes     bool
}

// NewConfirm creates a new ConfirmModel with the given warning message.
func NewConfirm(message string) ConfirmModel {
	return ConfirmModel{message: message}
}

// Update handles y/n key presses.
// Returns confirmed=true when the user presses y or Y.
func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, bool, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			m.yes = true
			return m, true, nil
		case "n", "N":
			m.yes = false
			return m, false, nil
		}
	}
	return m, false, nil
}

// View renders the confirmation dialog.
func (m ConfirmModel) View() string {
	warning := ui.StyleWarning.Render("⚠  " + m.message)
	yes := "  [ Yes ]  "
	no := "  [ No  ]  "
	yes = ui.StyleError.Render(yes)
	body := warning + "\n\n" + yes + no + "\n\n" +
		ui.StyleMuted.Render("y/Enter to confirm • n to cancel")
	return ui.StyleBorder.Render(body)
}
