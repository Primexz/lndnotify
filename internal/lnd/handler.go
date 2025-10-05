package lnd

import (
	"time"

	"github.com/Primexz/lndnotify/internal/events"
	"github.com/lightningnetwork/lnd/lnrpc"
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

	ev, err := c.client.SubscribePeerEvents(c.ctx, &lnrpc.PeerEventSubscription{})
	if err != nil {
		log.WithError(err).Error("error subscribing to peer events")
		return
	}

	for {
		peerEvent, err := ev.Recv()
		if err != nil {
			log.WithError(err).Error("error receiving peer event")
			return
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
}

// handleChannelEvents handles channel open and close events
func (c *Client) handleChannelEvents() {
	log.Debug("starting channel event handler")
	defer c.wg.Done()

	ev, err := c.client.SubscribeChannelEvents(c.ctx, &lnrpc.ChannelEventSubscription{})
	if err != nil {
		log.WithError(err).Error("error subscribing to channel events")
		return
	}

	for {
		peerEvent, err := ev.Recv()
		if err != nil {
			log.WithError(err).Error("error receiving channel event")
			return
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
}
