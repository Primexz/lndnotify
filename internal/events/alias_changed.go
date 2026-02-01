package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"golang.org/x/text/language"
)

type AliasChangedEvent struct {
	timestamp time.Time
	oldAlias  string
	newAlias  string
}

type AliasChangedTemplate struct {
	OldAlias string
	NewAlias string
}

func NewAliasChangedEvent(pubkey string, oldAlias string, newAlias string) *AliasChangedEvent {
	return &AliasChangedEvent{
		timestamp: time.Now(),
		oldAlias:  oldAlias,
		newAlias:  newAlias,
	}
}

func (e *AliasChangedEvent) Type() EventType {
	return Event_ALIAS_CHANGED
}

func (e *AliasChangedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *AliasChangedEvent) GetTemplateData(lang language.Tag) interface{} {
	return &AliasChangedTemplate{
		OldAlias: e.oldAlias,
		NewAlias: e.newAlias,
	}
}

func (e *AliasChangedEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.AliasChangedEvents
}
