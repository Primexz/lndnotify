package events

import (
	"time"

	"github.com/Primexz/lnd-notify/pkg/format"
)

type ChannelOpenEvent struct {
	Alias     string
	Capacity  int64
	timestamp time.Time
}

type ChannelOpenTemplate struct {
	PeerAlias string
	Capacity  string
}

func NewChannelOpenEvent(alias string, capacity int64) *ChannelOpenEvent {
	return &ChannelOpenEvent{
		Alias:     alias,
		Capacity:  capacity,
		timestamp: time.Now(),
	}
}

func (e *ChannelOpenEvent) Type() string {
	return "channel_open_event"
}

func (e *ChannelOpenEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChannelOpenEvent) GetTemplateData() interface{} {
	return &ChannelOpenTemplate{
		PeerAlias: e.Alias,
		Capacity:  format.FormatSats(float64(e.Capacity)),
	}
}
