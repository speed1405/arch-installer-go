package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/speed1405/arch-installer-go/messages"
	"github.com/speed1405/arch-installer-go/ui"
)

// ProgressModel shows a progress bar and a live scrolling log view.
type ProgressModel struct {
	bar      progress.Model
	lines    []string
	maxLines int
	percent  float64
	width    int
	height   int
}

// NewProgress creates a new ProgressModel.
func NewProgress() ProgressModel {
	b := progress.New(progress.WithDefaultGradient())
	return ProgressModel{
		bar:      b,
		maxLines: 10,
	}
}

// SetSize updates terminal dimensions.
func (m *ProgressModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.bar.Width = w - 8
}

// Update handles CmdOutputMsg to append log lines and advance the bar.
func (m ProgressModel) Update(msg tea.Msg) (ProgressModel, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.CmdOutputMsg:
		m.lines = append(m.lines, msg.Line)
		// Advance bar by a small increment per line received.
		if m.percent < 0.95 {
			m.percent += 0.01
		}
		return m, nil
	}
	return m, nil
}

// View renders the progress bar and log.
func (m ProgressModel) View() string {
	title := ui.StyleTitle.Render("Working…")
	bar := m.bar.ViewAs(m.percent)

	var logSb strings.Builder
	start := 0
	if len(m.lines) > m.maxLines {
		start = len(m.lines) - m.maxLines
	}
	for _, l := range m.lines[start:] {
		logSb.WriteString(ui.StyleLog.Render(l) + "\n")
	}

	return ui.StyleBorder.Render(
		title + "\n\n" + bar + "\n\n" + logSb.String(),
	)
}
