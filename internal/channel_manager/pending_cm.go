package channelmanager

import (
	"context"
	"sync"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type PendingChannelManager struct {
	client lnrpc.LightningClient

	// Sets of channel points we have seen over the lifetime of the manager.
	chanPointsOpening map[string]struct{}
	chanPointsClosing map[string]struct{}

	firstPollDone  bool
	pendingUpdates chan proto.Message
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup

	refreshInterval time.Duration
	refreshDelay    time.Duration
}

func NewPendingChannelManager(client lnrpc.LightningClient, pendingUpdates chan proto.Message) *PendingChannelManager {

	ctx, cancel := context.WithCancel(context.Background())
	return &PendingChannelManager{
		client:            client,
		chanPointsOpening: make(map[string]struct{}),
		chanPointsClosing: make(map[string]struct{}),
		pendingUpdates:    pendingUpdates,
		ctx:               ctx,
		cancel:            cancel,
		refreshInterval:   1 * time.Minute,
		refreshDelay:      5 * time.Second,
	}
}

func (cm *PendingChannelManager) Start() error {
	log.Debug("starting pending channel manager")

	if err := cm.refreshChannels(); err != nil {
		return err
	}

	cm.wg.Add(1)
	go cm.refreshLoop()

	return nil
}

func (cm *PendingChannelManager) Stop() {
	log.Debug("stopping pending channel manager")
	cm.cancel()
	cm.wg.Wait()
}

// RefreshDelayed refreshes the pending channels after a short delay to avoid
// data inconsistencies that may occur when called immediately after a channel
// event is received (e.g., missing closing txid and hex for closed channels).
func (cm *PendingChannelManager) RefreshDelayed() {
	cm.wg.Add(1)
	go func() {
		defer cm.wg.Done()
		select {
		case <-cm.ctx.Done():
			return
		case <-time.After(cm.refreshDelay):
			if err := cm.refreshChannels(); err != nil {
				log.WithError(err).Error("failed to refresh pending channels")
			}
		}
	}()
}

func (cm *PendingChannelManager) refreshLoop() {
	defer cm.wg.Done()

	ticker := time.NewTicker(cm.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-ticker.C:
			if err := cm.refreshChannels(); err != nil {
				log.WithError(err).Error("failed to refresh pending channels")
			}
		}
	}
}

func (cm *PendingChannelManager) refreshChannels() error {
	log.Debug("refreshing pending channel state")

	resp, err := cm.client.PendingChannels(cm.ctx, &lnrpc.PendingChannelsRequest{
		IncludeRawTx: true,
	})
	if err != nil {
		log.WithError(err).Error("error fetching pending channels")
		return err
	}

	log.WithFields(log.Fields{
		"opening_count": len(resp.PendingOpenChannels),
		"closing_count": len(resp.WaitingCloseChannels),
	}).Debug("fetched pending channels")

	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, channel := range resp.PendingOpenChannels {
		if channel == nil || channel.Channel == nil {
			continue
		}

		chanPoint := channel.Channel.ChannelPoint

		if _, ok := cm.chanPointsOpening[chanPoint]; !ok && cm.firstPollDone {
			cm.pendingUpdates <- channel
		}
		cm.chanPointsOpening[chanPoint] = struct{}{}
	}

	// WaitingCloseChannels includes both cooperatively and force-closed channels
	// which are not confirmed on-chain yet.
	// But force-closed channels only if they are initiated by the local node.
	// LND is not scanning the mempool, hence remote force-closed channels
	// are not seen until they are confirmed on-chain.
	for _, channel := range resp.WaitingCloseChannels {
		if channel == nil || channel.Channel == nil {
			continue
		}

		chanPoint := channel.Channel.ChannelPoint

		if _, ok := cm.chanPointsClosing[chanPoint]; !ok && cm.firstPollDone {
			cm.pendingUpdates <- channel
		}
		cm.chanPointsClosing[chanPoint] = struct{}{}
	}
	cm.firstPollDone = true

	return nil
}
