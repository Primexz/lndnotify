package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type PeerOfflineEvent struct {
	NodeInfo  *lnrpc.NodeInfo
	Event     *lnrpc.PeerEvent
	timestamp time.Time
}

type PeerOfflineTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
}

func NewPeerOfflineEvent(nodeInfo *lnrpc.NodeInfo, event *lnrpc.PeerEvent) *PeerOfflineEvent {
	return &PeerOfflineEvent{
		NodeInfo:  nodeInfo,
		Event:     event,
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
	var alias string
	if e.NodeInfo != nil {
		alias = e.NodeInfo.Node.Alias
	}

	return &PeerOfflineTemplate{
		PeerAlias:       alias,
		PeerPubKey:      e.Event.GetPubKey(),
		PeerPubkeyShort: format.FormatPubKey(e.Event.GetPubKey()),
	}
}

func (e *PeerOfflineEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.PeerEvents
}
