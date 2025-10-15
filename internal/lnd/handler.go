package lnd

import (
	"context"
	"time"

	"github.com/Primexz/lndnotify/internal/events"
	"github.com/Primexz/lndnotify/pkg/format"
	"github.com/cenkalti/backoff/v5"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	log "github.com/sirupsen/logrus"
)

// handleForwards polls for forwarding events
func (c *Client) handleForwards() {
	log.Debug("starting forward event handler")
	defer c.wg.Done()

	var lastOffset uint32
	start := time.Now()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			log.WithFields(log.Fields{
				"since":       start,
				"last_offset": lastOffset,
			}).Debug("polling for forwarding events")

			resp, err := c.client.ForwardingHistory(c.ctx, &lnrpc.ForwardingHistoryRequest{
				StartTime:       uint64(start.Unix()), // #nosec G115
				PeerAliasLookup: true,
				IndexOffset:     lastOffset,
			})
			if err != nil {
				log.WithError(err).Error("error fetching forwarding history")
				continue
			}

			forwards := resp.GetForwardingEvents()
			for _, fwd := range forwards {
				c.eventSub <- events.NewForwardEvent(fwd)
			}

			// push last offset for next request. lnd will return the current offset
			// if no new events are available.
			lastOffset = resp.LastOffsetIndex
		}
	}
}

// handlePeerEvents handles peer connection and disconnection events
// Deprecated: Replace with channel_status
func (c *Client) handlePeerEvents() {
	log.Debug("starting peer event handler")
	defer c.wg.Done()

	retry(c.ctx, "peer event subscription", func() (string, error) {
		ev, err := c.client.SubscribePeerEvents(c.ctx, &lnrpc.PeerEventSubscription{})
		if err != nil {
			return "", err
		}

		log.Debug("peer event subscription established")

		for {
			select {
			case <-c.ctx.Done():
				return "", nil
			default:
			}

			peerEvent, err := ev.Recv()
			if err != nil {
				return "", err // Return error to trigger retry
			}

			nodeInfo, err := c.client.GetNodeInfo(c.ctx, &lnrpc.NodeInfoRequest{
				PubKey: peerEvent.GetPubKey(),
			})
			if err != nil {
				log.WithField("pubkey", peerEvent.GetPubKey()).WithError(err).Warn("error fetching node info")
			}

			switch peerEvent.GetType() {
			case lnrpc.PeerEvent_PEER_ONLINE:
				c.eventSub <- events.NewPeerOnlineEvent(nodeInfo, peerEvent)
			case lnrpc.PeerEvent_PEER_OFFLINE:
				c.eventSub <- events.NewPeerOfflineEvent(nodeInfo, peerEvent)
			}
		}
	})
}

// handleChannelEvents handles channel open and close events
func (c *Client) handleChannelEvents() {
	log.Debug("starting channel event handler")
	defer c.wg.Done()

	retry(c.ctx, "channel event subscription", func() (string, error) {
		ev, err := c.client.SubscribeChannelEvents(c.ctx, &lnrpc.ChannelEventSubscription{})
		if err != nil {
			return "", err
		}

		log.Debug("channel event subscription established")

		for {
			select {
			case <-c.ctx.Done():
				return "", nil
			default:
			}

			chanEvent, err := ev.Recv()
			if err != nil {
				return "", err // Return error to trigger retry
			}

			log.WithField("channel_event", chanEvent).Trace("received channel event")

			switch chanEvent.GetType() {
			case lnrpc.ChannelEventUpdate_PENDING_OPEN_CHANNEL:
				c.pendChanManager.RefreshDelayed()
			case lnrpc.ChannelEventUpdate_OPEN_CHANNEL:
				channel := chanEvent.GetOpenChannel()
				nodeInfo, err := c.client.GetNodeInfo(c.ctx, &lnrpc.NodeInfoRequest{
					PubKey: channel.RemotePubkey,
				})
				if err != nil {
					log.WithError(err).Error("error fetching node info")
					continue
				}

				c.eventSub <- events.NewChannelOpenEvent(nodeInfo.Node, channel)
			case lnrpc.ChannelEventUpdate_CLOSED_CHANNEL:
				channel := chanEvent.GetClosedChannel()
				nodeInfo, err := c.client.GetNodeInfo(c.ctx, &lnrpc.NodeInfoRequest{
					PubKey: channel.RemotePubkey,
				})
				if err != nil {
					log.WithError(err).Error("error fetching node info")
					continue
				}

				c.eventSub <- events.NewChannelCloseEvent(nodeInfo.Node, channel)
			}
		}
	})
}

