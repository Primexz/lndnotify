package format

import (
	"testing"
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
		got := FormatDetailed(tt.value)
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
		got := FormatBasic(tt.value)
		if got != tt.expected {
			t.Errorf("TestFormatBasic(%v) = %q; want %q", tt.value, got, tt.expected)
		}
	}
}
