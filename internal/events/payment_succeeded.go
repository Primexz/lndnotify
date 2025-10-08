package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type PaymentSucceededEvent struct {
	Payment     *lnrpc.Payment
	PayReq      *lnrpc.PayReq
	RecHopAlias string
	timestamp   time.Time
}

type PaymentSucceededTemplate struct {
	PaymentHash string
	Value       string
	Fee         string
	RecHopAlias string
	Memo        string
}

func NewPaymentSucceededEvent(payment *lnrpc.Payment, payReq *lnrpc.PayReq,
	recHopAlias string) *PaymentSucceededEvent {
	return &PaymentSucceededEvent{
		Payment:     payment,
		PayReq:      payReq,
		RecHopAlias: recHopAlias,
		timestamp:   time.Now(),
	}
}

func (e *PaymentSucceededEvent) Type() EventType {
	return Event_PAYMENT_SUCCEEDED
}

func (e *PaymentSucceededEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *PaymentSucceededEvent) GetTemplateData() interface{} {
	valueSats := float64(e.Payment.ValueMsat) / 1000
	feeSats := float64(e.Payment.FeeMsat) / 1000

	// Get memo from PayReq
	var memo string
	if e.PayReq != nil {
		memo = e.PayReq.Description
	}

	return &PaymentSucceededTemplate{
		PaymentHash: e.Payment.PaymentHash,
		Value:       format.FormatBasic(valueSats),
		Fee:         format.FormatDetailed(feeSats),
		RecHopAlias: e.RecHopAlias,
		Memo:        memo,
	}
}

func (e *PaymentSucceededEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.PaymentEvents
}
