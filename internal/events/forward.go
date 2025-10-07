package events

import (
	"time"

	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type ForwardEvent struct {
	Forward   *lnrpc.ForwardingEvent
	timestamp time.Time
}

type ForwardTemplate struct {
	PeerAliasIn  string
	PeerAliasOut string
	Amount       string
	AmountOut    string
	Fee          string
}

func NewForwardEvent(forward *lnrpc.ForwardingEvent) *ForwardEvent {
	return &ForwardEvent{
		Forward:   forward,
		timestamp: time.Now(),
	}
}

func (e *ForwardEvent) Type() EventType {
	return Event_FORWARD
}

func (e *ForwardEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ForwardEvent) GetTemplateData() interface{} {
	amtInSats := float64(e.Forward.AmtInMsat) / 1000
	amtOutSats := float64(e.Forward.AmtOutMsat) / 1000
	feeSats := float64(e.Forward.FeeMsat) / 1000

	return &ForwardTemplate{
		PeerAliasIn:  e.Forward.PeerAliasIn,
		PeerAliasOut: e.Forward.PeerAliasOut,
		Amount:       format.FormatBasic(amtInSats),
		AmountOut:    format.FormatBasic(amtOutSats),
		Fee:          format.FormatDetailed(feeSats),
	}
}
