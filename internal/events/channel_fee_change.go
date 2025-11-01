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

	OldFeeRate           string
	NewFeeRate           string
	FeeRateChange        string
	FeeRateChangePercent string

	OldBaseFee           string
	NewBaseFee           string
	BaseFeeChange        string
	BaseFeeChangePercent string

	OldInboundFeeRate           string
	NewInboundFeeRate           string
	InboundFeeRateChange        string
	InboundFeeRateChangePercent string

	OldInboundBaseFee           string
	NewInboundBaseFee           string
	InboundBaseFeeChange        string
	InboundBaseFeeChangePercent string
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
	feeChange := e.FeeChange

	const millisatsToSats = 1000
	oldBaseFeeInSats := float64(feeChange.OldBaseFee) / 1000
	newBaseFeeInSats := float64(feeChange.NewBaseFee) / 1000
	oldInboundBaseFeeInSats := float64(feeChange.OldInboundBaseFee) / 1000
	newInboundBaseFeeInSats := float64(feeChange.NewInboundBaseFee) / 1000

	return &ChannelFeeChangeTemplate{
		PeerAlias:       e.getAlias(ch.RemotePubkey),
		PeerPubKey:      ch.RemotePubkey,
		PeerPubkeyShort: format.FormatPubKey(ch.RemotePubkey),
		ChannelPoint:    ch.ChannelPoint,
		Capacity:        format.FormatBasic(float64(ch.Capacity), lang),
		ChanId:          ch.ChanId,

		OldFeeRate:           format.FormatBasic(float64(feeChange.OldFeeRate), lang),
		NewFeeRate:           format.FormatBasic(float64(feeChange.NewFeeRate), lang),
		FeeRateChange:        format.CalculateAbsoluteChange(feeChange.OldFeeRate, feeChange.NewFeeRate),
		FeeRateChangePercent: format.CalculatePercentageChange(feeChange.OldFeeRate, feeChange.NewFeeRate),

		OldBaseFee:           format.FormatBasic(oldBaseFeeInSats, lang),
		NewBaseFee:           format.FormatBasic(newBaseFeeInSats, lang),
		BaseFeeChange:        format.CalculateAbsoluteChange(int64(oldBaseFeeInSats), int64(newBaseFeeInSats)),
		BaseFeeChangePercent: format.CalculatePercentageChange(int64(oldBaseFeeInSats), int64(newBaseFeeInSats)),

		OldInboundFeeRate:           format.FormatBasic(float64(feeChange.OldInboundFeeRate), lang),
		NewInboundFeeRate:           format.FormatBasic(float64(feeChange.NewInboundFeeRate), lang),
		InboundFeeRateChange:        format.CalculateAbsoluteChange(int64(feeChange.OldInboundFeeRate), int64(feeChange.NewInboundFeeRate)),
		InboundFeeRateChangePercent: format.CalculatePercentageChange(int64(feeChange.OldInboundFeeRate), int64(feeChange.NewInboundFeeRate)),

		OldInboundBaseFee:           format.FormatBasic(oldInboundBaseFeeInSats, lang),
		NewInboundBaseFee:           format.FormatBasic(newInboundBaseFeeInSats, lang),
		InboundBaseFeeChange:        format.CalculateAbsoluteChange(int64(oldInboundBaseFeeInSats), int64(newInboundBaseFeeInSats)),
		InboundBaseFeeChangePercent: format.CalculatePercentageChange(int64(oldInboundBaseFeeInSats), int64(newInboundBaseFeeInSats)),
	}
}

func (e *ChannelFeeChangeEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChannelFeeEvents
}
