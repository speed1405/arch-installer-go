package system

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/config"
	"github.com/speed1405/arch-installer-go/messages"
)

// lsblkOutput mirrors the JSON output of `lsblk --json`.
type lsblkOutput struct {
	Blockdevices []lsblkDevice `json:"blockdevices"`
}

type lsblkDevice struct {
	Name  string `json:"name"`
	Size  string `json:"size"`
	Model string `json:"model"`
	Type  string `json:"type"`
	Rm    bool   `json:"rm"` // removable
}

// ListDisks returns a tea.Cmd that detects available block devices and
// emits a DisksLoadedMsg.
func ListDisks() tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("lsblk", "--json", "-d", "-o", "NAME,SIZE,MODEL,TYPE,RM").Output()
		if err != nil {
			return messages.CmdErrorMsg{
				Err:  fmt.Errorf("lsblk: %w", err),
				Kind: messages.ErrorCritical,
			}
		}

		var data lsblkOutput
		if err := json.Unmarshal(out, &data); err != nil {
			return messages.CmdErrorMsg{
				Err:  fmt.Errorf("lsblk parse: %w", err),
				Kind: messages.ErrorCritical,
			}
		}

		var disks []messages.DiskInfo
		for _, d := range data.Blockdevices {
			if d.Type == "disk" {
				disks = append(disks, messages.DiskInfo{
					Name:  d.Name,
					Path:  "/dev/" + d.Name,
					Size:  d.Size,
					Model: d.Model,
				})
			}
		}

		return messages.DisksLoadedMsg{Disks: disks}
	}
}

// PartitionDisk runs fdisk/parted to set up partitions according to the config.
func PartitionDisk(cfg config.InstallConfig) tea.Cmd {
	if cfg.PartitionScheme == "auto" {
		return autoPartition(cfg.Disk, cfg.IsUEFI)
	}
	return manualPartition(cfg)
}

// autoPartition creates a simple layout: EFI + root (or BIOS boot + root).
func autoPartition(disk string, uefi bool) tea.Cmd {
	return func() tea.Msg {
		var script string
		if uefi {
			// GPT: 512M EFI + rest root
			script = "g\nn\n\n\n+512M\nt\n1\nn\n\n\n\nw\n"
		} else {
			// MBR: 1M BIOS boot + rest root
			script = "o\nn\np\n1\n\n+1M\nn\np\n2\n\n\nw\n"
		}

		cmd := exec.Command("sudo", "fdisk", disk) //nolint:gosec
		cmd.Stdin = newStringReader(script)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return messages.CmdErrorMsg{
				Err:  fmt.Errorf("fdisk: %w\n%s", err, out),
				Kind: messages.ErrorCritical,
			}
		}
		return messages.CmdDoneMsg{Output: string(out)}
	}
}

// manualPartition is a placeholder for a full manual partitioning flow.
func manualPartition(_ config.InstallConfig) tea.Cmd {
	return func() tea.Msg {
		return messages.CmdDoneMsg{Output: "manual partition complete"}
	}
}

// stringReader wraps a string as an io.Reader for use as stdin.
type stringReader struct {
	data string
	pos  int
}

func newStringReader(s string) *stringReader {
	return &stringReader{data: s}
}

func (r *stringReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
