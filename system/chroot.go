package system

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

const mountpoint = "/mnt"

// GenerateFstab runs genfstab and writes /mnt/etc/fstab.
func GenerateFstab() tea.Cmd {
	return Execute(context.Background(), true,
		"sh", "-c", "genfstab -U "+mountpoint+" >> "+mountpoint+"/etc/fstab",
	)
}

// SetTimezone configures the system timezone inside the chroot.
func SetTimezone(tz string) tea.Cmd {
	if tz == "" {
		tz = "UTC"
	}
	return ExecuteChroot(context.Background(), mountpoint,
		"ln", "-sf", "/usr/share/zoneinfo/"+tz, "/etc/localtime",
	)
}

// SetLocale configures the locale inside the chroot.
func SetLocale(locale string) tea.Cmd {
	if locale == "" {
		locale = "en_US.UTF-8"
	}
	return ExecuteChroot(context.Background(), mountpoint,
		"sh", "-c",
		"echo '"+locale+" UTF-8' > /etc/locale.gen && locale-gen && echo 'LANG="+locale+"' > /etc/locale.conf",
	)
}

// SetHostname writes /etc/hostname and /etc/hosts.
func SetHostname(hostname string) tea.Cmd {
	if hostname == "" {
		hostname = "archlinux"
	}
	script := "echo '" + hostname + "' > /etc/hostname && " +
		"echo '127.0.0.1 localhost' > /etc/hosts && " +
		"echo '::1 localhost' >> /etc/hosts && " +
		"echo '127.0.1.1 " + hostname + ".localdomain " + hostname + "' >> /etc/hosts"
	return ExecuteChroot(context.Background(), mountpoint, "sh", "-c", script)
}
