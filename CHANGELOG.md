# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## [Unreleased]
### Added
### Deprecated
### Fixed
### Changed
### Removed


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