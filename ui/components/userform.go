package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/ui"
)

// UserFormResult holds the collected user configuration values.
type UserFormResult struct {
	Hostname     string
	Username     string
	Password     string
	RootPassword string
}

// UserFormModel collects hostname, username, and passwords.
type UserFormModel struct {
	hostname     textinput.Model
	username     textinput.Model
	password     textinput.Model
	rootPassword textinput.Model
	focused      int
}

// NewUserForm creates a new UserFormModel.
func NewUserForm() UserFormModel {
	host := textinput.New()
	host.Placeholder = "archlinux"
	host.CharLimit = 64
	host.Focus()

	user := textinput.New()
	user.Placeholder = "username"
	user.CharLimit = 32

	pass := textinput.New()
	pass.Placeholder = "password"
	pass.EchoMode = textinput.EchoPassword
	pass.CharLimit = 128

	root := textinput.New()
	root.Placeholder = "root password"
	root.EchoMode = textinput.EchoPassword
	root.CharLimit = 128

	return UserFormModel{
		hostname:     host,
		username:     user,
		password:     pass,
		rootPassword: root,
	}
}

// Update handles keyboard input.
func (m UserFormModel) Update(msg tea.Msg) (UserFormModel, bool, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.focused = (m.focused + 1) % 4
			m.refocusUser()
		case "shift+tab", "up":
			m.focused = (m.focused + 3) % 4
			m.refocusUser()
		case "enter":
			if m.focused == 3 {
				return m, true, nil
			}
			m.focused = (m.focused + 1) % 4
			m.refocusUser()
		}
	}

	var cmd tea.Cmd
	switch m.focused {
	case 0:
		m.hostname, cmd = m.hostname.Update(msg)
	case 1:
		m.username, cmd = m.username.Update(msg)
	case 2:
		m.password, cmd = m.password.Update(msg)
	case 3:
		m.rootPassword, cmd = m.rootPassword.Update(msg)
	}
	return m, false, cmd
}

func (m *UserFormModel) refocusUser() {
	m.hostname.Blur()
	m.username.Blur()
	m.password.Blur()
	m.rootPassword.Blur()
	switch m.focused {
	case 0:
		m.hostname.Focus()
	case 1:
		m.username.Focus()
	case 2:
		m.password.Focus()
	case 3:
		m.rootPassword.Focus()
	}
}

// Result returns the form values.
func (m UserFormModel) Result() UserFormResult {
	return UserFormResult{
		Hostname:     m.hostname.Value(),
		Username:     m.username.Value(),
		Password:     m.password.Value(),
		RootPassword: m.rootPassword.Value(),
	}
}

// View renders the user config form.
func (m UserFormModel) View() string {
	title := ui.StyleTitle.Render("User Configuration")
	fields := "\n" +
		userField("Hostname", m.hostname.View(), m.focused == 0) + "\n" +
		userField("Username", m.username.View(), m.focused == 1) + "\n" +
		userField("Password", m.password.View(), m.focused == 2) + "\n" +
		userField("Root Password", m.rootPassword.View(), m.focused == 3) + "\n\n" +
		ui.StyleMuted.Render("Tab/↓ next • Shift+Tab/↑ prev • Enter on last field to confirm")
	return ui.StyleBorder.Render(title + fields)
}

func userField(label, input string, focused bool) string {
	lbl := ui.StyleMuted.Render(label + ": ")
	if focused {
		lbl = ui.StyleHighlight.Render(label + ": ")
	}
	return lbl + input
}
