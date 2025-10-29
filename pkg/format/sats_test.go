package format

import (
	"testing"

	"golang.org/x/text/language"
)

func TestFormatDetailed(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{12, "12"},
		{123, "123"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
		{1234567.0, "1,234,567"},
		{1234567.89, "1,234,567.89"},
		{1234567.8912, "1,234,567.891"},
		{1234567.8999, "1,234,567.899"},
		{1000000, "1,000,000"},
		{1000000.1234, "1,000,000.123"},
		{999.9999, "999.999"},
		{123.1, "123.1"},
		{123.10, "123.1"},
		{123.100, "123.1"},
		{123.001, "123.001"},
		{100.000, "100"},
	}

	for _, tt := range tests {
		got := FormatDetailed(tt.value, language.English)
		if got != tt.expected {
			t.Errorf("TestFormatDetailed(%v) = %q; want %q", tt.value, got, tt.expected)
		}
	}
}

func TestFormatBasic(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{12, "12"},
		{123, "123"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
		{1234567.0, "1,234,567"},
		{1234567.49, "1,234,567"}, // rounds down
		{1234567.5, "1,234,568"},  // rounds up
		{1234567.89, "1,234,568"}, // rounds up
		{1000000, "1,000,000"},
		{1000000.1234, "1,000,000"}, // rounds down
		{999.9999, "1,000"},         // rounds up
		{123.1, "123"},
		{123.6, "124"},
		{-123.4, "-123"}, // negative rounds toward zero
		{-123.6, "-124"}, // negative rounds away from zero
	}

	for _, tt := range tests {
		got := FormatBasic(tt.value, language.English)
		if got != tt.expected {
			t.Errorf("TestFormatBasic(%v) = %q; want %q", tt.value, got, tt.expected)
		}
	}
}

func TestFormatRatePPM(t *testing.T) {
	tests := []struct {
		numerator, denominator float64
		expected               string
	}{
		{0, 100, "0"},          // Zero numerator
		{50, 0, "0"},           // Zero denominator
		{1, 1000, "1,000"},     // 1000 ppm
		{5, 1000, "5,000"},     // 5000 ppm
		{10, 100000, "100"},    // 100 ppm
		{0.5, 1000, "500"},     // Fractional
		{1.7, 3000, "567"},     // Rounds up
		{50, 1000, "50,000"},   // 5%
		{100, 1000, "100,000"}, // 10%
		{2, 10000, "200"},      // Typical LN fee
		{1, 100000000, "0"},    // Very small rate
		{90, 100, "900,000"},   // 90%
		{-5, 1000, "-5,000"},   // Negative
		{5, -1000, "0"},        // Negative denominator
	}

	for i, tt := range tests {
		got := FormatRatePPM(tt.numerator, tt.denominator, language.English)
		if got != tt.expected {
			t.Errorf("Test %d: FormatRatePPM(%.1f, %.0f) = %s; want %s",
				i+1, tt.numerator, tt.denominator, got, tt.expected)
		}
	}
}

func TestFormatDetailedWithDifferentLanguages(t *testing.T) {
	tests := []struct {
		value    float64
		lang     language.Tag
		expected string
	}{
		// English
		{1234567.89, language.English, "1,234,567.89"},
		{1000000, language.English, "1,000,000"},
		{1234567.123, language.English, "1,234,567.123"},

		// German
		{1234567.89, language.German, "1.234.567,89"},
		{1000000, language.German, "1.000.000"},
		{1234567.123, language.German, "1.234.567,123"},
		{1.001, language.German, "1,001"},
	}

	for _, tt := range tests {
		got := FormatDetailed(tt.value, tt.lang)
		if got != tt.expected {
			t.Errorf("FormatDetailed(%v, %v) = %q; want %q", tt.value, tt.lang, got, tt.expected)
		}
	}
}

func TestFormatBasicWithDifferentLanguages(t *testing.T) {
	tests := []struct {
		value    float64
		lang     language.Tag
		expected string
	}{
		// English
		{1234567, language.English, "1,234,567"},
		{1000000, language.English, "1,000,000"},
		{1234567.123, language.English, "1,234,567"},

		// German
		{1234567, language.German, "1.234.567"},
		{1000000, language.German, "1.000.000"},
		{1234567.123, language.German, "1.234.567"},
	}

	for _, tt := range tests {
		got := FormatBasic(tt.value, tt.lang)
		if got != tt.expected {
			t.Errorf("FormatBasic(%v, %v) = %q; want %q", tt.value, tt.lang, got, tt.expected)
		}
	}
}

func TestFormatRatePPMWithDifferentLanguages(t *testing.T) {
	tests := []struct {
		numerator, denominator float64
		lang                   language.Tag
		expected               string
	}{
		// English
		{50, 1000, language.English, "50,000"},
		{2, 10000, language.English, "200"},
		{100, 100000, language.English, "1,000"},

		// German
		{50, 1000, language.German, "50.000"},
		{2, 10000, language.German, "200"},
		{100, 100000, language.German, "1.000"},
	}

	for _, tt := range tests {
		got := FormatRatePPM(tt.numerator, tt.denominator, tt.lang)
		if got != tt.expected {
			t.Errorf("FormatRatePPM(%v, %v, %v) = %q; want %q", tt.numerator, tt.denominator, tt.lang, got, tt.expected)
		}
	}
}
