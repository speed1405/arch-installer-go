package system

import (
"context"
"fmt"
"regexp"

tea "github.com/charmbracelet/bubbletea"
"github.com/speed1405/arch-installer-go/messages"
)

const mountpoint = "/mnt"

// validLocale matches strings like "en_US.UTF-8".
var validLocale = regexp.MustCompile(`^[a-zA-Z]{2,3}_[A-Z]{2,3}(\.[A-Z0-9\-]+)?(@[a-zA-Z]+)?$`)

// validHostname matches RFC 1123 hostnames.
var validHostname = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)

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
// The locale string is validated before use to prevent shell injection.
func SetLocale(locale string) tea.Cmd {
if locale == "" {
locale = "en_US.UTF-8"
}
return func() tea.Msg {
if !validLocale.MatchString(locale) {
return messages.CmdErrorMsg{
Err:  fmt.Errorf("invalid locale: %q", locale),
Kind: messages.ErrorCritical,
}
}
return ExecuteChrootWithStdin(
context.Background(), mountpoint,
locale+" UTF-8\n",
"sh", "-c",
"cat > /etc/locale.gen && locale-gen && echo 'LANG="+locale+"' > /etc/locale.conf",
)()
}
}

// SetHostname writes /etc/hostname and /etc/hosts.
// The hostname is validated against RFC 1123 rules before use.
func SetHostname(hostname string) tea.Cmd {
if hostname == "" {
hostname = "archlinux"
}
return func() tea.Msg {
if !validHostname.MatchString(hostname) {
return messages.CmdErrorMsg{
Err:  fmt.Errorf("invalid hostname: %q", hostname),
Kind: messages.ErrorCritical,
}
}
hostsContent := fmt.Sprintf(
"127.0.0.1 localhost\n::1 localhost\n127.0.1.1 %s.localdomain %s\n",
hostname, hostname,
)
steps := []tea.Cmd{
ExecuteChrootWithStdin(context.Background(), mountpoint, hostname+"\n",
"sh", "-c", "cat > /etc/hostname"),
ExecuteChrootWithStdin(context.Background(), mountpoint, hostsContent,
"sh", "-c", "cat > /etc/hosts"),
}
for _, step := range steps {
result := step()
if errMsg, ok := result.(messages.CmdErrorMsg); ok {
return errMsg
}
}
return messages.CmdDoneMsg{Output: "hostname configured"}
}
}