func (c *Client) handleInvoiceEvents() {
	log.Debug("starting invoice event handler")
	defer c.wg.Done()

	retry(c.ctx, "invoice event subscription", func() (string, error) {
		ev, err := c.client.SubscribeInvoices(c.ctx, &lnrpc.InvoiceSubscription{})
		if err != nil {
			return "", err
		}

		log.Debug("invoice event subscription established")

		for {
			select {
			case <-c.ctx.Done():
				return "", nil
			default:
			}

			invoice, err := ev.Recv()
			if err != nil {
				return "", err // Return error to trigger retry
			}

			switch invoice.GetState() {
			case lnrpc.Invoice_SETTLED:
				// We check if there is a payment with this hash in our lnd instance.
				// If yes, it is a rebalancing payment, so we do not send an invoice event.
				ctx, cancel := context.WithCancel(c.ctx)
				stream, err := c.router.TrackPaymentV2(ctx, &routerrpc.TrackPaymentRequest{
					PaymentHash: invoice.RHash,
				})

				// If an error occurs here, we assume that there is no payment with this hash.
				if err != nil {
					c.eventSub <- events.NewInvoiceSettledEvent(invoice)
					cancel()
					continue
				}

				// The rpc error "payment isn't initiated" is returned, when fetching the first
				// element from the stream.
				if _, err := stream.Recv(); err != nil {
					c.eventSub <- events.NewInvoiceSettledEvent(invoice)
				}
				cancel()
			}
		}
	})
}

func (c *Client) handleFailedHtlcEvents() {
	log.Debug("starting failed htlc event handler")
	defer c.wg.Done()

	retry(c.ctx, "htlc event subscription", func() (string, error) {
		ev, err := c.router.SubscribeHtlcEvents(c.ctx, &routerrpc.SubscribeHtlcEventsRequest{})
		if err != nil {
			log.WithError(err).Error("error subscribing to failed htlc events")
			return "", err
		}

		for {
			select {
			case <-c.ctx.Done():
				return "", nil
			default:
			}

			htlcEvent, err := ev.Recv()
			if err != nil {
				return "", err
			}

			if htlcEvent.GetEventType() != routerrpc.HtlcEvent_FORWARD {
				log.WithField("htlc_event", htlcEvent).Trace("ignoring non-forward htlc event")
				continue
			}

			linkFailEvent := htlcEvent.GetLinkFailEvent()
			if linkFailEvent != nil {
				c.eventSub <- events.NewFailedHtlcLinkEvent(htlcEvent, linkFailEvent, c.channelManager)
			} else {
				log.WithField("htlc_event", htlcEvent).Trace("unhandled htlc event")
			}
		}
	})
}

func (c *Client) handleKeysendEvents() {
	log.Debug("keysend event handler")
	defer c.wg.Done()

	retry(c.ctx, "keysend event subscription", func() (string, error) {
		ev, err := c.client.SubscribeInvoices(c.ctx, &lnrpc.InvoiceSubscription{})
		if err != nil {
			return "", err
		}

		log.Debug("keysend (invoice) event subscription established")

		for {
			select {
			case <-c.ctx.Done():
				return "", nil
			default:
			}

			invoice, err := ev.Recv()
			if err != nil {
				return "", err // Return error to trigger retry
			}

			if invoice.GetState() != lnrpc.Invoice_SETTLED {
				continue
			}

			htlcs := invoice.GetHtlcs()
			for _, htlc := range htlcs {
				records := htlc.GetCustomRecords()

				// https://github.com/satoshisstream/satoshis.stream/blob/main/TLV_registry.md
				if msgBuf, ok := records[34349334]; ok {
					channel := c.channelManager.GetChannelById(htlc.ChanId)
					msg := string(msgBuf)

					c.eventSub <- events.NewKeysendEvent(msg, channel, htlc)
					break
				}
			}
		}
	})
}

