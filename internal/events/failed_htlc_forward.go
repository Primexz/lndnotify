package events

import (
	"time"

	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
)

type FailedHtlcForwardEvent struct {
	HtlcEvent *routerrpc.HtlcEvent
	timestamp time.Time
}

type FailedHtlcForwardTemplate struct {
}

func NewFailedHtlcForwardEvent(htlcEvent *routerrpc.HtlcEvent) *FailedHtlcForwardEvent {
	return &FailedHtlcForwardEvent{
		HtlcEvent: htlcEvent,
		timestamp: time.Now(),
	}
}

func (e *FailedHtlcForwardEvent) Type() EventType {
	return Event_FAILED_HTLC
}

func (e *FailedHtlcForwardEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *FailedHtlcForwardEvent) GetTemplateData() interface{} {
	return &FailedHtlcForwardTemplate{}
}
