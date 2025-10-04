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

func (e *PeerOfflineEvent) Type() string {
	return "peer_offline_event"
}

func (e *PeerOfflineEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *PeerOfflineEvent) GetTemplateData() interface{} {
	return &PeerOfflineTemplate{
		PeerAlias: e.Alias,
	}
}
