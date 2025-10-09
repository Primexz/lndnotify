package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
)

// Event is the base interface for all Lightning Network events
type Event interface {
	Type() EventType
	Timestamp() time.Time
	GetTemplateData() interface{}
	ShouldProcess(cfg *config.Config) bool
}

type EventType string

const (
	Event_FORWARD         EventType = "forward_event"
	Event_PEER_ONLINE     EventType = "peer_online_event"
	Event_PEER_OFFLINE    EventType = "peer_offline_event"
	Event_CHANNEL_OPEN    EventType = "channel_open_event"
	Event_CHANNEL_CLOSE   EventType = "channel_close_event"
	Event_INVOICE_SETTLED EventType = "invoice_settled_event"
	Event_FAILED_HTLC     EventType = "failed_htlc_event"
)

func (et EventType) String() string {
	return string(et)
}
