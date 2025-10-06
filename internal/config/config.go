package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the root configuration structure
type Config struct {
	LND           LNDConfig          `yaml:"lnd" validate:"required"`
	Notifications NotificationConfig `yaml:"notifications" validate:"required"`
	Events        EventConfig        `yaml:"events"`
}

// LNDConfig holds the LND node connection settings
type LNDConfig struct {
	Host         string `yaml:"host" validate:"required"`
	Port         int    `yaml:"port" validate:"required"`
	TLSCertPath  string `yaml:"tls_cert_path" validate:"required,file"`
	MacaroonPath string `yaml:"macaroon_path" validate:"required,file"`
}

// NotificationConfig holds notification service settings
type NotificationConfig struct {
	Providers []ProviderConfig     `yaml:"providers" validate:"required,min=1"`
	Templates NotificationTemplate `yaml:"templates"`
}

// ProviderConfig represents a single notification provider configuration
type ProviderConfig struct {
	URL  string `yaml:"url" validate:"required"`
	Name string `yaml:"name"`
}

// NotificationTemplate holds customizable message templates
type NotificationTemplate struct {
	Forward        string `yaml:"forward_event"`
	PeerOnline     string `yaml:"peer_online_event"`
	PeerOffline    string `yaml:"peer_offline_event"`
	ChannelOpen    string `yaml:"channel_open_event"`
	ChannelClose   string `yaml:"channel_close_event"`
	InvoiceSettled string `yaml:"invoice_settled_event"`
	FailedHtlc     string `yaml:"failed_htlc_event"`
}

// EventConfig controls which events to monitor
type EventConfig struct {
	ForwardEvents bool `yaml:"forward_events"`
	PeerEvents    bool `yaml:"peer_events"`
	ChannelEvents bool `yaml:"channel_events"`
	InvoiceEvents bool `yaml:"invoice_events"`
	FailedHtlc    bool `yaml:"failed_htlc_events"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	// #nosec G304 -- Just a config path from the user
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	cfg.setDefaults()
	return &cfg, nil
}

// Validate the configuration fields
func (c *Config) validate() error {
	// Basic validation
	if c.LND.Host == "" {
		return fmt.Errorf("LND host is required")
	}
	if c.LND.Port == 0 {
		return fmt.Errorf("LND port is required")
	}
	if c.LND.TLSCertPath == "" {
		return fmt.Errorf("LND TLS certificate path is required")
	}
	if c.LND.MacaroonPath == "" {
		return fmt.Errorf("LND macaroon path is required")
	}
	if len(c.Notifications.Providers) == 0 {
		return fmt.Errorf("at least one notification provider is required")
	}
	for _, p := range c.Notifications.Providers {
		if p.URL == "" {
			return fmt.Errorf("notification provider URL is required")
		}
	}

	return nil
}

// Set default templates if not specified
func (c *Config) setDefaults() {
	if c.Notifications.Templates.Forward == "" {
		c.Notifications.Templates.Forward = "💰 Forwarded {{.Amount}} sats, {{.PeerAliasIn}} -> {{.PeerAliasOut}}, earned {{.Fee}} sats"
	}
	if c.Notifications.Templates.PeerOnline == "" {
		c.Notifications.Templates.PeerOnline = "✅ Peer {{.PeerAlias}} ({{.PeerPubkeyShort}}) is online"
	}
	if c.Notifications.Templates.PeerOffline == "" {
		c.Notifications.Templates.PeerOffline = "⚠️ Peer {{.PeerAlias}} ({{.PeerPubkeyShort}}) went offline"
	}
	if c.Notifications.Templates.ChannelOpen == "" {
		c.Notifications.Templates.ChannelOpen = "🚀 Channel opened with {{.PeerAlias}}, capacity {{.Capacity}} sats"
	}
	if c.Notifications.Templates.ChannelClose == "" {
		c.Notifications.Templates.ChannelClose = "🔒 Channel closed with {{.PeerAlias}}, settled balance {{.SettledBalance}} sats"
	}
	if c.Notifications.Templates.InvoiceSettled == "" {
		c.Notifications.Templates.InvoiceSettled = "💵 Invoice settled: {{or .Memo \"No Memo\"}} for {{.Value}} sats"
	}
	if c.Notifications.Templates.FailedHtlc == "" {
		c.Notifications.Templates.FailedHtlc = "❌ Failed HTLC of {{.Amount}} sats\n{{.InChanAlias}} -> {{.OutChanAlias}}\nReason: {{.WireFailure}} ({{.FailureDetail}})\nActual Outbound: {{.OutChanLiquidity}} sats\nMissed Fee: {{.MissedFee}} sats"
	}
}
