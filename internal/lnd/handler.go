package lnd

import (
	"fmt"
	"time"

	"github.com/Primexz/lndnotify/internal/events"
	"github.com/lightningnetwork/lnd/lnrpc"
)

// handleForwards polls for forwarding events
func (c *Client) handleForwards() {
	defer c.wg.Done()

	start := time.Now()
	for range time.Tick(time.Minute * 1) {
		resp, err := c.client.ForwardingHistory(c.ctx, &lnrpc.ForwardingHistoryRequest{
			StartTime:       uint64(start.Unix()), // #nosec G115
			PeerAliasLookup: true,
		})
		if err != nil {
			fmt.Printf("Error fetching forwarding history: %v\n", err)
			continue
		}

		forwards := resp.GetForwardingEvents()
		for _, fwd := range forwards {
			c.eventSub <- events.NewForwardEvent(
				fwd.PeerAliasIn,
				fwd.PeerAliasOut,
				fwd.AmtInMsat,
				fwd.AmtOutMsat,
				fwd.FeeMsat,
			)
		}

		// push start time forward
		if len(forwards) > 0 {
			start = time.Now()
		}
	}
}

// handlePeerEvents handles peer connection and disconnection events
func (c *Client) handlePeerEvents() {
	defer c.wg.Done()

	ev, err := c.client.SubscribePeerEvents(c.ctx, &lnrpc.PeerEventSubscription{})
	if err != nil {
		fmt.Printf("Error subscribing to peer events: %v\n", err)
		return
	}

	for {
		peerEvent, err := ev.Recv()
		if err != nil {
			fmt.Printf("Error receiving peer event: %v\n", err)
			return
		}

		nodeInfo, err := c.client.GetNodeInfo(c.ctx, &lnrpc.NodeInfoRequest{
			PubKey: peerEvent.GetPubKey(),
		})
		if err != nil {
			fmt.Printf("Error fetching node info: %v\n", err)
			continue
		}

		switch peerEvent.GetType() {
		case lnrpc.PeerEvent_PEER_ONLINE:
			c.eventSub <- events.NewPeerOnlineEvent(nodeInfo.Node.Alias)
		case lnrpc.PeerEvent_PEER_OFFLINE:
			c.eventSub <- events.NewPeerOfflineEvent(nodeInfo.Node.Alias)
		}
	}
}

// handleChannelEvents handles channel open and close events
func (c *Client) handleChannelEvents() {
	defer c.wg.Done()

	ev, err := c.client.SubscribeChannelEvents(c.ctx, &lnrpc.ChannelEventSubscription{})
	if err != nil {
		fmt.Printf("Error subscribing to peer events: %v\n", err)
		return
	}

	for {
		peerEvent, err := ev.Recv()
		if err != nil {
			fmt.Printf("Error receiving peer event: %v\n", err)
			return
		}

		switch peerEvent.GetType() {
		case lnrpc.ChannelEventUpdate_OPEN_CHANNEL:
			channel := peerEvent.GetOpenChannel()
			nodeInfo, err := c.client.GetNodeInfo(c.ctx, &lnrpc.NodeInfoRequest{
				PubKey: channel.RemotePubkey,
			})
			if err != nil {
				fmt.Printf("Error fetching node info: %v\n", err)
				continue
			}

			c.eventSub <- events.NewChannelOpenEvent(nodeInfo.Node.Alias, channel.Capacity)
		case lnrpc.ChannelEventUpdate_CLOSED_CHANNEL:
			channel := peerEvent.GetClosedChannel()
			nodeInfo, err := c.client.GetNodeInfo(c.ctx, &lnrpc.NodeInfoRequest{
				PubKey: channel.RemotePubkey,
			})
			if err != nil {
				fmt.Printf("Error fetching node info: %v\n", err)
				continue
			}

			c.eventSub <- events.NewChannelCloseEvent(nodeInfo.Node.Alias, channel.SettledBalance)
		}
	}
}
