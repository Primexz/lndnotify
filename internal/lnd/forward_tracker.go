package lnd

import (
	"fmt"
	"sync"

	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
)

type ForwardTracker struct {
	mu         sync.RWMutex
	forwardMap map[string]*routerrpc.HtlcEvent
}

func NewForwardTracker() *ForwardTracker {
	return &ForwardTracker{
		forwardMap: make(map[string]*routerrpc.HtlcEvent),
	}
}

func (ft *ForwardTracker) AddForward(event *routerrpc.HtlcEvent) {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	key := ft.getHtlcKey(event)
	ft.forwardMap[key] = event
}

func (ft *ForwardTracker) RemoveForward(event *routerrpc.HtlcEvent) bool {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	key := ft.getHtlcKey(event)
	if _, exists := ft.forwardMap[key]; exists {
		delete(ft.forwardMap, key)
		return true
	}
	return false
}

func (ft *ForwardTracker) GetForward(event *routerrpc.HtlcEvent) (*routerrpc.HtlcEvent, bool) {
	ft.mu.RLock()
	defer ft.mu.RUnlock()

	key := ft.getHtlcKey(event)
	forwardEvent, exists := ft.forwardMap[key]
	return forwardEvent, exists
}

// Count returns the number of forwards currently being tracked
func (ft *ForwardTracker) Count() int {
	ft.mu.RLock()
	defer ft.mu.RUnlock()
	return len(ft.forwardMap)
}

func (ft *ForwardTracker) getHtlcKey(event *routerrpc.HtlcEvent) string {
	return fmt.Sprintf("%d%d%d%d", event.IncomingChannelId, event.OutgoingChannelId, event.IncomingHtlcId, event.OutgoingHtlcId)
}
