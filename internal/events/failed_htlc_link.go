package events

import (
	"time"

	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
)

type FailedHtlcLinkEvent struct {
	HtlcEvent *routerrpc.HtlcEvent
	timestamp time.Time
}

type FailedHtlcLinkTemplate struct {
}

func NewFailedHtlcEvent(htlcEvent *routerrpc.HtlcEvent) *FailedHtlcLinkEvent {
	return &FailedHtlcLinkEvent{
		HtlcEvent: htlcEvent,
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
	return &FailedHtlcLinkTemplate{}
}
