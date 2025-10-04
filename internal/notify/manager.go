package notify

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/Primexz/lnd-notify/internal/config"
	"github.com/nicholas-fedor/shoutrrr"
	"github.com/nicholas-fedor/shoutrrr/pkg/router"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// ManagerConfig holds the configuration for the notification manager
type ManagerConfig struct {
	Providers []config.ProviderConfig
	Templates config.NotificationTemplate
	RateLimit config.RateLimitConfig
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

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	MaxNotificationsPerMinute int
	BatchSimilarEvents        bool
	BatchWindowSeconds        int
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
			// Log error but continue with other providers
			fmt.Printf("Error initializing provider %s: %v\n", p.Name, err)
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
	templates := map[string]string{
		"forward_event": m.cfg.Templates.Forward,
	}

	for name, text := range templates {
		if text == "" {
			continue
		}
		tmpl, err := template.New(name).Parse(text)
		if err != nil {
			fmt.Printf("Error parsing template %s: %v\n", name, err)
			continue
		}
		m.templates[name] = tmpl
	}
}

// RenderTemplate renders a notification template with the provided data
func (m *Manager) RenderTemplate(name string, data interface{}) (string, error) {
	tmpl, ok := m.templates[name]
	if !ok {
		return "", fmt.Errorf("template not found: %s", name)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// Send sends a notification to all configured providers
func (m *Manager) Send(message string) error {
	if message == "" {
		return fmt.Errorf("empty message")
	}

	if !m.checkRateLimit() {
		return fmt.Errorf("rate limit exceeded")
	}

	var lastErr error
	for name, provider := range m.providers {
		if err := provider.Send(message, &types.Params{}); len(err) > 0 {
			lastErr = fmt.Errorf("provider %s: %w", name, err)
			continue
		}
	}

	return lastErr
}

// SendBatch sends multiple notifications as a batch
func (m *Manager) SendBatch(messages []string) error {
	if len(messages) == 0 {
		return fmt.Errorf("empty message batch")
	}

	// Join messages with newlines
	message := ""
	for i, msg := range messages {
		if msg == "" {
			return fmt.Errorf("empty message in batch at index %d", i)
		}
		if i > 0 {
			message += "\n"
		}
		message += msg
	}

	return m.Send(message)
}

// checkRateLimit checks if sending a notification would exceed the rate limit
func (m *Manager) checkRateLimit() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if now.Sub(m.lastReset) >= time.Minute {
		m.sent = 0
		m.lastReset = now
	}

	if m.cfg.RateLimit.MaxNotificationsPerMinute > 0 &&
		m.sent >= m.cfg.RateLimit.MaxNotificationsPerMinute {
		return false
	}

	m.sent++
	return true
}
