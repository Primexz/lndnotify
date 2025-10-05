package format

import (
	"fmt"
	"strconv"
)

// FormatSats formats a satoshi amount with appropriate decimal places based on the count
// For low numbers (< 1000), it uses 3 decimal places
// For high numbers (>= 1000), it uses 0 decimal places with thousand separators
func FormatSats(sats float64) string {
	// if sats is a whole number, format without decimal places
	if sats == float64(int64(sats)) {
		return strconv.FormatInt(int64(sats), 10)
	}

	if sats < 1000 {
		return fmt.Sprintf("%.3f", float64(sats))
	}

	// Convert to string and add thousand separators
	str := strconv.FormatInt(int64(sats), 10)
	n := len(str)
	if n <= 3 {
		return str
	}

	// Build the formatted string with thousand separators
	var result []byte
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, str[i])
	}

	return string(result)
}
