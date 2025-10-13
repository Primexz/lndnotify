package notify

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/internal/events"
	"github.com/Primexz/lndnotify/pkg/uploader"
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

type Provider struct {
	Sender   *router.ServiceRouter
	Uploader uploader.Uploader
}

// Manager handles notification delivery
type Manager struct {
	cfg       *ManagerConfig
	providers map[string]Provider
	templates map[string]*template.Template
	mu        sync.Mutex
	sent      int
	lastReset time.Time
}

// NewManager creates a new notification manager
func NewManager(cfg *ManagerConfig) *Manager {
	m := &Manager{
		cfg:       cfg,
		providers: make(map[string]Provider),
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

		name, url, err := sender.ExtractServiceName(p.URL)
		if err != nil {
			log.WithField("provider", p.Name).WithError(err).Error("cannot initialize uploader, invalid URL")
			m.providers[p.Name] = Provider{Sender: sender, Uploader: nil}
			continue
		}
		upl := uploader.NewUploader(name, url)
		m.providers[p.Name] = Provider{Sender: sender, Uploader: upl}
	}

	// Initialize templates
	m.parseTemplates()

	return m
}

// parseTemplates parses all notification templates
func (m *Manager) parseTemplates() {
	templates := map[events.EventType]string{
		events.Event_BACKUP_MULTI:          m.cfg.Templates.BackupMulti,
		events.Event_FORWARD:               m.cfg.Templates.Forward,
		events.Event_PEER_OFFLINE:          m.cfg.Templates.PeerOffline,
		events.Event_PEER_ONLINE:           m.cfg.Templates.PeerOnline,
		events.Event_CHAIN_SYNC_LOST:       m.cfg.Templates.ChainSyncLost,
		events.Event_CHAIN_SYNC_RESTORED:   m.cfg.Templates.ChainSyncRestored,
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

	for name, p := range m.providers {
		logger := log.WithField("provider", name).WithField("message", message)

		m.sendRouter(message, p.Sender, logger)
	}
}

// sendRouter sends a notification using a specific router
func (m *Manager) sendRouter(message string, router *router.ServiceRouter, logger *log.Entry) {
	logger.Info("sending notification")

	errs := router.Send(message, &types.Params{})
	for _, err := range errs {
		if err == nil {
			continue
		}
		logger.WithError(err).Error("error sending notification")
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

func (m *Manager) UploadFile(message string, file *uploader.File) {
	if file == nil {
		m.Send(message)
		return
	}

	for name, p := range m.providers {
		logger := log.WithFields(log.Fields{
			"provider": name,
			"filename": file.Filename,
			"message":  message,
			"size":     len(file.Data),
		})

		// fallback sends the message without attachment via shoutrrr
		fallback := func(err error) {
			msg := message
			msg += "\n\n⚠️ Attachment removed"
			if err != nil {
				msg += fmt.Sprintf("\n🚨 Upload error: %v", err)
			} else {
				msg += " (file upload not supported for this provider)"
			}
			m.sendRouter(msg, p.Sender, logger)
		}

		if p.Uploader == nil {
			fallback(nil)
			continue
		}
		logger.Info("uploading file")

		err := p.Uploader.Upload(message, file)
		if err != nil {
			logger.WithError(err).Error("error uploading file, trying fallback")
			fallback(err)
		}
	}
}
