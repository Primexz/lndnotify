package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/chainutil"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"
)

type HTLCExpirationEvent struct {
	timestamp       time.Time
	htlc            *lnrpc.HTLC
	channel         *lnrpc.Channel
	remainingBlocks int32
}

type HTLCExpirationTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
	ChannelPoint    string
	HTLCAmount      string
	RemainingBlocks int32
	RemainingTime   time.Duration
}

func NewHTLCExpirationEvent(htlc *lnrpc.HTLC, channel *lnrpc.Channel, remainingBlocks int32) *HTLCExpirationEvent {
	return &HTLCExpirationEvent{
		htlc:            htlc,
		channel:         channel,
		remainingBlocks: remainingBlocks,
		timestamp:       time.Now(),
	}
}

func (e *HTLCExpirationEvent) Type() EventType {
	return Event_HTLC_EXPIRATION
}

func (e *HTLCExpirationEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *HTLCExpirationEvent) GetTemplateData(lang language.Tag) interface{} {
	return &HTLCExpirationTemplate{
		PeerAlias:       e.channel.PeerAlias,
		PeerPubKey:      e.channel.RemotePubkey,
		PeerPubkeyShort: format.FormatPubKey(e.channel.RemotePubkey),
		ChannelPoint:    e.channel.ChannelPoint,
		HTLCAmount:      format.FormatBasic(float64(e.htlc.Amount), lang),
		RemainingBlocks: e.remainingBlocks,
		RemainingTime:   format.FormatDuration(chainutil.BlockCountToDuration(e.remainingBlocks)),
	}
}

func (e *HTLCExpirationEvent) ShouldProcess(cfg *config.Config) bool {
	if !cfg.Events.HTLCExpirationEvents {
		return false
	}

	return e.remainingBlocks <= cfg.EventConfig.HTLCExpirationEvent.RemainingBlocks
}
