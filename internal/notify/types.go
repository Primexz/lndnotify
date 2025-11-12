package notify

import (
	"sync"
	"text/template"
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/uploader"
	"github.com/nicholas-fedor/shoutrrr/pkg/router"
)

// ManagerConfig holds the configuration for the notification manager
type ManagerConfig struct {
	Providers []config.ProviderConfig
	Templates config.NotificationTemplate
	Batching  config.BatchingConfig
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

// Provider represents a notification provider with sender and uploader
type Provider struct {
	Sender   *router.ServiceRouter
	Uploader uploader.Uploader
}

// QueuedNotification represents a notification waiting to be sent
type QueuedNotification struct {
	Message string
	File    *uploader.File
}

// Manager handles notification delivery
type Manager struct {
	cfg       *ManagerConfig
	providers map[string]Provider
	templates map[string]*template.Template
	lastReset time.Time

	// Batching fields
	batchQueue []QueuedNotification
	batchMu    sync.Mutex
	flushTimer *time.Timer
}
