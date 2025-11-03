# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## [Unreleased]
### Added
### Fixed
### Changed
### Removed
### Deprecated


## [1.11.0] - 2025-11-03

### Added
- Added channel fee change event: Notifications for changes in channel fee policies including inbound fee rates and base fees. (@Primexz)
- Added block explorer URL in on-chain transaction events. (@Primexz)

### Changed
- Updated [shoutrrr](https://github.com/nicholas-fedor/shoutrrr/releases/tag/v0.11.1)


## [1.10.1] - 2025-10-29

### Fixed
- Fixed issue where formatting of satoshi amounts with decimal places was incorrect in some locales. (@Primexz)


## [1.10.0] - 2025-10-28

### Added
- Added LND update available event. (@Primexz)

### Changed
- Added **dev** docker image. (@Primexz)


## [1.9.0] - 2025-10-26

### Added
- Added LND health event. (@Primexz)

### Fixed
- Fixed issue where events were not being sent correctly. (@Primexz)

### Changed
- Skip keysend invoices by default in invoice settled event. This can be configured via `event_config.invoice_event.skip_keysend` option. (@Primexz)
- Added VSCode code format settings. (@Primexz)


## [1.8.0] - 2025-10-24

### Added
- Added wallet state event. (@Primexz)

### Changed
- Reduced binary size with [UPX compression](https://github.com/upx/upx) for Linux builds. (@Primexz)


## [1.7.0] - 2025-10-21

### Added
- Added TLS certificate expiry event: Notifications for upcoming LND TLS certificate expiration with configurable threshold before notification. (@Primexz)
- Notification batching: Support for batching notifications to reduce frequency while ensuring important events are delivered promptly. Configurable flush interval and maximum batch size. (@Primexz)
- Added fee rate (ppm) to forward event template. (@Primexz)


## [1.6.3] - 2025-10-19

### Added
- MacOS binary build in release process (@Primexz)

### Fixed
- Prevented multiple notifications from being sent when lndClient.SubscribeEvents fails (e.g., during lnd startup when server is not active) and lndnotify is retried externally. (@Feelancer21)
- Fixed issue where channel up events and chain sync lost events were incorrectly triggered. (@Primexz)

### Changed

- Updated [shoutrrr](https://github.com/nicholas-fedor/shoutrrr/releases/tag/v0.10.3)
- Updated Golang to version 1.25.3 (@Primexz)

## [1.6.2] - 2025-10-18

### Changed
- Improved notification templates (@Primexz)
- Fixed notification template for backup event (@Primexz)


## [1.6.1] - 2025-10-15

### Fixed
- Fixed issue where channel up events were incorrectly triggered without a corresponding channel down event being sent previously. (@Primexz)


## [1.6.0] - 2025-10-15

### Added
- Channel status event: Notifications for channel up/down events with configurable minimum downtime before notification. (@Primexz)

### Deprecated
- Peer events: Peer online/offline events are deprecated and will be removed in a future release. Please use channel status events instead. (@Primexz)

### Fixed
- Fixed issue where chain sync lost event was getting triggered incorrectly. (@Primexz)


## [1.5.0] - 2025-10-14

### Added
- New Event: Chain sync lost notification (@Primexz)
- New event: multi-channel backup notification. Backup is uploaded as file if ntfy is used as provider (@Feelancer21)


## [1.4.0] - 2025-10-12

### Added
- New event: opening channel notification (@Feelancer21)
- New event: closing channel notification (@Feelancer21)
- Configuration option for number formatting locale (e.g., "en-US" for English, "de-DE" for German) (@Primexz)
- New template variables `{{.CloseInitiator}}` and `{{.CloseType}}` for `channel_close_event`. (@Feelancer21)


## [1.3.1] - 2025-10-11

### Fixed
- LND Connection attempts are now performed indefinitely (Previously, it was automatically abandoned after 15 minutes) (@Primexz)


## [1.3.0] - 2025-10-11

### Added
- New event: on chain transaction (@Primexz)


## [1.2.1] - 2025-10-09

### Added
- ``--version`` command to get the current version of the app (@Primexz)


## [1.2.0] - 2025-10-09

### Added
- New event: payment succeeded and rebalancing succeeded notifications (@feelancer21)
- New event: keysend notifications (@Primexz)
- Filter for `forward_min_amount` and `invoice_min_amount` to skip events with
low values (@feelancer21)
- New template variables `{{.FeeRate}}` for `forward_event`. See
[TEMPLATES.md](TEMPLATES.md) for details. (@feelancer21)

### Fixed

- Peer event: use pubkey as fallback (@Primexz)

### Changed

- Internal: new ppm rate formatting (@feelancer21)