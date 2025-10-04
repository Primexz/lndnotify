package events

import (
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

type ForwardTemplate struct {
	PeerAliasIn  string
	PeerAliasOut string
	Amount       string
	Fee          string
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

func (e *ForwardEvent) GetTemplateData() interface{} {
	amtSats := float64(e.AmtInMsat) / 1000
	feeSats := float64(e.FeeMsat) / 1000

	return &ForwardTemplate{
		PeerAliasIn:  e.PeerAliasIn,
		PeerAliasOut: e.PeerAliasOut,
		Amount:       format.FormatSats(amtSats),
		Fee:          format.FormatSats(feeSats),
	}
}
