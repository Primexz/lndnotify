package format

import (
	"testing"
)

func TestFormatSats(t *testing.T) {
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
		{1234567.89, "1,234,567.890"},
		{1234567.8912, "1,234,567.891"},
		{1234567.8999, "1,234,567.899"},
		{1000000, "1,000,000"},
		{1000000.1234, "1,000,000.123"},
		{999.9999, "999.999"},
	}

	for _, tt := range tests {
		got := FormatSats(tt.value)
		if got != tt.expected {
			t.Errorf("FormatSats(%v) = %q; want %q", tt.value, got, tt.expected)
		}
	}
}
