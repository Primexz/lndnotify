package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type ChannelClosingEvent struct {
	Channel   *lnrpc.PendingChannelsResponse_WaitingCloseChannel
	getAlias  func(pubKey string) string
	timestamp time.Time
}

type ChannelClosingTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
	ChannelPoint    string
	Capacity        string
	LimboBalance    string
	ClosingTxid     string
	ClosingTxHex    string
}

func NewChannelClosingEvent(channel *lnrpc.PendingChannelsResponse_WaitingCloseChannel, getAlias func(pubKey string) string) *ChannelClosingEvent {
	return &ChannelClosingEvent{
		Channel:   channel,
		getAlias:  getAlias,
		timestamp: time.Now(),
	}
}

func (e *ChannelClosingEvent) Type() EventType {
	return Event_CHANNEL_CLOSING
}

func (e *ChannelClosingEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChannelClosingEvent) GetTemplateData() interface{} {
	remotePubkey := e.Channel.Channel.RemoteNodePub

	return &ChannelClosingTemplate{
		PeerAlias:       e.getAlias(remotePubkey),
		PeerPubKey:      remotePubkey,
		PeerPubkeyShort: format.FormatPubKey(remotePubkey),
		ChannelPoint:    e.Channel.Channel.ChannelPoint,
		Capacity:        format.FormatBasic(float64(e.Channel.Channel.Capacity)),
		LimboBalance:    format.FormatBasic(float64(e.Channel.LimboBalance)),
		ClosingTxid:     e.Channel.ClosingTxid,
		ClosingTxHex:    e.Channel.ClosingTxHex,
	}
}

func (e *ChannelClosingEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChannelEvents
}
