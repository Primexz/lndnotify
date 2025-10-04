package events

import (
	"time"
)

// Event is the base interface for all Lightning Network events
type Event interface {
	Type() string
	Timestamp() time.Time
	GetTemplateData() interface{}
}
