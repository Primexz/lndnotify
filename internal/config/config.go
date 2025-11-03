package config

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/text/language"
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
	Providers  []ProviderConfig     `yaml:"providers" validate:"required,min=1"`
	Templates  NotificationTemplate `yaml:"templates"`
	Formatting Formatting           `yaml:"formatting"`
	Batching   BatchingConfig       `yaml:"batching"`
}

// BatchingConfig holds batching configuration
type BatchingConfig struct {
	Enabled       bool          `yaml:"enabled"`
	FlushInterval time.Duration `yaml:"flush_interval"`
	MaxSize       int           `yaml:"max_size"`
}

type Formatting struct {
	Locale LanguageTag `yaml:"locale"`
}

// ProviderConfig represents a single notification provider configuration
type ProviderConfig struct {
	URL  string `yaml:"url" validate:"required"`
	Name string `yaml:"name"`
}

// NotificationTemplate holds customizable message templates
// NOTE: Keep fields in alphabetical order to prevent merge conflicts when adding new events
type NotificationTemplate struct {
	BackupMulti          string `yaml:"backup_multi_event"`
	ChainSyncLost        string `yaml:"chain_sync_lost_event"`
	ChainSyncRestored    string `yaml:"chain_sync_restored_event"`
	ChannelClose         string `yaml:"channel_close_event"`
	ChannelClosing       string `yaml:"channel_closing_event"`
	ChannelFeeChange     string `yaml:"channel_fee_change_event"`
	ChannelOpen          string `yaml:"channel_open_event"`
	ChannelOpening       string `yaml:"channel_opening_event"`
	ChannelStatusUp      string `yaml:"channel_status_up_event"`
	ChannelStatusDown    string `yaml:"channel_status_down_event"`
	FailedHtlc           string `yaml:"failed_htlc_event"`
	Forward              string `yaml:"forward_event"`
	Healthy              string `yaml:"healthy_event"`
	Unhealthy            string `yaml:"unhealthy_event"`
	InvoiceSettled       string `yaml:"invoice_settled_event"`
	Keysend              string `yaml:"keysend_event"`
	OnChainConfirmed     string `yaml:"on_chain_confirmed_event"`
	OnChainMempool       string `yaml:"on_chain_mempool_event"`
	PaymentSucceeded     string `yaml:"payment_succeeded_event"`
	PeerOffline          string `yaml:"peer_offline_event"`
	PeerOnline           string `yaml:"peer_online_event"`
	RebalancingSucceeded string `yaml:"rebalancing_succeeded_event"`
	TLSCertExpiry        string `yaml:"tls_cert_expiry_event"`
	WalletState          string `yaml:"wallet_state_event"`
	LndUpdateAvailable   string `yaml:"lnd_update_available_event"`
}

