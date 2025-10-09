package format

import (
	"math"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// FormatSats formats a number with thousand separators and up to 3 decimal places.
// Should be used for fees or other fractional amounts.
func FormatDetailed(value float64) string {
	p := message.NewPrinter(language.English)

	if value == math.Floor(value) {
		return p.Sprintf("%.0f", value)
	}

	// Truncate to 3 decimal places
	truncated := math.Floor(math.Abs(value)*1000) / 1000
	if value < 0 {
		truncated = -truncated
	}

	formatted := p.Sprintf("%.3f", truncated)
	// Remove trailing zeros and decimal point if necessary
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")

	return formatted
}

// FormatWholeNumber formats a number with thousand separators, rounding to the nearest integer.
// Should be used for high-value amounts where fractional parts are not relevant.
func FormatBasic(value float64) string {
	p := message.NewPrinter(language.English)
	// Round to nearest integer
	rounded := math.Round(value)
	return p.Sprintf("%.0f", rounded)
}

// FormatRatePPM formats a rate as parts per million (ppm), rounding to the nearest integer.
// If the total is zero, it returns "0".
func FormatRatePPM(value float64, total float64) string {
	var rate float64

	if total > 0 {
		rate = value * 1e6 / total
	}
	return FormatBasic(rate)
}
