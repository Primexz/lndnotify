package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"golang.org/x/text/language"
)

type TLSCertExpiryEvent struct {
	ExpiryDate time.Time
	timestamp  time.Time
}

type TLSEventTemplate struct {
	ExpiryDate      time.Time
	TimeUntilExpiry time.Duration
}

func NewTLSCertExpiryEvent(expiryDate time.Time) *TLSCertExpiryEvent {
	return &TLSCertExpiryEvent{
		ExpiryDate: expiryDate,
		timestamp:  time.Now(),
	}
}

func (e *TLSCertExpiryEvent) Type() EventType {
	return Event_TLS_CERT_EXPIRY
}

func (e *TLSCertExpiryEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *TLSCertExpiryEvent) GetTemplateData(lang language.Tag) interface{} {
	return &TLSEventTemplate{
		ExpiryDate:      e.ExpiryDate,
		TimeUntilExpiry: time.Until(e.ExpiryDate),
	}
}

func (e *TLSCertExpiryEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.TLSCertExpiryEvents
}
