// Taken from https://github.com/leg100
// All credits to leg100 for the border title function in bubble tea

package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color constants for borders
var (
	Blue                  = lipgloss.Color("#1077e5")
	InactivePreviewBorder = lipgloss.Color("240")
)

// pane models expose metadata to be embedded in certain positions in borders:
// options:
// * one method per position (Title(), TopLeft(), etc)
// * one method returning positions: Metadata() map[border-position]string

// inputs:
// * content to wrap with border, from which width and height is determined
// * metadata to place within certain positions in border:
// map[border-position]string; or use functional options,
// WithBottomLeft(string), etc.
// * border style; alternatively let method define style based on a
// border-active input
// output:
// * content wrapped with border.

// positions:
// * Title (table row info)
// * TopLeft (task info; resource info; log msg info; task group info)
// * BottomLeft (task status; task summary; state summary)
// * BottomRight (scroll percentage)
//

type BorderPosition int

const (
	TopLeftBorder BorderPosition = iota
	TopMiddleBorder
	TopRightBorder
	BottomLeftBorder
	BottomMiddleBorder
	BottomRightBorder
)

func borderize(content string, active bool, embeddedText map[BorderPosition]string) string {
	if embeddedText == nil {
		embeddedText = make(map[BorderPosition]string)
	}

	// Always use rounded borders for consistent minimalist aesthetic
	border := lipgloss.RoundedBorder()

	// Use theme colors
	theme := DefaultTheme()
	color := map[bool]lipgloss.TerminalColor{
		true:  theme.BorderActive,
		false: theme.Border,
	}

	style := lipgloss.NewStyle().Foreground(color[active])
	width := lipgloss.Width(content)

	encloseInSquareBrackets := func(text string) string {
		if text != "" {
			return fmt.Sprintf("%s%s%s",
				style.Render(border.TopRight),
				text,
				style.Render(border.TopLeft),
			)
		}
		return text
	}
	buildHorizontalBorder := func(leftText, middleText, rightText, leftCorner, inbetween, rightCorner string) string {
		leftText = encloseInSquareBrackets(leftText)
		middleText = encloseInSquareBrackets(middleText)
		rightText = encloseInSquareBrackets(rightText)
		// Calculate length of border between embedded texts
		remaining := max(0, width-lipgloss.Width(leftText)-lipgloss.Width(middleText)-lipgloss.Width(rightText))
		leftBorderLen := max(0, (width/2)-lipgloss.Width(leftText)-(lipgloss.Width(middleText)/2))
		rightBorderLen := max(0, remaining-leftBorderLen)
		// Then construct border string
		s := leftText +
			style.Render(strings.Repeat(inbetween, leftBorderLen)) +
			middleText +
			style.Render(strings.Repeat(inbetween, rightBorderLen)) +
			rightText
		// Make it fit in the space available between the two corners.
		s = lipgloss.NewStyle().
			Inline(true).
			MaxWidth(width).
			Render(s)
		// Add the corners
		return style.Render(leftCorner) + s + style.Render(rightCorner)
	}
	// Stack top border, content and horizontal borders, and bottom border.
	return strings.Join([]string{
		buildHorizontalBorder(
			embeddedText[TopLeftBorder],
			embeddedText[TopMiddleBorder],
			embeddedText[TopRightBorder],
			border.TopLeft,
			border.Top,
			border.TopRight,
		),
		lipgloss.NewStyle().
			BorderForeground(color[active]).
			Border(border, false, true, false, true).Render(content),
		buildHorizontalBorder(
			embeddedText[BottomLeftBorder],
			embeddedText[BottomMiddleBorder],
			embeddedText[BottomRightBorder],
			border.BottomLeft,
			border.Bottom,
			border.BottomRight,
		),
	}, "\n")
}
