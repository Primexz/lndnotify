package events

import (
	"sync"
	"time"

	"github.com/Primexz/lnd-notify/internal/config"
)

// ProcessorConfig holds the configuration for event processing
type ProcessorConfig struct {
	EnabledEvents config.EventConfig
	RateLimit     config.RateLimitConfig
}

// Processor handles event processing and filtering
type Processor struct {
	cfg       *ProcessorConfig
	eventChan chan Event
	batchMu   sync.Mutex
	batches   map[string][]Event
	lastFlush time.Time
	running   bool
	runningMu sync.Mutex
}

// NewProcessor creates a new event processor
func NewProcessor(cfg *ProcessorConfig) *Processor {
	return &Processor{
		cfg:       cfg,
		eventChan: make(chan Event, 100),
		batches:   make(map[string][]Event),
		lastFlush: time.Now(),
	}
}

// Start begins processing events
func (p *Processor) Start() error {
	p.runningMu.Lock()
	defer p.runningMu.Unlock()

	if p.running {
		return nil
	}

	p.running = true
	go p.processBatches()
	return nil
}

// Stop halts event processing
func (p *Processor) Stop() {
	p.runningMu.Lock()
	defer p.runningMu.Unlock()

	if !p.running {
		return
	}

	p.running = false
	p.flushBatches()
}

// ProcessEvent handles a single event
func (p *Processor) ProcessEvent(event Event) error {
	if event == nil {
		return nil
	}

	if !p.shouldProcess(event) {
		return nil
	}

	if p.cfg.RateLimit.BatchSimilarEvents {
		p.batchEvent(event)
	} else {
		p.eventChan <- event
	}

	return nil
}

// shouldProcess checks if an event type is enabled
func (p *Processor) shouldProcess(event Event) bool {
	switch event.Type() {
	case "forward_event":
		return p.cfg.EnabledEvents.ForwardEvents
	case "peer_event":
		return p.cfg.EnabledEvents.PeerEvents
	default:
		return false
	}
}

// batchEvent adds an event to its type batch
func (p *Processor) batchEvent(event Event) {
	p.batchMu.Lock()
	defer p.batchMu.Unlock()

	eventType := event.Type()
	p.batches[eventType] = append(p.batches[eventType], event)

	if time.Since(p.lastFlush) >= time.Duration(p.cfg.RateLimit.BatchWindowSeconds)*time.Second {
		p.flushBatchesLocked()
	}
}

// processBatches periodically flushes event batches
func (p *Processor) processBatches() {
	ticker := time.NewTicker(time.Duration(p.cfg.RateLimit.BatchWindowSeconds) * time.Second)
	defer ticker.Stop()

	for {
		p.runningMu.Lock()
		if !p.running {
			p.runningMu.Unlock()
			return
		}
		p.runningMu.Unlock()

		<-ticker.C
		p.flushBatches()
	}
}

// flushBatches sends all batched events
func (p *Processor) flushBatches() {
	p.batchMu.Lock()
	defer p.batchMu.Unlock()
	p.flushBatchesLocked()
}

// flushBatchesLocked sends all batched events (must hold batchMu)
func (p *Processor) flushBatchesLocked() {
	for eventType, events := range p.batches {
		if len(events) > 0 {
			// Send batched events
			for _, event := range events {
				p.eventChan <- event
			}
			// Clear batch
			delete(p.batches, eventType)
		}
	}
	p.lastFlush = time.Now()
}
