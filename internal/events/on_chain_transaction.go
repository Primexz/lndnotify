package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type OnChainTransactionEvent struct {
	Event     *lnrpc.Transaction
	timestamp time.Time
}

type OnChainTransactionTemplate struct {
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

	return &OnChainTransactionTemplate{}
}

func (e *OnChainTransactionEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.OnChainEvents
}
