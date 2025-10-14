package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"
)

type ChannelStatusDownEvent struct {
	Channel   *lnrpc.Channel
	getAlias  func(pubKey string) string
	timestamp time.Time
}

type ChannelStatusDownEventTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
	ChannelPoint    string
	Capacity        string
}

func NewChannelStatusDownEvent(channel *lnrpc.Channel, getAlias func(pubKey string) string) *ChannelStatusDownEvent {
	return &ChannelStatusDownEvent{
		Channel:   channel,
		getAlias:  getAlias,
		timestamp: time.Now(),
	}
}

func (e *ChannelStatusDownEvent) Type() EventType {
	return Event_CHANNEL_STATUS_DOWN
}

func (e *ChannelStatusDownEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChannelStatusDownEvent) GetTemplateData(lang language.Tag) interface{} {
	remotePubkey := e.Channel.RemotePubkey

	return &ChannelClosingTemplate{
		PeerAlias:       e.getAlias(remotePubkey),
		PeerPubKey:      remotePubkey,
		PeerPubkeyShort: format.FormatPubKey(remotePubkey),
		ChannelPoint:    e.Channel.ChannelPoint,
		Capacity:        format.FormatBasic(float64(e.Channel.Capacity), lang),
	}
}

func (e *ChannelStatusDownEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChannelStatusEvents
}
