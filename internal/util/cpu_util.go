package cpu_util

import (
	"fmt"
)

func FormatGHz(g float64) string {
	return fmt.Sprintf("%.1fGHz", g/1000.0)
}

// TODO: move it to another util
func FormatMemory(kb uint64) string {
	const (
		mb = 1024
		gb = 1024 * 1024
	)

	const format = "%6.2f%s"

	if kb > gb {
		return fmt.Sprintf(format, float64(kb)/float64(gb), "G")
	}

	if kb > mb {
		return fmt.Sprintf(format, float64(kb)/float64(mb), "M")
	}

	return fmt.Sprintf(format, float64(kb), "K")
}
