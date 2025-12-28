package ui

import (
	"fmt"
	"strings"

	"wrenchi/internal/cpu"
	"wrenchi/internal/process"
	"wrenchi/internal/system"
	cpu_util "wrenchi/internal/util"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type layoutMode int

const (
	layoutVertical   layoutMode = iota // system summary | cores | processes (vertical)
	layoutHorizontal                   // system summary | (cores + processes side-by-side)
	layoutCompact                      // system summary + minimal info
)

// determines which layout to use based on terminal dimensions
func (m model) getLayoutMode() layoutMode {
	if m.height < 30 {
		return layoutHorizontal
	}

	if m.height < 15 {
		return layoutCompact
	}
	return layoutVertical
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	defer func() {
		if m.cacheDirty {
			m.cacheDirty = false
		}
	}()

	mode := m.getLayoutMode()

	switch mode {
	case layoutCompact:
		return m.renderCompactView()
	case layoutHorizontal:
		return m.renderHorizontalView()
	default:
		return m.renderVerticalView()
	}
}

// renders only the summary for small screens
func (m model) renderCompactView() string {
	return m.renderSystemSummary(m.height)
}

// renders the normal stacked layout
func (m model) renderVerticalView() string {
	summaryHeight := int(float64(m.height) * 0.2)
	if summaryHeight < 8 {
		summaryHeight = 8
	}

	coresHeight := int(float64(m.height) * 0.3)
	if coresHeight < 8 {
		coresHeight = 8
	}

	processHeight := m.height - summaryHeight - coresHeight
	if processHeight < 10 {
		processHeight = 10
	}

	summary := m.renderSystemSummary(summaryHeight)
	cores := m.renderPerCoreAdaptiveGrid(coresHeight, m.width)

	m.processTable.SetHeight(processHeight - 3)
	processes := m.renderProcessTable(processHeight, m.width)

	return lipgloss.JoinVertical(lipgloss.Left, summary, cores, processes)
}

// renders cores and processes side-by-side for small screens
func (m model) renderHorizontalView() string {
	summaryHeight := int(float64(m.height) * 0.3)
	if summaryHeight < 10 {
		summaryHeight = 10
	}

	bottomHeight := m.height - summaryHeight

	// needs at least 6 lines to show anything useful in bottom component
	if bottomHeight < 6 {
		bottomHeight = 6
		summaryHeight = m.height - bottomHeight
		if summaryHeight < 8 {
			summaryHeight = 8
		}
	}

	// don't exceed terminal height
	if summaryHeight+bottomHeight > m.height {
		ratio := float64(m.height) / float64(summaryHeight+bottomHeight)
		summaryHeight = int(float64(summaryHeight) * ratio)
		bottomHeight = m.height - summaryHeight
	}

	summary := m.renderSystemSummary(summaryHeight)

	coresWidth := int(float64(m.width) * 0.4)
	if coresWidth < 30 {
		coresWidth = 30
	}
	processWidth := m.width - coresWidth
	if processWidth < 40 {
		processWidth = 40
		coresWidth = m.width - processWidth
	}

	cores := m.renderPerCoreAdaptiveGrid(bottomHeight, coresWidth)

	processes := m.renderProcessTable(bottomHeight, processWidth)

	bottom := lipgloss.JoinHorizontal(lipgloss.Top, cores, processes)

	return lipgloss.JoinVertical(lipgloss.Left, summary, bottom)
}

// renders system summary with optional spacing control
func (m model) renderSystemSummary(maxHeight int) string {
	var s strings.Builder
	contentHeight := maxHeight - 2 // -2 for top and bottom borders

	// hostname | cores | uptime
	headerParts := []string{m.hostname}
	headerParts = append(headerParts, fmt.Sprintf("%d cores", len(m.cpuStats.PerCore)))

	uptime, err := system.Uptime()
	if err == nil {
		headerParts = append(headerParts, fmt.Sprintf("up %s", system.FormatUptime(uptime)))
	}

	headerStyle := lipgloss.NewStyle().Foreground(m.theme.TextDim)
	s.WriteString(headerStyle.Render(strings.Join(headerParts, " │ ")))
	s.WriteString("\n\n")

	if m.cpuStats.CPU != "" && contentHeight > 5 {
		cpuName := m.cpuStats.CPU
		maxCPUNameWidth := m.width - 10
		if len(cpuName) > maxCPUNameWidth {
			cpuName = cpuName[:maxCPUNameWidth-3] + "..."
		}
		cpuNameStyle := lipgloss.NewStyle().Foreground(m.theme.Text).Bold(false)
		s.WriteString(cpuNameStyle.Render(cpuName))
		s.WriteString("\n\n")

	}

	cpuUsage := cpu.CalculateCPUUsage(m.cpuPrev.Total, m.cpuStats.Total)

	// cpu usage bar
	barWidth := m.width - 30
	if barWidth < 20 {
		barWidth = 20
	}
	if barWidth > 50 {
		barWidth = 50
	}

	cpuBar := renderProgressBar(cpuUsage, barWidth, m.theme)
	usageStyle := m.theme.UsageStyle(cpuUsage)
	s.WriteString(fmt.Sprintf("CPU  %s %s\n",
		cpuBar,
		usageStyle.Render(fmt.Sprintf("%5.1f%% | %s", cpuUsage, m.cpuStats.AvgFreq))))

	// memory usage bar
	memUsed := m.memoryStats.Total - m.memoryStats.Available
	memUsage := (float64(memUsed) / float64(m.memoryStats.Total)) * 100.0

	memBar := renderProgressBar(memUsage, barWidth, m.theme)
	memStyle := lipgloss.NewStyle().Foreground(m.theme.MemoryUsed).Bold(true)
	s.WriteString(fmt.Sprintf("MEM  %s %s (%s/%s)\n",
		memBar,
		memStyle.Render(fmt.Sprintf("%5.1f%%", memUsage)),
		formatMemoryShort(memUsed),
		formatMemoryShort(m.memoryStats.Total)))

	content := s.String()
	constrainedStyle := lipgloss.NewStyle().Height(contentHeight).Width(m.width - 2)
	content = constrainedStyle.Render(content)

	return borderize(content, true, map[BorderPosition]string{
		TopLeftBorder: "wrenchi",
	}, m.theme)
}

// creates a colored progress bar using block characters
func renderProgressBar(percent float64, width int, theme Theme) string {
	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	// Use block characters for a clean look
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	// Color based on usage
	style := lipgloss.NewStyle().Foreground(theme.UsageColor(percent))
	return style.Render(bar)
}

// TODO: create utils package or move somewhere else
// formats memory in a compact way (e.g., "8.4G")
func formatMemoryShort(kb uint64) string {
	const (
		mb = 1024
		gb = 1024 * 1024
	)

	if kb > gb {
		return fmt.Sprintf("%.1fG", float64(kb)/float64(gb))
	}
	if kb > mb {
		return fmt.Sprintf("%.1fM", float64(kb)/float64(mb))
	}
	return fmt.Sprintf("%dK", kb)
}

// creates a sparkline with semantic color coding
func renderBrailleSparklineColored(history []float64, width int, theme Theme) string {
	heights := []rune{
		'⣀', // 0-2%
		'⣀', // >3%
		'⣄', // ~28%
		'⣆', // ~42%
		'⣇', // ~57%
		'⣧', // ~71%
		'⣷', // ~85%
		'⣿', // 100%
	}

	var result strings.Builder

	paddingCount := width - len(history)
	if paddingCount < 0 {
		paddingCount = 0
	}

	dimStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
	for i := 0; i < paddingCount; i++ {
		result.WriteString(dimStyle.Render("⣀"))
	}

	startIdx := 0
	if len(history) > width {
		startIdx = len(history) - width
	}

	for i := startIdx; i < len(history); i++ {
		usage := history[i]

		if usage < 0 {
			usage = 0
		}
		if usage > 100 {
			usage = 100
		}

		// map usage to height index
		heightIdx := int(usage / 100.0 * float64(len(heights)-1))
		if heightIdx >= len(heights) {
			heightIdx = len(heights) - 1
		}

		char := string(heights[heightIdx])

		style := lipgloss.NewStyle().Foreground(theme.UsageColor(usage))
		result.WriteString(style.Render(char))
	}

	return result.String()
}

// renders the process table with a specific height and width
func (m model) renderProcessTable(height, width int) string {
	// PID(8) + Name(20) + CPU%(8) + MEM%(8) = 44
	const (
		fixedColumnsWidth = 44
		tablePadding      = 10
		minCommandWidth   = 14
	)

	commandWidth := width - fixedColumnsWidth - tablePadding
	if commandWidth < minCommandWidth {
		commandWidth = minCommandWidth
	}

	columns := []table.Column{
		{Title: "PID", Width: 8},
		{Title: "Name", Width: 20},
		{Title: "CPU%", Width: 8},
		{Title: "MEM%", Width: 8},
		{Title: "Command", Width: commandWidth},
	}
	m.processTable.SetColumns(columns)

	tableView := m.processTable.View()

	scrollInfo := fmt.Sprintf(" %d/%d ",
		m.processTable.Cursor()+1,
		len(m.processTable.Rows()))

	constrainedStyle := lipgloss.NewStyle().MaxWidth(width - 2).MaxHeight(height - 2).Width(width)
	tableView = constrainedStyle.Render(tableView)

	helpStyle := lipgloss.NewStyle().Foreground(m.theme.TextDim)
	helpText := helpStyle.Render(" ↑↓ scroll │ q quit ")

	return borderize(tableView, true, map[BorderPosition]string{
		TopLeftBorder:     "processes",
		BottomLeftBorder:  helpText,
		BottomRightBorder: scrollInfo,
	}, m.theme)
}

// returns the number of columns for per-core usage
// based on current component alloted width
func gridColumns(componentWidth int) int {
	switch {
	case componentWidth >= 200:
		return 5
	case componentWidth >= 150:
		return 4
	case componentWidth >= 100:
		return 3
	case componentWidth >= 65:
		return 2
	default:
		return 1
	}
}

// renders all cores in an adaptive multi-column grid
func (m model) renderPerCoreAdaptiveGrid(maxHeight, width int) string {
	cores := m.cpuStats.PerCore
	if len(cores) == 0 {
		return borderize("No CPU cores detected", true, map[BorderPosition]string{
			TopLeftBorder: "cpu cores",
		}, m.theme)
	}

	contentHeight := maxHeight - 2
	columns := gridColumns(width)
	cellWidth := (width - 2) / columns

	// TODO: make it a constant
	// 6 (core name) + 7 (percentage) + 8 (frequency) + 4 (spacing/markers) = 25
	var sparkWidth int = cellWidth - 25
	if sparkWidth <= 5 {
		sparkWidth = 0
	}

	var maxUsage float64
	for i, core := range cores {
		var usage float64
		if i < len(m.cpuPrev.PerCore) {
			usage = cpu.CalculateCPUUsage(m.cpuPrev.PerCore[i], core)
		}
		if usage > maxUsage {
			maxUsage = usage
		}
	}

	var gridRows []string

	rows := (len(cores) + columns - 1) / columns
	for row := 0; row < rows && row < contentHeight; row++ {
		var rowCells []string
		for col := 0; col < columns; col++ {
			coreIdx := row*columns + col
			if coreIdx >= len(cores) {
				emptyCellStyle := lipgloss.NewStyle().Width(cellWidth)
				rowCells = append(rowCells, emptyCellStyle.Render(""))
				continue
			}

			cell := m.renderCoreCell(coreIdx, cellWidth, sparkWidth)
			rowCells = append(rowCells, cell)
		}
		gridRows = append(gridRows, lipgloss.JoinHorizontal(lipgloss.Top, rowCells...))
	}

	content := strings.Join(gridRows, "\n")

	constrainedStyle := lipgloss.NewStyle().Height(contentHeight).Width(width - 2)
	content = constrainedStyle.Render(content)

	statsStyle := lipgloss.NewStyle().Foreground(m.theme.TextDim)

	borderInfo := fmt.Sprintf(" peak: %.1f%% ", maxUsage)
	borderLabels := map[BorderPosition]string{
		TopLeftBorder:  "cpu cores",
		TopRightBorder: statsStyle.Render(borderInfo),
	}

	return borderize(content, true, borderLabels, m.theme)
}

// renders a single core cell with optional sparkline and frequency
func (m model) renderCoreCell(coreIdx int, cellWidth int, sparklineWidth int) string {
	if coreIdx >= len(m.cpuStats.PerCore) {
		return lipgloss.NewStyle().Width(cellWidth).Render("")
	}

	core := m.cpuStats.PerCore[coreIdx]
	var usage float64
	if coreIdx < len(m.cpuPrev.PerCore) {
		usage = cpu.CalculateCPUUsage(m.cpuPrev.PerCore[coreIdx], core)
	}

	// use cache if available
	cacheKey := coreIdx
	if !m.cacheDirty {
		if cached, ok := m.coreViewCache[cacheKey]; ok {
			return lipgloss.NewStyle().Width(cellWidth).Render(cached)
		}
	}

	var parts []string

	coreName := lipgloss.NewStyle().
		Foreground(m.theme.UsageColor(usage)).
		Render(fmt.Sprintf("%-5s", core.Name))
	parts = append(parts, coreName)

	usageStr := m.theme.UsageStyle(usage).Render(fmt.Sprintf("%5.1f%%", usage))
	parts = append(parts, usageStr)

	if sparklineWidth > 0 && coreIdx < len(m.cpuHistory) {
		sparkline := renderBrailleSparklineColored(m.cpuHistory[coreIdx], sparklineWidth, m.theme)
		parts = append(parts, sparkline)
	}

	if usage > 90 {
		hotIndicator := lipgloss.NewStyle().
			Foreground(m.theme.UsageHigh).
			Bold(true).
			Render("!")
		parts = append(parts, hotIndicator)
	}

	cellContent := strings.Join(parts, " ")

	if m.cacheDirty {
		if m.coreViewCache == nil {
			m.coreViewCache = make(map[int]string)
		}
		m.coreViewCache[cacheKey] = cellContent
	}

	return lipgloss.NewStyle().Width(cellWidth).Render(cellContent)
}

// converts process info to table rows
func processesToRows(processes []process.Process, cpuUsage float64) []table.Row {
	var rows []table.Row
	for _, p := range processes {
		pid := fmt.Sprintf("%d", p.PID)
		var actualPercent float64
		if cpuUsage > 0 {
			actualPercent = (p.CPUPercent / 100.0) * cpuUsage
		}
		cpuPercent := fmt.Sprintf("%.1f%%", actualPercent)
		memPercent := fmt.Sprintf("%s", cpu_util.FormatMemory(p.MemKB))

		// get command up to a space
		cmd := strings.Fields(p.Command)[0]

		rows = append(rows, table.Row{pid, p.Name, cpuPercent, memPercent, cmd})
	}

	return rows
}
