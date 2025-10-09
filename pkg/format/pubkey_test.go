package format

import "testing"

func TestFormatPubKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abcdef", "abcdef"},                       // less than 8 chars
		{"12345678", "12345678"},                   // exactly 8 chars
		{"123456789", "12345678"},                  // exactly 9 chars
		{"abcdefghijklmnopqrstuvwxyz", "abcdefgh"}, // more than 9 chars
		{"", ""}, // empty string
	}

	for _, tt := range tests {
		result := FormatPubKey(tt.input)
		if result != tt.expected {
			t.Errorf("FormatPubKey(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}
