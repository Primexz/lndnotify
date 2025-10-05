package format

import "testing"

func TestFormatPubKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abcdef", "abcdef"},                       // less than 10 chars
		{"123456789", "123456789"},                 // exactly 9 chars
		{"1234567890", "12345678"},                 // exactly 10 chars
		{"abcdefghijklmnopqrstuvwxyz", "abcdefgh"}, // more than 10 chars
		{"", ""}, // empty string
	}

	for _, tt := range tests {
		result := FormatPubKey(tt.input)
		if result != tt.expected {
			t.Errorf("FormatPubKey(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}
