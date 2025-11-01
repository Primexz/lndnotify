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
	OldFeeRate      string
	NewFeeRate      string
	OldBaseFee      string
	NewBaseFee      string
	FeeRateChange   string // "increased" or "decreased"
	BaseFeeChange   string // "increased" or "decreased"
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

	feeRateChange := "unchanged"
	if e.FeeChange.NewFeeRate > e.FeeChange.OldFeeRate {
		feeRateChange = "increased"
	} else if e.FeeChange.NewFeeRate < e.FeeChange.OldFeeRate {
		feeRateChange = "decreased"
	}

	baseFeeChange := "unchanged"
	if e.FeeChange.NewBaseFee > e.FeeChange.OldBaseFee {
		baseFeeChange = "increased"
	} else if e.FeeChange.NewBaseFee < e.FeeChange.OldBaseFee {
		baseFeeChange = "decreased"
	}

	return &ChannelFeeChangeTemplate{
		PeerAlias:       e.getAlias(ch.RemotePubkey),
		PeerPubKey:      ch.RemotePubkey,
		PeerPubkeyShort: format.FormatPubKey(ch.RemotePubkey),
		ChannelPoint:    ch.ChannelPoint,
		Capacity:        format.FormatBasic(float64(ch.Capacity), lang),
		ChanId:          ch.ChanId,
		OldFeeRate:      format.FormatBasic(float64(e.FeeChange.OldFeeRate)/1000, lang),
		NewFeeRate:      format.FormatBasic(float64(e.FeeChange.NewFeeRate)/1000, lang),
		OldBaseFee:      format.FormatBasic(float64(e.FeeChange.OldBaseFee)/1000, lang),
		NewBaseFee:      format.FormatBasic(float64(e.FeeChange.NewBaseFee)/1000, lang),
		FeeRateChange:   feeRateChange,
		BaseFeeChange:   baseFeeChange,
	}
}

func (e *ChannelFeeChangeEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChannelFeeEvents
}
