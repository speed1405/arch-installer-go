package system

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/messages"
)

// CreateUser sets up the root password, creates a regular user, and configures sudo.
func CreateUser(username, password, rootPassword string) tea.Cmd {
	return func() tea.Msg {
		steps := []tea.Cmd{
			setPassword(context.Background(), "root", rootPassword),
			addUser(context.Background(), username),
			setPassword(context.Background(), username, password),
			addToWheelGroup(context.Background(), username),
			configureSudoers(context.Background()),
		}
		for _, step := range steps {
			result := step()
			switch r := result.(type) {
			case messages.CmdErrorMsg:
				return r
			}
		}
		return messages.CmdDoneMsg{Output: fmt.Sprintf("User %s created", username)}
	}
}

func setPassword(ctx context.Context, user, password string) tea.Cmd {
	script := fmt.Sprintf("echo '%s:%s' | chpasswd", user, password)
	return ExecuteChroot(ctx, mountpoint, "sh", "-c", script)
}

func addUser(ctx context.Context, username string) tea.Cmd {
	return ExecuteChroot(ctx, mountpoint,
		"useradd", "-m", "-G", "audio,video,optical,storage", "-s", "/bin/bash", username)
}

func addToWheelGroup(ctx context.Context, username string) tea.Cmd {
	return ExecuteChroot(ctx, mountpoint, "usermod", "-aG", "wheel", username)
}

func configureSudoers(ctx context.Context) tea.Cmd {
	script := "sed -i 's/^# %wheel ALL=(ALL:ALL) ALL/%wheel ALL=(ALL:ALL) ALL/' /etc/sudoers"
	return ExecuteChroot(ctx, mountpoint, "sh", "-c", script)
}
