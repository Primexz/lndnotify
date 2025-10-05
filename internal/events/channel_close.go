package events

import (
	"time"

	"github.com/Primexz/lndnotify/pkg/format"
)

type ChannelCloseEvent struct {
	Alias          string
	SettledBalance int64
	timestamp      time.Time
}

type ChannelCloseTemplate struct {
	PeerAlias      string
	SettledBalance string
}

func NewChannelCloseEvent(alias string, settledBalance int64) *ChannelCloseEvent {
	return &ChannelCloseEvent{
		Alias:          alias,
		SettledBalance: settledBalance,
		timestamp:      time.Now(),
	}
}

func (e *ChannelCloseEvent) Type() EventType {
	return Event_CHANNEL_CLOSE
}

func (e *ChannelCloseEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChannelCloseEvent) GetTemplateData() interface{} {
	return &ChannelCloseTemplate{
		PeerAlias:      e.Alias,
		SettledBalance: format.FormatSats(float64(e.SettledBalance)),
	}
}
