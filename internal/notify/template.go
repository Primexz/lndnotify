package notify

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Primexz/lndnotify/internal/events"
	log "github.com/sirupsen/logrus"
)

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
		events.Event_CHANNEL_STATUS_DOWN:   m.cfg.Templates.ChannelStatusDown,
		events.Event_CHANNEL_STATUS_UP:     m.cfg.Templates.ChannelStatusUp,
		events.Event_TLS_CERT_EXPIRY:       m.cfg.Templates.TLSCertExpiry,
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
