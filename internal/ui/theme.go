package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	// accent color for highlights and active elements
	Accent lipgloss.Color

	// usage levels
	UsageIdle   lipgloss.Color // <5% (dim, barely running)
	UsageLow    lipgloss.Color // <20% (cool, idle)
	UsageMedium lipgloss.Color // 20-70% (warm, active)
	UsageHigh   lipgloss.Color // >70% (hot, critical)

	// frequency states
	FreqBase  lipgloss.Color // at or below base frequency
	FreqBoost lipgloss.Color // turbo boost active

	// memory colors
	MemoryUsed lipgloss.Color
	MemoryFree lipgloss.Color

	// UI element colors
	Border       lipgloss.Color
	BorderActive lipgloss.Color
	Text         lipgloss.Color
	TextDim      lipgloss.Color
	TextBold     lipgloss.Color

	// border style
	BorderStyle lipgloss.Border
}

func DefaultTheme() Theme {
	return Theme{
		Accent:       lipgloss.Color("#00D4FF"),
		UsageIdle:    lipgloss.Color("240"),     // Dim gray for idle
		UsageLow:     lipgloss.Color("#00D4FF"), // Cyan for low usage
		UsageMedium:  lipgloss.Color("#FFDD00"), // Yellow for medium usage
		UsageHigh:    lipgloss.Color("#FF0088"), // Pink/Magenta for high usage
		FreqBase:     lipgloss.Color("255"),     // Normal white for base freq
		FreqBoost:    lipgloss.Color("#00D4FF"), // Cyan accent for turbo boost
		MemoryUsed:   lipgloss.Color("#00FF88"), // Green for memory
		MemoryFree:   lipgloss.Color("240"),     // Dim gray
		Border:       lipgloss.Color("240"),
		BorderActive: lipgloss.Color("#00D4FF"),
		Text:         lipgloss.Color("255"),
		TextDim:      lipgloss.Color("240"),
		TextBold:     lipgloss.Color("255"),
		BorderStyle:  lipgloss.RoundedBorder(),
	}
}

// returns the appropriate color based on usage percentage
// refined 3-level scheme: <20% low (cyan), 20-70% medium (yellow), >70% high (pink)
func (t Theme) UsageColor(usage float64) lipgloss.Color {
	switch {
	case usage < 5:
		return t.UsageIdle
	case usage < 20:
		return t.UsageLow
	case usage < 70:
		return t.UsageMedium
	default:
		return t.UsageHigh
	}
}

// returns a styled string for usage percentage with semantic color
func (t Theme) UsageStyle(usage float64) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.UsageColor(usage)).
		Bold(true)
}

// returns a status label for usage percentage
func UsageLabel(usage float64) string {
	switch {
	case usage < 10:
		return "idle"
	case usage < 30:
		return "low"
	case usage < 70:
		return "active"
	case usage < 90:
		return "busy"
	default:
		return "HOT"
	}
}
