package events

import (
	"fmt"
	"math"
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type PaymentSucceededEvent struct {
	Payment       *lnrpc.Payment
	PayReq        *lnrpc.PayReq
	IsRebalancing bool
	getAlias      func(pubKey string) string
	timestamp     time.Time
}

type PaymentHopInfo struct {
	Pubkey  string
	Alias   string
	Amount  string
	Fee     string
	FeeRate string
}

type PaymentHtlcInfo struct {
	FirstHop  string
	PenultHop string
	HopInfo   []PaymentHopInfo
	Amount    string
	Fee       string
	FeeRate   string
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

func (e *PaymentSucceededEvent) GetTemplateData() interface{} {
	amountSats := float64(e.Payment.ValueMsat) / 1000
	feeSats := float64(e.Payment.FeeMsat) / 1000
	feeRate := FormattedFeeRate(feeSats, amountSats)

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

		hopInfo := make([]PaymentHopInfo, 0, len(htlc.Route.Hops))
		for _, hop := range htlc.Route.Hops {
			feeSats := float64(hop.FeeMsat) / 1000
			amountSats := float64(hop.AmtToForwardMsat) / 1000
			feeRate := FormattedFeeRate(feeSats, amountSats)

			hopInfo = append(hopInfo, PaymentHopInfo{
				Pubkey:  hop.PubKey,
				Alias:   e.getAlias(hop.PubKey),
				Amount:  format.FormatBasic(amountSats),
				Fee:     format.FormatDetailed(feeSats),
				FeeRate: feeRate,
			})
		}

		firstHop := htlc.Route.Hops[0]
		feeSats := float64(htlc.Route.TotalFeesMsat) / 1000
		amountSats := float64(htlc.Route.TotalAmtMsat)/1000 - feeSats
		feeRate := FormattedFeeRate(feeSats, amountSats)

		var penultHop string
		if len(htlc.Route.Hops) > 1 {
			penultHop = e.getAlias(htlc.Route.Hops[len(htlc.Route.Hops)-2].PubKey)
		}

		htlcInfo = append(htlcInfo, PaymentHtlcInfo{
			FirstHop:  e.getAlias(firstHop.PubKey),
			PenultHop: penultHop,
			HopInfo:   hopInfo,
			Fee:       format.FormatDetailed(feeSats),
			FeeRate:   feeRate,
			Amount:    format.FormatBasic(amountSats),
		})
	}

	return &PaymentSucceededTemplate{
		PaymentHash: e.Payment.PaymentHash,
		Amount:      format.FormatBasic(amountSats),
		Fee:         format.FormatDetailed(feeSats),
		FeeRate:     feeRate,
		HtlcInfo:    htlcInfo,
		Receiver:    receiver,
		Memo:        memo,
	}
}

func (e *PaymentSucceededEvent) ShouldProcess(cfg *config.Config) bool {
	if e.Type() == Event_REBALANCING_SUCCEEDED {
		return cfg.Events.RebalancingEvents
	}
	return cfg.Events.PaymentEvents
}

func FormattedFeeRate(fee, amount float64) string {
	var rate int32

	// Since the amounts in LN are rounded down, we have to round
	// commercially in order to reconstruct a correct ppm.
	if amount > 0 {
		rate = int32(math.Round(fee * 1e6 / amount))
	}
	return fmt.Sprintf("%d", rate)
}
