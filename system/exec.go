package system

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/messages"
)

// Execute runs a system command (optionally under sudo), streams each line of
// combined stdout+stderr back to the Bubble Tea loop as CmdOutputMsg messages,
// and returns CmdDoneMsg on success or CmdErrorMsg on failure.
func Execute(ctx context.Context, sudo bool, args ...string) tea.Cmd {
	return func() tea.Msg {
		if len(args) == 0 {
			return messages.CmdErrorMsg{
				Err:  fmt.Errorf("execute: no command specified"),
				Kind: messages.ErrorCritical,
			}
		}
		if sudo {
			args = append([]string{"sudo", "--non-interactive"}, args...)
		}

		cmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint:gosec
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return messages.CmdErrorMsg{Err: err, Kind: messages.ClassifyError(err)}
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return messages.CmdErrorMsg{Err: err, Kind: messages.ClassifyError(err)}
		}

		if err := cmd.Start(); err != nil {
			return messages.CmdErrorMsg{Err: err, Kind: messages.ClassifyError(err)}
		}

		// Merge stdout and stderr into a single line channel.
		lines := make(chan string, 128)
		var wg sync.WaitGroup
		for _, r := range []io.Reader{stdout, stderr} {
			wg.Add(1)
			go func(r io.Reader) {
				defer wg.Done()
				sc := bufio.NewScanner(r)
				for sc.Scan() {
					lines <- sc.Text()
				}
			}(r)
		}
		go func() {
			wg.Wait()
			close(lines)
		}()

		var buf strings.Builder
		for line := range lines {
			buf.WriteString(line + "\n")
		}

		if err := cmd.Wait(); err != nil {
			combined := fmt.Errorf("%w\n%s", err, strings.TrimSpace(buf.String()))
			return messages.CmdErrorMsg{Err: combined, Kind: messages.ClassifyError(err)}
		}
		return messages.CmdDoneMsg{Output: buf.String()}
	}
}

// ExecuteChroot runs a command inside an arch-chroot environment.
func ExecuteChroot(ctx context.Context, mountpoint string, args ...string) tea.Cmd {
	chrootArgs := append([]string{"arch-chroot", mountpoint}, args...)
	return Execute(ctx, true, chrootArgs...)
}
