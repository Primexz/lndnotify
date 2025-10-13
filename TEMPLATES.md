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

## Channel Opening Event
Triggered when a channel is in the process of being opened (pending).

| Variable | Description |
|----------|-------------|
| `{{.PeerAlias}}` | The alias of the peer with whom the channel is being opened |
| `{{.PeerPubKey}}` | The full public key of the peer |
| `{{.PeerPubkeyShort}}` | A shortened version of the peer's public key |
| `{{.ChannelPoint}}` | The channel point (funding transaction ID and output index) |
| `{{.Capacity}}` | The total capacity of the channel in satoshis (formatted) |
| `{{.Initiator}}` | Boolean indicating if the channel was initiated by your node |
| `{{.IsPrivate}}` | Boolean indicating if this is a private channel |

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

## Channel Closing Event
Triggered when a channel close is in progress (waiting for confirmation).
It is sent when a cooperative or force close is initiated. Force closes are only sent when
they are initiated by your node.

| Variable | Description |
|----------|-------------|
| `{{.PeerAlias}}` | The alias of the peer with whom the channel is being closed |
| `{{.PeerPubKey}}` | The full public key of the peer |
| `{{.PeerPubkeyShort}}` | A shortened version of the peer's public key |
| `{{.ChannelPoint}}` | The channel point (funding transaction ID and output index) |
| `{{.Capacity}}` | The total capacity of the channel in satoshis (formatted) |
| `{{.LimboBalance}}` | The balance in satoshis encumbered in this pending close (formatted) |
| `{{.ClosingTxid}}` | The transaction ID of the closing transaction |
| `{{.ClosingTxHex}}` | The full hex of the closing transaction |

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
| `{{.CloseInitiator}}` | Boolean indicating if the channel close was initiated by your node |
| `{{.CloseType}}` | Integer indicating the type of close: 0=Cooperative, 1=Local Force, 2=Remote Force, 3=Breach, 4=Funding Canceled, 5=Abandoned |

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

For rebalancing payments, the last hop (which is always your own node) is excluded from the `HopInfo` list.

## On-Chain Transaction Event
Triggered when an on-chain transaction is detected involving your node's wallet.

| Variable | Description |
|----------|-------------|
| `{{.TxHash}}` | The transaction hash (txid) of the on-chain transaction |
| `{{.RawTxHex}}` | The raw transaction data in hexadecimal format |
| `{{.Amount}}` | The net amount of the transaction in satoshis (formatted) |
| `{{.TotalFees}}` | The total fees paid for the transaction in satoshis (formatted) |
| `{{.Outputs}}` | List of transaction outputs (see below) |

### Transaction Output Information ({{.Outputs}})
Each output in the transaction contains:

| Variable | Description |
|----------|-------------|
| `{{.Amount}}` | The amount sent to this output in satoshis (formatted) |
| `{{.Address}}` | The destination address for this output |
| `{{.OutputType}}` | The type of output |
| `{{.IsOurAddress}}` | Boolean indicating if this address belongs to your wallet |

## Chain Sync Lost Event
Triggered when your node loses chain synchronization with the Bitcoin network.

| Variable | Description |
|----------|-------------|
| `{{.Duration}}` | The duration for which the chain sync was lost |

## Chain Sync Restored Event
Triggered when your node regains chain synchronization with the Bitcoin network after being out of sync.

| Variable | Description |
|----------|-------------|
| `{{.Duration}}` | The duration for which the chain sync was lost before being restored |

## Channel Backup (Multi) Event
Triggered when a new multi-channel backup is created. This event includes a file attachment.

| Variable | Description |
|----------|-------------|
| `{{.ChanPoints}}` | A list of channel points included in the backup. |
| `{{.NumChanPoints}}` | The total number of channel points in the backup. |
| `{{.Filename}}` | The filename of the backup file. |
| `{{.Sha256Sum}}` | The SHA256 checksum of the backup file. |

## Example Usage

You can use these variables in your notification templates in the config.yaml file. For example:

```yaml
notifications:
  telegram:
    forward_event:
      template: "⚡ New Forward: {{.Amount}} sats\nFee earned: {{.Fee}} sats\nRoute: {{.PeerAliasIn}} -> {{.PeerAliasOut}}"
```

Note: All amounts (satoshis) are automatically formatted with proper separators for readability.