package channelmanager

import (
	"context"
	"sync"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	log "github.com/sirupsen/logrus"
)

// FeeChangeEvent represents a change in a channel's fee policy
type FeeChangeEvent struct {
	Channel     *lnrpc.Channel
	ChannelEdge *lnrpc.ChannelEdge
	Timestamp   time.Time

	OldFeeRate int64
	NewFeeRate int64

	OldBaseFee int64
	NewBaseFee int64

	OldInboundFeeRate int32
	NewInboundFeeRate int32

	OldInboundBaseFee int32
	NewInboundBaseFee int32
}

type ChannelManager struct {
	client       lnrpc.LightningClient
	channels     map[uint64]*lnrpc.Channel
	channelEdges map[uint64]*lnrpc.ChannelEdge
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup

	refreshInterval time.Duration

	feeChangeCh chan FeeChangeEvent
}

func NewChannelManager(client lnrpc.LightningClient) *ChannelManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ChannelManager{
		client:          client,
		channels:        make(map[uint64]*lnrpc.Channel),
		channelEdges:    make(map[uint64]*lnrpc.ChannelEdge),
		ctx:             ctx,
		cancel:          cancel,
		refreshInterval: 5 * time.Minute,
		feeChangeCh:     make(chan FeeChangeEvent, 100), // Buffered channel for fee changes
	}
}

// Start initializes the channel manager and begins periodic refreshes
func (cm *ChannelManager) Start() error {
	log.Debug("starting channel manager")

	if err := cm.refreshChannels(); err != nil {
		return err
	}

	cm.wg.Add(1)
	go cm.refreshLoop()

	return nil
}

// Stop stops the channel manager
func (cm *ChannelManager) Stop() {
	log.Debug("stopping channel manager")
	cm.cancel()
	cm.wg.Wait()
	close(cm.feeChangeCh)
}

// GetChannelById retrieves a channel by its ID
func (cm *ChannelManager) GetChannelById(chanId uint64) *lnrpc.Channel {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	ch, exists := cm.channels[chanId]
	if !exists {
		return nil
	}
	return ch
}

// GetAllChannels returns a slice of all managed channels
func (cm *ChannelManager) GetAllChannels() []*lnrpc.Channel {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	channels := make([]*lnrpc.Channel, 0, len(cm.channels))
	for _, ch := range cm.channels {
		channels = append(channels, ch)
	}
	return channels
}

// RefreshNow triggers an immediate refresh of channel states
func (cm *ChannelManager) RefreshNow() error {
	return cm.refreshChannels()
}

// SetRefreshInterval sets the interval for periodic channel state refreshes
func (cm *ChannelManager) SetRefreshInterval(interval time.Duration) {
	cm.refreshInterval = interval
}

// GetFeeChangeChannel returns the channel for receiving fee change events
func (cm *ChannelManager) GetFeeChangeChannel() <-chan FeeChangeEvent {
	return cm.feeChangeCh
}

func (cm *ChannelManager) refreshLoop() {
	defer cm.wg.Done()

	ticker := time.NewTicker(cm.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-ticker.C:
			if err := cm.refreshChannels(); err != nil {
				log.WithError(err).Error("failed to refresh channels")
			}
		}
	}
}

func (cm *ChannelManager) refreshChannels() error {
	log.Debug("refreshing channel state")

	resp, err := cm.client.ListChannels(cm.ctx, &lnrpc.ListChannelsRequest{
		PeerAliasLookup: true,
	})
	if err != nil {
		return err
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.channels = make(map[uint64]*lnrpc.Channel)
	for _, ch := range resp.Channels {
		logger := log.WithFields(log.Fields{
			"channel_id": ch.ChanId,
			"peer":       ch.RemotePubkey,
		})
		chanEdge, err := cm.client.GetChanInfo(cm.ctx, &lnrpc.ChanInfoRequest{
			ChanId: ch.ChanId,
		})
		if err != nil {
			logger.WithError(err).Warn("failed to get channel info for fee change detection")
			continue
		}

		oldEdge := cm.channelEdges[ch.ChanId]
		cm.channelEdges[ch.ChanId] = chanEdge
		cm.channels[ch.ChanId] = ch

		cm.checkFeeChanges(ch, chanEdge, oldEdge)
	}

	log.WithField("channel_count", len(cm.channels)).Debug("channel state refreshed")
	return nil
}

// checkFeeChanges compares old and new channel states to detect fee changes
func (cm *ChannelManager) checkFeeChanges(ch *lnrpc.Channel, newEdge *lnrpc.ChannelEdge, oldEdge *lnrpc.ChannelEdge) {
	if oldEdge == nil || newEdge == nil {
		return
	}

	logger := log.WithFields(log.Fields{
		"channel_id": ch.ChanId,
		"peer":       ch.RemotePubkey,
	})

	var oldRemotePolicy, newRemotePolicy *lnrpc.RoutingPolicy
	if newEdge.Node1Pub == ch.RemotePubkey {
		newRemotePolicy = newEdge.Node1Policy
		oldRemotePolicy = oldEdge.Node1Policy
	} else if newEdge.Node2Pub == ch.RemotePubkey {
		newRemotePolicy = newEdge.Node2Policy
		oldRemotePolicy = oldEdge.Node2Policy
	} else {
		logger.Warn("could not identify remote peer in channel edge")
		return
	}

	if oldRemotePolicy == nil || newRemotePolicy == nil {
		logger.Warn("missing routing policy for fee change detection")
		return
	}

	if oldRemotePolicy.FeeRateMilliMsat == newRemotePolicy.FeeRateMilliMsat &&
		oldRemotePolicy.FeeBaseMsat == newRemotePolicy.FeeBaseMsat &&
		oldRemotePolicy.InboundFeeRateMilliMsat == newRemotePolicy.InboundFeeRateMilliMsat &&
		oldRemotePolicy.InboundFeeBaseMsat == newRemotePolicy.InboundFeeBaseMsat {
		logger.WithFields(log.Fields{
			"rate": oldRemotePolicy.FeeRateMilliMsat,
			"base": oldRemotePolicy.FeeBaseMsat,
		}).Trace("no change in remote peer fees detected")
		return
	}

	logger.WithFields(log.Fields{
		"old_fee_rate": oldRemotePolicy.FeeRateMilliMsat,
		"new_fee_rate": newRemotePolicy.FeeRateMilliMsat,
		"old_base_fee": oldRemotePolicy.FeeBaseMsat,
		"new_base_fee": newRemotePolicy.FeeBaseMsat,
	}).Debug("detected remote peer fee change")

	cm.feeChangeCh <- FeeChangeEvent{
		Channel:           ch,
		ChannelEdge:       newEdge,
		OldFeeRate:        oldRemotePolicy.FeeRateMilliMsat,
		NewFeeRate:        newRemotePolicy.FeeRateMilliMsat,
		OldBaseFee:        oldRemotePolicy.FeeBaseMsat,
		NewBaseFee:        newRemotePolicy.FeeBaseMsat,
		OldInboundFeeRate: oldRemotePolicy.InboundFeeRateMilliMsat,
		NewInboundFeeRate: newRemotePolicy.InboundFeeRateMilliMsat,
		OldInboundBaseFee: oldRemotePolicy.InboundFeeBaseMsat,
		NewInboundBaseFee: newRemotePolicy.InboundFeeBaseMsat,
		Timestamp:         time.Now(),
	}
}
