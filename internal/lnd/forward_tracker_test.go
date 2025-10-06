package lnd

import (
	"testing"

	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/stretchr/testify/assert"
)

func TestNewForwardTracker(t *testing.T) {
	ft := NewForwardTracker()

	assert.NotNil(t, ft)
	assert.NotNil(t, ft.forwardMap)
	assert.Equal(t, 0, len(ft.forwardMap))
}

func TestForwardTracker_AddForward(t *testing.T) {
	ft := NewForwardTracker()

	event := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    2,
	}

	ft.AddForward(event)

	// Verify the event was added
	assert.Equal(t, 1, len(ft.forwardMap))

	// Verify we can retrieve it
	retrievedEvent, exists := ft.GetForward(event)
	assert.True(t, exists)
	assert.Equal(t, event, retrievedEvent)
}

func TestForwardTracker_AddMultipleForwards(t *testing.T) {
	ft := NewForwardTracker()

	event1 := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    2,
	}

	event2 := &routerrpc.HtlcEvent{
		IncomingChannelId: 111111111,
		OutgoingChannelId: 222222222,
		IncomingHtlcId:    3,
		OutgoingHtlcId:    4,
	}

	ft.AddForward(event1)
	ft.AddForward(event2)

	assert.Equal(t, 2, len(ft.forwardMap))

	// Verify both events can be retrieved
	retrievedEvent1, exists1 := ft.GetForward(event1)
	assert.True(t, exists1)
	assert.Equal(t, event1, retrievedEvent1)

	retrievedEvent2, exists2 := ft.GetForward(event2)
	assert.True(t, exists2)
	assert.Equal(t, event2, retrievedEvent2)
}

func TestForwardTracker_AddSameForwardTwice(t *testing.T) {
	ft := NewForwardTracker()

	event := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    2,
	}

	ft.AddForward(event)
	ft.AddForward(event) // Add the same event again

	// Should still only have one entry (overwritten)
	assert.Equal(t, 1, len(ft.forwardMap))

	retrievedEvent, exists := ft.GetForward(event)
	assert.True(t, exists)
	assert.Equal(t, event, retrievedEvent)
}

func TestForwardTracker_RemoveForward(t *testing.T) {
	ft := NewForwardTracker()

	event := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    2,
	}

	// Add the event first
	ft.AddForward(event)
	assert.Equal(t, 1, len(ft.forwardMap))

	// Remove the event
	removed := ft.RemoveForward(event)
	assert.True(t, removed)
	assert.Equal(t, 0, len(ft.forwardMap))

	// Verify it's no longer retrievable
	_, exists := ft.GetForward(event)
	assert.False(t, exists)
}

func TestForwardTracker_RemoveNonExistentForward(t *testing.T) {
	ft := NewForwardTracker()

	event := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    2,
	}

	// Try to remove an event that was never added
	removed := ft.RemoveForward(event)
	assert.False(t, removed)
	assert.Equal(t, 0, len(ft.forwardMap))
}

func TestForwardTracker_GetNonExistentForward(t *testing.T) {
	ft := NewForwardTracker()

	event := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    2,
	}

	// Try to get an event that was never added
	retrievedEvent, exists := ft.GetForward(event)
	assert.False(t, exists)
	assert.Nil(t, retrievedEvent)
}

func TestForwardTracker_getHtlcKey(t *testing.T) {
	ft := NewForwardTracker()

	event := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    2,
	}

	key := ft.getHtlcKey(event)
	expectedKey := "12345678998765432112"
	assert.Equal(t, expectedKey, key)
}

func TestForwardTracker_getHtlcKeyUniqueness(t *testing.T) {
	ft := NewForwardTracker()

	event1 := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    2,
	}

	event2 := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    3, // Different outgoing HTLC ID
	}

	key1 := ft.getHtlcKey(event1)
	key2 := ft.getHtlcKey(event2)

	assert.NotEqual(t, key1, key2, "Keys should be unique for different events")
}

func TestForwardTracker_Count(t *testing.T) {
	ft := NewForwardTracker()

	// Initially should be empty
	assert.Equal(t, 0, ft.Count())

	// Add one event
	event1 := &routerrpc.HtlcEvent{
		IncomingChannelId: 123456789,
		OutgoingChannelId: 987654321,
		IncomingHtlcId:    1,
		OutgoingHtlcId:    2,
	}
	ft.AddForward(event1)
	assert.Equal(t, 1, ft.Count())

	// Add another event
	event2 := &routerrpc.HtlcEvent{
		IncomingChannelId: 111111111,
		OutgoingChannelId: 222222222,
		IncomingHtlcId:    3,
		OutgoingHtlcId:    4,
	}
	ft.AddForward(event2)
	assert.Equal(t, 2, ft.Count())

	// Remove one event
	removed := ft.RemoveForward(event1)
	assert.True(t, removed)
	assert.Equal(t, 1, ft.Count())

	// Remove the last event
	removed = ft.RemoveForward(event2)
	assert.True(t, removed)
	assert.Equal(t, 0, ft.Count())
}
