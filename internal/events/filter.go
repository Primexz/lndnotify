package events

import (
	"github.com/Primexz/lndnotify/internal/config"
	log "github.com/sirupsen/logrus"
)

type ProcessorConfig struct {
	EnabledEvents config.EventConfig
}

type Processor struct {
	cfg *ProcessorConfig
}

// NewProcessor creates a new event processor
func NewProcessor(cfg *ProcessorConfig) *Processor {
	return &Processor{
		cfg: cfg,
	}
}

// shouldProcess checks if an event type is enabled
func (p *Processor) ShouldProcess(event Event) bool {
	switch event.Type() {
	case Event_FORWARD:
		return p.cfg.EnabledEvents.ForwardEvents
	case Event_PEER_ONLINE, Event_PEER_OFFLINE:
		return p.cfg.EnabledEvents.PeerEvents
	case Event_CHANNEL_OPEN, Event_CHANNEL_CLOSE:
		return p.cfg.EnabledEvents.ChannelEvents
	case Event_INVOICE_SETTLED:
		return p.cfg.EnabledEvents.InvoiceEvents
	case Event_FAILED_HTLC:
		return p.cfg.EnabledEvents.FailedHtlc
	case Event_KEYSEND:
		return p.cfg.EnabledEvents.KeysendEvents
	default:
		log.WithField("event_type", event.Type()).Warn("unknown event type")
		return false
	}
}
