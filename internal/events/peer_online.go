package events

import (
	"time"
)

type PeerOnlineEvent struct {
	Alias     string
	timestamp time.Time
}

type PeerOnlineTemplate struct {
	PeerAlias string
}

func NewPeerOnlineEvent(alias string) *PeerOnlineEvent {
	return &PeerOnlineEvent{
		Alias:     alias,
		timestamp: time.Now(),
	}
}

func (e *PeerOnlineEvent) Type() string {
	return "peer_online_event"
}

func (e *PeerOnlineEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *PeerOnlineEvent) GetTemplateData() interface{} {
	return &PeerOnlineTemplate{
		PeerAlias: e.Alias,
	}
}
