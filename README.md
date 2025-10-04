# LND Notify

A notification system for Lightning Network nodes that monitors and notifies about important events

## Features

- Real-time monitoring of LND node events
- Configurable notifications for:
  - Payment forwardings
- Multiple notification providers support via shoutrrr
- Customizable message templates
- Event filtering and rate limiting
- Environment variable configuration support

## Prerequisites

- Go 1.25 or later
- Running LND node with gRPC access
- LND TLS certificate
- LND macaroon file

## Installation

```bash
go install github.com/Primexz/lnd-notify@latest
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
    forward: "💸 Forwarded {{.Amount}} sats, earned {{.Fee}} sats"

# Event settings
events:
  forward_events: true

# Rate limiting settings
rate_limiting:
  max_notifications_per_minute: 60
  batch_similar_events: true
  batch_window_seconds: 30
```

### Environment Variables

You can also configure the program using environment variables:

```bash
export LND_HOST=localhost
export LND_PORT=10009
export LND_TLS_CERT_PATH=/path/to/tls.cert
export LND_MACAROON_PATH=/path/to/admin.macaroon
export NOTIFICATION_URL="discord://token@channel"
export ENABLED_EVENTS="forwards"
```

### Notification Providers

The program uses shoutrrr for notifications, which supports various services:

- Discord: `discord://token@channel`
- Telegram: `telegram://token@telegram?channels=channel-1`
- Slack: `slack://token@channel`
- Generic Webhook: `generic://example.com/webhook`

## Usage

```bash
# Using config file
lnd-notify -config config.yaml

# Using environment variables
lnd-notify
```

## Development

### Building from Source

```bash
git clone https://github.com/Primexz/lnd-notify.git
cd lnd-notify
go mod download
go build -o lnd-notify cmd/lnd-notify/main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.