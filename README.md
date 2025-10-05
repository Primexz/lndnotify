# ⚡️ LND Notify

A notification system for Lightning Network nodes that monitors and notifies about important events.

This project is heavily inspired by [balanceofsatoshis](https://github.com/alexbosworth/balanceofsatoshis), but with the aim of offering much greater customisation and a wider range of notification destinations.

![CI](https://img.shields.io/github/actions/workflow/status/primexz/lndnotify/ci.yml)
![License](https://img.shields.io/github/license/primexz/lndnotify)


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

- Running LND node with gRPC access
- LND TLS certificate
- LND macaroon file

## Installation

### 🐳 Run with Docker

#### Docker-Compose

```bash
vim docker-compose.yml
```

```yaml
version: "3.8"
services:
  lndnotify:
    image: ghcr.io/primexz/lndnotify:latest
    container_name: lndnotify
    restart: always
```

### 💻 Run without Docker
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
    forward_event: "💰 Forwarded {{.Amount}} sats, {{.PeerAliasIn}} -> {{.PeerAliasOut}}, earned {{.Fee}} sats"
    peer_online_event: "✅ Peer {{.PeerAlias}} ({{.PeerPubkeyShort}}) is online"
    peer_offline_event: "⚠️ Peer {{.PeerAlias}} ({{.PeerPubkeyShort}}) went offline"
    channel_open_event: "🚀 Channel opened with {{.PeerAlias}}, capacity {{.Capacity}} sats"
    channel_close_event: "🔒 Channel closed with {{.PeerAlias}}, settled balance {{.SettledBalance}} sats"

# Event settings
events:
  forward_events: true
  peer_events: true
  channel_events: true
```


### Notification Providers

The program uses [shoutrrr](https://github.com/nicholas-fedor/shoutrrr) for notifications, which supports various services:

- Discord: `discord://token@channel`
- Telegram: `telegram://token@telegram?channels=channel-1`
- Slack: `slack://token@channel`
- Generic Webhook: `generic://example.com/webhook`

To see the full list of supported providers, check out the [official documentation](https://shoutrrr.nickfedor.com/v0.10.1/services/overview/).

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