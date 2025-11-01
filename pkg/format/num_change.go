package format

import "fmt"

// calculatePercentageChange calculates the percentage change between old and new values
func CalculatePercentageChange(oldVal int64, newVal int64) string {
	if oldVal == 0 {
		if newVal == 0 {
			return "0%"
		}
		return "∞%"
	}

	change := ((float64(newVal) - float64(oldVal)) / float64(oldVal)) * 100
	if change > 0 {
		return fmt.Sprintf("+%.1f%%", change)
	}
	return fmt.Sprintf("%.1f%%", change)
}

// calculateAbsoluteChange calculates the absolute change between old and new values
func CalculateAbsoluteChange(oldVal int64, newVal int64) string {
	change := newVal - oldVal
	if change > 0 {
		return fmt.Sprintf("+%d", change)
	}
	return fmt.Sprintf("%d", change)
}
