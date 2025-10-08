package events

import (
	"fmt"
	"math"
	"time"

	"github.com/Primexz/lndnotify/internal/config"
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
	FeeRate      string
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

	var feeRatePpm int32
	// Since the amounts in LN are rounded down, we have to round
	// commercially in order to reconstruct a correct ppm.
	if amtOutSats > 0 {
		feeRatePpm = int32(math.Round(feeSats * 1e6 / amtOutSats))
	}

	return &ForwardTemplate{
		PeerAliasIn:  e.Forward.PeerAliasIn,
		PeerAliasOut: e.Forward.PeerAliasOut,
		Amount:       format.FormatBasic(amtInSats),
		AmountOut:    format.FormatBasic(amtOutSats),
		Fee:          format.FormatDetailed(feeSats),
		FeeRate:      fmt.Sprintf("%d", feeRatePpm),
	}
}

func (e *ForwardEvent) ShouldProcess(cfg *config.Config) bool {
	if !cfg.Events.ForwardEvents {
		return false
	}
	return e.Forward.AmtOut >= cfg.EventConfig.ForwardEvent.MinAmount
}
