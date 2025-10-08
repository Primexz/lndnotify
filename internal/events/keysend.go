package events

import (
	"time"

	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type KeysendEvent struct {
	Msg       string
	Channel   *lnrpc.Channel
	Htlc      *lnrpc.InvoiceHTLC
	timestamp time.Time
}

type KeysendTemplate struct {
	Msg         string
	InChanAlias string
	InChanId    uint64
	Amount      string
}

func NewKeysendEvent(msg string, channel *lnrpc.Channel, htlc *lnrpc.InvoiceHTLC) *KeysendEvent {
	return &KeysendEvent{
		Msg:       msg,
		Channel:   channel,
		Htlc:      htlc,
		timestamp: time.Now(),
	}
}

func (e *KeysendEvent) Type() EventType {
	return Event_KEYSEND
}

func (e *KeysendEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *KeysendEvent) GetTemplateData() interface{} {
	var inChanAlias string
	if e.Channel != nil {
		inChanAlias = e.Channel.PeerAlias
	}

	return &KeysendTemplate{
		Msg:         e.Msg,
		InChanAlias: inChanAlias,
		InChanId:    e.Htlc.ChanId,
		Amount:      format.FormatDetailed(float64(e.Htlc.AmtMsat / 1000)),
	}
}
