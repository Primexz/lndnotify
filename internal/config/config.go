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
	Events        EventFlags         `yaml:"events"`
	EventConfig   EventConfig        `yaml:"event_config"`
	LogLevel      string             `yaml:"log_level" validate:"omitempty,oneof=panic fatal error warn info debug trace"`
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
// NOTE: Keep fields in alphabetical order to prevent merge conflicts when adding new events
type NotificationTemplate struct {
	ChannelClose         string `yaml:"channel_close_event"`
	ChannelOpen          string `yaml:"channel_open_event"`
	FailedHtlc           string `yaml:"failed_htlc_event"`
	Forward              string `yaml:"forward_event"`
	InvoiceSettled       string `yaml:"invoice_settled_event"`
	Keysend              string `yaml:"keysend_event"`
	OnChain              string `yaml:"on_chain_event"`
	PaymentSucceeded     string `yaml:"payment_succeeded_event"`
	PeerOffline          string `yaml:"peer_offline_event"`
	PeerOnline           string `yaml:"peer_online_event"`
	RebalancingSucceeded string `yaml:"rebalancing_succeeded_event"`
}

// EventFlags controls which events to monitor (feature flags)
// NOTE: Keep fields in alphabetical order to prevent merge conflicts when adding new events
type EventFlags struct {
	ChannelEvents     bool `yaml:"channel_events"`
	FailedHtlc        bool `yaml:"failed_htlc_events"`
	ForwardEvents     bool `yaml:"forward_events"`
	InvoiceEvents     bool `yaml:"invoice_events"`
	KeysendEvents     bool `yaml:"keysend_events"`
	OnChainEvents     bool `yaml:"on_chain_events"`
	PaymentEvents     bool `yaml:"payment_events"`
	PeerEvents        bool `yaml:"peer_events"`
	RebalancingEvents bool `yaml:"rebalancing_events"`
	StatusEvents      bool `yaml:"status_events"`
}

// EventConfig contains specific configuration for each event type
// NOTE: Keep fields in alphabetical order to prevent merge conflicts when adding new events
type EventConfig struct {
	FailedHtlcEvent struct {
		MinAmount uint64 `yaml:"min_amount"`
	} `yaml:"failed_htlc_event"`
	ForwardEvent struct {
		MinAmount uint64 `yaml:"min_amount"`
	} `yaml:"forward_event"`
	InvoiceEvent struct {
		MinAmount uint64 `yaml:"min_amount"`
	} `yaml:"invoice_event"`
	PaymentEvent struct {
		MinAmount uint64 `yaml:"min_amount"`
	} `yaml:"payment_event"`
	RebalancingEvent struct {
		MinAmount uint64 `yaml:"min_amount"`
	} `yaml:"rebalancing_event"`
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
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}

	// Set default templates in alphabetical order to prevent merge conflicts
	if c.Notifications.Templates.ChannelClose == "" {
		c.Notifications.Templates.ChannelClose = "🔒 Channel closed with {{.PeerAlias}}, settled balance {{.SettledBalance}} sats"
	}
	if c.Notifications.Templates.ChannelOpen == "" {
		c.Notifications.Templates.ChannelOpen = "🚀 Channel opened with {{.PeerAlias}}, capacity {{.Capacity}} sats"
	}
	if c.Notifications.Templates.FailedHtlc == "" {
		c.Notifications.Templates.FailedHtlc = "❌ Failed HTLC of {{.Amount}} sats\n{{.InChanAlias}} -> {{.OutChanAlias}}\nReason: {{.WireFailure}} ({{.FailureDetail}})\nActual Outbound: {{.OutChanLiquidity}} sats\nMissed Fee: {{.MissedFee}} sats"
	}
	if c.Notifications.Templates.Forward == "" {
		c.Notifications.Templates.Forward = "💰 Forwarded {{.Amount}} sats, {{.PeerAliasIn}} -> {{.PeerAliasOut}}, earned {{.Fee}} sats"
	}
	if c.Notifications.Templates.InvoiceSettled == "" {
		c.Notifications.Templates.InvoiceSettled = "💵 Invoice settled: {{or .Memo \"No Memo\"}} for {{.Value}} sats"
	}
	if c.Notifications.Templates.Keysend == "" {
		c.Notifications.Templates.Keysend = "📨 Keysend received:\n\n{{.Msg}}\n\nChannel In: {{.InChanAlias}} ({{.InChanId}})"
	}
	if c.Notifications.Templates.OnChain == "" {
		c.Notifications.Templates.OnChain = "🔗 On-chain transaction: {{.Amount}} sats\nFee: {{.TotalFees}} sats\n\nOutputs:\n{{range .Outputs}}- {{.Amount}} sats to {{.Address}} ({{.OutputType}}{{if .IsOurAddress}}, ours{{end}})\n{{end}}\nTxID: {{.TxHash}}\nRaw TX: {{.RawTxHex}}"
	}
	if c.Notifications.Templates.PaymentSucceeded == "" {
		c.Notifications.Templates.PaymentSucceeded = "⚡️ Payment: {{.Amount}} sats (fee: {{.Fee}}) to {{.Receiver}}{{if .Memo}} - {{.Memo}}{{end}}{{range .HtlcInfo}}\n  HTLC: {{.Amount}} via {{.FirstHop}} (fee: {{.Fee}}){{end}}\nHash: {{.PaymentHash}}"
	}
	if c.Notifications.Templates.PeerOffline == "" {
		c.Notifications.Templates.PeerOffline = `{{if .PeerAlias}}⚠️ Peer {{.PeerAlias}} ({{.PeerPubkeyShort}}) went offline{{else}}⚠️ Peer {{.PeerPubKey}} went offline{{end}}`
	}
	if c.Notifications.Templates.PeerOnline == "" {
		c.Notifications.Templates.PeerOnline = `{{if .PeerAlias}}✅ Peer {{.PeerAlias}} ({{.PeerPubkeyShort}}) is online{{else}}✅ Peer {{.PeerPubKey}} is online{{end}}`
	}
	if c.Notifications.Templates.RebalancingSucceeded == "" {
		c.Notifications.Templates.RebalancingSucceeded = "{{range .HtlcInfo}}☯️ Rebalanced {{.Amount}} sats {{.FirstHop}} → {{.PenultHop}}\nFee: {{.Fee}} sats ({{.FeeRate}} ppm)\nRoute: {{range $i, $hop := .HopInfo}}{{if $i}} -> {{end}}{{$hop.Alias}} ({{$hop.FeeRate}} ppm){{end}}\n{{end}}"
	}
}
