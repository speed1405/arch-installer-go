package system

import (
	"net"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/messages"
)

// DetectEnvironment checks for UEFI support and internet connectivity,
// then returns an EnvDetectedMsg.
func DetectEnvironment() tea.Cmd {
	return func() tea.Msg {
		isUEFI := checkUEFI()
		internet := checkInternet()
		return messages.EnvDetectedMsg{
			IsUEFI:   isUEFI,
			Internet: internet,
		}
	}
}

// checkUEFI returns true if /sys/firmware/efi exists (indicating UEFI boot).
func checkUEFI() bool {
	_, err := os.Stat("/sys/firmware/efi")
	return err == nil
}

// checkInternet attempts a TCP dial to a well-known address to verify connectivity.
func checkInternet() bool {
	conn, err := net.DialTimeout("tcp", "archlinux.org:443", 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
