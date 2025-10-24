package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"
)

type WalletStateEvent struct {
	OldState  lnrpc.WalletState
	NewState  lnrpc.WalletState
	timestamp time.Time
}

type WalletStateTemplate struct {
	OldState string
	NewState string
}

func NewWalletStateEvent(oldState lnrpc.WalletState, newState lnrpc.WalletState) *WalletStateEvent {
	return &WalletStateEvent{
		OldState:  oldState,
		NewState:  newState,
		timestamp: time.Now(),
	}
}

func (e *WalletStateEvent) Type() EventType {
	return Event_WALLET_STATE
}

func (e *WalletStateEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *WalletStateEvent) GetTemplateData(lang language.Tag) interface{} {
	return &WalletStateTemplate{
		OldState: e.OldState.String(),
		NewState: e.NewState.String(),
	}
}

func (e *WalletStateEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.WalletStateEvents
}
