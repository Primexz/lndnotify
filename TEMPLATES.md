# Template Variables Documentation

This document describes all available template variables that can be used in notification templates for each event type in lndnotify.

Templates use Go's `text/template` syntax. For more information on template syntax, functions, and advanced usage, see the [official Go documentation](https://pkg.go.dev/text/template).

## Forward Event
Triggered when a payment is forwarded through your node.

| Variable | Description |
|----------|-------------|
| `{{.PeerAliasIn}}` | The alias of the peer that sent the payment to your node |
| `{{.PeerAliasOut}}` | The alias of the peer that received the payment from your node |
| `{{.Amount}}` | The incoming amount of the forward in satoshis (formatted) |
| `{{.AmountOut}}` | The outgoing amount of the forward in satoshis (formatted) |
| `{{.Fee}}` | The fee earned from forwarding the payment in satoshis (formatted) |
| `{{.FeeRate}}` | The fee rate in ppm earned from forwarding the payment (formatted) |

## Invoice Settled Event
Triggered when an invoice is paid/settled.

| Variable | Description |
|----------|-------------|
| `{{.Memo}}` | The memo/description attached to the invoice |
| `{{.Value}}` | The value of the invoice in satoshis (formatted) |
| `{{.IsKeysend}}` | Boolean indicating if this was a keysend payment |
| `{{.PaymentRequest}}` | The original payment request string (invoice) |

## Keysend Event
Triggered when a keysend payment is received through your node.

| Variable | Description |
|----------|-------------|
| `{{.Msg}}` | The message attached to the keysend payment |
| `{{.InChanAlias}}` | The alias of the peer on the incoming channel |
| `{{.InChanId}}` | The ID of the incoming channel |
| `{{.Amount}}` | The amount of the keysend payment in satoshis (formatted) |

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

## Payment Succeeded Event
Triggered when an outgoing payment is successfully completed.

| Variable | Description |
|----------|-------------|
| `{{.PaymentHash}}` | The payment hash of the completed payment |
| `{{.Amount}}` | The total amount of the payment in satoshis (formatted) |
| `{{.Fee}}` | The total fee paid for the payment in satoshis (formatted) |
| `{{.FeeRate}}` | The total fee rate of the payment in ppm |
| `{{.Receiver}}` | The alias of the receiving node (final destination) |
| `{{.Memo}}` | The memo/description from the payment request |
| `{{.HtlcInfo}}` | List of HTLC information (see below) |

### HTLC Information ({{.HtlcInfo}})
Each HTLC in the list contains:

| Variable | Description |
|----------|-------------|
| `{{.FirstHop}}` | The alias of the first hop in this HTLC route |
| `{{.PenultHop}}` | The alias of the penultimate (second-to-last) hop |
| `{{.Amount}}` | The amount sent via this HTLC in satoshis (formatted) |
| `{{.Fee}}` | The fee paid for this HTLC in satoshis (formatted) |
| `{{.FeeRate}}` | The fee rate for this HTLC in ppm |
| `{{.HopInfo}}` | List of individual hop information (see below) |

### Hop Information ({{.HopInfo}})
Each hop in an HTLC route contains:

| Variable | Description |
|----------|-------------|
| `{{.Pubkey}}` | The public key of this hop |
| `{{.Alias}}` | The alias of this hop |
| `{{.Amount}}` | The amount forwarded to this hop in satoshis (formatted) |
| `{{.Fee}}` | The fee paid to this hop in satoshis (formatted) |
| `{{.FeeRate}}` | The fee rate for this hop in ppm |

## Example Usage

You can use these variables in your notification templates in the config.yaml file. For example:

```yaml
notifications:
  telegram:
    forward_event:
      template: "⚡ New Forward: {{.Amount}} sats\nFee earned: {{.Fee}} sats\nRoute: {{.PeerAliasIn}} -> {{.PeerAliasOut}}"
```

Note: All amounts (satoshis) are automatically formatted with proper separators for readability.