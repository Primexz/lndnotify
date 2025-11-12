package events

import (
	"bytes"
	"text/template"
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"

	log "github.com/sirupsen/logrus"
)

type OnChainTransactionEvent struct {
	Event     *lnrpc.Transaction
	cfg       *config.Config
	timestamp time.Time
}

type OnChainTransactionTemplate struct {
	TxHash         string
	RawTxHex       string
	Amount         string
	TotalFees      string
	Confirmed      bool
	Outputs        []OnChainOutput
	TransactionURL string
}

type OnChainOutput struct {
	Amount       string
	Address      string
	OutputType   string
	IsOurAddress bool
}

func NewOnChainTransactionEvent(event *lnrpc.Transaction, cfg *config.Config) *OnChainTransactionEvent {
	return &OnChainTransactionEvent{
		Event:     event,
		cfg:       cfg,
		timestamp: time.Now(),
	}
}

func (e *OnChainTransactionEvent) Type() EventType {
	if e.Event.GetNumConfirmations() > 0 {
		return Event_ONCHAIN_CONFIRMED
	}

	return Event_ONCHAIN_MEMPOOL
}

func (e *OnChainTransactionEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *OnChainTransactionEvent) GetTemplateData(lang language.Tag) interface{} {
	outputs := make([]OnChainOutput, 0, len(e.Event.OutputDetails))
	for _, output := range e.Event.OutputDetails {
		outputs = append(outputs, OnChainOutput{
			Amount:       format.FormatBasic(float64(output.Amount), lang),
			OutputType:   output.OutputType.String(),
			IsOurAddress: output.IsOurAddress,
			Address:      output.Address,
		})
	}

	transactionURL, err := e.generateTransactionURL(e.Event.TxHash)
	if err != nil {
		log.WithError(err).WithField("tx_hash", e.Event.TxHash).Error("failed to generate transaction URL")
		transactionURL = "error generating URL (see logs)"
	}

	return &OnChainTransactionTemplate{
		TxHash:         e.Event.TxHash,
		RawTxHex:       e.Event.RawTxHex,
		Outputs:        outputs,
		Amount:         format.FormatBasic(float64(e.Event.Amount), lang),
		TotalFees:      format.FormatDetailed(float64(e.Event.TotalFees), lang),
		Confirmed:      e.Event.NumConfirmations > 0,
		TransactionURL: transactionURL,
	}
}

func (e *OnChainTransactionEvent) ShouldProcess(cfg *config.Config) bool {
	if !cfg.Events.OnChainEvents {
		return false
	}

	return uint64(e.Event.Amount) >= cfg.EventConfig.OnChainEvent.MinAmount
}

func (e *OnChainTransactionEvent) generateTransactionURL(txHash string) (string, error) {
	tmpl, err := template.New("transaction_url").Parse(e.cfg.EventConfig.OnChainEvent.TransactionUrlTemplate)
	if err != nil {
		return "", err
	}

	data := map[string]string{"TxHash": txHash}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
