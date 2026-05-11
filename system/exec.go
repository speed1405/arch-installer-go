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

// Execute runs a system command (optionally under sudo) and streams each line
// of combined stdout+stderr back to the Bubble Tea loop as CmdOutputMsg
// messages. Each CmdOutputMsg carries a Next tea.Cmd that the Update function
// must dispatch to read the following line. When all output is consumed the
// final message is CmdDoneMsg (or CmdErrorMsg on failure).
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

		// Merge stdout and stderr into a single buffered line channel.
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

		// done receives the process exit error once all output has been read.
		done := make(chan error, 1)
		go func() {
			wg.Wait()
			close(lines)
			done <- cmd.Wait()
		}()

		// Return the first polling command to the Bubble Tea loop.
		return readNextLine(lines, done)()
	}
}

// readNextLine returns a tea.Cmd that reads one line from the channel and
// emits a CmdOutputMsg (with the next polling command attached) or a
// CmdDoneMsg / CmdErrorMsg when the stream is exhausted.
func readNextLine(lines <-chan string, done <-chan error) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-lines
		if !ok {
			// Channel closed — wait for process exit.
			if err := <-done; err != nil {
				return messages.CmdErrorMsg{
					Err:  err,
					Kind: messages.ClassifyError(err),
				}
			}
			return messages.CmdDoneMsg{}
		}
		return messages.CmdOutputMsg{
			Line: line,
			Next: readNextLine(lines, done),
		}
	}
}

// ExecuteChroot runs a command inside an arch-chroot environment.
func ExecuteChroot(ctx context.Context, mountpoint string, args ...string) tea.Cmd {
	chrootArgs := append([]string{"arch-chroot", mountpoint}, args...)
	return Execute(ctx, true, chrootArgs...)
}

// ExecuteChrootWithStdin runs a command inside arch-chroot with the given string
// piped to the process's stdin. This avoids passing secrets through shell arguments.
func ExecuteChrootWithStdin(ctx context.Context, mountpoint, stdinData string, args ...string) tea.Cmd {
	return func() tea.Msg {
		chrootArgs := append(
			[]string{"sudo", "--non-interactive", "arch-chroot", mountpoint},
			args...,
		)
		cmd := exec.CommandContext(ctx, chrootArgs[0], chrootArgs[1:]...) //nolint:gosec
		cmd.Stdin = strings.NewReader(stdinData)
		out, err := cmd.CombinedOutput()
		if err != nil {
			combined := fmt.Errorf("%w\n%s", err, strings.TrimSpace(string(out)))
			return messages.CmdErrorMsg{Err: combined, Kind: messages.ClassifyError(err)}
		}
		return messages.CmdDoneMsg{Output: string(out)}
	}
}
