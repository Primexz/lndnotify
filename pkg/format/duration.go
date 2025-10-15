package format

import (
	"time"
)

func FormatDuration(d time.Duration) time.Duration {
	return time.Duration.Round(d, time.Second)
}
