package system

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/config"
	"github.com/speed1405/arch-installer-go/messages"
)

// InstallBootloader selects and runs the appropriate bootloader installer.
func InstallBootloader(cfg config.InstallConfig) tea.Cmd {
	if cfg.Bootloader == "systemd-boot" {
		return installSystemdBoot(cfg)
	}
	return installGrub(cfg)
}

// installGrub installs GRUB for both BIOS and UEFI targets.
func installGrub(cfg config.InstallConfig) tea.Cmd {
	return func() tea.Msg {
		var cmds []tea.Cmd
		if cfg.IsUEFI {
			cmds = []tea.Cmd{
				ExecuteChroot(context.Background(), mountpoint,
					"pacman", "-S", "--noconfirm", "grub", "efibootmgr"),
				ExecuteChroot(context.Background(), mountpoint,
					"grub-install", "--target=x86_64-efi",
					"--efi-directory=/boot", "--bootloader-id=GRUB"),
				ExecuteChroot(context.Background(), mountpoint,
					"grub-mkconfig", "-o", "/boot/grub/grub.cfg"),
			}
		} else {
			cmds = []tea.Cmd{
				ExecuteChroot(context.Background(), mountpoint,
					"pacman", "-S", "--noconfirm", "grub"),
				ExecuteChroot(context.Background(), mountpoint,
					"grub-install", "--target=i386-pc", cfg.Disk),
				ExecuteChroot(context.Background(), mountpoint,
					"grub-mkconfig", "-o", "/boot/grub/grub.cfg"),
			}
		}

		for _, c := range cmds {
			result := c()
			switch r := result.(type) {
			case messages.CmdErrorMsg:
				return r
			}
		}
		return messages.CmdDoneMsg{Output: "GRUB installed"}
	}
}

// installSystemdBoot installs systemd-boot (UEFI only).
func installSystemdBoot(cfg config.InstallConfig) tea.Cmd {
	return func() tea.Msg {
		if !cfg.IsUEFI {
			return messages.CmdErrorMsg{
				Err:  fmt.Errorf("systemd-boot requires UEFI; detected BIOS system"),
				Kind: messages.ErrorCritical,
			}
		}
		cmds := []tea.Cmd{
			ExecuteChroot(context.Background(), mountpoint, "bootctl", "install"),
		}
		for _, c := range cmds {
			result := c()
			switch r := result.(type) {
			case messages.CmdErrorMsg:
				return r
			}
		}
		return messages.CmdDoneMsg{Output: "systemd-boot installed"}
	}
}
