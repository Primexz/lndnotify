package events

import (
	"time"

	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type PeerOfflineEvent struct {
	Node      *lnrpc.LightningNode
	timestamp time.Time
}

type PeerOfflineTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
}

func NewPeerOfflineEvent(node *lnrpc.LightningNode) *PeerOfflineEvent {
	return &PeerOfflineEvent{
		Node:      node,
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
		PeerAlias:       e.Node.Alias,
		PeerPubKey:      e.Node.PubKey,
		PeerPubkeyShort: format.FormatPubKey(e.Node.PubKey),
	}
}
