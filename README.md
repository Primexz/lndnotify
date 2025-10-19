# ⚡️ LND Notify

A notification system for Lightning Network nodes that monitors and notifies about important events.

This project is heavily inspired by [balanceofsatoshis](https://github.com/alexbosworth/balanceofsatoshis), but with the aim of offering much greater customisation and a wider range of notification destinations.

![CI](https://img.shields.io/github/actions/workflow/status/primexz/lndnotify/ci.yml)
![License](https://img.shields.io/github/license/primexz/lndnotify)


## Features

- Real-time monitoring of LND node events
- Configurable notifications for:
  - Payment forward
  - Channel Opening (pending)
  - Channel Open
  - Channel Closing (pending)
  - Channel Close
  - Channel Backup (Multi) (File uploads only supported via ntfy at the moment)
  - Channel Status Change (Up/Down)
  - Settled Invoice
  - Failed HTLCs
  - Payment Succeeded
  - Rebalancing Succeeded 
  - On-Chain Transactions
  - On-Chain Sync Lost
- Multiple notification providers support via [shoutrrr](https://github.com/nicholas-fedor/shoutrrr)
- Customizable message templates ([see all template variables](TEMPLATES.md))
- Customizable notification formatting (e.g., number formatting based on locale)
- Optional notification batching with configurable intervals
- Event filtering

## Prerequisites

- Running LND node with gRPC access
- LND TLS certificate
- LND readonly macaroon file
- Basic Understanding of Linux, Docker and Docker Compose

## Installation

### 🐳 Run with Docker

#### Docker-Compose

```bash
vim docker-compose.yml
```

```yaml
services:
  lndnotify:
    image: ghcr.io/primexz/lndnotify:latest
    container_name: lndnotify
    volumes:
      - HOST_LND_PATH:/root/.lnd:ro
      - ./lndnotify/config.yaml:/data/config.yaml
    command: -config /data/config.yaml
    networks:
      - LND_NETWORK
    restart: always
```

- Adjust ``HOST_LND_PATH``and ``LND_NETWORK``
- Add ``config.yaml`` file to the ``lndnotify`` directory and adjust the configuration

### 💻 Run without Docker
```bash
git clone https://github.com/Primexz/lndnotify.git
cd lndnotify
go build .
./lndnotify
```

## Configuration

> [!TIP] 
> The full configuration options can be found in the [example config file](config.example.yaml).

Create a basic configuration file `config.yaml`:

```yaml
# LND connection settings
lnd:
  host: "localhost"
  port: 10009
  tls_cert_path: "~/.lnd/tls.cert"
  macaroon_path: "~/.lnd/data/chain/bitcoin/mainnet/readonly.macaroon"

# Notification settings
notifications:
  providers:
    - url: "discord://token@channel?SplitLines=false"  # Discord webhook URL
      name: "main-discord"
  formatting:
    locale: "en-US"  # Language for number formatting (e.g. "en-US" for English, "de-DE" for German)

# Event settings
events:
  backup_events: true
  channel_events: true
  failed_htlc_events: true
  forward_events: true
  invoice_events: true
  keysend_events: true
  payment_events: true
  rebalancing_events: true
  status_events: true
  on_chain_events: true
  chain_sync_events: true
  channel_status_events: true
  tls_cert_expiry_events: true
  peer_events: false

# Event-specific configuration
event_config:
  failed_htlc_event:
    min_amount: 0
  forward_event:
    min_amount: 0
  invoice_event:
    min_amount: 0
  payment_event:
    min_amount: 0
  rebalancing_event:
    min_amount: 0
  on_chain_event:
    min_amount: 0
  chain_lost_event:
    threshold: 5m 
    warning_interval: 15m 
  channel_status_event:
    min_downtime: 10m

```

### Notification Batching

LND Notify supports batching notifications to reduce the frequency of messages while ensuring important events are still delivered promptly. This is particularly useful for high-traffic nodes that might generate many notifications.

When batching is enabled, notifications are collected and sent together based on two criteria:
- **Time-based flushing**: Notifications are sent after a configurable interval (default: 5 seconds)
- **Size-based flushing**: Notifications are sent immediately when the batch reaches a maximum size (default: 10 notifications)

#### Configuration

```yaml
notifications:
  batching:
    enabled: true  # Enable notification batching
    flush_interval: "10s"  # Send batched notifications every 10 seconds
    max_size: 15  # Send immediately when 15 notifications are queued
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
go install ./cmd/lndnotify
```

## Contributing

> [!NOTE] 
> This project uses [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for commit messages. Please follow the guidelines when contributing.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the [MIT License](LICENSE).