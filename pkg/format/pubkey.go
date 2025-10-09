package format

// FormatPubKey formats a pubkey to a shorter version
func FormatPubKey(value string) string {
	// Check just for safety preventing panics.
	if len(value) < 9 {
		return value
	}

	return value[:8]
}
