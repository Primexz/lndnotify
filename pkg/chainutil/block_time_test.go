package chainutil

import (
	"testing"
	"time"
)

func TestBlockCountToDuration(t *testing.T) {
	tests := []struct {
		name  string
		input int32
		want  time.Duration
	}{
		{name: "144 blocks => 24h", input: 144, want: 144 * 10 * time.Minute},
		{name: "72 blocks => 12h", input: 72, want: 72 * 10 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BlockCountToDuration(tt.input)
			if got != tt.want {
				t.Fatalf("BlockCountToDuration(%d) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}
