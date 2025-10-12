package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"
)

type PaymentSucceededEvent struct {
	Payment       *lnrpc.Payment
	PayReq        *lnrpc.PayReq
	IsRebalancing bool
	getAlias      func(pubKey string) string
	timestamp     time.Time
}

type PaymentSucceededTemplate struct {
	PaymentHash string
	Amount      string
	Fee         string
	FeeRate     string
	HtlcInfo    []PaymentHtlcInfo
	Receiver    string
	Memo        string
}

type PaymentHtlcInfo struct {
	FirstHop  string
	PenultHop string
	HopInfo   []PaymentHopInfo
	Amount    string
	Fee       string
	FeeRate   string
}

type PaymentHopInfo struct {
	Pubkey  string
	Alias   string
	Amount  string
	Fee     string
	FeeRate string
}

func NewPaymentSucceededEvent(payment *lnrpc.Payment, payReq *lnrpc.PayReq,
	isRebalancing bool, getAlias func(pubKey string) string) *PaymentSucceededEvent {

	return &PaymentSucceededEvent{
		Payment:       payment,
		PayReq:        payReq,
		IsRebalancing: isRebalancing,
		getAlias:      getAlias,
		timestamp:     time.Now(),
	}
}

func (e *PaymentSucceededEvent) Type() EventType {
	if e.IsRebalancing {
		return Event_REBALANCING_SUCCEEDED
	}
	return Event_PAYMENT_SUCCEEDED
}

func (e *PaymentSucceededEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *PaymentSucceededEvent) GetTemplateData(lang language.Tag) interface{} {
	amountSats := float64(e.Payment.ValueMsat) / 1000
	feeSats := float64(e.Payment.FeeMsat) / 1000

	// Get memo from PayReq
	var memo string
	if e.PayReq != nil {
		memo = e.PayReq.Description
	}

	var receiver string
	if len(e.Payment.Htlcs) > 0 {
		lastHtlc := e.Payment.Htlcs[len(e.Payment.Htlcs)-1]
		receiver = e.getAlias(lastHtlc.Route.Hops[len(lastHtlc.Route.Hops)-1].PubKey)
	}

	var htlcInfo []PaymentHtlcInfo
	for _, htlc := range e.Payment.Htlcs {
		if len(htlc.Route.Hops) == 0 {
			continue
		}

		// We exclude the last hop for rebalancing payments, as it is
		// always our own node.
		hopsToProcess := htlc.Route.Hops
		if e.IsRebalancing && len(hopsToProcess) > 0 {
			hopsToProcess = hopsToProcess[:len(hopsToProcess)-1]
		}

		hopInfo := make([]PaymentHopInfo, 0, len(hopsToProcess))
		for _, hop := range hopsToProcess {
			feeSats := float64(hop.FeeMsat) / 1000
			amountSats := float64(hop.AmtToForwardMsat) / 1000

			hopInfo = append(hopInfo, PaymentHopInfo{
				Pubkey:  hop.PubKey,
				Alias:   e.getAlias(hop.PubKey),
				Amount:  format.FormatBasic(amountSats, lang),
				Fee:     format.FormatDetailed(feeSats, lang),
				FeeRate: format.FormatRatePPM(feeSats, amountSats, lang),
			})
		}

		firstHop := htlc.Route.Hops[0]
		feeSats := float64(htlc.Route.TotalFeesMsat) / 1000
		amountSats := float64(htlc.Route.TotalAmtMsat)/1000 - feeSats

		var penultHop string
		if len(htlc.Route.Hops) > 1 {
			penultHop = e.getAlias(htlc.Route.Hops[len(htlc.Route.Hops)-2].PubKey)
		}

		htlcInfo = append(htlcInfo, PaymentHtlcInfo{
			FirstHop:  e.getAlias(firstHop.PubKey),
			PenultHop: penultHop,
			HopInfo:   hopInfo,
			Fee:       format.FormatDetailed(feeSats, lang),
			FeeRate:   format.FormatRatePPM(feeSats, amountSats, lang),
			Amount:    format.FormatBasic(amountSats, lang),
		})
	}

	return &PaymentSucceededTemplate{
		PaymentHash: e.Payment.PaymentHash,
		Amount:      format.FormatBasic(amountSats, lang),
		Fee:         format.FormatDetailed(feeSats, lang),
		FeeRate:     format.FormatRatePPM(feeSats, amountSats, lang),
		HtlcInfo:    htlcInfo,
		Receiver:    receiver,
		Memo:        memo,
	}
}

func (e *PaymentSucceededEvent) ShouldProcess(cfg *config.Config) bool {
	switch e.Type() {
	case Event_REBALANCING_SUCCEEDED:
		if !cfg.Events.RebalancingEvents {
			return false
		}
		return uint64(e.Payment.ValueSat) >= cfg.EventConfig.RebalancingEvent.MinAmount // #nosec G115

	case Event_PAYMENT_SUCCEEDED:
		if !cfg.Events.PaymentEvents {
			return false
		}
		return uint64(e.Payment.ValueSat) >= cfg.EventConfig.PaymentEvent.MinAmount // #nosec G115

	default:
		return false
	}
}
