# LND Notify

A notification system for Lightning Network nodes that monitors and notifies about important events.
This project is heavily inspured by balanceofsatoshis(https://github.com/alexbosworth/balanceofsatoshis), but with the aim of offering much greater customisation and a wider range of notification destinations.

## Features

- Real-time monitoring of LND node events
- Configurable notifications for:
  - Payment forward
  - Channel Open
  - Channel Close
  - Peer Online
  - Peer Offline
- Multiple notification providers support via [shoutrrr](https://github.com/nicholas-fedor/shoutrrr)
- Customizable message templates
- Event filtering

## Prerequisites

- Go 1.25 or later
- Running LND node with gRPC access
- LND TLS certificate
- LND macaroon file

## Installation

```bash
go install github.com/Primexz/lndnotify@latest
```

## Configuration

Create a configuration file `config.yaml`:

```yaml
# LND connection settings
lnd:
  host: "localhost"
  port: 10009
  tls_cert_path: "~/.lnd/tls.cert"
  macaroon_path: "~/.lnd/data/chain/bitcoin/mainnet/admin.macaroon"

# Notification settings
notifications:
  providers:
    - url: "discord://token@channel"  # Discord webhook URL
      name: "main-discord"
  templates:
    forward_event: "💸 Forwarded {{.Amount}} sats, earned {{.Fee}} sats"
    peer_online_event: "✅ Peer {{.PeerAlias}} is online"
    peer_offline_event: "⚠️ Peer {{.PeerAlias}} went offline"
    channel_open_event: "🚀 Channel opened with {{.PeerAlias}}, capacity {{.Capacity}} sats"
    channel_close_event: "🔒 Channel closed with {{.PeerAlias}}, settled balance {{.SettledBalance}} sats"

# Event settings
events:
  forward_events: true
  peer_events: true
  channel_events: true

# Rate limiting settings
rate_limiting:
  max_notifications_per_minute: 60
  batch_similar_events: true
  batch_window_seconds: 30
```


### Notification Providers

The program uses [shoutrrr](https://github.com/nicholas-fedor/shoutrrr) for notifications, which supports various services:

- Discord: `discord://token@channel`
- Telegram: `telegram://token@telegram?channels=channel-1`
- Slack: `slack://token@channel`
- Generic Webhook: `generic://example.com/webhook`

To see the full list of supported providers, check out the [official list](https://github.com/nicholas-fedor/shoutrrr#supported-services).

## Usage

```bash
lndnotify -config config.yaml
```

## Development

### Building from Source

```bash
git clone https://github.com/Primexz/lndnotify.git
cd lndnotify
go build -o lndnotify cmd/lndnotify/main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.