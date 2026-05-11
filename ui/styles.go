package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	ColorPrimary   = lipgloss.Color("#7C3AED") // violet
	ColorSecondary = lipgloss.Color("#06B6D4") // cyan
	ColorSuccess   = lipgloss.Color("#22C55E") // green
	ColorWarning   = lipgloss.Color("#F59E0B") // amber
	ColorError     = lipgloss.Color("#EF4444") // red
	ColorMuted     = lipgloss.Color("#6B7280") // gray
	ColorWhite     = lipgloss.Color("#F9FAFB")
	ColorBg        = lipgloss.Color("#111827")
)

// Base styles
var (
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Italic(true)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning)

	StyleMuted = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	StyleHeader = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 2).
			Width(80)

	StyleFooter = lipgloss.NewStyle().
			Foreground(ColorMuted).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(ColorMuted).
			Padding(0, 2)

	StyleBody = lipgloss.NewStyle().
			Padding(1, 2)

	StyleHighlight = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	StyleLog = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 1)
)
