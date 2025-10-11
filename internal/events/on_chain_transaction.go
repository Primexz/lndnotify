package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
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
	return Event_ONCHAIN
}

func (e *OnChainTransactionEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *OnChainTransactionEvent) GetTemplateData() interface{} {
	outputs := make([]OnChainOutput, 0, len(e.Event.OutputDetails))
	for _, output := range e.Event.OutputDetails {
		outputs = append(outputs, OnChainOutput{
			Amount:       format.FormatBasic(float64(output.Amount)),
			OutputType:   output.OutputType.String(),
			IsOurAddress: output.IsOurAddress,
			Address:      output.Address,
		})
	}

	return &OnChainTransactionTemplate{
		TxHash:    e.Event.TxHash,
		RawTxHex:  e.Event.RawTxHex,
		Outputs:   outputs,
		Amount:    format.FormatBasic(float64(e.Event.Amount)),
		TotalFees: format.FormatDetailed(float64(e.Event.TotalFees)),
	}
}

func (e *OnChainTransactionEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.OnChainEvents
}
