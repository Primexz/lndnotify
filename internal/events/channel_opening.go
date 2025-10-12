package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"
)

type ChannelOpeningEvent struct {
	Channel   *lnrpc.PendingChannelsResponse_PendingOpenChannel
	getAlias  func(pubKey string) string
	timestamp time.Time
}

type ChannelOpeningTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
	ChannelPoint    string
	Capacity        string
	Initiator       bool
	IsPrivate       bool
}

func NewChannelOpeningEvent(channel *lnrpc.PendingChannelsResponse_PendingOpenChannel,
	getAlias func(pubKey string) string) *ChannelOpeningEvent {

	return &ChannelOpeningEvent{
		Channel:   channel,
		getAlias:  getAlias,
		timestamp: time.Now(),
	}
}

func (e *ChannelOpeningEvent) Type() EventType {
	return Event_CHANNEL_OPENING
}

func (e *ChannelOpeningEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChannelOpeningEvent) GetTemplateData(lang language.Tag) interface{} {
	remotePubkey := e.Channel.Channel.RemoteNodePub
	initiator := e.Channel.Channel.Initiator == lnrpc.Initiator_INITIATOR_LOCAL

	return &ChannelOpeningTemplate{
		PeerAlias:       e.getAlias(remotePubkey),
		PeerPubKey:      remotePubkey,
		PeerPubkeyShort: format.FormatPubKey(remotePubkey),
		ChannelPoint:    e.Channel.Channel.ChannelPoint,
		Capacity:        format.FormatBasic(float64(e.Channel.Channel.Capacity), lang),
		Initiator:       initiator,
		IsPrivate:       e.Channel.Channel.Private,
	}
}

func (e *ChannelOpeningEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChannelEvents
}