// EventFlags controls which events to monitor (feature flags)
// NOTE: Keep fields in alphabetical order to prevent merge conflicts when adding new events
type EventFlags struct {
	BackupEvents        bool `yaml:"backup_events"`
	ChainSyncEvents     bool `yaml:"chain_sync_events"`
	ChannelEvents       bool `yaml:"channel_events"`
	ChannelFeeEvents    bool `yaml:"channel_fee_events"`
	ChannelStatusEvents bool `yaml:"channel_status_events"`
	FailedHtlc          bool `yaml:"failed_htlc_events"`
	ForwardEvents       bool `yaml:"forward_events"`
	HealthEvents        bool `yaml:"health_events"`
	InvoiceEvents       bool `yaml:"invoice_events"`
	KeysendEvents       bool `yaml:"keysend_events"`
	OnChainEvents       bool `yaml:"on_chain_events"`
	PaymentEvents       bool `yaml:"payment_events"`
	PeerEvents          bool `yaml:"peer_events"`
	RebalancingEvents   bool `yaml:"rebalancing_events"`
	StatusEvents        bool `yaml:"status_events"`
	TLSCertExpiryEvents bool `yaml:"tls_cert_expiry_events"`
	WalletStateEvents   bool `yaml:"wallet_state_events"`
	LndUpdateEvents     bool `yaml:"lnd_update_events"`
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
		MinAmount   uint64 `yaml:"min_amount"`
		SkipKeysend *bool  `yaml:"skip_keysend"`
	} `yaml:"invoice_event"`
	PaymentEvent struct {
		MinAmount uint64 `yaml:"min_amount"`
	} `yaml:"payment_event"`
	RebalancingEvent struct {
		MinAmount uint64 `yaml:"min_amount"`
	} `yaml:"rebalancing_event"`
	OnChainEvent struct {
		MinAmount              uint64 `yaml:"min_amount"`
		TransactionUrlTemplate string `yaml:"transaction_url_template"`
	} `yaml:"on_chain_event"`
	ChainLostEvent struct {
		Threshold       time.Duration `yaml:"threshold"`
		WarningInterval time.Duration `yaml:"warning_interval"`
	} `yaml:"chain_lost_event"`
	ChannelStatusEvent struct {
		MinDowntime time.Duration `yaml:"min_downtime"`
	} `yaml:"channel_status_event"`
	TLSCertExpiryEvent struct {
		Threshold time.Duration `yaml:"threshold"`
	} `yaml:"tls_cert_expiry_event"`
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
	if c.Notifications.Formatting.Locale.Tag == language.Und {
		c.Notifications.Formatting.Locale.Tag = language.English
	}
	if c.EventConfig.OnChainEvent.TransactionUrlTemplate == "" {
		c.EventConfig.OnChainEvent.TransactionUrlTemplate = "https://mempool.space/tx/{{.TxHash}}"
	}

	// Set default templates in alphabetical order to prevent merge conflicts
	if c.Notifications.Templates.BackupMulti == "" {
		c.Notifications.Templates.BackupMulti = "❗️ Channel backup received for {{.NumChanPoints}} channels\n\nChannel Points:\n{{range .ChanPoints}}- {{.}}\n{{end}}\nFilename: {{.Filename}}\n SHA256: {{.Sha256Sum}}"
	}
	if c.Notifications.Templates.ChannelClose == "" {
		c.Notifications.Templates.ChannelClose = "🔒 Channel closed with {{.PeerAlias}}\nCapacity {{.Capacity}} sats\nSettled balance {{.SettledBalance}} sats\n\nChannel Point: {{.ChannelPoint}}\nClose Type: {{if eq .CloseType 0}}🤝 Cooperatively {{if .CloseInitiator}}Local{{else}}Remote{{end}}{{else if eq .CloseType 1}}🔴 Force Local{{else if eq .CloseType 2}}🔴 Force Remote{{else if eq .CloseType 3}}🚨 Breach{{else}}💀 Other{{end}}"
	}
	if c.Notifications.Templates.ChannelClosing == "" {
		c.Notifications.Templates.ChannelClosing = "⏳ Closing channel with {{.PeerAlias}}\nCapacity {{.Capacity}} sats\nLimbo: {{.LimboBalance}} sats\n\nClosing TxID: {{.ClosingTxid}}\nRaw TX: {{.ClosingTxHex}}"
	}
	if c.Notifications.Templates.ChannelFeeChange == "" {
		c.Notifications.Templates.ChannelFeeChange = "✏️ Fee change detected on channel with {{.PeerAlias}} ({{.PeerPubkeyShort}})\nCapacity: {{.Capacity}} sats\n\nFee Rate: {{if ne .OldFeeRate .NewFeeRate}}{{.OldFeeRate}} -> {{.NewFeeRate}} ({{.FeeRateChange}} ppm, {{.FeeRateChangePercent}}){{else}}{{.OldFeeRate}}{{end}} ppm\nBase Fee: {{if ne .OldBaseFee .NewBaseFee}}{{.OldBaseFee}} -> {{.NewBaseFee}} ({{.BaseFeeChange}} sats, {{.BaseFeeChangePercent}}){{else}}{{.OldBaseFee}}{{end}} sats\nInbound Fee Rate: {{if ne .OldInboundFeeRate .NewInboundFeeRate}}{{.OldInboundFeeRate}} -> {{.NewInboundFeeRate}} ({{.InboundFeeRateChange}} ppm, {{.InboundFeeRateChangePercent}}){{else}}{{.OldInboundFeeRate}}{{end}} ppm\nInbound Base Fee: {{if ne .OldInboundBaseFee .NewInboundBaseFee}}{{.OldInboundBaseFee}} -> {{.NewInboundBaseFee}} ({{.InboundBaseFeeChange}} sats, {{.InboundBaseFeeChangePercent}}){{else}}{{.OldInboundBaseFee}}{{end}} sats"
	}
	if c.Notifications.Templates.ChannelOpen == "" {
		c.Notifications.Templates.ChannelOpen = "🚀 Channel opened with {{.PeerAlias}}\nCapacity {{.Capacity}} sats"
	}
	if c.Notifications.Templates.ChannelOpening == "" {
		c.Notifications.Templates.ChannelOpening = "{{if .Initiator}}⏳ Opening new {{.Capacity}} sats channel to {{.PeerAlias}}{{else}}⏳ Accepting new {{.Capacity}} sats channel from {{.PeerAlias}}{{end}}"
	}
	if c.Notifications.Templates.FailedHtlc == "" {
		c.Notifications.Templates.FailedHtlc = "❌ Failed HTLC of {{.Amount}} sats\n{{.InChanAlias}} -> {{.OutChanAlias}}\nReason: {{.WireFailure}} ({{.FailureDetail}})\nActual Outbound: {{.OutChanLiquidity}} sats\nMissed Fee: {{.MissedFee}} sats\nLocal liquidity failure: {{if .IsLocalLiquidityFailure}}✅{{else}}❌{{end}}"
	}
	if c.Notifications.Templates.Forward == "" {
		c.Notifications.Templates.Forward = "💰 Forwarded {{.Amount}} sats\n{{.PeerAliasIn}} -> {{.PeerAliasOut}}\nEarned {{.Fee}} sats ({{.FeeRate}} ppm)"
	}
	if c.Notifications.Templates.InvoiceSettled == "" {
		c.Notifications.Templates.InvoiceSettled = "💵 Invoice settled: {{or .Memo \"No Memo\"}} for {{.Value}} sats"
	}
	if c.Notifications.Templates.Keysend == "" {
		c.Notifications.Templates.Keysend = "📨 Keysend received:\n\n{{.Msg}}\n\nChannel In: {{.InChanAlias}} ({{.InChanId}})"
	}
	if c.Notifications.Templates.ChainSyncLost == "" {
		c.Notifications.Templates.ChainSyncLost = "⚠️ Chain is out of sync since {{.Duration}}"
	}
	if c.Notifications.Templates.ChainSyncRestored == "" {
		c.Notifications.Templates.ChainSyncRestored = "✅ Chain is back in sync after {{.Duration}}"
	}
	if c.Notifications.Templates.OnChainMempool == "" {
		c.Notifications.Templates.OnChainMempool = "🔗 Discovered On-Chain transaction in mempool: {{.Amount}} sats\nFee: {{.TotalFees}} sats\n\nOutputs:\n{{range .Outputs}}- {{.Amount}} sats to {{.Address}} ({{.OutputType}}{{if .IsOurAddress}}, ours{{end}})\n{{end}}\nView on explorer: {{.TransactionURL}}\nTxID: {{.TxHash}}"
	}
	if c.Notifications.Templates.OnChainConfirmed == "" {
		c.Notifications.Templates.OnChainConfirmed = "🔗 Confirmed On-Chain transaction: {{.Amount}} sats\nFee: {{.TotalFees}} sats\n\nView on explorer: {{.TransactionURL}}\nTxID: {{.TxHash}}"
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
		c.Notifications.Templates.RebalancingSucceeded = "{{range .HtlcInfo}}☯️ Rebalanced {{.Amount}} sats {{.FirstHop}} → {{.PenultHop}}\nFee: {{.Fee}} sats ({{.FeeRate}} ppm)\nRoute: {{range $i, $hop := .HopInfo}}{{if $i}} -> {{end}}{{$hop.Alias}} ({{$hop.FeeRate}} ppm){{end}}\n\n{{end}}"
	}
	if c.Notifications.Templates.ChannelStatusUp == "" {
		c.Notifications.Templates.ChannelStatusUp = "🟢 Channel with {{.PeerAlias}} ({{.PeerPubkeyShort}}) is back online after {{.Duration}}\nCapacity {{.Capacity}} sats"
	}
	if c.Notifications.Templates.ChannelStatusDown == "" {
		c.Notifications.Templates.ChannelStatusDown = "🔴 Channel with {{.PeerAlias}} ({{.PeerPubkeyShort}}) is down since {{.Duration}}\nCapacity {{.Capacity}} sats"
	}
	if c.Notifications.Templates.TLSCertExpiry == "" {
		c.Notifications.Templates.TLSCertExpiry = "⚠️ LND TLS certificate is expiring soon on {{.ExpiryDate}} (in {{.TimeUntilExpiry}})"
	}
	if c.Notifications.Templates.WalletState == "" {
		c.Notifications.Templates.WalletState = "👛 Wallet state changed: {{.OldState}} -> {{.NewState}}"
	}
	if c.Notifications.Templates.Healthy == "" {
		c.Notifications.Templates.Healthy = "✅ LND node is healthy again"
	}
	if c.Notifications.Templates.Unhealthy == "" {
		c.Notifications.Templates.Unhealthy = "❌ LND node is unhealthy\nError: {{.Err}}"
	}
	if c.Notifications.Templates.LndUpdateAvailable == "" {
		c.Notifications.Templates.LndUpdateAvailable = "⬆️ New LND version available: {{.LatestVersion}}\nYou are currently running version: {{.CurrentVersion}}"
	}

	if c.EventConfig.ChainLostEvent.Threshold == 0 {
		c.EventConfig.ChainLostEvent.Threshold = 5 * time.Minute
	}
	if c.EventConfig.ChainLostEvent.WarningInterval == 0 {
		c.EventConfig.ChainLostEvent.WarningInterval = 15 * time.Minute
	}
	if c.EventConfig.ChannelStatusEvent.MinDowntime == 0 {
		c.EventConfig.ChannelStatusEvent.MinDowntime = 10 * time.Minute
	}
	if c.EventConfig.TLSCertExpiryEvent.Threshold == 0 {
		c.EventConfig.TLSCertExpiryEvent.Threshold = 7 * 24 * time.Hour
	}
	if c.EventConfig.InvoiceEvent.SkipKeysend == nil {
		defaultSkip := true
		c.EventConfig.InvoiceEvent.SkipKeysend = &defaultSkip
	}

	// Set default batching configuration
	if c.Notifications.Batching.FlushInterval == 0 {
		c.Notifications.Batching.FlushInterval = 5 * time.Second
	}
	if c.Notifications.Batching.MaxSize == 0 {
		c.Notifications.Batching.MaxSize = 10
	}
}
