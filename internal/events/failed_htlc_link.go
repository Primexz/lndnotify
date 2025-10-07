package events

import (
	"time"

	"github.com/Primexz/lndnotify/pkg/channel"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	log "github.com/sirupsen/logrus"
)

type FailedHtlcLinkEvent struct {
	HtlcEvent *routerrpc.HtlcEvent
	FailEvent *routerrpc.LinkFailEvent
	Channels  []*lnrpc.Channel
	timestamp time.Time
}

type FailedHtlcLinkTemplate struct {
	OutChanId               uint64
	InChanId                uint64
	InChanAlias             string
	OutChanAlias            string
	OutChanLiquidity        string
	MissingOutChanLiquidity string
	IsLocalLiquidityFailure bool
	Amount                  string
	WireFailure             string
	FailureDetail           string
	MissedFee               string
}

func NewFailedHtlcLinkEvent(htlcEvent *routerrpc.HtlcEvent, failEvent *routerrpc.LinkFailEvent, channels []*lnrpc.Channel) *FailedHtlcLinkEvent {
	return &FailedHtlcLinkEvent{
		HtlcEvent: htlcEvent,
		FailEvent: failEvent,
		Channels:  channels,
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

	if outChan := channel.GetChannelById(e.Channels, outChanId); outChan != nil {
		outChanAlias = outChan.PeerAlias
		outChanLiquidity = outChan.GetLocalBalance() - outChan.GetLocalChanReserveSat()
	} else {
		log.WithField("chan_id", outChanId).Warn("could not find outgoing channel")
	}

	if inChan := channel.GetChannelById(e.Channels, inChanId); inChan != nil {
		inChanAlias = inChan.PeerAlias
	} else {
		log.WithField("chan_id", inChanId).Warn("could not find incoming channel")
	}

	return &FailedHtlcLinkTemplate{
		InChanId:                inChanId,
		OutChanId:               outChanId,
		InChanAlias:             inChanAlias,
		OutChanAlias:            outChanAlias,
		OutChanLiquidity:        format.FormatBasic(float64(outChanLiquidity)),
		MissingOutChanLiquidity: format.FormatBasic(float64(failInfo.GetOutgoingAmtMsat())/1000 - float64(outChanLiquidity)),
		IsLocalLiquidityFailure: float64(failInfo.GetOutgoingAmtMsat()/1000) > float64(outChanLiquidity),
		Amount:                  format.FormatBasic(float64(failInfo.GetOutgoingAmtMsat()) / 1000),
		WireFailure:             e.FailEvent.GetWireFailure().String(),
		FailureDetail:           e.FailEvent.GetFailureDetail().String(),
		MissedFee:               format.FormatDetailed((float64(failInfo.GetIncomingAmtMsat() - failInfo.GetOutgoingAmtMsat())) / 1000),
	}
}
