package notify

import (
	"fmt"
	"text/template"
	"time"

	"github.com/Primexz/lndnotify/pkg/uploader"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	log "github.com/sirupsen/logrus"
)

// NewManager creates a new notification manager
func NewManager(cfg *ManagerConfig) *Manager {
	m := &Manager{
		cfg:       cfg,
		providers: make(map[string]Provider),
		templates: make(map[string]*template.Template),
		lastReset: time.Now(),
	}

	// Initialize providers and templates
	m.initializeProviders()
	m.parseTemplates()

	return m
}

// Send sends a notification to all configured providers
func (m *Manager) send(message string) {
	if message == "" {
		return
	}

	for name, p := range m.providers {
		logger := log.WithField("provider", name).WithField("message", message)

		logger.Info("sending notification")
		errs := p.Sender.Send(message, &types.Params{})
		for _, err := range errs {
			if err == nil {
				continue
			}
			logger.WithError(err).Error("error sending notification")
		}
	}
}

func (m *Manager) uploadFile(message string, file *uploader.File) {
	if file == nil {
		m.send(message)
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
			logger.Info("sending notification")

			errs := p.Sender.Send(msg, &types.Params{})
			for _, err := range errs {
				if err == nil {
					continue
				}
				logger.WithError(err).Error("error sending notification")
			}
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

// SendNotification sends a notification, either immediately or adds to batch
func (m *Manager) SendNotification(message string) {
	if m.cfg.Batching.Enabled {
		m.addToBatch(message, nil)
	} else {
		m.send(message)
	}
}

// SendNotificationWithFile sends a notification with file, either immediately or adds to batch
func (m *Manager) SendNotificationWithFile(message string, file *uploader.File) {
	if m.cfg.Batching.Enabled {
		m.addToBatch(message, file)
	} else {
		m.uploadFile(message, file)
	}
}

// Stop gracefully stops the notification manager and flushes any pending batches
func (m *Manager) Stop() {
	if m.cfg.Batching.Enabled {
		log.Info("flushing pending notification batch before shutdown")
		m.flushBatch()
	}
}
