package notify

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/internal/events"
	"github.com/nicholas-fedor/shoutrrr"
	"github.com/nicholas-fedor/shoutrrr/pkg/router"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	log "github.com/sirupsen/logrus"
)

// ManagerConfig holds the configuration for the notification manager
type ManagerConfig struct {
	Providers []config.ProviderConfig
	Templates config.NotificationTemplate
}

// ProviderConfig holds the configuration for a notification provider
type ProviderConfig struct {
	URL  string
	Name string
}

// NotificationTemplates holds message templates for different event types
type NotificationTemplates struct {
	Forward string
}

// Manager handles notification delivery
type Manager struct {
	cfg       *ManagerConfig
	providers map[string]*router.ServiceRouter
	templates map[string]*template.Template
	mu        sync.Mutex
	sent      int
	lastReset time.Time
}

// NewManager creates a new notification manager
func NewManager(cfg *ManagerConfig) *Manager {
	m := &Manager{
		cfg:       cfg,
		providers: make(map[string]*router.ServiceRouter),
		templates: make(map[string]*template.Template),
		lastReset: time.Now(),
	}

	// Initialize providers
	for _, p := range cfg.Providers {
		sender, err := shoutrrr.CreateSender(p.URL)
		if err != nil {
			log.WithField("provider", p.Name).WithError(err).Error("error creating sender")
			continue
		}
		m.providers[p.Name] = sender
	}

	// Initialize templates
	m.parseTemplates()

	return m
}

// parseTemplates parses all notification templates
func (m *Manager) parseTemplates() {
	templates := map[events.EventType]string{
		events.Event_FORWARD:               m.cfg.Templates.Forward,
		events.Event_PEER_OFFLINE:          m.cfg.Templates.PeerOffline,
		events.Event_PEER_ONLINE:           m.cfg.Templates.PeerOnline,
		events.Event_CHANNEL_OPEN:          m.cfg.Templates.ChannelOpen,
		events.Event_CHANNEL_OPENING:       m.cfg.Templates.ChannelOpening,
		events.Event_CHANNEL_CLOSE:         m.cfg.Templates.ChannelClose,
		events.Event_CHANNEL_CLOSING:       m.cfg.Templates.ChannelClosing,
		events.Event_INVOICE_SETTLED:       m.cfg.Templates.InvoiceSettled,
		events.Event_FAILED_HTLC:           m.cfg.Templates.FailedHtlc,
		events.Event_KEYSEND:               m.cfg.Templates.Keysend,
		events.Event_ONCHAIN_CONFIRMED:     m.cfg.Templates.OnChainConfirmed,
		events.Event_ONCHAIN_MEMPOOL:       m.cfg.Templates.OnChainMempool,
		events.Event_PAYMENT_SUCCEEDED:     m.cfg.Templates.PaymentSucceeded,
		events.Event_REBALANCING_SUCCEEDED: m.cfg.Templates.RebalancingSucceeded,
	}

	for name, text := range templates {
		if text == "" {
			continue
		}
		tmpl, err := template.New(name.String()).Parse(text)
		if err != nil {
			log.WithField("template", name).WithError(err).Error("error parsing template")
			continue
		}
		m.templates[name.String()] = tmpl
	}
}

// RenderTemplate renders a notification template with the provided data
func (m *Manager) RenderTemplate(name string, data interface{}) (string, error) {
	tmpl, ok := m.templates[name]
	if !ok {
		return "", fmt.Errorf("template not found: %s", name)
	}

	log.WithFields(log.Fields{
		"template": name,
		"data":     data,
	}).Debug("rendering template")

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// Send sends a notification to all configured providers
func (m *Manager) Send(message string) {
	if message == "" {
		return
	}

	for name, provider := range m.providers {
		logger := log.WithField("provider", name).WithField("message", message)
		logger.Info("sending notification")

		errs := provider.Send(message, &types.Params{})
		for _, err := range errs {
			if err == nil {
				continue
			}

			logger.WithError(err).Error("error sending notification")
		}
	}
}

// SendBatch sends multiple notifications as a batch
func (m *Manager) SendBatch(messages []string) {
	if len(messages) == 0 {
		return
	}

	// Join messages with newlines
	message := ""
	for i, msg := range messages {
		if msg == "" {
			return
		}
		if i > 0 {
			message += "\n"
		}
		message += msg
	}

	m.Send(message)
}
