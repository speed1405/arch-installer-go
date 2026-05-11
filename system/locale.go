package system

import (
	"context"
	"fmt"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/messages"
)

// validKeymap matches simple keyboard layout identifiers like "us", "de", "fr".
var validKeymap = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

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
// The layout is validated against a safe pattern before use.
func ConfigureKeyboard(layout string) tea.Cmd {
	if layout == "" {
		layout = "us"
	}
	return func() tea.Msg {
		if !validKeymap.MatchString(layout) {
			return messages.CmdErrorMsg{
				Err:  fmt.Errorf("invalid keyboard layout: %q", layout),
				Kind: messages.ErrorCritical,
			}
		}
		return ExecuteChrootWithStdin(
			context.Background(), mountpoint,
			"KEYMAP="+layout+"\n",
			"sh", "-c", "cat > /etc/vconsole.conf",
		)()
	}
}

// SyncClock runs hwclock to synchronize hardware clock.
func SyncClock() tea.Cmd {
	return func() tea.Msg {
		result := ExecuteChroot(context.Background(), mountpoint, "hwclock", "--systohc")()
		if errMsg, ok := result.(messages.CmdErrorMsg); ok {
			// non-fatal: log and continue
			return messages.CmdOutputMsg{Line: "hwclock: " + errMsg.Err.Error()}
		}
		return result
	}
}
