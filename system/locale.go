package system

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/messages"
)

// AvailableLocales returns a set of commonly used locales.
func AvailableLocales() []string {
	return []string{
		"en_US.UTF-8",
		"en_GB.UTF-8",
		"de_DE.UTF-8",
		"fr_FR.UTF-8",
		"es_ES.UTF-8",
		"pt_BR.UTF-8",
		"ja_JP.UTF-8",
		"zh_CN.UTF-8",
	}
}

// AvailableTimezones returns a short list of common timezones.
func AvailableTimezones() []string {
	return []string{
		"UTC",
		"America/New_York",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"Europe/London",
		"Europe/Berlin",
		"Europe/Paris",
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Australia/Sydney",
	}
}

// ConfigureKeyboard sets the keyboard layout inside chroot.
func ConfigureKeyboard(layout string) tea.Cmd {
	if layout == "" {
		layout = "us"
	}
	return func() tea.Msg {
		script := fmt.Sprintf("echo 'KEYMAP=%s' > /etc/vconsole.conf", layout)
		return ExecuteChroot(context.Background(), mountpoint, "sh", "-c", script)()
	}
}

// SyncClock runs hwclock to synchronize hardware clock.
func SyncClock() tea.Cmd {
	return func() tea.Msg {
		result := ExecuteChroot(context.Background(), mountpoint, "hwclock", "--systohc")()
		switch r := result.(type) {
		case messages.CmdErrorMsg:
			// non-fatal
			return messages.CmdOutputMsg{Line: "hwclock: " + r.Err.Error()}
		}
		return result
	}
}
