package channel

import "github.com/lightningnetwork/lnd/lnrpc"

func GetChannelById(channels []*lnrpc.Channel, chanId uint64) *lnrpc.Channel {
	for _, ch := range channels {
		if ch.ChanId == chanId {
			return ch
		}
	}
	return nil
}
