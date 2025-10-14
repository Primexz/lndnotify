package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"
)

type ChannelStatusUpEvent struct {
	Channel   *lnrpc.Channel
	getAlias  func(pubKey string) string
	timestamp time.Time
}

type ChannelStatusUpTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
	ChannelPoint    string
	Capacity        string
}

func NewChannelStatusUpEvent(channel *lnrpc.Channel, getAlias func(pubKey string) string) *ChannelStatusUpEvent {
	return &ChannelStatusUpEvent{
		Channel:   channel,
		getAlias:  getAlias,
		timestamp: time.Now(),
	}
}

func (e *ChannelStatusUpEvent) Type() EventType {
	return Event_CHANNEL_STATUS_UP
}

func (e *ChannelStatusUpEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChannelStatusUpEvent) GetTemplateData(lang language.Tag) interface{} {
	remotePubkey := e.Channel.RemotePubkey

	return &ChannelClosingTemplate{
		PeerAlias:       e.getAlias(remotePubkey),
		PeerPubKey:      remotePubkey,
		PeerPubkeyShort: format.FormatPubKey(remotePubkey),
		ChannelPoint:    e.Channel.ChannelPoint,
		Capacity:        format.FormatBasic(float64(e.Channel.Capacity), lang),
	}
}

func (e *ChannelStatusUpEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChannelStatusEvents
}
