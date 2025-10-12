package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"golang.org/x/text/language"
)

type ChainSyncLostEvent struct {
	Duration  time.Duration
	timestamp time.Time
}

type ChainSyncLostTemplate struct {
	Duration time.Duration
}

func NewChainSyncLostEvent(duration time.Duration) *ChainSyncLostEvent {
	return &ChainSyncLostEvent{
		Duration:  duration,
		timestamp: time.Now(),
	}
}

func (e *ChainSyncLostEvent) Type() EventType {
	return Event_CHAIN_SYNC_LOST
}

func (e *ChainSyncLostEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *ChainSyncLostEvent) GetTemplateData(lang language.Tag) interface{} {
	return &ChainSyncLostTemplate{
		Duration: e.Duration,
	}
}

func (e *ChainSyncLostEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.ChainSyncEvents
}
