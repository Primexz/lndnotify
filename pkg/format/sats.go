package format

import (
	"math"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// FormatSats formats a number with thousand separators and up to 3 decimal places.
// Should be used for fees or other fractional amounts.
func FormatDetailed(value float64, lang language.Tag) string {
	p := message.NewPrinter(lang)

	if value == math.Floor(value) {
		return p.Sprintf("%.0f", value)
	}

	multiplied := math.Abs(value) * 1000
	truncated := math.Floor(multiplied+1e-9) / 1000
	if value < 0 {
		truncated = -truncated
	}

	formatted := p.Sprintf("%.3f", truncated)
	formatted = strings.TrimRight(formatted, "0")
	if strings.HasSuffix(formatted, ".") || strings.HasSuffix(formatted, ",") {
		formatted = strings.TrimRight(formatted, ".,")
	}

	return formatted
}

// FormatWholeNumber formats a number with thousand separators, rounding to the nearest integer.
// Should be used for high-value amounts where fractional parts are not relevant.
func FormatBasic(value float64, lang language.Tag) string {
	p := message.NewPrinter(lang)
	// Round to nearest integer
	rounded := math.Round(value)
	return p.Sprintf("%.0f", rounded)
}

// FormatRatePPM formats a rate as parts per million (ppm), rounding to the nearest integer.
// If the total is zero, it returns "0".
func FormatRatePPM(value float64, total float64, lang language.Tag) string {
	var rate float64

	if total > 0 {
		rate = value * 1e6 / total
	}
	return FormatBasic(rate, lang)
}
