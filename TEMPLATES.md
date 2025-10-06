# Template Variables Documentation

This document describes all available template variables that can be used in notification templates for each event type in lndnotify.

## Forward Event
Triggered when a payment is forwarded through your node.

| Variable | Description |
|----------|-------------|
| `{{.PeerAliasIn}}` | The alias of the peer that sent the payment to your node |
| `{{.PeerAliasOut}}` | The alias of the peer that received the payment from your node |
| `{{.Amount}}` | The amount of the payment in satoshis (formatted) |
| `{{.Fee}}` | The fee earned from forwarding the payment (formatted) |

## Invoice Settled Event
Triggered when an invoice is paid/settled.

| Variable | Description |
|----------|-------------|
| `{{.Memo}}` | The memo/description attached to the invoice |
| `{{.Value}}` | The value of the invoice in satoshis (formatted) |
| `{{.IsKeysend}}` | Boolean indicating if this was a keysend payment |
| `{{.PaymentRequest}}` | The original payment request string (invoice) |

## Peer Online Event
Triggered when a peer connects to your node.

| Variable | Description |
|----------|-------------|
| `{{.PeerAlias}}` | The alias of the peer that came online |
| `{{.PeerPubKey}}` | The full public key of the peer |
| `{{.PeerPubkeyShort}}` | A shortened version of the peer's public key |

## Peer Offline Event
Triggered when a peer disconnects from your node.

| Variable | Description |
|----------|-------------|
| `{{.PeerAlias}}` | The alias of the peer that went offline |
| `{{.PeerPubKey}}` | The full public key of the peer |
| `{{.PeerPubkeyShort}}` | A shortened version of the peer's public key |

## Channel Open Event
Triggered when a new channel is opened with your node.

| Variable | Description |
|----------|-------------|
| `{{.PeerAlias}}` | The alias of the peer with whom the channel was opened |
| `{{.PeerPubKey}}` | The full public key of the peer |
| `{{.PeerPubkeyShort}}` | A shortened version of the peer's public key |
| `{{.ChanId}}` | The numeric channel ID |
| `{{.ChannelPoint}}` | The channel point (funding transaction ID and output index) |
| `{{.RemotePubkey}}` | The public key of the remote peer |
| `{{.Capacity}}` | The total capacity of the channel in satoshis (formatted) |

## Channel Close Event
Triggered when a channel is closed.

| Variable | Description |
|----------|-------------|
| `{{.PeerAlias}}` | The alias of the peer with whom the channel was closed |
| `{{.PeerPubKey}}` | The full public key of the peer |
| `{{.PeerPubkeyShort}}` | A shortened version of the peer's public key |
| `{{.ChanId}}` | The numeric channel ID |
| `{{.ChannelPoint}}` | The channel point (funding transaction ID and output index) |
| `{{.RemotePubkey}}` | The public key of the remote peer |
| `{{.Capacity}}` | The total capacity of the channel in satoshis (formatted) |
| `{{.SettledBalance}}` | The final settled balance in satoshis (formatted) |

## Failed HTLC Event
Triggered when an HTLC (Hash Time Locked Contract) fails during routing.

| Variable | Description |
|----------|-------------|
| `{{.InChanId}}` | The ID of the incoming channel |
| `{{.OutChanId}}` | The ID of the outgoing channel |
| `{{.InChanAlias}}` | The alias of the peer on the incoming channel |
| `{{.OutChanAlias}}` | The alias of the peer on the outgoing channel |
| `{{.OutChanLiquidity}}` | The available local balance in the outgoing channel (formatted) |
| `{{.Amount}}` | The amount that was attempted to be forwarded (formatted) |
| `{{.WireFailure}}` | The type of wire failure that occurred |
| `{{.FailureDetail}}` | Detailed description of the failure |
| `{{.MissedFee}}` | The routing fee that was missed due to the failure (formatted) |

## Example Usage

You can use these variables in your notification templates in the config.yaml file. For example:

```yaml
notifications:
  telegram:
    forward_event:
      template: "⚡ New Forward: {{.Amount}} sats\nFee earned: {{.Fee}} sats\nRoute: {{.PeerAliasIn}} -> {{.PeerAliasOut}}"
```

Note: All amounts (satoshis) are automatically formatted with proper separators for readability.