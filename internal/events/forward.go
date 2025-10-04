package events

import (
	"fmt"
	"time"

	"github.com/Primexz/lnd-notify/pkg/format"
)

type ForwardEvent struct {
	PeerAliasIn  string
	PeerAliasOut string
	AmtInMsat    uint64
	AmtOutMsat   uint64
	FeeMsat      uint64

	timestamp time.Time
}

func NewForwardEvent(peerAliasIn, peerAliasOut string, amtIn, amtOut, fee uint64) *ForwardEvent {
	return &ForwardEvent{
		PeerAliasIn:  peerAliasIn,
		PeerAliasOut: peerAliasOut,
		AmtInMsat:    amtIn,
		AmtOutMsat:   amtOut,
		FeeMsat:      fee,
		timestamp:    time.Now(),
	}
}

func (e *ForwardEvent) Type() string {
	return "forward_event"
}

func (e *ForwardEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ForwardEvent) ToNotification() string {
	amtSats := float64(e.AmtInMsat) / 1000
	feeSats := float64(e.FeeMsat) / 1000

	return fmt.Sprintf("💰 Forwarded %s sats %s -> %s, earning %s sats fee", format.FormatSats(amtSats), e.PeerAliasIn, e.PeerAliasOut, format.FormatSats(feeSats))
}
