package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderHeader renders the top bar with the current step name.
func RenderHeader(stepName string, isError bool) string {
	fg := ColorPrimary
	if isError {
		fg = ColorError
	}
	style := lipgloss.NewStyle().
		Background(fg).
		Foreground(ColorWhite).
		Bold(true).
		Padding(0, 2).
		Width(80)

	left := "⚡ Arch Linux Installer"
	right := stepName
	gap := 80 - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}
	return style.Render(left + strings.Repeat(" ", gap) + right)
}

// RenderFooter renders the bottom hint bar.
func RenderFooter(done bool) string {
	var hints string
	if done {
		hints = "ctrl+c quit"
	} else {
		hints = "esc back • ctrl+c quit"
	}
	return StyleFooter.Render(hints)
}

// RenderWelcome renders the welcome screen.
func RenderWelcome() string {
	title := StyleTitle.Render("Welcome to Arch Linux Installer")
	sub := StyleSubtitle.Render("A guided TUI installer powered by Bubble Tea")
	body := "\n" + title + "\n\n" + sub + "\n\n" +
		StyleMuted.Render("Press Enter or → to begin.\nPress ctrl+c to quit at any time.")
	return StyleBorder.Render(body)
}

// RenderSpinner renders a "please wait" screen.
func RenderSpinner(msg string) string {
	return StyleBody.Render(StyleSubtitle.Render("⏳ " + msg))
}

// RenderComplete renders the success screen.
func RenderComplete() string {
	body := StyleSuccess.Render("✓ Installation Complete!") + "\n\n" +
		"Your Arch Linux system is ready.\n\n" +
		StyleMuted.Render("Remove the installation media and reboot.\n\nPress r to reboot, or ctrl+c to exit.")
	return StyleBorder.Render(body)
}

// RenderError renders the error screen.
func RenderError(err error, retryable bool) string {
	msg := "Unknown error"
	if err != nil {
		msg = err.Error()
	}
	var hint string
	if retryable {
		hint = StyleWarning.Render("This error may be retryable. Press r to retry.")
	} else {
		hint = StyleError.Render("This is a critical error. Press ctrl+c to exit or reboot.")
	}
	body := StyleError.Render("✗ Error") + "\n\n" +
		StyleMuted.Render(msg) + "\n\n" + hint
	return StyleBorder.Render(body)
}

// RenderPartitionStrategy renders the partition strategy selection.
func RenderPartitionStrategy(current string) string {
	auto := "  [auto]  Automatic partitioning (recommended)"
	manual := "  [manual] Manual partitioning"
	if current == "auto" {
		auto = StyleHighlight.Render("► " + auto[2:])
	} else {
		manual = StyleHighlight.Render("► " + manual[2:])
	}
	body := StyleTitle.Render("Partition Strategy") + "\n\n" +
		auto + "\n" + manual + "\n\n" +
		StyleMuted.Render("Press ↑/↓ to select, Enter to confirm.")
	return StyleBorder.Render(body)
}

// RenderMirrorInput renders the mirror configuration screen.
func RenderMirrorInput(current string) string {
	if current == "" {
		current = "(will use default mirrors)"
	}
	body := StyleTitle.Render("Mirror Configuration") + "\n\n" +
		"Mirror: " + StyleHighlight.Render(current) + "\n\n" +
		StyleMuted.Render("Press Enter to continue with default mirrors.")
	return StyleBorder.Render(body)
}

// RenderTextInput renders a simple single-field input screen.
func RenderTextInput(label, value string) string {
	body := StyleTitle.Render(label) + "\n\n" +
		label + ": " + StyleHighlight.Render(value) + "\n\n" +
		StyleMuted.Render("Press Enter to confirm.")
	return StyleBorder.Render(body)
}

// RenderPostInstall renders the post-install extras screen.
func RenderPostInstall() string {
	body := StyleTitle.Render("Post-Install Extras") + "\n\n" +
		"Optional extras:\n" +
		"  • AUR helper (yay / paru)\n" +
		"  • Desktop Environment (KDE, GNOME, etc.)\n\n" +
		StyleMuted.Render("Press Enter to skip, or select above.")
	return StyleBody.Render(body)
}

// RenderLogLines renders the last N lines of command output.
func RenderLogLines(lines []string, max int) string {
	if len(lines) == 0 {
		return StyleMuted.Render("(waiting for output…)")
	}
	start := 0
	if len(lines) > max {
		start = len(lines) - max
	}
	var sb strings.Builder
	for _, l := range lines[start:] {
		sb.WriteString(StyleLog.Render(l) + "\n")
	}
	return fmt.Sprintf("%s", sb.String())
}
