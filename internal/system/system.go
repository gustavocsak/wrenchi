package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// reads system uptime from /proc/uptime
func Uptime() (time.Duration, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0, fmt.Errorf("invalid uptime format")
	}

	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(seconds * float64(time.Second)), nil
}

// formats a duration into a human-readable uptime string
func FormatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// returns the system hostname
func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
