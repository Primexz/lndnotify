package events

import (
	"time"

	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type ChannelCloseEvent struct {
	Node      *lnrpc.LightningNode
	Channel   *lnrpc.ChannelCloseSummary
	timestamp time.Time
}

type ChannelCloseTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
	SettledBalance  string
	ChanId          uint64
	ChannelPoint    string
	RemotePubkey    string
	Capacity        string
}

func NewChannelCloseEvent(node *lnrpc.LightningNode, channel *lnrpc.ChannelCloseSummary) *ChannelCloseEvent {
	return &ChannelCloseEvent{
		Node:      node,
		Channel:   channel,
		timestamp: time.Now(),
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
		PeerAlias:       e.Node.Alias,
		PeerPubKey:      e.Node.PubKey,
		PeerPubkeyShort: format.FormatPubKey(e.Node.PubKey),
		ChanId:          e.Channel.ChanId,
		ChannelPoint:    e.Channel.ChannelPoint,
		Capacity:        format.FormatSats(float64(e.Channel.Capacity)),
		RemotePubkey:    e.Channel.RemotePubkey,
		SettledBalance:  format.FormatSats(float64(e.Channel.SettledBalance)),
	}
}
