package events

import (
	"time"

	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type PeerOnlineEvent struct {
	NodeInfo  *lnrpc.NodeInfo
	Event     *lnrpc.PeerEvent
	timestamp time.Time
}

type PeerOnlineTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
}

func NewPeerOnlineEvent(nodeInfo *lnrpc.NodeInfo, event *lnrpc.PeerEvent) *PeerOnlineEvent {
	return &PeerOnlineEvent{
		NodeInfo:  nodeInfo,
		Event:     event,
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
	var alias string
	if e.NodeInfo != nil {
		alias = e.NodeInfo.Node.Alias
	}

	return &PeerOnlineTemplate{
		PeerAlias:       alias,
		PeerPubKey:      e.Event.GetPubKey(),
		PeerPubkeyShort: format.FormatPubKey(e.Event.GetPubKey()),
	}
}
