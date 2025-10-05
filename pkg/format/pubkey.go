package format

// FormatSats formats a pubkey to a shorter version
func FormatPubKey(value string) string {
	if len(value) < 10 {
		return value
	}

	return value[:8]
}
