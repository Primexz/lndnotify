package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"golang.org/x/text/language"
)

type UnhealthyEvent struct {
	Err       error
	timestamp time.Time
}

type UnhealthyTemplate struct {
	Err string
}

func NewLndUnhealthyEvent(err error) *UnhealthyEvent {
	return &UnhealthyEvent{
		Err:       err,
		timestamp: time.Now(),
	}
}

func (e *UnhealthyEvent) Type() EventType {
	return Event_UNHEALTHY
}

func (e *UnhealthyEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *UnhealthyEvent) GetTemplateData(lang language.Tag) interface{} {
	return &UnhealthyTemplate{
		Err: e.Err.Error(),
	}
}

func (e *UnhealthyEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.HealthEvents
}
