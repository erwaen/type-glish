package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Dark fantasy color palette
	ColorPrimary    = lipgloss.Color("#C9A227") // Gold/amber
	ColorSecondary  = lipgloss.Color("#7A7A7A") // Steel grey
	ColorTertiary   = lipgloss.Color("#A0A0A0") // Light grey
	ColorBackground = lipgloss.Color("#1A1A1A") // Near black
	ColorText       = lipgloss.Color("#E0E0E0") // Off-white
	ColorSubtext    = lipgloss.Color("#666666") // Dim grey
	ColorSuccess    = lipgloss.Color("#4A9F4A") // Muted green
	ColorError      = lipgloss.Color("#A03030") // Dark red
	ColorWarning    = lipgloss.Color("#B87333") // Copper/bronze

	// Layout
	WidthMain = 70

	// Styles
	StyleBase = lipgloss.NewStyle().
			Foreground(ColorText).
			Padding(1, 2)

	StyleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 3).
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
			Foreground(ColorPrimary).
			Bold(true).
			SetString("> ")

	StyleUnselected = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			SetString("  ")

	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Italic(true).
			MarginTop(1)

	StyleDamageDealt = lipgloss.NewStyle().
				Foreground(ColorSuccess).
				Bold(true)

	StyleDamageReceived = lipgloss.NewStyle().
				Foreground(ColorError).
				Bold(true)

	StyleEnemyName = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	StyleLocation = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Italic(true)
)

// Box returns a styled box with a title and content
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

// CenteredView centers a box in the terminal
func CenteredView(title string, content string, isActive bool, width, height int) string {
	box := Box(title, content, isActive)

	// Center horizontally and vertically
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		box,
		lipgloss.WithWhitespaceBackground(ColorBackground),
	)
}

func RenderMenuItem(text string, selected bool) string {
	style := StyleUnselected
	if selected {
		style = StyleSelected
	}
	return style.Render(text)
}

// RenderHPBar renders an HP bar with filled and empty segments
func RenderHPBar(current, max int, label string, barWidth int) string {
	if max <= 0 {
		max = 1
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	percent := float64(current) / float64(max)
	filled := int(percent * float64(barWidth))
	empty := barWidth - filled

	// Choose color based on HP percentage
	barColor := ColorSuccess
	if percent <= 0.25 {
		barColor = ColorError
	} else if percent <= 0.5 {
		barColor = ColorWarning
	}

	filledStyle := lipgloss.NewStyle().Foreground(barColor)
	damagedStyle := lipgloss.NewStyle().Foreground(barColor).Faint(true)
	emptyStyle := lipgloss.NewStyle().Foreground(ColorSubtext)

	// Use different characters for visual feedback
	bar := filledStyle.Render(strings.Repeat("█", filled))
	if empty > 0 {
		// Show first few empty as "damaged" (▓) and rest as truly empty (░)
		damaged := min(empty, 2)
		bar += damagedStyle.Render(strings.Repeat("▓", damaged))
		bar += emptyStyle.Render(strings.Repeat("░", empty-damaged))
	}

	percentStyle := lipgloss.NewStyle().Bold(true).Foreground(barColor)
	return fmt.Sprintf("[%s]: %s %s", label, bar, percentStyle.Render(fmt.Sprintf("(%d%%)", int(percent*100))))
}

// RenderStatusBar renders a compact status bar with HP, Gold, and XP
func RenderStatusBar(hp, maxHP, gold, xp int) string {
	// HP portion
	percent := float64(hp) / float64(maxHP)
	hpColor := ColorSuccess
	if percent <= 0.25 {
		hpColor = ColorError
	} else if percent <= 0.5 {
		hpColor = ColorWarning
	}

	hpStyle := lipgloss.NewStyle().Foreground(hpColor).Bold(true)
	goldStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	xpStyle := lipgloss.NewStyle().Foreground(ColorTertiary)
	labelStyle := lipgloss.NewStyle().Foreground(ColorSubtext)

	// Build compact HP bar (10 chars)
	barWidth := 10
	filled := int(percent * float64(barWidth))
	empty := barWidth - filled

	bar := hpStyle.Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(ColorSubtext).Render(strings.Repeat("░", empty))

	return fmt.Sprintf("%s %s %s  %s %s  %s %s",
		labelStyle.Render("HP:"), bar, hpStyle.Render(fmt.Sprintf("%d/%d", hp, maxHP)),
		labelStyle.Render("Gold:"), goldStyle.Render(fmt.Sprintf("%d", gold)),
		labelStyle.Render("XP:"), xpStyle.Render(fmt.Sprintf("%d", xp)),
	)
}

// RenderCombatHeader renders the location and enemy info header
func RenderCombatHeader(location, enemyName string) string {
	loc := StyleLocation.Render(location)
	enemy := StyleEnemyName.Render(enemyName)
	return fmt.Sprintf("LOCATION: %s    ENEMY: %s", loc, enemy)
}
