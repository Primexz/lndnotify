package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type ChannelOpenEvent struct {
	Node      *lnrpc.LightningNode
	Channel   *lnrpc.Channel
	timestamp time.Time
}

type ChannelOpenTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
	SettledBalance  string
	ChanId          uint64
	ChannelPoint    string
	RemotePubkey    string
	Capacity        string
}

func NewChannelOpenEvent(node *lnrpc.LightningNode, channel *lnrpc.Channel) *ChannelOpenEvent {
	return &ChannelOpenEvent{
		Node:      node,
		Channel:   channel,
		timestamp: time.Now(),
	}
}

func (e *ChannelOpenEvent) Type() EventType {
	return Event_CHANNEL_OPEN
}

func (e *ChannelOpenEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChannelOpenEvent) GetTemplateData() interface{} {
	return &ChannelOpenTemplate{
		PeerAlias:       e.Node.Alias,
		PeerPubKey:      e.Node.PubKey,
		PeerPubkeyShort: format.FormatPubKey(e.Node.PubKey),
		ChanId:          e.Channel.ChanId,
		ChannelPoint:    e.Channel.ChannelPoint,
		Capacity:        format.FormatBasic(float64(e.Channel.Capacity)),
		RemotePubkey:    e.Channel.RemotePubkey,
	}
}

func (e *ChannelOpenEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChannelEvents
}
