package events

import (
	"time"

	channelmanager "github.com/Primexz/lndnotify/internal/channel_manager"
	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"golang.org/x/text/language"
)

type ChannelFeeChangeEvent struct {
	FeeChange channelmanager.FeeChangeEvent
	getAlias  func(pubKey string) string
	timestamp time.Time
}

type ChannelFeeChangeTemplate struct {
	PeerAlias       string
	PeerPubKey      string
	PeerPubkeyShort string
	ChannelPoint    string
	Capacity        string
	ChanId          uint64

	OldFeeRate string
	NewFeeRate string

	OldBaseFee string
	NewBaseFee string

	OldInboundFeeRate string
	NewInboundFeeRate string

	OldInboundBaseFee string
	NewInboundBaseFee string
}

func NewChannelFeeChangeEvent(feeChange channelmanager.FeeChangeEvent, getAlias func(pubKey string) string) *ChannelFeeChangeEvent {
	return &ChannelFeeChangeEvent{
		FeeChange: feeChange,
		getAlias:  getAlias,
		timestamp: time.Now(),
	}
}

func (e *ChannelFeeChangeEvent) Type() EventType {
	return Event_CHANNEL_FEE_CHANGE
}

func (e *ChannelFeeChangeEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChannelFeeChangeEvent) GetTemplateData(lang language.Tag) interface{} {
	ch := e.FeeChange.Channel

	return &ChannelFeeChangeTemplate{
		PeerAlias:         e.getAlias(ch.RemotePubkey),
		PeerPubKey:        ch.RemotePubkey,
		PeerPubkeyShort:   format.FormatPubKey(ch.RemotePubkey),
		ChannelPoint:      ch.ChannelPoint,
		Capacity:          format.FormatBasic(float64(ch.Capacity), lang),
		ChanId:            ch.ChanId,
		OldFeeRate:        format.FormatBasic(float64(e.FeeChange.OldFeeRate), lang),
		NewFeeRate:        format.FormatBasic(float64(e.FeeChange.NewFeeRate), lang),
		OldBaseFee:        format.FormatBasic(float64(e.FeeChange.OldBaseFee)/1000, lang),
		NewBaseFee:        format.FormatBasic(float64(e.FeeChange.NewBaseFee)/1000, lang),
		OldInboundFeeRate: format.FormatBasic(float64(e.FeeChange.OldInboundFeeRate), lang),
		NewInboundFeeRate: format.FormatBasic(float64(e.FeeChange.NewInboundFeeRate), lang),
		OldInboundBaseFee: format.FormatBasic(float64(e.FeeChange.OldInboundBaseFee)/1000, lang),
		NewInboundBaseFee: format.FormatBasic(float64(e.FeeChange.NewInboundBaseFee)/1000, lang),
	}
}

func (e *ChannelFeeChangeEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChannelFeeEvents
}
