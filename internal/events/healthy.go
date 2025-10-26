package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"golang.org/x/text/language"
)

type HealthyEvent struct {
	timestamp time.Time
}

type HealthyTemplate struct {
}

func NewLndHealthyEvent() *HealthyEvent {
	return &HealthyEvent{
		timestamp: time.Now(),
	}
}

func (e *HealthyEvent) Type() EventType {
	return Event_HEALTHY
}

func (e *HealthyEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *HealthyEvent) GetTemplateData(lang language.Tag) interface{} {
	return &HealthyTemplate{}
}

func (e *HealthyEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.HealthEvents
}
