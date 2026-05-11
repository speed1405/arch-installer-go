# arch-installer-go

A modular, TUI-driven Arch Linux installer written in Go, powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **18-step state machine** from welcome → disk detection → partitioning → pacstrap → chroot config → bootloader → complete
- **Live log view** with streaming command output piped back to the TUI
- **Responsive UI** — long-running commands (pacstrap, grub-install, etc.) run in the background via `tea.Cmd`; the TUI never blocks
- **Error classification** — Critical vs. Retryable errors with dedicated recovery screens
- **Back navigation** blocked past destructive operations (partitioning, pacstrap)
- **UEFI / BIOS** auto-detection at startup
- **GRUB** and **systemd-boot** supported

## Project Structure

```
arch-installer-go/
├── main.go                     # tea.Program entry point
├── config/                     # InstallConfig schema + sensible defaults
├── messages/                   # Shared Bubble Tea message types
├── installer/                  # 18-step state machine, root Model/Update/View, error classification
├── ui/
│   ├── styles.go               # Lipgloss colour palette + base styles
│   ├── layout.go               # Header, footer, and screen renderers
│   └── components/             # DiskList, PartForm, UserForm, ConfirmDialog, ProgressBar
└── system/                     # OS interaction: disk, pacstrap, chroot, locale, bootloader, user
```

## Building

```bash
go build -o arch-installer .
```

Requires Go 1.22+.

## Running

**Must be run as root (or with sudo) on a live Arch ISO.**

```bash
sudo ./arch-installer
```

## Dependencies

| Library | Purpose |
|---|---|
| `github.com/charmbracelet/bubbletea` | The Elm Architecture TUI framework |
| `github.com/charmbracelet/bubbles` | Text inputs, progress bar components |
| `github.com/charmbracelet/lipgloss` | Terminal layout and styling |

## Architecture Notes

### State Machine

Steps are defined in `installer/steps.go` as `Step` (`int`) constants.  
`nextStep()` provides the linear transition table.  
`destructiveSteps` marks steps that block backward navigation once entered.

### Concurrency Model

All system commands run via `system.Execute()` which returns a `tea.Cmd`.  
Output lines stream back as `messages.CmdOutputMsg`; completion signals `messages.CmdDoneMsg`.  
The progress screen (`ui/components/progress.go`) consumes these messages in real-time.

### Security

- Passwords are **never** written to disk or shell command lines — passed via `stdin` pipe to `chpasswd`
- User-supplied strings (locale, hostname, keyboard layout) are **validated with regexp** before use in shell commands
- `sudo --non-interactive` is used to prevent interactive prompts from blocking the TUI
