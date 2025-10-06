package events

import (
	"time"

	"github.com/Primexz/lndnotify/pkg/channel"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
)

type FailedHtlcLinkEvent struct {
	HtlcEvent *routerrpc.HtlcEvent
	FailEvent *routerrpc.LinkFailEvent
	Channels  []*lnrpc.Channel
	timestamp time.Time
}

func NewFailedHtlcLinkEvent(htlcEvent *routerrpc.HtlcEvent, failEvent *routerrpc.LinkFailEvent, channels []*lnrpc.Channel) *FailedHtlcLinkEvent {
	return &FailedHtlcLinkEvent{
		HtlcEvent: htlcEvent,
		FailEvent: failEvent,
		timestamp: time.Now(),
	}
}

func (e *FailedHtlcLinkEvent) Type() EventType {
	return Event_FAILED_HTLC
}

func (e *FailedHtlcLinkEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *FailedHtlcLinkEvent) GetTemplateData() interface{} {
	failInfo := e.FailEvent.GetInfo()
	inChanId := e.HtlcEvent.GetIncomingChannelId()
	outChanId := e.HtlcEvent.GetOutgoingChannelId()
	inChanAlias := "unknown"
	outChanAlias := "unknown"
	outChanLiquidity := int64(0)
	amount := float64(failInfo.GetOutgoingAmtMsat()) / 1000
	wireFailure := e.FailEvent.GetWireFailure().String()
	failureDetail := e.FailEvent.GetFailureDetail().String()
	missedFee := (float64(failInfo.GetIncomingAmtMsat() - failInfo.GetOutgoingAmtMsat())) / 1000

	if outChan := channel.GetChannelById(e.Channels, outChanId); outChan != nil {
		outChanAlias = outChan.PeerAlias
		outChanLiquidity = outChan.GetLocalBalance() - outChan.GetLocalChanReserveSat()
	}

	if inChan := channel.GetChannelById(e.Channels, inChanId); inChan != nil {
		inChanAlias = inChan.PeerAlias
	}

	return &FailedHtlcTemplate{
		InChanId:         inChanId,
		OutChanId:        outChanId,
		InChanAlias:      inChanAlias,
		OutChanAlias:     outChanAlias,
		OutChanLiquidity: format.FormatSats(float64(outChanLiquidity)),
		Amount:           format.FormatSats(amount),
		WireFailure:      wireFailure,
		FailureDetail:    failureDetail,
		MissedFee:        format.FormatSats(missedFee),
	}
}
