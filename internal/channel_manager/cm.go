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
	OldFeeRate  int64
	NewFeeRate  int64
	OldBaseFee  int64
	NewBaseFee  int64
	Timestamp   time.Time
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

func (cm *ChannelManager) Start() error {
	log.Debug("starting channel manager")

	if err := cm.refreshChannels(); err != nil {
		return err
	}

	cm.wg.Add(1)
	go cm.refreshLoop()

	return nil
}

func (cm *ChannelManager) Stop() {
	log.Debug("stopping channel manager")
	cm.cancel()
	cm.wg.Wait()
	close(cm.feeChangeCh)
}

func (cm *ChannelManager) GetChannelById(chanId uint64) *lnrpc.Channel {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	ch, exists := cm.channels[chanId]
	if !exists {
		return nil
	}
	return ch
}

func (cm *ChannelManager) GetAllChannels() []*lnrpc.Channel {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	channels := make([]*lnrpc.Channel, 0, len(cm.channels))
	for _, ch := range cm.channels {
		channels = append(channels, ch)
	}
	return channels
}

func (cm *ChannelManager) RefreshNow() error {
	return cm.refreshChannels()
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
		cm.channels[ch.ChanId] = ch
		cm.checkFeeChanges(ch)
	}

	log.WithField("channel_count", len(cm.channels)).Debug("channel state refreshed")
	return nil
}

func (cm *ChannelManager) SetRefreshInterval(interval time.Duration) {
	cm.refreshInterval = interval
}

// GetFeeChangeChannel returns the channel for receiving fee change events
func (cm *ChannelManager) GetFeeChangeChannel() <-chan FeeChangeEvent {
	return cm.feeChangeCh
}

// checkFeeChanges compares old and new channel states to detect fee changes
func (cm *ChannelManager) checkFeeChanges(ch *lnrpc.Channel) {
	logger := log.WithFields(log.Fields{
		"channel_id": ch.ChanId,
		"peer":       ch.RemotePubkey,
	})

	// Get current channel edge information to check fee policies
	chanInfo, err := cm.client.GetChanInfo(cm.ctx, &lnrpc.ChanInfoRequest{
		ChanId: ch.ChanId,
	})
	if err != nil {
		logger.WithError(err).Warn("failed to get channel info for fee change detection")
		return
	}

	oldEdge := cm.channelEdges[ch.ChanId]
	cm.channelEdges[ch.ChanId] = chanInfo

	// If we don't have previous edge info, this is the first time we see this channel
	if oldEdge == nil {
		return
	}

	var oldRemotePolicy, newRemotePolicy *lnrpc.RoutingPolicy
	if chanInfo.Node1Pub == ch.RemotePubkey {
		newRemotePolicy = chanInfo.Node1Policy
		oldRemotePolicy = oldEdge.Node1Policy
	} else if chanInfo.Node2Pub == ch.RemotePubkey {
		newRemotePolicy = chanInfo.Node2Policy
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
		oldRemotePolicy.FeeBaseMsat == newRemotePolicy.FeeBaseMsat {
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
		Channel:     ch,
		ChannelEdge: chanInfo,
		OldFeeRate:  oldRemotePolicy.FeeRateMilliMsat,
		NewFeeRate:  newRemotePolicy.FeeRateMilliMsat,
		OldBaseFee:  oldRemotePolicy.FeeBaseMsat,
		NewBaseFee:  newRemotePolicy.FeeBaseMsat,
		Timestamp:   time.Now(),
	}
}
