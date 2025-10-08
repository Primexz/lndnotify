package lnd

import (
	"context"
	"time"

	"github.com/Primexz/lndnotify/internal/events"
	"github.com/cenkalti/backoff/v5"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	log "github.com/sirupsen/logrus"
)

// handleForwards polls for forwarding events
func (c *Client) handleForwards() {
	log.Debug("starting forward event handler")
	defer c.wg.Done()

	start := time.Now()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			log.WithField("since", start).Debug("polling for forwarding events")
			resp, err := c.client.ForwardingHistory(c.ctx, &lnrpc.ForwardingHistoryRequest{
				StartTime:       uint64(start.Unix()), // #nosec G115
				PeerAliasLookup: true,
			})
			if err != nil {
				log.WithError(err).Error("error fetching forwarding history")
				continue
			}

			forwards := resp.GetForwardingEvents()
			for _, fwd := range forwards {
				c.eventSub <- events.NewForwardEvent(fwd)
			}

			// push start time forward
			if len(forwards) > 0 {
				start = time.Now()
			}
		}
	}
}

// handlePeerEvents handles peer connection and disconnection events
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
				log.WithField("pubkey", peerEvent.GetPubKey()).WithError(err).Error("error fetching node info")
				continue
			}

			switch peerEvent.GetType() {
			case lnrpc.PeerEvent_PEER_ONLINE:
				c.eventSub <- events.NewPeerOnlineEvent(nodeInfo.Node)
			case lnrpc.PeerEvent_PEER_OFFLINE:
				c.eventSub <- events.NewPeerOfflineEvent(nodeInfo.Node)
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

			peerEvent, err := ev.Recv()
			if err != nil {
				return "", err // Return error to trigger retry
			}

			switch peerEvent.GetType() {
			case lnrpc.ChannelEventUpdate_OPEN_CHANNEL:
				channel := peerEvent.GetOpenChannel()
				nodeInfo, err := c.client.GetNodeInfo(c.ctx, &lnrpc.NodeInfoRequest{
					PubKey: channel.RemotePubkey,
				})
				if err != nil {
					log.WithError(err).Error("error fetching node info")
					continue
				}

				c.eventSub <- events.NewChannelOpenEvent(nodeInfo.Node, channel)
			case lnrpc.ChannelEventUpdate_CLOSED_CHANNEL:
				channel := peerEvent.GetClosedChannel()
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
				c.eventSub <- events.NewInvoiceSettledEvent(invoice)
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
				log.WithError(err).Error("error receiving failed htlc event")
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

func retry(ctx context.Context, name string, operation backoff.Operation[string]) {
	logger := log.WithField("name", name)
	notify := func(err error, duration time.Duration) {
		logger.WithError(err).WithField("next_retry_in", duration).WithError(err).Warn("operation failed, retrying")
	}

	_, err := backoff.Retry(ctx, operation, backoff.WithNotify(notify))
	if err != nil {
		if ctx.Err() != nil {
			logger.Debug("context cancelled, stopping retry")
		} else {
			logger.WithError(err).Error("operation failed permanently")
		}
	}
}
