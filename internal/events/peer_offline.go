package events

import (
	"time"
)

type PeerOfflineEvent struct {
	Alias     string
	timestamp time.Time
}

type PeerOfflineTemplate struct {
	PeerAlias string
}

func NewPeerOfflineEvent(alias string) *PeerOfflineEvent {
	return &PeerOfflineEvent{
		Alias:     alias,
		timestamp: time.Now(),
	}
}

func (e *PeerOfflineEvent) Type() EventType {
	return Event_PEER_OFFLINE
}

func (e *PeerOfflineEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *PeerOfflineEvent) GetTemplateData() interface{} {
	return &PeerOfflineTemplate{
		PeerAlias: e.Alias,
	}
}
