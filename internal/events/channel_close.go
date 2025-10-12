package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"
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
	CloseInitiator  bool
	CloseType       int32
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

func (e *ChannelCloseEvent) GetTemplateData(lang language.Tag) interface{} {
	return &ChannelCloseTemplate{
		PeerAlias:       e.Node.Alias,
		PeerPubKey:      e.Node.PubKey,
		PeerPubkeyShort: format.FormatPubKey(e.Node.PubKey),
		ChanId:          e.Channel.ChanId,
		ChannelPoint:    e.Channel.ChannelPoint,
		Capacity:        format.FormatBasic(float64(e.Channel.Capacity), lang),
		RemotePubkey:    e.Channel.RemotePubkey,
		SettledBalance:  format.FormatBasic(float64(e.Channel.SettledBalance), lang),
		CloseInitiator:  e.Channel.CloseInitiator == lnrpc.Initiator_INITIATOR_LOCAL,
		CloseType:       int32(e.Channel.CloseType),
	}
}

func (e *ChannelCloseEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChannelEvents
}
