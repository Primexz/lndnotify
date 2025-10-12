package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"golang.org/x/text/language"
)

type ChainSyncRestoredEvent struct {
	Duration  time.Duration
	timestamp time.Time
}

type ChainSyncRestoredTemplate struct {
	Duration time.Duration
}

func NewChainSyncRestoredEvent(duration time.Duration) *ChainSyncRestoredEvent {
	return &ChainSyncRestoredEvent{
		Duration:  duration,
		timestamp: time.Now(),
	}
}

func (e *ChainSyncRestoredEvent) Type() EventType {
	return Event_CHAIN_SYNC_RESTORED
}

func (e *ChainSyncRestoredEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChainSyncRestoredEvent) GetTemplateData(lang language.Tag) interface{} {
	return &ChainSyncRestoredTemplate{
		Duration: e.Duration,
	}
}

func (e *ChainSyncRestoredEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChainSyncEvents
}
