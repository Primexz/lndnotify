package chainutil

import (
	"time"
)

// BlockCountToDuration converts a Bitcoin block count to a time.Duration.
func BlockCountToDuration(blockCount int) time.Duration {
	// avg block time is 10 minutes
	blockTime := 10 * time.Minute
	return time.Duration(blockCount) * blockTime
}
