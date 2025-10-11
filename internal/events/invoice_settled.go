package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type InvoiceSettledEvent struct {
	Invoice   *lnrpc.Invoice
	timestamp time.Time
}

type InvoiceSettledTemplate struct {
	Memo           string
	Value          string
	IsKeysend      bool
	PaymentRequest string
}

func NewInvoiceSettledEvent(invoice *lnrpc.Invoice) *InvoiceSettledEvent {
	return &InvoiceSettledEvent{
		Invoice:   invoice,
		timestamp: time.Now(),
	}
}

func (e *InvoiceSettledEvent) Type() EventType {
	return Event_INVOICE_SETTLED
}

func (e *InvoiceSettledEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *InvoiceSettledEvent) GetTemplateData(cfg *config.Config) interface{} {
	langTag := cfg.Formatting.Language.Tag

	return &InvoiceSettledTemplate{
		Memo:           e.Invoice.Memo,
		Value:          format.FormatBasic(float64(e.Invoice.Value), langTag),
		IsKeysend:      e.Invoice.IsKeysend,
		PaymentRequest: e.Invoice.PaymentRequest,
	}
}

func (e *InvoiceSettledEvent) ShouldProcess(cfg *config.Config) bool {
	if !cfg.Events.InvoiceEvents {
		return false
	}
	return uint64(e.Invoice.Value) >= cfg.EventConfig.InvoiceEvent.MinAmount // #nosec G115
}