func (c *Client) handlePaymentEvents() {
	log.Debug("starting payment event handler")
	defer c.wg.Done()

	retry(c.ctx, "payment event subscription", func() (string, error) {
		// Pubkey of the local node to distinguish between rebalancing and external payment
		var localPubkey string
		if info, err := c.client.GetInfo(c.ctx, &lnrpc.GetInfoRequest{}); err == nil {
			localPubkey = info.IdentityPubkey
		} else {
			return "", err
		}

		ev, err := c.router.TrackPayments(c.ctx, &routerrpc.TrackPaymentsRequest{})
		if err != nil {
			return "", err
		}

		log.Debug("payment event subscription established")

		for {
			select {
			case <-c.ctx.Done():
				return "", nil
			default:
			}

			payment, err := ev.Recv()
			if err != nil {
				return "", err // Return error to trigger retry
			}

			switch payment.Status {
			case lnrpc.Payment_SUCCEEDED:
				var payReq *lnrpc.PayReq
				if payment.PaymentRequest != "" {
					if decoded, err := c.client.DecodePayReq(c.ctx, &lnrpc.PayReqString{
						PayReq: payment.PaymentRequest,
					}); err == nil {
						payReq = decoded
					}
				}

				var recPubkey string
				// The pub_key of the last hop (receiver of the payment) has to be identical
				// for all htlcs. Hence we use the first htlc.
				if len(payment.Htlcs) > 0 && len(payment.Htlcs[0].Route.Hops) > 0 {
					lastHop := payment.Htlcs[0].Route.Hops[len(payment.Htlcs[0].Route.Hops)-1]
					recPubkey = lastHop.PubKey
				}

				isRebalancing := recPubkey == localPubkey
				c.eventSub <- events.NewPaymentSucceededEvent(payment, payReq, isRebalancing, c.getAlias)
			}
		}
	})
}

func (c *Client) handleOnChainEvents() {
	log.Debug("starting on chain event handler")
	defer c.wg.Done()

	retry(c.ctx, "on chain event subscription", func() (string, error) {
		ev, err := c.client.SubscribeTransactions(c.ctx, &lnrpc.GetTransactionsRequest{})
		if err != nil {
			return "", err
		}

		log.Debug("on chain event subscription established")

		for {
			select {
			case <-c.ctx.Done():
				return "", nil
			default:
			}

			event, err := ev.Recv()
			if err != nil {
				return "", err // Return error to trigger retry
			}

			confirmCnt := event.GetNumConfirmations()
			if confirmCnt == 0 || confirmCnt == 1 {
				c.eventSub <- events.NewOnChainTransactionEvent(event)
				c.pendChanManager.RefreshDelayed()
			}
		}
	})
}

func (c *Client) handlePendingChannels() {
	log.Debug("starting pending channel event handler")
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		case update := <-c.pendChanUpdates:
			switch ev := update.(type) {
			case *lnrpc.PendingChannelsResponse_PendingOpenChannel:
				c.eventSub <- events.NewChannelOpeningEvent(ev, c.getAlias)
			case *lnrpc.PendingChannelsResponse_WaitingCloseChannel:
				c.eventSub <- events.NewChannelClosingEvent(ev, c.getAlias)
			default:
				log.WithField("update", update).Warn("unknown pending channel update type")
			}
		}
	}
}

func (c *Client) handleChainSyncState() {
	log.Debug("starting sync state event handler")
	defer c.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	unsyncedThreshold := c.cfg.EventConfig.ChainLostEvent.Threshold
	warningInterval := c.cfg.EventConfig.ChainLostEvent.WarningInterval

	var lastUnsyncedTime *time.Time
	var lastWarningTime *time.Time

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			log.Debug("polling for sync state")

			info, err := c.client.GetInfo(c.ctx, &lnrpc.GetInfoRequest{})
			if err != nil {
				log.WithError(err).Error("error fetching node info")
				continue
			}

			if info.GetSyncedToChain() {
				if lastWarningTime != nil {
					log.Debug("chain sync restored")
					c.eventSub <- events.NewChainSyncRestoredEvent(time.Since(*lastUnsyncedTime))
					lastUnsyncedTime = nil
					lastWarningTime = nil
				}
			} else {
				now := time.Now()
				if lastUnsyncedTime == nil {
					// first time we detect chain is not synced
					lastUnsyncedTime = &now
					log.Debug("chain sync lost, starting timer")
				} else {
					unsyncedDuration := now.Sub(*lastUnsyncedTime)
					if unsyncedDuration >= unsyncedThreshold {
						shouldWarn := false
						if lastWarningTime == nil {
							// initial warning after threshold
							shouldWarn = true
						} else if now.Sub(*lastWarningTime) >= warningInterval {
							// oh no, it's time for another warning
							shouldWarn = true
						}

						if shouldWarn {
							// chain has been unsynced for longer than threshold
							c.eventSub <- events.NewChainSyncLostEvent(unsyncedDuration)
							lastWarningTime = &now
						}
					}
				}
			}
		}
	}
}

