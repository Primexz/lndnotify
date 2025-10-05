package format

import (
	"math"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// FormatSats formats a float64 with thousand separators.
// Whole numbers are displayed without decimals, others show up to 3 decimal places.
func FormatSats(value float64) string {
	p := message.NewPrinter(language.English)

	if value == math.Floor(value) {
		return p.Sprintf("%.0f", value)
	}

	// Truncate to 3 decimal places
	truncated := math.Floor(math.Abs(value)*1000) / 1000
	if value < 0 {
		truncated = -truncated
	}

	return p.Sprintf("%.3f", truncated)
}
