package format

import "testing"

func TestCalculatePercentageChange(t *testing.T) {
	tests := []struct {
		name     string
		oldVal   int64
		newVal   int64
		expected string
	}{
		{
			name:     "both values zero",
			oldVal:   0,
			newVal:   0,
			expected: "0%",
		},
		{
			name:     "old value zero, new value positive",
			oldVal:   0,
			newVal:   100,
			expected: "∞%",
		},
		{
			name:     "old value zero, new value negative",
			oldVal:   0,
			newVal:   -100,
			expected: "∞%",
		},
		{
			name:     "positive change",
			oldVal:   100,
			newVal:   150,
			expected: "+50.0%",
		},
		{
			name:     "negative change",
			oldVal:   100,
			newVal:   50,
			expected: "-50.0%",
		},
		{
			name:     "no change",
			oldVal:   100,
			newVal:   100,
			expected: "0.0%",
		},
		{
			name:     "large positive change",
			oldVal:   10,
			newVal:   100,
			expected: "+900.0%",
		},
		{
			name:     "small positive change",
			oldVal:   1000,
			newVal:   1001,
			expected: "+0.1%",
		},
		{
			name:     "small negative change",
			oldVal:   1000,
			newVal:   999,
			expected: "-0.1%",
		},
		{
			name:     "positive to negative",
			oldVal:   100,
			newVal:   -100,
			expected: "-200.0%",
		},
		{
			name:     "both negative values",
			oldVal:   -200,
			newVal:   -100,
			expected: "-50.0%",
		},
		{
			name:     "fractional percentage",
			oldVal:   333,
			newVal:   334,
			expected: "+0.3%",
		},
		{
			name:     "large numbers",
			oldVal:   1000000,
			newVal:   1500000,
			expected: "+50.0%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePercentageChange(tt.oldVal, tt.newVal)
			if result != tt.expected {
				t.Errorf("CalculatePercentageChange(%d, %d) = %s, want %s", tt.oldVal, tt.newVal, result, tt.expected)
			}
		})
	}
}

func TestCalculateAbsoluteChange(t *testing.T) {
	tests := []struct {
		name     string
		oldVal   int64
		newVal   int64
		expected string
	}{
		{
			name:     "both values zero",
			oldVal:   0,
			newVal:   0,
			expected: "0",
		},
		{
			name:     "positive change",
			oldVal:   100,
			newVal:   150,
			expected: "+50",
		},
		{
			name:     "negative change",
			oldVal:   100,
			newVal:   50,
			expected: "-50",
		},
		{
			name:     "no change",
			oldVal:   100,
			newVal:   100,
			expected: "0",
		},
		{
			name:     "large positive change",
			oldVal:   10,
			newVal:   1000,
			expected: "+990",
		},
		{
			name:     "large negative change",
			oldVal:   1000,
			newVal:   10,
			expected: "-990",
		},
		{
			name:     "negative to positive",
			oldVal:   -100,
			newVal:   100,
			expected: "+200",
		},
		{
			name:     "positive to negative",
			oldVal:   100,
			newVal:   -100,
			expected: "-200",
		},
		{
			name:     "both negative values - increase",
			oldVal:   -200,
			newVal:   -100,
			expected: "+100",
		},
		{
			name:     "both negative values - decrease",
			oldVal:   -100,
			newVal:   -200,
			expected: "-100",
		},
		{
			name:     "from zero to positive",
			oldVal:   0,
			newVal:   100,
			expected: "+100",
		},
		{
			name:     "from zero to negative",
			oldVal:   0,
			newVal:   -100,
			expected: "-100",
		},
		{
			name:     "to zero from positive",
			oldVal:   100,
			newVal:   0,
			expected: "-100",
		},
		{
			name:     "to zero from negative",
			oldVal:   -100,
			newVal:   0,
			expected: "+100",
		},
		{
			name:     "maximum int64 values",
			oldVal:   9223372036854775806, // max int64 - 1
			newVal:   9223372036854775807, // max int64
			expected: "+1",
		},
		{
			name:     "minimum int64 values",
			oldVal:   -9223372036854775807, // min int64 + 1
			newVal:   -9223372036854775808, // min int64
			expected: "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAbsoluteChange(tt.oldVal, tt.newVal)
			if result != tt.expected {
				t.Errorf("CalculateAbsoluteChange(%d, %d) = %s, want %s", tt.oldVal, tt.newVal, result, tt.expected)
			}
		})
	}
}