func (c *Client) handleBackupEvents() {
	log.Debug("starting backup event handler")
	defer c.wg.Done()

	retry(c.ctx, "channel backup subscription", func() (string, error) {
		ev, err := c.client.SubscribeChannelBackups(c.ctx, &lnrpc.ChannelBackupSubscription{})
		if err != nil {
			return "", err
		}

		log.Debug("channel backup subscription established")

		for {
			select {
			case <-c.ctx.Done():
				return "", nil
			default:
			}

			backup, err := ev.Recv()
			if err != nil {
				return "", err // Return error to trigger retry
			}

			if multiBackup := backup.GetMultiChanBackup(); multiBackup != nil {
				c.eventSub <- events.NewBackupMultiEvent(multiBackup)
			}
		}
	})
}

func (c *Client) handleChannelStatusEvents() {
	downChannelMap := make(map[uint64]time.Time)
	notifiedChannels := make(map[uint64]bool)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	log.Debug("starting channel status event handler")
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			log.Debug("checking channel status")

			channels := c.channelManager.GetAllChannels()
			now := time.Now()

			for _, channel := range channels {
				chanId := channel.ChanId
				logger := log.WithField("channel_id", chanId)

				if channel.GetActive() {
					// Channel is active, remove from downChannelMap if present
					if _, exists := downChannelMap[chanId]; exists {
						downDuration := now.Sub(downChannelMap[chanId])

						logger.Debug("channel is back up")

						delete(downChannelMap, chanId)
						delete(notifiedChannels, chanId) // Clear notification flag
						c.eventSub <- events.NewChannelStatusUpEvent(channel, downDuration, c.getAlias)
					}
				} else {
					// Channel is inactive, track downtime
					if downStartTime, exists := downChannelMap[channel.ChanId]; exists {
						downDuration := now.Sub(downStartTime)

						logger.Debug("channel is still down, checking duration")

						if downDuration >= c.cfg.EventConfig.ChannelStatusEvent.MinDowntime && !notifiedChannels[chanId] {
							logger.Debug("channel has been down long enough, sending event")
							c.eventSub <- events.NewChannelStatusDownEvent(channel, downDuration, c.getAlias)
							notifiedChannels[chanId] = true
						}
					} else {
						// First time seeing this channel as down, start tracking
						downChannelMap[chanId] = now
						notifiedChannels[chanId] = false
						logger.Debug("channel is down, starting downtime tracking")
					}
				}
			}
		}
	}
}

// getAlias returns the alias for a given pubkey. If an error occurs, it returns the first
// 8 characters of the pubkey.
func (c *Client) getAlias(pubkey string) string {
	if nodeInfo, err := c.client.GetNodeInfo(c.ctx, &lnrpc.NodeInfoRequest{
		PubKey: pubkey,
	}); err == nil {
		return nodeInfo.Node.Alias
	}
	return format.FormatPubKey(pubkey)
}

func retry(ctx context.Context, name string, operation backoff.Operation[string]) {
	logger := log.WithField("name", name)
	notify := func(err error, duration time.Duration) {
		logger.WithError(err).WithField("next_retry_in", duration).Warn("operation failed, retrying")
	}

	_, err := backoff.Retry(ctx, operation, backoff.WithNotify(notify), backoff.WithMaxElapsedTime(0))
	if err != nil {
		if ctx.Err() != nil {
			logger.Debug("context cancelled, stopping retry")
		} else {
			logger.WithError(err).Error("operation failed permanently")
		}
	}
}
