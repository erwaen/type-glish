package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimary    = lipgloss.Color("#FF79C6") // Pink
	ColorSecondary  = lipgloss.Color("#BD93F9") // Purple
	ColorTertiary   = lipgloss.Color("#8BE9FD") // Cyan
	ColorBackground = lipgloss.Color("#282A36") // Dark Grey
	ColorText       = lipgloss.Color("#F8F8F2") // White
	ColorSubtext    = lipgloss.Color("#6272A4") // Greyish Blue
	ColorSuccess    = lipgloss.Color("#50FA7B") // Green
	ColorError      = lipgloss.Color("#FF5555") // Red
	ColorWarning    = lipgloss.Color("#FFB86C") // Orange

	// Layout
	WidthMain = 80

	// Styles
	StyleBase = lipgloss.NewStyle().
			Foreground(ColorText).
			Padding(1, 2)

	StyleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2).
			Width(WidthMain)

	StyleBoxActive = StyleBox.
			BorderForeground(ColorPrimary)

	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorBackground).
			Background(ColorPrimary).
			Padding(0, 1).
			Bold(true)

	StyleSubTitle = lipgloss.NewStyle().
			Foreground(ColorTertiary).
			Italic(true)

	StyleSelected = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true).
			SetString("> ")

	StyleUnselected = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			SetString("  ")

	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Italic(true).
			MarginTop(1)
)

// Box returns a styled string with a title and content
func Box(title string, content string, isActive bool) string {
	boxStyle := StyleBox
	if isActive {
		boxStyle = StyleBoxActive
	}

	return boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			StyleTitle.Render(title),
			"\n"+content,
		),
	)
}

func RenderMenuItem(text string, selected bool) string {
	style := StyleUnselected
	if selected {
		style = StyleSelected
	}
	return style.Render(text)
}
