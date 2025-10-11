package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"
)

type OnChainTransactionEvent struct {
	Event     *lnrpc.Transaction
	timestamp time.Time
}

type OnChainTransactionTemplate struct {
	TxHash    string
	RawTxHex  string
	Amount    string
	TotalFees string
	Confirmed bool
	Outputs   []OnChainOutput
}

type OnChainOutput struct {
	Amount       string
	Address      string
	OutputType   string
	IsOurAddress bool
}

func NewOnChainTransactionEvent(event *lnrpc.Transaction) *OnChainTransactionEvent {
	return &OnChainTransactionEvent{
		Event:     event,
		timestamp: time.Now(),
	}
}

func (e *OnChainTransactionEvent) Type() EventType {
	if e.Event.GetNumConfirmations() > 0 {
		return Event_ONCHAIN_CONFIRMED
	}

	return Event_ONCHAIN_MEMPOOL
}

func (e *OnChainTransactionEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *OnChainTransactionEvent) GetTemplateData(langTag language.Tag) interface{} {
	outputs := make([]OnChainOutput, 0, len(e.Event.OutputDetails))
	for _, output := range e.Event.OutputDetails {
		outputs = append(outputs, OnChainOutput{
			Amount:       format.FormatBasic(float64(output.Amount), langTag),
			OutputType:   output.OutputType.String(),
			IsOurAddress: output.IsOurAddress,
			Address:      output.Address,
		})
	}

	return &OnChainTransactionTemplate{
		TxHash:    e.Event.TxHash,
		RawTxHex:  e.Event.RawTxHex,
		Outputs:   outputs,
		Amount:    format.FormatBasic(float64(e.Event.Amount), langTag),
		TotalFees: format.FormatDetailed(float64(e.Event.TotalFees), langTag),
		Confirmed: e.Event.NumConfirmations > 0,
	}
}

func (e *OnChainTransactionEvent) ShouldProcess(cfg *config.Config) bool {
	if !cfg.Events.OnChainEvents {
		return false
	}

	return uint64(e.Event.Amount) >= cfg.EventConfig.OnChainEvent.MinAmount // #nosec G115
}
