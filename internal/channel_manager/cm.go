package channelmanager

import (
	"context"
	"sync"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	log "github.com/sirupsen/logrus"
)

type ChannelManager struct {
	client   lnrpc.LightningClient
	channels map[uint64]*lnrpc.Channel
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup

	refreshInterval time.Duration
}

func NewChannelManager(client lnrpc.LightningClient) *ChannelManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ChannelManager{
		client:          client,
		channels:        make(map[uint64]*lnrpc.Channel),
		ctx:             ctx,
		cancel:          cancel,
		refreshInterval: 5 * time.Minute,
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
}

func (cm *ChannelManager) GetChannelById(chanID uint64) *lnrpc.Channel {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	ch, exists := cm.channels[chanID]
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
	}

	log.WithField("channel_count", len(cm.channels)).Debug("channel state refreshed")
	return nil
}

func (cm *ChannelManager) SetRefreshInterval(interval time.Duration) {
	cm.refreshInterval = interval
}
