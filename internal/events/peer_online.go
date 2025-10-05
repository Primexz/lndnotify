package events

import (
	"time"

	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type PeerOnlineEvent struct {
	Node      *lnrpc.LightningNode
	timestamp time.Time
}

type PeerOnlineTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
}

func NewPeerOnlineEvent(node *lnrpc.LightningNode) *PeerOnlineEvent {
	return &PeerOnlineEvent{
		Node:      node,
		timestamp: time.Now(),
	}
}

func (e *PeerOnlineEvent) Type() EventType {
	return Event_PEER_ONLINE
}

func (e *PeerOnlineEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *PeerOnlineEvent) GetTemplateData() interface{} {
	return &PeerOnlineTemplate{
		PeerAlias:       e.Node.Alias,
		PeerPubKey:      e.Node.PubKey,
		PeerPubkeyShort: format.FormatPubKey(e.Node.PubKey),
	}
}
