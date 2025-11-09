package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/uploader"
	"golang.org/x/text/language"
)

// Event is the base interface for all Lightning Network events
type Event interface {
	Type() EventType
	Timestamp() time.Time
	GetTemplateData(lang language.Tag) interface{}
	ShouldProcess(cfg *config.Config) bool
}

// FileSource is an interface for types that can provide a file
type FileSource interface {
	GetFile() *uploader.File
}

type EventType string

// NOTE: Keep constants in alphabetical order to prevent merge conflicts when adding new events
const (
	Event_BACKUP_MULTI          EventType = "backup_multi_event"
	Event_CHAIN_SYNC_LOST       EventType = "chain_sync_lost_event"
	Event_CHAIN_SYNC_RESTORED   EventType = "chain_sync_restored_event"
	Event_CHANNEL_CLOSE         EventType = "channel_close_event"
	Event_CHANNEL_CLOSING       EventType = "channel_closing_event"
	Event_CHANNEL_FEE_CHANGE    EventType = "channel_fee_change_event"
	Event_CHANNEL_OPEN          EventType = "channel_open_event"
	Event_CHANNEL_OPENING       EventType = "channel_opening_event"
	Event_CHANNEL_STATUS_UP     EventType = "channel_status_up_event"
	Event_CHANNEL_STATUS_DOWN   EventType = "channel_status_down_event"
	Event_FAILED_HTLC           EventType = "failed_htlc_event"
	Event_FORWARD               EventType = "forward_event"
	Event_HEALTHY               EventType = "healthy_event"
	Event_UNHEALTHY             EventType = "unhealthy_event"
	Event_INVOICE_SETTLED       EventType = "invoice_settled_event"
	Event_KEYSEND               EventType = "keysend_event"
	Event_ONCHAIN_MEMPOOL       EventType = "on_chain_event"
	Event_ONCHAIN_CONFIRMED     EventType = "on_chain_event_confirmed"
	Event_PAYMENT_SUCCEEDED     EventType = "payment_succeeded_event"
	Event_PEER_OFFLINE          EventType = "peer_offline_event"
	Event_PEER_ONLINE           EventType = "peer_online_event"
	Event_REBALANCING_SUCCEEDED EventType = "rebalancing_succeeded_event"
	Event_TLS_CERT_EXPIRY       EventType = "tls_cert_expiry_event"
	Event_WALLET_STATE          EventType = "wallet_state_event"
	Event_LND_UPDATE_AVAILABLE  EventType = "lnd_update_available_event"
	Event_HTLC_EXPIRATION       EventType = "htlc_expiration_event"
)

func (et EventType) String() string {
	return string(et)
}
